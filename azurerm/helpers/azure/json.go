package azure

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// Deprecated: moved to utils and will be removed in 3.0
func NormalizeJson(jsonString interface{}) string {
	return utils.NormalizeJson(jsonString)
}
