package cue

import (
	"encoding/json"
	"fmt"

	"cuelang.org/go/cue"
	"github.com/k1LoW/tbls/schema"
)

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

func Analyze(v cue.Value) (_ *schema.SchemaJSON, err error) {
	s := &schema.SchemaJSON{}
	s.Name, err = v.LookupPath(cue.MakePath(cue.Str("Namespace"))).String()
	if err != nil {
		return nil, err
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
		t.Columns = append(t.Columns, columns...)

		s.Tables = append(s.Tables, t)
	}
	return s, nil
}

func analyzeFields(s *schema.SchemaJSON, t *schema.TableJSON, fields Fields, prefix string) ([]*schema.ColumnJSON, error) {
	var columns []*schema.ColumnJSON
	for name, field := range fields {
		if prefix != "" {
			name = fmt.Sprintf("%s.%s", prefix, name)
		}
		// TODO: Fix order
		// TODO: Array field, nested field
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
			sourceIdName := field.SourceId
			if prefix != "" {
				sourceIdName = fmt.Sprintf("%s.%s", prefix, sourceIdName)
			}

			parentTable := c.Type
			rel := &schema.RelationJSON{
				Table:       t.Name,
				Columns:     []string{sourceIdName, name},
				ParentTable: parentTable,
			}
			// Check sourceId
			sourceId, ok := fields[field.SourceId]
			if !ok {
				return nil, fmt.Errorf("sourceId %s not found", sourceIdName)
			}
			if sourceId.ForeignKey {
				if sourceId.ForeignKeyField != "" {
					rel.ParentColumns = []string{sourceId.ForeignKeyField}
				} else {
					rel.ParentColumns = []string{"id"}
				}
				rel.Def = fmt.Sprintf("ForeignKeyType: %s", sourceId.ForeignKeyType)

				t.Constraints = append(t.Constraints, &schema.Constraint{
					Name:    fmt.Sprintf("ForeignKey for %s to %s", name, parentTable),
					Type:    "FOREIGN KEY",
					Def:     fmt.Sprintf("ForeignKeyType: %s", parentTable),
					Table:   &t.Name,
					Columns: []string{sourceIdName, name},
				})

			} else {
				rel.ParentColumns = []string{"id"}
				rel.Def = fmt.Sprintf("Source: %s", parentTable)
			}
			s.Relations = append(s.Relations, rel)
		}

		if field.Index {
			t.Indexes = append(t.Indexes, &schema.Index{
				Name:    fmt.Sprintf("Index for %s", c.Name),
				Def:     "Index: true",
				Table:   &t.Name,
				Columns: []string{c.Name},
			})
		}

		if field.Unique {
			t.Indexes = append(t.Indexes, &schema.Index{
				Name:    fmt.Sprintf("Unique for %s", c.Name),
				Def:     "Unique: true",
				Table:   &t.Name,
				Columns: []string{c.Name},
			})
			t.Constraints = append(t.Constraints, &schema.Constraint{
				Name:    fmt.Sprintf("Unique for %s", c.Name),
				Type:    "UNIQUE",
				Def:     "Unique: true",
				Table:   &t.Name,
				Columns: []string{c.Name},
			})
		}
	}
	return columns, nil
}
