package sdk

import "github.com/hashicorp/terraform-plugin-framework/attr"

func HasChange(planVal, stateVal attr.Value) bool {
	// planned value is unknown + !staveVal.IsUnknown means computed with no default
	return !(planVal.IsUnknown() && !stateVal.IsUnknown()) || planVal.Equal(stateVal)
}
