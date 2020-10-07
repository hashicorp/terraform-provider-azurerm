package set

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/set"
)

// Deprecated: moved to internal and will be removed in 3.0
func HashInt(v interface{}) int {
	return set.HashInt(v)
}

// Deprecated: moved to internal and will be removed in 3.0
func HashStringIgnoreCase(v interface{}) int {
	return set.HashStringIgnoreCase(v)
}

// Deprecated: moved to internal and will be removed in 3.0
func FromStringSlice(slice []string) *schema.Set {
	return set.FromStringSlice(slice)
}
