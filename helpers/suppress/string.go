package suppress

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
)

// Deprecated: has been moved to internal and will be removed sooner then later
func CaseDifference(k, old, new string, d *pluginsdk.ResourceData) bool {
	return suppress.CaseDifference(k, old, new, d)
}
