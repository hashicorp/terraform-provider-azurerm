// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/routes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	waf "github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2025-03-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func flattenTransformSlice(input *[]waf.TransformType) []interface{} {
	result := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return result
	}

	for _, item := range *input {
		result = append(result, string(item))
	}

	return result
}

func flattenFrontendEndpointLinkSlice(input *[]waf.FrontendEndpointLink) []interface{} {
	result := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return result
	}

	for _, item := range *input {
		if item.Id == nil {
			continue
		}

		result = append(result, *item.Id)
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

func expandCustomDomainActivatedResourceArray(input []interface{}) *[]routes.ActivatedResourceReference {
	results := make([]routes.ActivatedResourceReference, 0)

	// NOTE: I have confirmed with the service team that this is required to be an explicit "nil" value, an empty
	// list will not work. I had to modify the SDK to allow for nil which in the API means disassociate the custom domains.
	if len(input) == 0 {
		return nil
	}

	for _, customDomain := range input {
		results = append(results, routes.ActivatedResourceReference{
			Id: pointer.To(customDomain.(string)),
		})
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

// determines if the slice contains the value case-insensitively
func routeSliceContains(input *[]routes.RouteId, value string) bool {
	if len(*input) == 0 || input == nil {
		return false
	}

	for _, key := range *input {
		v := key.ID()
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

func getRouteProperties(ctx context.Context, meta interface{}, id *routes.RouteId, resourceName string) ([]interface{}, *routes.RouteProperties, error) {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: retrieving %s: %+v", resourceName, id, err)
	}

	if resp.Model == nil {
		return nil, nil, fmt.Errorf("%s: retrieving %s: model was nil", resourceName, id)
	}

	if resp.Model.Properties == nil {
		return nil, nil, fmt.Errorf("%s: retrieving %s: properties was nil", resourceName, id)
	}
	props := resp.Model.Properties

	customDomains, err := flattenCdnFrontDoorRouteCustomDomainActivatedResourceArray(props.CustomDomains)
	if err != nil {
		return nil, nil, err
	}

	return customDomains, props, nil
}

func removeCustomDomainAssociationFromRoutes(ctx context.Context, meta interface{}, routes *[]routes.RouteId, customDomainID *afddomains.CustomDomainId) error {
	if routes != nil && len(*routes) != 0 {
		for _, route := range *routes {
			// lock the route resource for update...
			locks.ByID(route.ID())
			defer locks.UnlockByID(route.ID())

			// Check to see if the route still exists and grab its properties...
			// NOTE: cdnFrontDoorRouteResourceName is defined in the "cdn_frontdoor_route_disable_link_to_default_domain_resource" file
			// ignore the error because that could just mean that the route has already been deleted...
			customDomains, props, err := getRouteProperties(ctx, meta, &route, cdnFrontDoorCustomDomainResourceName)
			if err == nil {
				// Check to make sure the custom domain is still associated with the route
				isAssociated := sliceContainsString(customDomains, customDomainID.ID())

				if isAssociated {
					// it is, now removed the association...
					newDomains := sliceRemoveString(customDomains, customDomainID.ID())
					if err := updateRouteAssociations(ctx, meta, &route, newDomains, props, customDomainID); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func updateRouteAssociations(ctx context.Context, meta interface{}, routeId *routes.RouteId, customDomains []interface{}, props *routes.RouteProperties, customDomainID *afddomains.CustomDomainId) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient

	updateParams := routes.RouteUpdateParameters{
		Properties: &routes.RouteUpdatePropertiesParameters{
			CustomDomains: expandCustomDomainActivatedResourceArray(customDomains),
		},
	}

	// NOTE: You must pull the Cache Configuration from the existing route else you will get a diff
	// because a nil value means disabled
	if props.CacheConfiguration != nil {
		updateParams.Properties.CacheConfiguration = props.CacheConfiguration
	}

	// NOTE: If there are no more custom domains associated with the route you must flip the
	// 'link to default domain' field to 'true' else the route will be in an invalid state...
	if len(customDomains) == 0 {
		updateParams.Properties.LinkToDefaultDomain = pointer.To(routes.LinkToDefaultDomainEnabled)
	}

	if err := client.UpdateThenPoll(ctx, *routeId, updateParams); err != nil {
		return fmt.Errorf("%s: updating the association with %s: %+v", *customDomainID, *routeId, err)
	}

	return nil
}

func validateRoutesCustomDomainProfile(customDomains []interface{}, routeProfile string) error {
	wrongProfile := make([]string, 0)

	if len(customDomains) != 0 {
		// Verify all of the custom domains belong to the same profile as the route...
		for _, v := range customDomains {
			customDomain, err := afddomains.ParseCustomDomainID(v.(string))
			if err != nil {
				return err
			}

			if customDomain.ProfileName != routeProfile {
				wrongProfile = append(wrongProfile, customDomain.ID())
			}
		}

		if len(wrongProfile) > 0 {
			return fmt.Errorf("the following CDN FrontDoor Custom Domain(s) do not belong to the expected CDN FrontDoor Profile(Name: %q). Please remove the following CDN FrontDoor Custom Domain(s) from your CDN FrontDoor Route configuration block: %s", routeProfile, strings.Join(wrongProfile, ", "))
		}
	}

	return nil
}

// Validates that the CDN FrontDoor Custom Domain can be associated with the CDN FrontDoor Route
func validateCustomDomainRoutes(input *[]routes.RouteId, customDomainID *afddomains.CustomDomainId) error {
	if input == nil || len(*input) == 0 {
		return nil
	}

	// check for duplicates...
	if err := routeSliceHasDuplicates(input, "CDN FrontDoor Route"); err != nil {
		return err
	}

	for i, route := range *input {
		// the route and custom domain profiles must match...
		if customDomainID.ProfileName != route.ProfileName {
			return fmt.Errorf("the CDN FrontDoor Custom Domain(Name: %q, Profile: %q) and the CDN FrontDoor Route(Name: %q, Profile: %q) must belong to the same CDN FrontDoor Profile", customDomainID.CustomDomainName, customDomainID.ProfileName, route.RouteName, route.ProfileName)
		}

		// validate all routes are using the same endpoint because a custom domain can not
		// be associated with routes that target two different endpoints...
		for t, v := range *input {
			if i == t {
				continue
			}

			if route.AfdEndpointName != v.AfdEndpointName {
				return fmt.Errorf("the CDN FrontDoor Route(Name: %q) and CDN FrontDoor Route(Name: %q) do not reference the same CDN FrontDoor Endpoint(Name: %q). All CDN FrontDoor Routes must reference the same CDN FrontDoor Endpoint %q to associate the CDN FrontDoor Custom Domain(Name: %q) with more than one CDN FrontDoor Route", route.RouteName, v.RouteName, route.AfdEndpointName, route.AfdEndpointName, customDomainID.CustomDomainName)
			}
		}
	}

	return nil
}

func routeSliceHasDuplicates(input *[]routes.RouteId, resourceName string) error {
	k := make(map[string]bool)
	if input == nil || len(*input) == 0 {
		return nil
	}

	for _, route := range *input {
		if _, d := k[strings.ToLower(route.ID())]; !d {
			k[strings.ToLower(route.ID())] = true
		} else {
			return fmt.Errorf("duplicate %[1]s detected, please remove all duplicate entries for the %[1]s(ID: %q) from your configuration block", resourceName, route.ID())
		}
	}

	return nil
}

// Determines what CDN FrontDoor Routes need to be removed/disassociated from this CDN FrontDoor Custom Domain
func routeDelta(oldRoutes *[]routes.RouteId, newRoutes *[]routes.RouteId) (*[]routes.RouteId, *[]routes.RouteId) {
	remove := make([]routes.RouteId, 0)
	shared := make([]routes.RouteId, 0)
	if newRoutes == nil || len(*newRoutes) == 0 {
		return oldRoutes, &shared
	}

	if oldRoutes == nil || len(*oldRoutes) == 0 {
		return &remove, &shared
	}

	// just find what old routes are not in the new route list...
	for _, oldRoute := range *oldRoutes {
		if !routeSliceContains(newRoutes, oldRoute.ID()) {
			remove = append(remove, oldRoute)
		} else {
			shared = append(shared, oldRoute)
		}
	}

	return &remove, &shared
}

func expandRoutes(input []interface{}) (*[]routes.RouteId, []interface{}, error) {
	out := make([]routes.RouteId, 0)
	config := make([]interface{}, 0)
	if len(input) == 0 || input == nil {
		return &out, config, nil
	}

	for _, v := range input {
		id, err := routes.ParseRouteID(v.(string))
		if err != nil {
			return nil, nil, err
		}

		out = append(out, *id)
		config = append(config, id.ID())
	}

	return &out, config, nil
}

const RuleCacheBehaviorDisabled = "Disabled"

func PossibleValuesForRuleCacheBehavior() []string {
	return []string{
		string(rules.RuleCacheBehaviorHonorOrigin),
		string(rules.RuleCacheBehaviorOverrideAlways),
		string(rules.RuleCacheBehaviorOverrideIfOriginMissing),
		// Exposed `Disabled` as a valid value for provider issue #19008.
		RuleCacheBehaviorDisabled,
	}
}
