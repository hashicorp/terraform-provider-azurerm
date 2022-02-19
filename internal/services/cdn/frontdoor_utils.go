package cdn

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdorigingroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdorigins"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/routes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ConvertFrontdoorProfileTags(tagMap *map[string]string) map[string]*string {
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

func ConvertBoolToOriginsEnabledState(isEnabled bool) *afdorigins.EnabledState {
	out := afdorigins.EnabledState(afdorigins.EnabledStateDisabled)

	if isEnabled {
		out = afdorigins.EnabledState(afdorigins.EnabledStateEnabled)
	}
	return &out
}

func ConvertOriginsEnabledStateToBool(enabledState *afdorigins.EnabledState) bool {
	if enabledState == nil {
		return false
	}

	return (*enabledState == afdorigins.EnabledState(afdorigins.EnabledStateEnabled))
}

func ConvertBoolToOriginGroupsEnabledState(isEnabled bool) *afdorigingroups.EnabledState {
	out := afdorigingroups.EnabledState(afdorigingroups.EnabledStateDisabled)

	if isEnabled {
		out = afdorigingroups.EnabledState(afdorigingroups.EnabledStateEnabled)
	}

	return &out
}

func ConvertOriginGroupsEnabledStateToBool(enabledState *afdorigingroups.EnabledState) bool {
	if enabledState == nil {
		return false
	}

	return (*enabledState == afdorigingroups.EnabledState(afdorigingroups.EnabledStateEnabled))
}

func ConvertBoolToRoutesEnabledState(isEnabled bool) *routes.EnabledState {
	out := routes.EnabledState(routes.EnabledStateDisabled)

	if isEnabled {
		out = routes.EnabledState(routes.EnabledStateEnabled)
	}

	return &out
}

func ConvertRoutesEnabledStateToBool(enabledState *routes.EnabledState) bool {
	if enabledState == nil {
		return false
	}

	return (*enabledState == routes.EnabledState(routes.EnabledStateEnabled))
}

func ConvertBoolToRouteHttpsRedirect(isEnabled bool) *routes.HttpsRedirect {
	out := routes.HttpsRedirect(routes.HttpsRedirectDisabled)

	if isEnabled {
		out = routes.HttpsRedirect(routes.HttpsRedirectEnabled)
	}

	return &out
}

func ConvertRouteHttpsRedirectToBool(httpsRedirect *routes.HttpsRedirect) bool {
	if httpsRedirect == nil {
		return false
	}

	return (*httpsRedirect == routes.HttpsRedirect(routes.HttpsRedirectEnabled))
}

func ConvertBoolToRouteLinkToDefaultDomain(isLinked bool) *routes.LinkToDefaultDomain {
	out := routes.LinkToDefaultDomain(routes.LinkToDefaultDomainDisabled)

	if isLinked {
		out = routes.LinkToDefaultDomain(routes.LinkToDefaultDomainEnabled)
	}

	return &out
}

func ConvertRouteLinkToDefaultDomainToBool(linkToDefaultDomain *routes.LinkToDefaultDomain) bool {
	if linkToDefaultDomain == nil {
		return false
	}

	return (*linkToDefaultDomain == routes.LinkToDefaultDomain(routes.LinkToDefaultDomainEnabled))
}

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
