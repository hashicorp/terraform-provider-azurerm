package suppress

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
)

// Deprecated: has been moved to internal and will be removed sooner then later
func CaseDifference(k, old, new string, d *pluginsdk.ResourceData) bool {
	return suppress.CaseDifference(k, old, new, d)
}
