// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontDoorRoute() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorRouteCreate,
		Read:   resourceCdnFrontDoorRouteRead,
		Update: resourceCdnFrontDoorRouteUpdate,
		Delete: resourceCdnFrontDoorRouteDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontDoorRouteID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorRouteName,
			},

			"cdn_frontdoor_endpoint_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorEndpointID,
			},

			"cdn_frontdoor_origin_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FrontDoorOriginGroupID,
			},

			// NOTE: These are not sent to the API, they are only here so Terraform
			// can provision/destroy the resources in the correct order.
			"cdn_frontdoor_origin_ids": {
				Type:     pluginsdk.TypeList,
				Required: true,

				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.FrontDoorOriginID,
				},
			},

			"cdn_frontdoor_custom_domain_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,

				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.FrontDoorCustomDomainID,
				},
			},

			"link_to_default_domain": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			// NOTE: Per the service team this cannot just be omitted it must explicitly be set to nil to disable caching
			"cache": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						// NOTE: CSV string to API
						"query_strings": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringDoesNotContainAny(","),
							},
						},

						"query_string_caching_behavior": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(cdn.AfdQueryStringCachingBehaviorIgnoreQueryString),
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.AfdQueryStringCachingBehaviorIgnoreQueryString),
								string(cdn.AfdQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings),
								string(cdn.AfdQueryStringCachingBehaviorIncludeSpecifiedQueryStrings),
								string(cdn.AfdQueryStringCachingBehaviorUseQueryString),
							}, false),
						},

						"compression_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"content_types_to_compress": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice(frontDoorContentTypes(), false),
							},
						},
					},
				},
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"forwarding_protocol": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(cdn.ForwardingProtocolMatchRequest),
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.ForwardingProtocolHTTPOnly),
					string(cdn.ForwardingProtocolHTTPSOnly),
					string(cdn.ForwardingProtocolMatchRequest),
				}, false),
			},

			"https_redirect_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"cdn_frontdoor_origin_path": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"patterns_to_match": {
				Type:     pluginsdk.TypeList,
				Required: true,

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"cdn_frontdoor_rule_set_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"supported_protocols": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MaxItems: 2,

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(cdn.AFDEndpointProtocolsHTTP),
						string(cdn.AFDEndpointProtocolsHTTPS),
					}, false),
				},
			},
		},
	}
}

func resourceCdnFrontDoorRouteCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	endpointRaw := d.Get("cdn_frontdoor_endpoint_id").(string)
	endpoint, err := parse.FrontDoorEndpointID(endpointRaw)
	if err != nil {
		return err
	}

	id := parse.NewFrontDoorRouteID(endpoint.SubscriptionId, endpoint.ResourceGroup, endpoint.ProfileName, endpoint.AfdEndpointName, d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cdn_frontdoor_route", id.ID())
	}

	var origins []interface{}
	var originGroup *cdn.ResourceReference

	protocolsRaw := d.Get("supported_protocols").(*pluginsdk.Set).List()
	originGroupRaw := d.Get("cdn_frontdoor_origin_group_id").(string)
	ruleSetIdsRaw := d.Get("cdn_frontdoor_rule_set_ids").(*pluginsdk.Set).List()
	originsRaw := d.Get("cdn_frontdoor_origin_ids").([]interface{})
	customDomainsRaw := d.Get("cdn_frontdoor_custom_domain_ids").(*pluginsdk.Set).List()
	httpsRedirect := d.Get("https_redirect_enabled").(bool)
	linkToDefaultDomain := d.Get("link_to_default_domain").(bool)

	// NOTE: If HTTPS Redirect is enabled the Supported Protocols must support both HTTP and HTTPS
	// This configuration does not cause an error when provisioned, however the http requests that
	// are supposed to be redirected to https remain http requests
	if httpsRedirect {
		if err := validate.SupportsBothHttpAndHttps(protocolsRaw, "https_redirect_enabled"); err != nil {
			return err
		}
	}

	normalizedCustomDomains, err := expandCustomDomains(customDomainsRaw)
	if err != nil {
		return err
	}

	if !linkToDefaultDomain && len(normalizedCustomDomains) == 0 {
		return fmt.Errorf("it is invalid to disable the 'LinkToDefaultDomain' for the CDN Front Door Route(Name: %s) since the route does not have any CDN Front Door Custom Domains associated with it", id.RouteName)
	} else if len(normalizedCustomDomains) != 0 {
		if err := sliceHasDuplicates(normalizedCustomDomains, "CDN FrontDoor Custom Domain"); err != nil {
			return err
		}

		if err := validateRoutesCustomDomainProfile(normalizedCustomDomains, id.ProfileName); err != nil {
			return err
		}
	}

	if originGroupRaw != "" {
		id, err := parse.FrontDoorOriginGroupID(originGroupRaw)
		if err != nil {
			return err
		}

		originGroup = expandResourceReference(id.ID())
	}

	normalizedRuleSets, err := expandRuleSetIds(ruleSetIdsRaw)
	if err != nil {
		return err
	}

	props := cdn.Route{
		RouteProperties: &cdn.RouteProperties{
			CustomDomains:       expandCustomDomainActivatedResourceArray(normalizedCustomDomains),
			CacheConfiguration:  expandCdnFrontdoorRouteCacheConfiguration(d.Get("cache").([]interface{})),
			EnabledState:        expandEnabledBool(d.Get("enabled").(bool)),
			ForwardingProtocol:  cdn.ForwardingProtocol(d.Get("forwarding_protocol").(string)),
			HTTPSRedirect:       expandEnabledBoolToRouteHttpsRedirect(httpsRedirect),
			LinkToDefaultDomain: expandEnabledBoolToLinkToDefaultDomain(linkToDefaultDomain),
			OriginGroup:         originGroup,
			PatternsToMatch:     utils.ExpandStringSlice(d.Get("patterns_to_match").([]interface{})),
			RuleSets:            expandRuleSetReferenceArray(normalizedRuleSets),
			SupportedProtocols:  expandEndpointProtocolsArray(protocolsRaw),
		},
	}

	if originPath := d.Get("cdn_frontdoor_origin_path").(string); originPath != "" {
		props.RouteProperties.OriginPath = utils.String(originPath)
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName, props)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// NOTE: These are not sent to the API, they are only here so Terraform
	// can provision/destroy the resources in the correct order.
	for _, origin := range originsRaw {
		id, err := parse.FrontDoorOriginID(origin.(string))
		if err != nil {
			return err
		}

		origins = append(origins, id.ID())
	}

	if len(origins) != 0 {
		d.Set("cdn_frontdoor_origin_ids", origins)
	}

	return resourceCdnFrontDoorRouteRead(d, meta)
}

func resourceCdnFrontDoorRouteRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorRouteID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	// NOTE: These are not sent to the API, they are only here so Terraform
	// can provision/destroy the resources in the correct order.
	if originIds := d.Get("cdn_frontdoor_origin_ids").([]interface{}); len(originIds) > 0 {
		d.Set("cdn_frontdoor_origin_ids", utils.ExpandStringSlice(originIds))
	}

	d.Set("name", id.RouteName)
	d.Set("cdn_frontdoor_endpoint_id", parse.NewFrontDoorEndpointID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.AfdEndpointName).ID())

	if props := resp.RouteProperties; props != nil {
		customDomains, err := flattenCustomDomainActivatedResourceArray(props.CustomDomains)
		if err != nil {
			return err
		}

		d.Set("cdn_frontdoor_custom_domain_ids", customDomains)
		d.Set("enabled", flattenEnabledBool(props.EnabledState))
		d.Set("forwarding_protocol", props.ForwardingProtocol)
		d.Set("https_redirect_enabled", flattenHttpsRedirectToBool(props.HTTPSRedirect))
		d.Set("cdn_frontdoor_origin_path", props.OriginPath)
		d.Set("patterns_to_match", props.PatternsToMatch)
		d.Set("link_to_default_domain", flattenLinkToDefaultDomainToBool(props.LinkToDefaultDomain))

		if err := d.Set("cache", flattenCdnFrontdoorRouteCacheConfiguration(props.CacheConfiguration)); err != nil {
			return fmt.Errorf("setting `cache`: %+v", err)
		}

		originGroupId, err := flattenOriginGroupResourceReference(props.OriginGroup)
		if err != nil {
			return fmt.Errorf("flattening `cdn_frontdoor_origin_group_id`: %+v", err)
		}

		if err := d.Set("cdn_frontdoor_origin_group_id", originGroupId); err != nil {
			return fmt.Errorf("setting `cdn_frontdoor_origin_group_id`: %+v", err)
		}

		if err := d.Set("cdn_frontdoor_rule_set_ids", flattenRuleSetResourceArray(props.RuleSets)); err != nil {
			return fmt.Errorf("setting `cdn_frontdoor_rule_set_ids`: %+v", err)
		}

		if err := d.Set("supported_protocols", flattenCdnFrontdoorRouteEndpointProtocolsArray(props.SupportedProtocols)); err != nil {
			return fmt.Errorf("setting `supported_protocols`: %+v", err)
		}
	}

	return nil
}

func resourceCdnFrontDoorRouteUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	workaroundsClient := azuresdkhacks.NewCdnFrontDoorRoutesWorkaroundClient(client)
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorRouteID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}

	if existing.RouteProperties == nil {
		return fmt.Errorf("retrieving existing %s: 'properties' was nil", *id)
	}

	// we need to lock the route for update because the custom domain
	// association may also be trying to update the route as well...
	locks.ByName(id.RouteName, cdnFrontDoorRouteResourceName)
	defer locks.UnlockByName(id.RouteName, cdnFrontDoorRouteResourceName)

	httpsRedirect := d.Get("https_redirect_enabled").(bool)
	protocolsRaw := d.Get("supported_protocols").(*pluginsdk.Set).List()
	customDomainsRaw := d.Get("cdn_frontdoor_custom_domain_ids").(*pluginsdk.Set).List()
	originGroupRaw := d.Get("cdn_frontdoor_origin_group_id").(string)
	ruleSetIdsRaw := d.Get("cdn_frontdoor_rule_set_ids").(*pluginsdk.Set).List()
	linkToDefaultDomain := d.Get("link_to_default_domain").(bool)

	// NOTE: If HTTPS Redirect is enabled the Supported Protocols must support both HTTP and HTTPS
	// This configuration does not cause an error when provisioned, however the http requests that
	// are supposed to be redirected to https remain http requests
	if httpsRedirect {
		if err := validate.SupportsBothHttpAndHttps(protocolsRaw, "https_redirect_enabled"); err != nil {
			return err
		}
	}

	originGroup, err := parse.FrontDoorOriginGroupID(originGroupRaw)
	if err != nil {
		return err
	}

	customDomains, err := expandCustomDomains(customDomainsRaw)
	if err != nil {
		return err
	}

	if !linkToDefaultDomain && len(customDomains) == 0 {
		return fmt.Errorf("it is invalid to disable the 'LinkToDefaultDomain' for the CDN Front Door Route(Name: %s) since the route does not have any CDN Front Door Custom Domains associated with it", id.RouteName)
	} else if len(customDomains) != 0 {
		if err := sliceHasDuplicates(customDomains, "CDN FrontDoor Custom Domain"); err != nil {
			return err
		}

		if err := validateRoutesCustomDomainProfile(customDomains, id.ProfileName); err != nil {
			return err
		}
	}

	// NOTE: You need to always pass these three on update else you will
	// disable your cache, disassociate your custom domains or remove your origin path...
	updateProps := azuresdkhacks.RouteUpdatePropertiesParameters{
		CustomDomains:      existing.RouteProperties.CustomDomains,
		CacheConfiguration: existing.RouteProperties.CacheConfiguration,
		OriginPath:         existing.RouteProperties.OriginPath,
	}

	if d.HasChange("cache") {
		updateProps.CacheConfiguration = expandCdnFrontdoorRouteCacheConfiguration(d.Get("cache").([]interface{}))
	}

	if d.HasChange("enabled") {
		updateProps.EnabledState = expandEnabledBool(d.Get("enabled").(bool))
	}

	if d.HasChange("forwarding_protocol") {
		updateProps.ForwardingProtocol = cdn.ForwardingProtocol(d.Get("forwarding_protocol").(string))
	}

	if d.HasChange("https_redirect_enabled") {
		updateProps.HTTPSRedirect = expandEnabledBoolToRouteHttpsRedirect(httpsRedirect)
	}

	if d.HasChange("link_to_default_domain") {
		updateProps.LinkToDefaultDomain = expandEnabledBoolToLinkToDefaultDomain(d.Get("link_to_default_domain").(bool))
	}

	if d.HasChange("cdn_frontdoor_custom_domain_ids") {
		updateProps.CustomDomains = expandCustomDomainActivatedResourceArray(customDomains)
	}

	if d.HasChange("cdn_frontdoor_origin_group_id") {
		updateProps.OriginGroup = expandResourceReference(originGroup.ID())
	}

	if d.HasChange("cdn_frontdoor_origin_path") {
		updateProps.OriginPath = nil

		originPath := d.Get("cdn_frontdoor_origin_path").(string)
		if originPath != "" {
			updateProps.OriginPath = utils.String(originPath)
		}
	}

	if d.HasChange("patterns_to_match") {
		updateProps.PatternsToMatch = utils.ExpandStringSlice(d.Get("patterns_to_match").([]interface{}))
	}

	if d.HasChange("cdn_frontdoor_rule_set_ids") {
		ruleSets, err := expandRuleSetIds(ruleSetIdsRaw)
		if err != nil {
			return err
		}

		updateProps.RuleSets = expandRuleSetReferenceArray(ruleSets)
	}

	if d.HasChange("supported_protocols") {
		updateProps.SupportedProtocols = expandEndpointProtocolsArray(protocolsRaw)
	}

	updateParams := azuresdkhacks.RouteUpdateParameters{
		RouteUpdatePropertiesParameters: &updateProps,
	}

	future, err := workaroundsClient.Update(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName, updateParams)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	// NOTE: These are not sent to the API, they are only here so Terraform
	// can provision/destroy the resources in the correct order.
	if originIds := d.Get("cdn_frontdoor_origin_ids").([]interface{}); len(originIds) > 0 {
		d.Set("cdn_frontdoor_origin_ids", utils.ExpandStringSlice(originIds))
	}

	return resourceCdnFrontDoorRouteRead(d, meta)
}

func resourceCdnFrontDoorRouteDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorRouteID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandEndpointProtocolsArray(input []interface{}) *[]cdn.AFDEndpointProtocols {
	results := make([]cdn.AFDEndpointProtocols, 0)

	for _, item := range input {
		results = append(results, cdn.AFDEndpointProtocols(item.(string)))
	}

	return &results
}

func expandRuleSetReferenceArray(input []interface{}) *[]cdn.ResourceReference {
	results := make([]cdn.ResourceReference, 0)

	// NOTE: The Frontdoor service, do not treat an empty object like an empty object
	// if it is not nil they assume it is fully defined and then end up throwing errors
	// when they attempt to get a value from one of the fields.
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	for _, item := range input {
		results = append(results, cdn.ResourceReference{
			ID: utils.String(item.(string)),
		})
	}

	return &results
}

func expandCdnFrontdoorRouteCacheConfiguration(input []interface{}) *cdn.AfdRouteCacheConfiguration {
	// NOTE: If this is not an explicit nil you will receive an "Unsupported QueryStringCachingBehavior type:
	// Property 'RouteV2.CacheConfiguration.QueryStringCachingBehavior' is required but it was not set" error.
	// The Frontdoor service treats empty slices as if they are fully defined unlike other services.
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	queryStringCachingBehaviorValue := cdn.AfdQueryStringCachingBehavior(v["query_string_caching_behavior"].(string))
	compressionEnabled := v["compression_enabled"].(bool)

	cacheConfiguration := &cdn.AfdRouteCacheConfiguration{
		CompressionSettings: &cdn.CompressionSettings{
			IsCompressionEnabled: utils.Bool(compressionEnabled),
		},
		QueryParameters:            expandStringSliceToCsvFormat(v["query_strings"].([]interface{})),
		QueryStringCachingBehavior: queryStringCachingBehaviorValue,
	}

	if contentTypes := v["content_types_to_compress"].([]interface{}); len(contentTypes) > 0 {
		cacheConfiguration.CompressionSettings.ContentTypesToCompress = utils.ExpandStringSlice(contentTypes)
	}

	return cacheConfiguration
}

func flattenRuleSetResourceArray(input *[]cdn.ResourceReference) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	// Normalize these values in the configuration file we know they are valid because they were set on the
	// resource... if these are modified in the portal they will all be lowercased...
	// Issue: https://github.com/Azure/azure-sdk-for-go/issues/19378
	for _, ruleSet := range *input {
		id, _ := parse.FrontDoorRuleSetIDInsensitively(*ruleSet.ID)
		results = append(results, id.ID())
	}

	return results
}

func flattenCdnFrontdoorRouteEndpointProtocolsArray(input *[]cdn.AFDEndpointProtocols) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, string(item))
	}

	return results
}

func flattenCdnFrontdoorRouteCacheConfiguration(input *cdn.AfdRouteCacheConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	queryParameters := make([]interface{}, 0)
	if input.QueryParameters != nil {
		queryParameters = flattenCsvToStringSlice(input.QueryParameters)
	}

	cachingBehavior := ""
	if input.QueryStringCachingBehavior != "" {
		cachingBehavior = string(input.QueryStringCachingBehavior)
	}

	compressionEnabled := false
	contentTypesToCompress := make([]interface{}, 0)
	if v := input.CompressionSettings; v != nil {
		if v.IsCompressionEnabled != nil {
			compressionEnabled = *v.IsCompressionEnabled
		}
		contentTypesToCompress = utils.FlattenStringSlice(v.ContentTypesToCompress)
	}

	return []interface{}{
		map[string]interface{}{
			"compression_enabled":           compressionEnabled,
			"content_types_to_compress":     contentTypesToCompress,
			"query_string_caching_behavior": cachingBehavior,
			"query_strings":                 queryParameters,
		},
	}
}
