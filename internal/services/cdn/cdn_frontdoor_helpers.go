package cdn

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"             // nolint: staticcheck
	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-11-01/frontdoor" // nolint: staticcheck
	dnsValidate "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func expandEnabledBool(input bool) cdn.EnabledState {
	if input {
		return cdn.EnabledStateEnabled
	}

	return cdn.EnabledStateDisabled
}

func expandEnabledBoolToRouteHttpsRedirect(input bool) cdn.HTTPSRedirect {
	if input {
		return cdn.HTTPSRedirectEnabled
	}

	return cdn.HTTPSRedirectDisabled
}

func expandEnabledBoolToLinkToDefaultDomain(input bool) cdn.LinkToDefaultDomain {
	if input {
		return cdn.LinkToDefaultDomainEnabled
	}

	return cdn.LinkToDefaultDomainDisabled
}

func flattenLinkToDefaultDomainToBool(input cdn.LinkToDefaultDomain) bool {
	if len(input) == 0 {
		return false
	}

	return input == cdn.LinkToDefaultDomainEnabled
}

func expandResourceReference(input string) *cdn.ResourceReference {
	if len(input) == 0 {
		return nil
	}

	return &cdn.ResourceReference{
		ID: utils.String(input),
	}
}

func flattenOriginGroupResourceReference(input *cdn.ResourceReference) (string, error) {
	if input != nil && input.ID != nil {
		id, err := parse.FrontDoorOriginGroupIDInsensitively(*input.ID)
		if err != nil {
			return "", err
		}

		return id.ID(), nil
	}

	return "", nil
}

func flattenSecretResourceReference(input *cdn.ResourceReference) (string, error) {
	if input != nil && input.ID != nil {
		id, err := parse.FrontDoorSecretIDInsensitively(*input.ID)
		if err != nil {
			return "", err
		}

		return id.ID(), nil
	}

	return "", nil
}

func flattenDNSZoneResourceReference(input *cdn.ResourceReference) (string, error) {
	if input != nil && input.ID != nil {
		id, err := dnsValidate.ParseDnsZoneIDInsensitively(*input.ID)
		if err != nil {
			return "", err
		}

		return id.ID(), nil
	}

	return "", nil
}

func flattenEnabledBool(input cdn.EnabledState) bool {
	return input == cdn.EnabledStateEnabled
}

func flattenHttpsRedirectToBool(input cdn.HTTPSRedirect) bool {
	return input == cdn.HTTPSRedirectEnabled
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
	if input == nil || len(*input) == 0 {
		return result
	}

	if input != nil {
		for _, item := range *input {
			result = append(result, string(item))
		}
	}

	return result
}

func flattenFrontendEndpointLinkSlice(input *[]frontdoor.FrontendEndpointLink) []interface{} {
	result := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return result
	}

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

// Takes a Slice of strings and transforms it into a CSV formatted string.
func expandStringSliceToCsvFormat(input []interface{}) *string {
	if len(input) == 0 {
		return nil
	}

	v := utils.ExpandStringSlice(input)
	csv := strings.Trim(fmt.Sprintf("[%s]", strings.Join(*v, ",")), "[]")

	return &csv
}

func expandCustomDomainActivatedResourceArray(input []interface{}) *[]cdn.ActivatedResourceReference {
	results := make([]cdn.ActivatedResourceReference, 0)

	// NOTE: I have confirmed with the service team that this is required to be an explicit "nil" value, an empty
	// list will not work. I had to modify the SDK to allow for nil which in the API means disassociate the custom domains.
	if len(input) == 0 {
		return nil
	}

	for _, customDomain := range input {
		if id, err := parse.FrontDoorCustomDomainID(customDomain.(string)); err == nil {
			results = append(results, cdn.ActivatedResourceReference{
				ID: utils.String(id.ID()),
			})
		}
	}

	return &results
}

// Takes a CSV formatted string and transforms it into a Slice of strings.
func flattenCsvToStringSlice(input *string) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	v := strings.Split(*input, ",")

	for _, s := range v {
		results = append(results, s)
	}

	return results
}

func flattenCustomDomainActivatedResourceArray(input *[]cdn.ActivatedResourceReference) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return results, nil
	}

	// Normalize these values in the configuration file we know they are valid because they were set on the
	// resource... if these are modified in the portal they will all be lowercased...
	for _, customDomain := range *input {
		if customDomain.ID == nil {
			continue
		}
		id, err := parse.FrontDoorCustomDomainIDInsensitively(*customDomain.ID)
		if err != nil {
			return nil, err
		}
		results = append(results, id.ID())
	}

	return results, nil
}

// determines if the slice contains the value case-insensitively
func sliceContainsString(input []interface{}, value string) bool {
	if len(input) == 0 {
		return false
	}

	for _, key := range input {
		v := key.(string)
		if strings.EqualFold(v, value) {
			return true
		}
	}

	return false
}

func validateCustomDomainLinkToDefaultDomainState(resourceCustomDomains []interface{}, routeCustomDomains []interface{}, routeName string, routeProfile string) error {
	// NOTE: Only used in the deprecated custom domain link to default domain resource
	if !features.FourPointOhBeta() {
		// Make all of the custom domains belong to the same profile as the route...
		wrongProfile := make([]string, 0)

		for _, v := range resourceCustomDomains {
			customDomain, err := parse.FrontDoorCustomDomainID(v.(string))
			if err != nil {
				return err
			}

			if customDomain.ProfileName != routeProfile {
				wrongProfile = append(wrongProfile, fmt.Sprintf("%q", customDomain.ID()))
			}
		}

		if len(wrongProfile) > 0 {
			return fmt.Errorf("the following CDN FrontDoor Custom Domain(s) do not belong to the expected CDN FrontDoor Profile(Name: %q). Please remove the following CDN FrontDoor Custom Domain(s) from your CDN Route Disable Link To Default Domain configuration: %s", routeProfile, strings.Join(wrongProfile, ", "))
		}

		// Make sure the resource is referencing all of the custom domains that are associated with the route...
		missingDomains := make([]string, 0)

		for _, v := range routeCustomDomains {
			// If this was updated by the portal, it lowercases to resource ID...
			customDomain, err := parse.FrontDoorCustomDomainID(v.(string))
			if err != nil {
				return fmt.Errorf("unable to parse %q: %+v", v.(string), err)
			}

			if !sliceContainsString(resourceCustomDomains, customDomain.ID()) {
				missingDomains = append(missingDomains, fmt.Sprintf("%q", customDomain.ID()))
			}
		}

		if len(missingDomains) > 0 {
			return fmt.Errorf("does not contain all of the CDN FrontDoor Custom Domains that are associated with the CDN FrontDoor Route(Name: %q). Please add the following CDN FrontDoor Custom Domain(s) to your CDN Route Disable Link To Default Domain configuration: %s", routeName, strings.Join(missingDomains, ", "))
		}

		// Make sure all of the custom domains that are referenced by the resource are actually associated with the route...
		notAssociated := make([]string, 0)

		for _, v := range resourceCustomDomains {
			customDomain, err := parse.FrontDoorCustomDomainID(v.(string))
			if err != nil {
				return fmt.Errorf("unable to parse %q: %+v", v.(string), err)
			}

			if !sliceContainsString(routeCustomDomains, customDomain.ID()) {
				notAssociated = append(notAssociated, fmt.Sprintf("%q", customDomain.ID()))
			}
		}

		if len(notAssociated) > 0 {
			return fmt.Errorf("contains CDN FrontDoor Custom Domains that are not associated with the CDN FrontDoor Route(Name: %q). Please remove the following CDN FrontDoor Custom Domain(s) from your CDN Route Disable Link To Default Domain configuration: %s", routeName, strings.Join(notAssociated, ", "))
		}
	}

	return nil
}

func validateRoutesCustomDomainProfile(customDomains []interface{}, routeProfile string) error {
	wrongProfile := make([]string, 0)

	if len(customDomains) != 0 {
		// Verify all of the custom domains belong to the same profile as the route...
		for _, v := range customDomains {
			customDomain, err := parse.FrontDoorCustomDomainID(v.(string))
			if err != nil {
				return err
			}

			if customDomain.ProfileName != routeProfile {
				wrongProfile = append(wrongProfile, fmt.Sprintf("%q", customDomain.ID()))
			}
		}

		if len(wrongProfile) > 0 {
			return fmt.Errorf("the following CDN FrontDoor Custom Domain(s) do not belong to the expected CDN FrontDoor Profile(Name: %q). Please remove the following CDN FrontDoor Custom Domain(s) from your CDN FrontDoor Route configuration block: %s", routeProfile, strings.Join(wrongProfile, ", "))
		}
	}

	return nil
}

// Checks to make sure the list of CDN FrontDoor Custom Domains does not contain duplicate entries
func sliceHasDuplicates(input []interface{}, resourceTxt string) error {
	k := make(map[string]bool)
	if len(input) == 0 || input == nil {
		return nil
	}

	for _, v := range input {
		if _, d := k[strings.ToLower(v.(string))]; !d {
			k[strings.ToLower(v.(string))] = true
		} else {
			return fmt.Errorf("duplicate %[1]s detected, please remove all duplicate entries for the %[1]s(ID: %q) from your configuration block", resourceTxt, v.(string))
		}
	}

	return nil
}

func expandRuleSetIds(input []interface{}) ([]interface{}, error) {
	out := make([]interface{}, 0)
	if len(input) == 0 || input == nil {
		return out, nil
	}

	for _, ruleSet := range input {
		id, err := parse.FrontDoorRuleSetID(ruleSet.(string))
		if err != nil {
			return nil, err
		}

		out = append(out, id.ID())
	}

	return out, nil
}
