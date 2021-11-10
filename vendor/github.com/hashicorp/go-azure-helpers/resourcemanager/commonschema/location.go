package commonschema

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func LocationSchema() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		ValidateFunc:     location.EnhancedValidate,
		StateFunc:        location.StateFunc,
		DiffSuppressFunc: location.DiffSuppressFunc,
	}
}

func LocationSchemaOptional() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		StateFunc:        location.StateFunc,
		DiffSuppressFunc: location.DiffSuppressFunc,
	}
}

func LocationSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
}

func LocationSchemaWithoutForceNew() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		ValidateFunc:     location.EnhancedValidate,
		StateFunc:        location.StateFunc,
		DiffSuppressFunc: location.DiffSuppressFunc,
	}
}
