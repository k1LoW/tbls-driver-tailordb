package tf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/k1LoW/tbls/schema"
	"github.com/zclconf/go-cty/cty"
)

const providerTailordb = "tailor"
const typeTailorWorkspace = "tailor_workspace"
const typeTailordb = "tailor_tailordb"
const typeTailordbType = "tailor_tailordb_type"

type Types struct {
	Types []Type `hcl:"resource,block"`
}

type Type struct {
	Name string `hcl:"name,label"`
}

func Analyze(dir string) (*schema.SchemaJSON, error) {
	parser := hclparse.NewParser()
	if err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".tf" {
			return nil
		}
		_, diags := parser.ParseHCLFile(path)
		if diags.HasErrors() {
			return diags
		}
		return nil
	}); err != nil {
		return nil, err
	}

	v := newVariables()

	for _, file := range parser.Files() {
		content, _, diags := file.Body.PartialContent(rootSchema)
		if diags.HasErrors() {
			return nil, diags
		}
		for _, block := range content.Blocks {
			var attrs hcl.Attributes
			labels := block.Labels
			switch block.Type {
			case "terraform":
				content, _, diags := block.Body.PartialContent(terraformSchema)
				if diags.HasErrors() {
					return nil, diags
				}
				attrs = content.Attributes
			case "resource", "data":
				rtype := block.Labels[0]
				name := block.Labels[1]
				switch rtype {
				case typeTailorWorkspace:
					content, _, diags := block.Body.PartialContent(workspaceSchema)
					if diags.HasErrors() {
						return nil, diags
					}
					// Add workspace id (Read-only)
					if err := v.set([]string{rtype, name, "id"}, cty.StringVal("[read-only]")); err != nil {
						return nil, err
					}
					attrs = content.Attributes
				case typeTailordb:
					content, _, diags := block.Body.PartialContent(tailordbSchema)
					if diags.HasErrors() {
						return nil, diags
					}
					attrs = content.Attributes
				case typeTailordbType:
					content, _, diags := block.Body.PartialContent(tailordbTypeSchema)
					if diags.HasErrors() {
						return nil, diags
					}
					attrs = content.Attributes
				default:
					attrs, diags = block.Body.JustAttributes()
					if diags.HasErrors() && !strings.Contains(diags.Error(), "Blocks are not allowed here") {
						return nil, diags
					}
				}
			case "locals":
				attrs, diags = block.Body.JustAttributes()
				if diags.HasErrors() {
					return nil, diags
				}
				labels = []string{"local"}
			default:
				attrs, diags = block.Body.JustAttributes()
				if diags.HasErrors() {
					return nil, diags
				}
			}
			for _, attr := range attrs {
				if len(attr.Expr.Variables()) > 0 {
					v.setExpr(append(labels, attr.Name), attr.Expr)
					continue
				}
				vv, diags := attr.Expr.Value(nil)
				if diags.HasErrors() {
					if strings.Contains(diags.Error(), "Functions may not be called here") {
						vv = cty.StringVal("[function]")
					} else {
						return nil, diags
					}
				}
				if err := v.set(append(labels, attr.Name), vv); err != nil {
					return nil, err
				}
			}
		}
	}

	ctx, err := v.ctx()
	if err != nil {
		return nil, err
	}

	s := &schema.SchemaJSON{}
	for _, file := range parser.Files() {
		content, _, diags := file.Body.PartialContent(rootSchema)
		if diags.HasErrors() {
			return nil, diags
		}
		for _, block := range content.Blocks {
			if block.Type != "resource" && block.Type != "data" {
				continue
			}
			rtype := block.Labels[0]
			if rtype != typeTailordbType {
				continue
			}

			// tables
			content, _, diags := block.Body.PartialContent(tailordbTypeSchema)
			if diags.HasErrors() {
				return nil, diags
			}
			t := &schema.TableJSON{
				Type: "TailorDB.Type",
			}
			if v, ok := content.Attributes["name"]; ok {
				vv, diags := v.Expr.Value(ctx)
				if diags.HasErrors() {
					return nil, diags
				}
				t.Name = vv.AsString()
			}
			if v, ok := content.Attributes["description"]; ok {
				vv, diags := v.Expr.Value(ctx)
				if diags.HasErrors() {
					return nil, diags
				}
				t.Comment = vv.AsString()
			}

			// indexes
			if v, ok := content.Attributes["indexes"]; ok {
				vv, diags := v.Expr.Value(ctx)
				if diags.HasErrors() {
					return nil, diags
				}
				for k, v := range vv.AsValueMap() {
					vvv := v.AsValueMap()
					var fields []string
					if v, ok := vvv["field_name"]; ok {
						vvvv := v.AsValueSlice()
						for _, vv := range vvvv {
							fields = append(fields, vv.AsString())
						}
					}
					var unique bool
					if v, ok := vvv["unique"]; ok && v.True() {
						unique = true
					}
					b, err := json.Marshal(struct {
						FieldNames []string `json:"field_names"`
						Unique     bool     `json:"unique"`
					}{
						FieldNames: fields,
						Unique:     unique,
					})
					if err != nil {
						return nil, err
					}
					def := string(b)
					t.Indexes = append(t.Indexes, &schema.Index{
						Name:    k,
						Def:     def,
						Table:   &t.Name,
						Columns: fields,
					})
					if unique {
						t.Constraints = append(t.Constraints, &schema.Constraint{
							Name:    k,
							Type:    "UNIQUE",
							Def:     def,
							Table:   &t.Name,
							Columns: fields,
						})
					}
				}
			}

			// columns
			// Add id column
			id := &schema.ColumnJSON{
				Name:     "id",
				Type:     "uuid",
				Nullable: false,
			}
			t.Columns = append(t.Columns, id)

			if v, ok := content.Attributes["fields"]; ok {
				vv, diags := v.Expr.Value(ctx)
				if diags.HasErrors() {
					return nil, diags
				}
				columns, err := analyzeFields(ctx, s, t, vv, "")
				if err != nil {
					return nil, err
				}
				t.Columns = append(t.Columns, columns...)
			}

			// sort workaroud
			sort.Slice(t.Columns, func(i, j int) bool {
				// id column first
				if t.Columns[i].Name == "id" {
					return true
				}
				if t.Columns[j].Name == "id" {
					return false
				}
				return t.Columns[i].Name < t.Columns[j].Name
			})

			sort.Slice(t.Indexes, func(i, j int) bool {
				if t.Indexes[i].Def == t.Indexes[j].Def {
					return t.Indexes[i].Name < t.Indexes[j].Name
				}
				return t.Indexes[i].Def < t.Indexes[j].Def
			})
			sort.Slice(t.Constraints, func(i, j int) bool {
				if t.Constraints[i].Def == t.Constraints[j].Def {
					return t.Constraints[i].Name < t.Constraints[j].Name
				}
				return t.Constraints[i].Def < t.Constraints[j].Def
			})

			s.Tables = append(s.Tables, t)
		}
	}

	// sort workaroud
	sort.Slice(s.Tables, func(i, j int) bool {
		return s.Tables[i].Name < s.Tables[j].Name
	})
	sort.Slice(s.Relations, func(i, j int) bool {
		if s.Relations[i].Def == s.Relations[j].Def {
			return s.Relations[i].ParentTable < s.Relations[j].ParentTable
		}
		return s.Relations[i].Def < s.Relations[j].Def
	})
	return s, nil
}

func analyzeFields(ctx *hcl.EvalContext, s *schema.SchemaJSON, t *schema.TableJSON, fields cty.Value, prefix string) ([]*schema.ColumnJSON, error) {
	var columns []*schema.ColumnJSON
	for name, v := range fields.AsValueMap() {
		if prefix != "" {
			name = fmt.Sprintf("%s.%s", prefix, name)
		}
		c := &schema.ColumnJSON{
			Name: name,
		}
		vv := v.AsValueMap()
		var fieldType string
		if v, ok := vv["type"]; ok {
			fieldType = v.AsString()
		}
		c.Type = fieldType
		if v, ok := vv["description"]; ok {
			c.Comment = v.AsString()
		}
		if v, ok := vv["required"]; ok {
			c.Nullable = !v.True()
		}
		columns = append(columns, c)

		if v, ok := vv["fields"]; ok && fieldType == "nested" {
			nestedColumns, err := analyzeFields(ctx, s, t, v, name)
			if err != nil {
				return nil, err
			}
			columns = append(columns, nestedColumns...)
		}

		if v, ok := vv["source"]; ok {
			sourceName := v.AsString()
			if prefix != "" {
				sourceName = fmt.Sprintf("%s.%s", prefix, sourceName)
			}

			parentTable := fieldType
			rel := &schema.RelationJSON{
				Table:       t.Name,
				Columns:     []string{sourceName, name},
				ParentTable: parentTable,
			}
			// Check source
			vv, ok := fields.AsValueMap()[sourceName]
			if !ok {
				return nil, fmt.Errorf("source %s not found", sourceName)
			}
			source := vv.AsValueMap()
			if v, ok := source["foreign_key"]; ok {
				vv := v.AsValueMap()
				if v, ok := vv["field"]; ok {
					rel.ParentColumns = []string{v.AsString()}
				} else {
					rel.ParentColumns = []string{"id"}
				}
				vvv, ok := vv["type"]
				if !ok {
					return nil, fmt.Errorf("foreign_key.type not found: %s", sourceName)
				}
				rel.Def = fmt.Sprintf("ForeignKeyType: %s", vvv.AsString())

				t.Constraints = append(t.Constraints, &schema.Constraint{
					Name:    fmt.Sprintf("ForeignKey for %s to %s", name, parentTable),
					Type:    "FOREIGN KEY",
					Def:     fmt.Sprintf("ForeignKeyType: %s", parentTable),
					Table:   &t.Name,
					Columns: []string{sourceName, name},
				})
			} else {
				rel.ParentColumns = []string{"id"}
				rel.Def = fmt.Sprintf("Source: %s", parentTable)
			}
			s.Relations = append(s.Relations, rel)
		}

		var (
			index  bool
			unique bool
		)
		if v, ok := vv["index"]; ok && v.True() {
			index = true
		}
		if v, ok := vv["unique"]; ok && v.True() {
			unique = true
		}

		switch {
		case index && !unique:
			t.Indexes = append(t.Indexes, &schema.Index{
				Name:    fmt.Sprintf("Index for %s", name),
				Def:     "Index: true",
				Table:   &t.Name,
				Columns: []string{name},
			})
		case unique:
			t.Indexes = append(t.Indexes, &schema.Index{
				Name:    fmt.Sprintf("Unique for %s", c.Name),
				Def:     "Unique: true / Index: true",
				Table:   &t.Name,
				Columns: []string{c.Name},
			})
		}
	}

	return columns, nil
}
