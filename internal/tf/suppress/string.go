package suppress

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

func CaseDifference(_, old, new string, _ *schema.ResourceData) bool {
	return strings.EqualFold(old, new)
}

// CaseDifferenceV2Only only suppress case difference for v2.0.
func CaseDifferenceV2Only(_, old, new string, _ *schema.ResourceData) bool {
	if features.ThreePointOhBeta() {
		return old == new
	}
	return strings.EqualFold(old, new)
}
