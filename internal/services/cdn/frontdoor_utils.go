package cdn

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ConvertFrontdoorTags(tagMap *map[string]string) map[string]*string {
	t := make(map[string]*string)

	if tagMap != nil {
		for k, v := range *tagMap {
			tagKey := k
			tagValue := v
			t[tagKey] = &tagValue
		}
	}

	return t
}

func ConvertBoolToEnabledState(isEnabled bool) track1.EnabledState {
	out := track1.EnabledState(track1.EnabledStateDisabled)

	if isEnabled {
		out = track1.EnabledState(track1.EnabledStateEnabled)
	}

	return out
}

func ConvertEnabledStateToBool(enabledState *track1.EnabledState) bool {
	if enabledState == nil {
		return false
	}

	return (*enabledState == track1.EnabledState(track1.EnabledStateEnabled))
}

// func ConvertBoolToOriginsEnabledState(isEnabled bool) *afdorigins.EnabledState {
// 	out := afdorigins.EnabledState(afdorigins.EnabledStateDisabled)

// 	if isEnabled {
// 		out = afdorigins.EnabledState(afdorigins.EnabledStateEnabled)
// 	}

// 	return &out
// }

// func ConvertOriginsEnabledStateToBool(enabledState *afdorigins.EnabledState) bool {
// 	if enabledState == nil {
// 		return false
// 	}

// 	return (*enabledState == afdorigins.EnabledState(afdorigins.EnabledStateEnabled))
// }

// func ConvertBoolToOriginGroupsEnabledState(isEnabled bool) *afdorigingroups.EnabledState {
// 	out := afdorigingroups.EnabledState(afdorigingroups.EnabledStateDisabled)

// 	if isEnabled {
// 		out = afdorigingroups.EnabledState(afdorigingroups.EnabledStateEnabled)
// 	}

// 	return &out
// }

// func ConvertOriginGroupsEnabledStateToBool(enabledState *afdorigingroups.EnabledState) bool {
// 	if enabledState == nil {
// 		return false
// 	}

// 	return (*enabledState == afdorigingroups.EnabledState(afdorigingroups.EnabledStateEnabled))
// }

// func ConvertBoolToRoutesEnabledState(isEnabled bool) *routes.EnabledState {
// 	out := routes.EnabledState(routes.EnabledStateDisabled)

// 	if isEnabled {
// 		out = routes.EnabledState(routes.EnabledStateEnabled)
// 	}

// 	return &out
// }

// func ConvertRoutesEnabledStateToBool(enabledState *routes.EnabledState) bool {
// 	if enabledState == nil {
// 		return false
// 	}

// 	return (*enabledState == routes.EnabledState(routes.EnabledStateEnabled))
// }

// func ConvertBoolToRouteHttpsRedirect(isEnabled bool) *routes.HttpsRedirect {
// 	out := routes.HttpsRedirect(routes.HttpsRedirectDisabled)

// 	if isEnabled {
// 		out = routes.HttpsRedirect(routes.HttpsRedirectEnabled)
// 	}

// 	return &out
// }

// func ConvertRouteHttpsRedirectToBool(httpsRedirect *routes.HttpsRedirect) bool {
// 	if httpsRedirect == nil {
// 		return false
// 	}

// 	return (*httpsRedirect == routes.HttpsRedirect(routes.HttpsRedirectEnabled))
// }

// func ConvertBoolToRouteLinkToDefaultDomain(isLinked bool) *routes.LinkToDefaultDomain {
// 	out := routes.LinkToDefaultDomain(routes.LinkToDefaultDomainDisabled)

// 	if isLinked {
// 		out = routes.LinkToDefaultDomain(routes.LinkToDefaultDomainEnabled)
// 	}

// 	return &out
// }

// func ConvertRouteLinkToDefaultDomainToBool(linkToDefaultDomain *routes.LinkToDefaultDomain) bool {
// 	if linkToDefaultDomain == nil {
// 		return false
// 	}

// 	return (*linkToDefaultDomain == routes.LinkToDefaultDomain(routes.LinkToDefaultDomainEnabled))
// }

// func ConvertBoolToEndpointsEnabledState(isEnabled bool) *track1.EnabledState {
// 	out := track1.EnabledState(track1.EnabledStateDisabled)

// 	if isEnabled {
// 		out = track1.EnabledState(track1.EnabledStateEnabled)
// 	}

// 	return &out
// }

// func ConvertBoolToEndpointsEnabledState(isEnabled bool) *track1.EnabledState {
// 	out := afdendpoints.EnabledState(afdendpoints.EnabledStateDisabled)

// 	if isEnabled {
// 		out = afdendpoints.EnabledState(afdendpoints.EnabledStateEnabled)
// 	}

// 	return &out
// }

// func ConvertEndpointsEnabledStateToBool(enabledState *afdendpoints.EnabledState) bool {
// 	if enabledState == nil {
// 		return false
// 	}

// 	return (*enabledState == afdendpoints.EnabledState(afdendpoints.EnabledStateEnabled))
// }

func IsValidDomain(i interface{}, k string) (warnings []string, errors []error) {
	if warn, err := validation.IsIPv6Address(i, k); len(err) == 0 {
		return warn, err
	}

	if warn, err := validation.IsIPv4Address(i, k); len(err) == 0 {
		return warn, err
	}

	// TODO: Figure out a better way to validate Doman Name if not and IP Address
	if warn, err := validation.StringIsNotEmpty(i, k); len(err) == 0 {
		return warn, err
	}

	return warnings, errors
}

func ValidateFrontdoorRuleSetName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `(^[a-zA-Z])([\da-zA-Z]{1,88})([a-zA-Z]$)`); !m {
		return nil, append(regexErrs, fmt.Errorf(`%q must be between 1 and 90 characters in length and begin with a letter, end with a letter and may contain only letters and numbers.`, k))
	}

	return nil, nil
}
