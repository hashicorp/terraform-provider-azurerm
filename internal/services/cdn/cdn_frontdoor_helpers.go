package cdn

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-11-01/frontdoor"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func expandEnabledBool(isEnabled bool) cdn.EnabledState {
	if isEnabled {
		return cdn.EnabledStateEnabled
	}

	return cdn.EnabledStateDisabled
}

func expandEnabledBoolToRouteHttpsRedirect(isEnabled bool) cdn.HTTPSRedirect {
	if isEnabled {
		return cdn.HTTPSRedirectEnabled
	}

	return cdn.HTTPSRedirectDisabled
}

func expandEnabledBoolToLinkToDefaultDomain(isEnabled bool) cdn.LinkToDefaultDomain {
	if isEnabled {
		return cdn.LinkToDefaultDomainEnabled
	}

	return cdn.LinkToDefaultDomainDisabled
}

func flattenLinkToDefaultDomainToBool(linkToDefaultDomain cdn.LinkToDefaultDomain) bool {
	if len(linkToDefaultDomain) == 0 {
		return false
	}

	return linkToDefaultDomain == cdn.LinkToDefaultDomainEnabled
}

func expandResourceReference(input string) *cdn.ResourceReference {
	if len(input) == 0 {
		return nil
	}

	return &cdn.ResourceReference{
		ID: utils.String(input),
	}
}

func flattenResourceReference(input *cdn.ResourceReference) string {
	if input != nil && input.ID != nil {
		return *input.ID
	}

	return ""
}

func flattenEnabledBool(input cdn.EnabledState) bool {
	return input == cdn.EnabledStateEnabled
}

func flattenRouteHttpsRedirectToBool(httpsRedirect cdn.HTTPSRedirect) bool {
	if len(httpsRedirect) == 0 {
		return false
	}

	return httpsRedirect == cdn.HTTPSRedirectEnabled
}

func expandFrontDoorTags(tagMap *map[string]string) map[string]*string {
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

func flattenFrontDoorTags(tagMap map[string]*string) *map[string]string {
	t := make(map[string]string)

	for k, v := range tagMap {
		tagKey := k
		tagValue := v
		if tagValue == nil {
			continue
		}
		t[tagKey] = *tagValue
	}

	return &t
}

func flattenTransformSlice(input *[]frontdoor.TransformType) []interface{} {
	result := make([]interface{}, 0)

	if input != nil {
		for _, item := range *input {
			result = append(result, string(item))
		}
	}
	return result
}

func flattenFrontendEndpointLinkSlice(input *[]frontdoor.FrontendEndpointLink) []interface{} {
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

func ruleHasDeliveryRuleConditions(conditions map[string]interface{}) bool {
	var hasConditions bool

	for _, condition := range conditions {
		if len(condition.([]interface{})) > 0 {
			hasConditions = true
			break
		}
	}

	return hasConditions
}

func frontDoorContentTypes() []string {
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

func routeSupportsHttpAndHttps(supportedProtocols []interface{}) bool {
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

// Takes a Slice of strings and transforms it into a CSV formatted string.
func expandStringSliceToCsvFormat(input []interface{}) *string {
	if len(input) == 0 {
		return nil
	}

	v := utils.ExpandStringSlice(input)
	csv := strings.Trim(fmt.Sprintf("[%s]", strings.Join(*v, ",")), "[]")

	return &csv
}

// Takes a CSV formatted string and transforms it into a Slice of strings.
func flattenCsvToStringSlice(input *string) []interface{} {
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
