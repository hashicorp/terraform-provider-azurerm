package sdk

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
)

func HasChange(planVal, stateval attr.Value) bool {
	// planned value is unknown + !staveVal.IsUnknown means computed with no default
	if (planVal.IsUnknown() && !stateval.IsUnknown()) || planVal.Equal(stateval) {
		return false
	}

	return true
}
