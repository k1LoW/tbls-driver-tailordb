package cue

type Type struct {
	Name        string  `json:"Name"`
	Description string  `json:"Description"`
	Indexes     Indexes `json:"Indexes,omitempty"`
	Fields      Fields  `json:"Fields"`
}

type Field struct {
	Type            string `json:"Type"`
	Description     string `json:"Description"`
	Required        bool   `json:"Required,omitempty"`
	SourceId        string `json:"SourceId,omitempty"`
	ForeignKey      bool   `json:"ForeignKey,omitempty"`
	ForeignKeyType  string `json:"ForeignKeyType,omitempty"`
	ForeignKeyField string `json:"ForeignKeyField,omitempty"`
	Index           bool   `json:"Index,omitempty"`
	Unique          bool   `json:"Unique,omitempty"`
}

type Fields map[string]*Field

type Index struct {
	FieldNames []string `json:"FieldNames"`
	Unique     bool     `json:"Unique"`
}

type Indexes map[string]*Index
