package cdn

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-11-01/frontdoor"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

// TODO: May not need these anymore... remove if the association resource tests work...
// func expandEnabledBoolToLinkToDefaultDomain(isEnabled bool) cdn.LinkToDefaultDomain {
// 	if isEnabled {
// 		return cdn.LinkToDefaultDomainEnabled
// 	}

// 	return cdn.LinkToDefaultDomainDisabled
// }

// func flattenLinkToDefaultDomainToBool(linkToDefaultDomain cdn.LinkToDefaultDomain) bool {
// 	if len(linkToDefaultDomain) == 0 {
// 		return false
// 	}

// 	return linkToDefaultDomain == cdn.LinkToDefaultDomainEnabled
// }

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

// determines if the slice contains the value case-insensitively
func sliceContainsString(input []interface{}, value string) bool {
	for _, key := range input {
		v := key.(string)
		if strings.EqualFold(v, value) {
			return true
		}
	}

	return false
}

// returns the slice with the value removed case-insensitively
func sliceRemoveString(input []interface{}, value string) []interface{} {
	out := make([]interface{}, 0)
	if len(input) == 0 {
		return out
	}

	for _, key := range input {
		v := key.(string)
		if strings.EqualFold(v, value) {
			continue
		}
		out = append(out, key)
	}

	return out
}

func checkIfRouteExists(d *pluginsdk.ResourceData, meta interface{}, id *parse.FrontDoorRouteId, resourceName string) ([]interface{}, *cdn.RouteProperties, error) {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	ctx, routeCancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer routeCancel()

	// Check to see if the route exists
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: retrieving existing %s: %+v", resourceName, *id, err)
	}

	props := resp.RouteProperties
	if props == nil {
		return nil, nil, fmt.Errorf("%s: %s properties are 'nil': %+v", resourceName, *id, err)
	}

	customDomains := flattenCdnFrontdoorRouteActivatedResourceArray(props.CustomDomains)

	return customDomains, props, nil
}

func addCustomDomainAssociationToRoute(d *pluginsdk.ResourceData, meta interface{}, routeId *parse.FrontDoorRouteId, customDomainID *parse.FrontDoorCustomDomainId) (bool, error) {
	var associationFailed bool

	// Check to see if the route still exists or not...
	customDomains, props, err := checkIfRouteExists(d, meta, routeId, cdnFrontDoorCustomDomainResourceName)
	if err != nil {
		return associationFailed, err
	}

	// Check to make sure the custom domain is not already associated with the route
	// if it is, then there is nothing for us to do...
	isAssociated := sliceContainsString(customDomains, customDomainID.ID())

	// if it is not associated update the route to add the association...
	if !isAssociated {
		customDomains = append(customDomains, customDomainID.ID())
		err := updateRouteAssociations(d, meta, routeId, customDomains, props, customDomainID)
		if err != nil {
			return associationFailed, err
		}
	}

	return !associationFailed, nil
}

func removeCustomDomainAssociationFromRoute(d *pluginsdk.ResourceData, meta interface{}, routeId *parse.FrontDoorRouteId, customDomainID *parse.FrontDoorCustomDomainId) error {
	// Check to see if the route still exists or not...
	customDomains, props, err := checkIfRouteExists(d, meta, routeId, cdnFrontDoorCustomDomainResourceName)
	if err != nil {
		return err
	}

	// Check to make sure the custom domain is still associated with the route
	isAssociated := sliceContainsString(customDomains, customDomainID.ID())

	if isAssociated {
		// it is, now remove the association...
		newDomains := sliceRemoveString(customDomains, customDomainID.ID())
		err := updateRouteAssociations(d, meta, routeId, newDomains, props, customDomainID)
		if err != nil {
			return err
		}

		// remove the field from state...
		d.Set("associate_with_cdn_frontdoor_route_id", "")
	}

	return nil
}

func updateRouteAssociations(d *pluginsdk.ResourceData, meta interface{}, routeId *parse.FrontDoorRouteId, customDomains []interface{}, props *cdn.RouteProperties, customDomainID *parse.FrontDoorCustomDomainId) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	workaroundsClient := azuresdkhacks.NewCdnFrontDoorRoutesWorkaroundClient(client)
	ctx, routeCancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer routeCancel()

	updateProps := azuresdkhacks.RouteUpdatePropertiesParameters{
		CustomDomains: expandCdnFrontdoorRouteActivatedResourceArray(customDomains),
	}

	// NOTE: You must pull the Cache Configuration from the existing route else you will get a diff
	// because a nil value means disabled
	if props.CacheConfiguration != nil {
		updateProps.CacheConfiguration = props.CacheConfiguration
	}

	// NOTE: If there are no more custom domains associated with the route you must flip the
	// 'link to default domain' field to 'true' else the route will be in an invalid state...
	if len(customDomains) == 0 {
		updateProps.LinkToDefaultDomain = cdn.LinkToDefaultDomainEnabled
	}

	updatePrarams := azuresdkhacks.RouteUpdateParameters{
		RouteUpdatePropertiesParameters: &updateProps,
	}

	future, err := workaroundsClient.Update(ctx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName, updatePrarams)
	if err != nil {
		return fmt.Errorf("%s: updating the association with %s: %+v", *customDomainID, *routeId, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("%s: waiting to update the association with %s: %+v", *customDomainID, *routeId, err)
	}

	return nil
}

func validateCustomDomanLinkToDefaultDomainState(resourceCustomDomains []interface{}, routeCustomDomains []interface{}, routeName string, routeProfile string) error {
	// Make all of the custom domains belong to the same profile as the route...
	wrongProfile := make([]string, 0)

	for _, v := range resourceCustomDomains {
		customDomain, err := parse.FrontDoorCustomDomainIDInsensitively(v.(string))
		if err != nil {
			return err
		}

		if customDomain.ProfileName != routeProfile {
			wrongProfile = append(wrongProfile, fmt.Sprintf("%q", customDomain.ID()))
		}
	}

	if len(wrongProfile) > 0 {
		return fmt.Errorf("the following CDN Front Door Custom Domain(s) do not belong to the expected CDN Front Door Profile(Name: %q). Please remove the following CDN Front Door Custom Domain(s) from your CDN Route Disable Link To Default Domain configuration: %s", routeProfile, strings.Join(wrongProfile, ", "))
	}

	// Make sure the resource is referencing all of the custom domains that are associated with the route...
	missingDomains := make([]string, 0)

	for _, v := range routeCustomDomains {
		// If this was updated by the portal, it lowercases to resource ID...
		customDomain, err := parse.FrontDoorCustomDomainIDInsensitively(v.(string))
		if err != nil {
			return fmt.Errorf("unable to parse %q: %+v", v.(string), err)
		}

		if !sliceContainsString(resourceCustomDomains, customDomain.ID()) {
			missingDomains = append(missingDomains, fmt.Sprintf("%q", customDomain.ID()))
		}
	}

	if len(missingDomains) > 0 {
		return fmt.Errorf("does not contain all of the CDN Front Door Custom Domains that are associated with the CDN Front Door Route(Name: %q). Please add the following CDN Front Door Custom Domain(s) to your CDN Route Disable Link To Default Domain configuration: %s", routeName, strings.Join(missingDomains, ", "))
	}

	// Make sure all of the custom domains that are referenced by the resource are actually associated with the route...
	notAssociated := make([]string, 0)

	for _, v := range resourceCustomDomains {
		customDomain, err := parse.FrontDoorCustomDomainIDInsensitively(v.(string))
		if err != nil {
			return fmt.Errorf("unable to parse %q: %+v", v.(string), err)
		}

		if !sliceContainsString(routeCustomDomains, customDomain.ID()) {
			notAssociated = append(notAssociated, fmt.Sprintf("%q", customDomain.ID()))
		}
	}

	if len(notAssociated) > 0 {
		return fmt.Errorf("contains CDN Front Door Custom Domains that are not associated with the CDN Front Door Route(Name: %q). Please remove the following CDN Front Door Custom Domain(s) from your CDN Route Disable Link To Default Domain configuration: %s", routeName, strings.Join(notAssociated, ", "))
	}

	return nil
}
