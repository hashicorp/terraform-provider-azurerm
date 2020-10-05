package suppress

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
)

// Deprecated: moved to internal and will be removed in 3.0
func XmlDiff(k, old, new string, d *schema.ResourceData) bool {
	return suppress.XmlDiff(k, old, new, d)
}
