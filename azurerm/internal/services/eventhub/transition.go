package eventhub

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

func expandTags(input map[string]interface{}) *map[string]string {
	output := make(map[string]string)
	for k, v := range input {
		output[k] = v.(string)
	}
	return &output
}

func flattenTags(input *map[string]string) map[string]*string {
	output := make(map[string]*string)

	if input != nil {
		for k, v := range *input {
			output[k] = utils.String(v)
		}
	}

	return output
}
