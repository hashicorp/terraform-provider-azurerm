package cdn

import "github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"

func expandEnabledBool(isEnabled bool) cdn.EnabledState {
	if isEnabled {
		return cdn.EnabledStateEnabled
	}

	return cdn.EnabledStateDisabled
}

func flattenEnabledBool(input cdn.EnabledState) bool {
	return input == cdn.EnabledStateEnabled
}
