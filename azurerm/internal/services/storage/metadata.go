package storage

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
)

func MetaDataSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeMap,
		Optional:     true,
		ValidateFunc: validate.MetaDataKeys,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func MetaDataComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeMap,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validate.MetaDataKeys,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func ExpandMetaData(input map[string]interface{}) map[string]string {
	output := make(map[string]string)

	for k, v := range input {
		output[k] = v.(string)
	}

	return output
}

func FlattenMetaData(input map[string]string) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range input {
		output[k] = v
	}

	return output
}
