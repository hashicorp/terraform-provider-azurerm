package cdn

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-11-01/frontdoor"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// TODO: the validation methods want moving into ./validation
// WS: Fixed
// TODO: the convert methods need to be made private (since they're only related to this package)
// WS: Fixed

func cdnFrontdoorPrivateLinkTargetTypes() []string {
	return []string{"blob", "blob_secondary", "sites", "web"}
}

func convertCdnFrontdoorTags(tagMap *map[string]string) map[string]*string {
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

func convertCdnFrontdoorTagsToTagsFlatten(tagMap map[string]*string) *map[string]string {
	t := make(map[string]string)

	for k, v := range tagMap {
		tagKey := k
		tagValue := v
		t[tagKey] = *tagValue
	}

	return &t
}

func convertCdnFrontdoorBoolToEnabledState(isEnabled bool) cdn.EnabledState {
	out := cdn.EnabledStateDisabled

	if isEnabled {
		out = cdn.EnabledStateEnabled
	}

	return out
}

func convertCdnFrontdoorEnabledStateToBool(enabledState *cdn.EnabledState) bool {
	if enabledState == nil {
		return false
	}

	return (*enabledState == cdn.EnabledStateEnabled)
}

func expandCdnFrontdoorResourceReference(input string) *cdn.ResourceReference {
	if len(input) == 0 {
		return nil
	}

	return &cdn.ResourceReference{
		ID: utils.String(input),
	}
}

// Takes a Slice of strings and transforms it into a CSV formatted string.
func expandCdnFrontdoorStringSliceToCsvFormat(input []interface{}) *string {
	if len(input) == 0 {
		return nil
	}

	v := utils.ExpandStringSlice(input)
	csv := strings.Trim(fmt.Sprintf("[%s]", strings.Join(*v, ",")), "[]")

	return &csv
}

// Takes a CSV formatted string and transforms it into a Slice of strings.
func flattenCdnFrontdoorCsvToStringSlice(input *string) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	v := strings.Split(*input, ",")

	for _, s := range v {
		results = append(results, s)
	}

	return results
}

func flattenCdnFrontdoorResourceReference(input *cdn.ResourceReference) string {
	result := ""
	if input == nil {
		return result
	}

	if input.ID != nil {
		result = *input.ID
	}

	return result
}

func flattenCdnFrontdoorTransformSlice(input *[]frontdoor.TransformType) []interface{} {
	result := make([]interface{}, 0)

	if input != nil {
		for _, item := range *input {
			result = append(result, string(item))
		}
	}
	return result
}

func flattenCdnFrontdoorFrontendEndpointLinkSlice(input *[]frontdoor.FrontendEndpointLink) []interface{} {
	result := make([]interface{}, 0)

	if input != nil {
		for _, item := range *input {
			if item.ID == nil {
				continue
			}

			result = append(result, *item.ID)
		}
	}
	return result
}

func convertCdnFrontdoorBoolToRouteHttpsRedirect(isEnabled bool) cdn.HTTPSRedirect {
	out := cdn.HTTPSRedirectDisabled

	if isEnabled {
		out = cdn.HTTPSRedirectEnabled
	}

	return out
}

func convertCdnFrontdoorRouteHttpsRedirectToBool(httpsRedirect *cdn.HTTPSRedirect) bool {
	if httpsRedirect == nil {
		return false
	}

	return (*httpsRedirect == cdn.HTTPSRedirectEnabled)
}

func convertCdnFrontdoorBoolToRouteLinkToDefaultDomain(isLinked bool) cdn.LinkToDefaultDomain {
	out := cdn.LinkToDefaultDomainDisabled

	if isLinked {
		out = cdn.LinkToDefaultDomainEnabled
	}

	return out
}

func convertCdnFrontdoorRouteLinkToDefaultDomainToBool(linkToDefaultDomain *cdn.LinkToDefaultDomain) bool {
	if linkToDefaultDomain == nil {
		return false
	}

	return (*linkToDefaultDomain == cdn.LinkToDefaultDomainEnabled)
}

func cdnFrontdoorContentTypes() []string {
	return []string{
		"application/eot",
		"application/font",
		"application/font-sfnt",
		"application/javascript",
		"application/json",
		"application/opentype",
		"application/otf",
		"application/pkcs7-mime",
		"application/truetype",
		"application/ttf",
		"application/vnd.ms-fontobject",
		"application/xhtml+xml",
		"application/xml",
		"application/xml+rss",
		"application/x-font-opentype",
		"application/x-font-truetype",
		"application/x-font-ttf",
		"application/x-httpd-cgi",
		"application/x-mpegurl",
		"application/x-opentype",
		"application/x-otf",
		"application/x-perl",
		"application/x-ttf",
		"application/x-javascript",
		"font/eot",
		"font/ttf",
		"font/otf",
		"font/opentype",
		"image/svg+xml",
		"text/css",
		"text/csv",
		"text/html",
		"text/javascript",
		"text/js",
		"text/plain",
		"text/richtext",
		"text/tab-separated-values",
		"text/xml",
		"text/x-script",
		"text/x-component",
		"text/x-java-source",
	}
}

func cdnFrontdoorRuleHasDeliveryRuleConditions(conditions map[string]interface{}) bool {
	var hasConditions bool

	for _, condition := range conditions {
		if len(condition.([]interface{})) > 0 {
			hasConditions = true
			break
		}
	}

	return hasConditions
}

func cdnFrontdoorRouteSupportsHttpHttps(supportedProtocols []interface{}) bool {
	var supportsBoth bool
	if len(supportedProtocols) == 0 {
		return supportsBoth
	}

	protocols := utils.ExpandStringSlice(supportedProtocols)
	if utils.SliceContainsValue(*protocols, string(cdn.AFDEndpointProtocolsHTTP)) && utils.SliceContainsValue(*protocols, string(cdn.AFDEndpointProtocolsHTTPS)) {
		supportsBoth = true
	}

	return supportsBoth
}
