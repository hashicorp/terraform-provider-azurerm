package shim

import "github.com/hashicorp/terraform-provider-azurerm/utils"

func mapStringPtrToMapString(input map[string]*string) map[string]string {
	output := make(map[string]string, 0)

	for k, v := range input {
		if v == nil {
			continue
		}

		output[k] = *v
	}

	return output
}

func mapStringToMapStringPtr(input map[string]string) map[string]*string {
	output := make(map[string]*string, 0)

	for k, v := range input {
		output[k] = utils.String(v)
	}

	return output
}
