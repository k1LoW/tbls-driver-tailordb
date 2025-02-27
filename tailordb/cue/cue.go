package cue

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	cueload "cuelang.org/go/cue/load"
	"github.com/k1LoW/tbls/schema"
)

const kindTailordb = "tailordb"

type Type struct {
	Name        string  `json:"Name"`
	Description string  `json:"Description"`
	Indexes     Indexes `json:"Indexes,omitempty"`
	Fields      Fields  `json:"Fields"`
}

type Field struct {
	Type            string `json:"Type"`
	Description     string `json:"Description"`
	SourceId        string `json:"SourceId,omitempty"`
	ForeignKey      bool   `json:"ForeignKey,omitempty"`
	ForeignKeyType  string `json:"ForeignKeyType,omitempty"`
	ForeignKeyField string `json:"ForeignKeyField,omitempty"`
	Required        bool   `json:"Required,omitempty"`
	Array           bool   `json:"Array,omitempty"`
	Index           bool   `json:"Index,omitempty"`
	Unique          bool   `json:"Unique,omitempty"`
	Fields          Fields `json:"Fields,omitempty"`
}

type Fields map[string]*Field

type Index struct {
	FieldNames []string `json:"FieldNames"`
	Unique     bool     `json:"Unique"`
}

type Indexes map[string]*Index

func LookupModRoot(root string) (string, error) {
	if fi, err := os.Stat(path.Join(root, "cue.mod")); err == nil && fi.IsDir() {
		return root, nil
	}
	if root == "." || root == "/" {
		return "", fmt.Errorf("module root not found")
	}
	return LookupModRoot(path.Dir(root))
}

func Analyze(root string) (_ *schema.SchemaJSON, err error) {
	v, err := load(root)
	if err != nil {
		return nil, err
	}
	return analyze(v)
}

func load(root string) (cue.Value, error) {
	modRoot, err := LookupModRoot(root)
	if err != nil {
		return cue.Value{}, err
	}

	ctx := cuecontext.New()
	cfg := &cueload.Config{
		ModuleRoot: modRoot,
	}
	insts := cueload.Instances([]string{root}, cfg)
	v := ctx.BuildInstance(insts[0])
	return detectTailordb(v)
}

func detectTailordb(v cue.Value) (cue.Value, error) {
	services := v.LookupPath(cue.MakePath(cue.Str("Services")))
	if !services.Exists() {
		kind, err := v.LookupPath(cue.MakePath(cue.Str("Kind"))).String()
		if err != nil {
			return cue.Value{}, err
		}
		if kind != kindTailordb {
			return cue.Value{}, fmt.Errorf("no tailordb found")
		}
		return v, nil
	}
	var databases []cue.Value
	iter, err := services.List()
	if err != nil {
		return cue.Value{}, err
	}
	for iter.Next() {
		value := iter.Value()
		kindValue := value.LookupPath(cue.MakePath(cue.Str("Kind")))
		if kindValue.Exists() {
			k, err := kindValue.String()
			if err != nil {
				return cue.Value{}, err
			}
			if k == kindTailordb {
				databases = append(databases, value)
			}
		}
	}
	switch len(databases) {
	case 0:
		return cue.Value{}, fmt.Errorf("no tailordb found")
	case 1:
	default:
		return cue.Value{}, fmt.Errorf("multiple tailordb found")
	}
	database := databases[0]
	return database, nil
}

func analyze(v cue.Value) (_ *schema.SchemaJSON, err error) {
	s := &schema.SchemaJSON{}
	s.Name, err = v.LookupPath(cue.MakePath(cue.Str("Namespace"))).String()
	if err != nil {
		return nil, err
	}
	vv := v.LookupPath(cue.MakePath(cue.Str("Types")))
	if !vv.Exists() {
		return nil, fmt.Errorf("no types found")
	}

	typesIter, err := v.LookupPath(cue.MakePath(cue.Str("Types"))).List()
	if err != nil {
		return nil, err
	}

	// tables
	for typesIter.Next() {
		v := typesIter.Value()
		typ := &Type{}
		if err := v.Decode(&typ); err != nil {
			return nil, err
		}

		t := &schema.TableJSON{
			Type: "TailorDB.Type",
		}
		t.Name = typ.Name
		t.Comment = typ.Description

		// indexes
		{
			b, err := json.Marshal(typ.Indexes)
			if err != nil {
				return nil, err
			}
			def := string(b)
			for name, idx := range typ.Indexes {
				t.Indexes = append(t.Indexes, &schema.Index{
					Name:    name,
					Def:     def,
					Table:   &t.Name,
					Columns: idx.FieldNames,
				})
				if idx.Unique {
					t.Constraints = append(t.Constraints, &schema.Constraint{
						Name:    name,
						Type:    "UNIQUE",
						Def:     def,
						Table:   &t.Name,
						Columns: idx.FieldNames,
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

		// Add fields
		columns, err := analyzeFields(s, t, typ.Fields, "")
		if err != nil {
			return nil, err
		}
		fieldsIter, err := v.LookupPath(cue.MakePath(cue.Str("Fields"))).Fields()
		if err != nil {
			return nil, err
		}
		fieldsLabels := fieldsOrder(fieldsIter)
		for _, label := range fieldsLabels {
			for _, c := range columns {
				if c.Name == label {
					t.Columns = append(t.Columns, c)
					break
				}
			}
		}

		// sort workarounds
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

	// sort workarounds
	sort.Slice(s.Relations, func(i, j int) bool {
		if s.Relations[i].Def == s.Relations[j].Def {
			if s.Relations[i].Table == s.Relations[j].Table {
				return s.Relations[i].Columns[0] < s.Relations[j].Columns[0]
			}
			return s.Relations[i].Table < s.Relations[j].Table
		}
		return s.Relations[i].Def < s.Relations[j].Def
	})

	return s, nil
}

func analyzeFields(s *schema.SchemaJSON, t *schema.TableJSON, fields Fields, prefix string) ([]*schema.ColumnJSON, error) {
	var columns []*schema.ColumnJSON
	for name, field := range fields {
		if prefix != "" {
			name = fmt.Sprintf("%s.%s", prefix, name)
		}
		c := &schema.ColumnJSON{}

		c.Name = name
		c.Type = field.Type
		if field.Array {
			c.Type = fmt.Sprintf("Array\\<%s\\>", field.Type)
		}
		c.Comment = field.Description
		c.Nullable = !field.Required
		columns = append(columns, c)

		if field.Type == "nested" {
			nestedColumns, err := analyzeFields(s, t, field.Fields, name)
			if err != nil {
				return nil, err
			}
			columns = append(columns, nestedColumns...)
		}

		if field.SourceId != "" {
			sourceName := field.SourceId
			if prefix != "" {
				sourceName = fmt.Sprintf("%s.%s", prefix, sourceName)
			}

			parentTable := field.Type
			rel := &schema.RelationJSON{
				Table:       t.Name,
				Columns:     []string{sourceName, name},
				ParentTable: parentTable,
			}
			// Check sourceId
			source, ok := fields[field.SourceId]
			if !ok {
				return nil, fmt.Errorf("source %s not found", sourceName)
			}
			if source.ForeignKey {
				if source.ForeignKeyField != "" {
					rel.ParentColumns = []string{source.ForeignKeyField}
				} else {
					rel.ParentColumns = []string{"id"}
				}
				rel.Def = fmt.Sprintf("ForeignKeyType: %s", source.ForeignKeyType)

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

		switch {
		case field.Index && !field.Unique:
			t.Indexes = append(t.Indexes, &schema.Index{
				Name:    fmt.Sprintf("Index for %s", name),
				Def:     "Index: true",
				Table:   &t.Name,
				Columns: []string{name},
			})
		case field.Unique:
			t.Indexes = append(t.Indexes, &schema.Index{
				Name:    fmt.Sprintf("Unique for %s", name),
				Def:     "Unique: true / Index: true",
				Table:   &t.Name,
				Columns: []string{name},
			})
			t.Constraints = append(t.Constraints, &schema.Constraint{
				Name:    fmt.Sprintf("Unique for %s", name),
				Type:    "UNIQUE",
				Def:     "Unique: true / Index: true",
				Table:   &t.Name,
				Columns: []string{name},
			})
		}
	}
	return columns, nil
}

func fieldsOrder(fieldsIter *cue.Iterator) []string {
	var fieldsLabels []string
	for fieldsIter.Next() {
		fieldsLabels = append(fieldsLabels, fieldsIter.Selector().String())
		v := fieldsIter.Value()
		// nested fields
		if iter, err := v.LookupPath(cue.MakePath(cue.Str("Fields"))).Fields(); err == nil {
			for _, label := range fieldsOrder(iter) {
				fieldsLabels = append(fieldsLabels, fmt.Sprintf("%s.%s", fieldsIter.Selector().String(), label))
			}
		}
	}
	return fieldsLabels
}
