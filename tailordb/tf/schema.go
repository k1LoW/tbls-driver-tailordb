package tf

import "github.com/hashicorp/hcl/v2"

var rootSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "terraform",
			LabelNames: nil,
		},
		{
			Type:       "variable",
			LabelNames: []string{"name"},
		},
		{
			Type:       "output",
			LabelNames: []string{"name"},
		},
		{
			Type:       "provider",
			LabelNames: []string{"name"},
		},
		{
			Type:       "resource",
			LabelNames: []string{"type", "name"},
		},
		{
			Type:       "data",
			LabelNames: []string{"type", "name"},
		},
		{
			Type:       "module",
			LabelNames: []string{"name"},
		},
		{
			Type:       "locals",
			LabelNames: nil,
		},
	},
}

var workspaceSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name:     "name",
			Required: true,
		},
		{
			Name:     "region",
			Required: true,
		},
		{
			Name: "delete_protection",
		},
		{
			Name: "folder_id",
		},
		{
			Name: "organization_id",
		},
	},
}

var tailordbSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name:     "namespace",
			Required: true,
		},
		{
			Name:     "workspace_id",
			Required: true,
		},
	},
}

var tailordbTypeSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name:     "name",
			Required: true,
		},
		{
			Name:     "namespace",
			Required: true,
		},
		{
			Name:     "workspace_id",
			Required: true,
		},
		{
			Name: "fields",
		},
		{
			Name: "description",
		},
		{
			Name: "directives",
		},
		{
			Name: "extends",
		},
		{
			Name: "indexes",
		},
		{
			Name: "record_permission",
		},
		{
			Name: "settings",
		},
		{
			Name: "type_permission",
		},
	},
}

var terraformSchema = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{
			Name: "required_version",
		},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type: "required_providers",
		},
	},
}
