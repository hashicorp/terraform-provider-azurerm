package suppress

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
)

// moved to internal and will be removed in 3.0 (nor marked as deprecated as it would affect many packages)
func CaseDifference(k, old, new string, d *schema.ResourceData) bool {
	return suppress.CaseDifference(k, old, new, d)
}
