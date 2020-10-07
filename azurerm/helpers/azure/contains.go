package azure

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

// Deprecated: moved to utils and will be remove din 3.0
func SliceContainsValue(input []string, value string) bool {
	return utils.SliceContainsValue(input, value)
}
