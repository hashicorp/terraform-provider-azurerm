package cdn

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

func expandResourceReference(input string) *track1.ResourceReference {
	if len(input) == 0 {
		return nil
	}

	return &track1.ResourceReference{
		ID: utils.String(input),
	}
}

func flattenResourceReference(input *track1.ResourceReference) string {
	result := ""
	if input == nil {
		return result
	}

	if input.ID != nil {
		result = *input.ID
	}

	return result
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

func ConvertBoolToRouteHttpsRedirect(isEnabled bool) track1.HTTPSRedirect {
	out := track1.HTTPSRedirect(track1.HTTPSRedirectDisabled)

	if isEnabled {
		out = track1.HTTPSRedirect(track1.HTTPSRedirectEnabled)
	}

	return out
}

func ConvertRouteHttpsRedirectToBool(httpsRedirect *track1.HTTPSRedirect) bool {
	if httpsRedirect == nil {
		return false
	}

	return (*httpsRedirect == track1.HTTPSRedirect(track1.HTTPSRedirectEnabled))
}

func ConvertBoolToRouteLinkToDefaultDomain(isLinked bool) track1.LinkToDefaultDomain {
	out := track1.LinkToDefaultDomain(track1.LinkToDefaultDomainDisabled)

	if isLinked {
		out = track1.LinkToDefaultDomain(track1.LinkToDefaultDomainEnabled)
	}

	return out
}

func ConvertRouteLinkToDefaultDomainToBool(linkToDefaultDomain *track1.LinkToDefaultDomain) bool {
	if linkToDefaultDomain == nil {
		return false
	}

	return (*linkToDefaultDomain == track1.LinkToDefaultDomain(track1.LinkToDefaultDomainEnabled))
}

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

func SchemaFrontdoorOperator() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeString,
		Required: true,
		ValidateFunc: validation.StringInSlice([]string{
			"Any",
			"Equal",
			"Contains",
			"BeginsWith",
			"EndsWith",
			"LessThan",
			"LessThanOrEqual",
			"GreaterThan",
			"GreaterThanOrEqual",
			"RegEx",
		}, false),
	}
}

func SchemaFrontdoorNegateCondition() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeBool,
		Optional: true,
		Default:  false,
	}
}

func SchemaFrontdoorMatchValues() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 10,

		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func SchemaFrontdoorRuleTransforms() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 6,

		Elem: &pluginsdk.Schema{
			Type:    pluginsdk.TypeString,
			Default: string(track1.TransformLowercase),
			ValidateFunc: validation.StringInSlice([]string{
				string(track1.TransformLowercase),
				string(track1.TransformRemoveNulls),
				string(track1.TransformTrim),
				string(track1.TransformUppercase),
				string(track1.TransformURLDecode),
				string(track1.TransformURLEncode),
			}, false),
		},
	}
}
