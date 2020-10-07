package azure

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// Deprecated: moved to utils and will be removed in 3.0
func NormalizeIPv6Address(ipv6 interface{}) string {
	return utils.NormalizeIPv6Address(ipv6)
}

// Deprecated: moved to internal and will be removed in 3.0
func HashIPv6Address(ipv6 interface{}) int {
	return set.HashIPv6Address(ipv6)
}
