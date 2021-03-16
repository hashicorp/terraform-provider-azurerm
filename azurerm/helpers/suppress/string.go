package suppress

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
)

// Deprecated: has been moved to internal and will be removed sooner then later
func CaseDifference(k, old, new string, d *schema.ResourceData) bool {
	return suppress.CaseDifference(k, old, new, d)
}
