package storage

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func MetaDataSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeMap,
		Optional:     true,
		ValidateFunc: validate.MetaDataKeys,
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
		},
	}
}

func MetaDataComputedSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeMap,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validate.MetaDataKeys,
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
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
