package storage

import "github.com/hashicorp/terraform/helper/schema"

func MetaDataSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
	}
}

func ExpandMetaData(input map[string]interface{}) map[string]string {
	output := make(map[string]string, 0)

	for k, v := range input {
		output[k] = v.(string)
	}

	return output
}

func FlattenMetaData(input map[string]string) map[string]interface{} {
	output := make(map[string]interface{}, 0)

	for k, v := range input {
		output[k] = v
	}

	return output
}
