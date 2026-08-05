// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdendpoints"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigingroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigins"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/routes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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
			Create: pluginsdk.DefaultTimeout(4 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(4 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := routes.ParseRouteID(id)
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
				ValidateFunc: routes.ValidateAfdEndpointID,
			},

			"cdn_frontdoor_origin_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: afdorigingroups.ValidateOriginGroupID,
			},

			// NOTE: These are not sent to the API, they are only here so Terraform
			// can provision/destroy the resources in the correct order.
			// Made this field optional to address comments in Issue #29063
			"cdn_frontdoor_origin_ids": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: afdorigins.ValidateOriginGroupOriginID,
				},
			},

			"cdn_frontdoor_custom_domain_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: afddomains.ValidateCustomDomainID,
				},
			},

			"link_to_default_domain": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"cache": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
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
							Default:  string(routes.AfdQueryStringCachingBehaviorIgnoreQueryString),
							ValidateFunc: validation.StringInSlice([]string{
								string(routes.AfdQueryStringCachingBehaviorIgnoreQueryString),
								string(routes.AfdQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings),
								string(routes.AfdQueryStringCachingBehaviorIncludeSpecifiedQueryStrings),
								string(routes.AfdQueryStringCachingBehaviorUseQueryString),
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
				Default:  string(routes.ForwardingProtocolMatchRequest),
				ValidateFunc: validation.StringInSlice([]string{
					string(routes.ForwardingProtocolHTTPOnly),
					string(routes.ForwardingProtocolHTTPSOnly),
					string(routes.ForwardingProtocolMatchRequest),
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
					Type:         pluginsdk.TypeString,
					ValidateFunc: rulesets.ValidateRuleSetID,
				},
			},

			"supported_protocols": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MaxItems: 2,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(routes.AFDEndpointProtocolsHTTP),
						string(routes.AFDEndpointProtocolsHTTPS),
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

	endpoint, err := routes.ParseAfdEndpointID(d.Get("cdn_frontdoor_endpoint_id").(string))
	if err != nil {
		return err
	}

	id := routes.NewRouteID(endpoint.SubscriptionId, endpoint.ResourceGroupName, endpoint.ProfileName, endpoint.AfdEndpointName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_route", id.ID())
		}
	}

	// NOTE: If HTTPS Redirect is enabled the Supported Protocols must support both HTTP and HTTPS
	// This configuration does not cause an error when provisioned, however the http requests that
	// are supposed to be redirected to https remain http requests
	httpsRedirect := d.Get("https_redirect_enabled").(bool)
	protocols := d.Get("supported_protocols").(*pluginsdk.Set).List()
	if httpsRedirect {
		if err := validate.SupportsBothHttpAndHttps(protocols, "https_redirect_enabled"); err != nil {
			return err
		}
	}

	customDomains := d.Get("cdn_frontdoor_custom_domain_ids").(*pluginsdk.Set).List()
	linkToDefaultDomain := d.Get("link_to_default_domain").(bool)
	if !linkToDefaultDomain && len(customDomains) == 0 {
		return errors.New("`link_to_default_domain` cannot be disabled when no `cdn_frontdoor_custom_domain_ids` are defined")
	}

	if len(customDomains) != 0 {
		if err := validateRoutesCustomDomainProfile(customDomains, id.ProfileName); err != nil {
			return err
		}
	}

	var originGroup *routes.ResourceReference
	if originGroupRaw := d.Get("cdn_frontdoor_origin_group_id").(string); originGroupRaw != "" {
		originGroup = expandCdnFrontDoorRouteResourceReference(originGroupRaw)
	}

	props := routes.Route{
		Properties: &routes.RouteProperties{
			CustomDomains:       expandCustomDomainActivatedResourceArray(customDomains),
			CacheConfiguration:  expandCdnFrontdoorRouteCacheConfiguration(d.Get("cache").([]interface{})),
			EnabledState:        expandCdnFrontDoorRouteEnabled(d.Get("enabled").(bool)),
			ForwardingProtocol:  pointer.ToEnum[routes.ForwardingProtocol](d.Get("forwarding_protocol").(string)),
			HTTPSRedirect:       expandCdnFrontDoorRouteHttpsRedirect(httpsRedirect),
			LinkToDefaultDomain: expandCdnFrontDoorRouteDefaultDomain(linkToDefaultDomain),
			OriginGroup:         originGroup,
			OriginPath:          pointer.ToOrNil(d.Get("cdn_frontdoor_origin_path").(string)),
			PatternsToMatch:     utils.ExpandStringSlice(d.Get("patterns_to_match").([]interface{})),
			RuleSets:            expandCdnFrontdoorRouteRuleSetReferenceArray(d.Get("cdn_frontdoor_rule_set_ids").(*pluginsdk.Set).List()),
			SupportedProtocols:  expandCdnFrontDoorRouteEndpointProtocolsArray(protocols),
		},
	}

	if err := client.CreateCallbackThenPoll(ctx, id, props, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// NOTE: These are not sent to the API, they are only here so Terraform
	// can provision/destroy the resources in the correct order.
	if origins := d.Get("cdn_frontdoor_origin_ids").([]interface{}); len(origins) != 0 {
		d.Set("cdn_frontdoor_origin_ids", origins)
	}

	return resourceCdnFrontDoorRouteRead(d, meta)
}

func resourceCdnFrontDoorRouteRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := routes.ParseRouteID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	// NOTE: These are not sent to the API, they are only here so Terraform
	// can provision/destroy the resources in the correct order.
	if originIds := d.Get("cdn_frontdoor_origin_ids").([]interface{}); len(originIds) > 0 {
		d.Set("cdn_frontdoor_origin_ids", originIds)
	}

	d.Set("name", id.RouteName)
	d.Set("cdn_frontdoor_endpoint_id", afdendpoints.NewAfdEndpointID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.AfdEndpointName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			customDomains, err := flattenCdnFrontDoorRouteCustomDomainActivatedResourceArray(props.CustomDomains)
			if err != nil {
				return err
			}

			d.Set("cdn_frontdoor_custom_domain_ids", customDomains)
			d.Set("enabled", pointer.From(props.EnabledState) == routes.EnabledStateEnabled)
			d.Set("forwarding_protocol", pointer.FromEnum(props.ForwardingProtocol))
			d.Set("https_redirect_enabled", pointer.From(props.HTTPSRedirect) == routes.HTTPSRedirectEnabled)
			d.Set("cdn_frontdoor_origin_path", props.OriginPath)
			d.Set("patterns_to_match", props.PatternsToMatch)
			d.Set("link_to_default_domain", pointer.From(props.LinkToDefaultDomain) == routes.LinkToDefaultDomainEnabled)

			if err := d.Set("cache", flattenCdnFrontDoorRouteCacheConfiguration(props.CacheConfiguration)); err != nil {
				return fmt.Errorf("setting `cache`: %+v", err)
			}

			originGroupId, err := flattenCdnFrontDoorRouteOriginGroupResourceReference(props.OriginGroup)
			if err != nil {
				return fmt.Errorf("flattening `cdn_frontdoor_origin_group_id`: %+v", err)
			}

			if err := d.Set("cdn_frontdoor_origin_group_id", originGroupId); err != nil {
				return fmt.Errorf("setting `cdn_frontdoor_origin_group_id`: %+v", err)
			}

			ruleSetIDs, err := flattenCdnFrontDoorRouteRuleSetResourceArray(props.RuleSets)
			if err != nil {
				return fmt.Errorf("flattening `cdn_frontdoor_rule_set_ids`: %+v", err)
			}

			if err := d.Set("cdn_frontdoor_rule_set_ids", ruleSetIDs); err != nil {
				return fmt.Errorf("setting `cdn_frontdoor_rule_set_ids`: %+v", err)
			}

			if err := d.Set("supported_protocols", pointer.FromEnumSlice(props.SupportedProtocols)); err != nil {
				return fmt.Errorf("setting `supported_protocols`: %+v", err)
			}
		}
	}

	return nil
}

func resourceCdnFrontDoorRouteUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := routes.ParseRouteID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: properties was nil", id)
	}
	props := existing.Model.Properties

	// we need to lock the route for update because the custom domain
	// association may also be trying to update the route as well...
	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	// NOTE: If HTTPS Redirect is enabled the Supported Protocols must support both HTTP and HTTPS
	// This configuration does not cause an error when provisioned, however the http requests that
	// are supposed to be redirected to https remain http requests
	httpsRedirect := d.Get("https_redirect_enabled").(bool)
	protocols := d.Get("supported_protocols").(*pluginsdk.Set).List()
	if httpsRedirect {
		if err := validate.SupportsBothHttpAndHttps(protocols, "https_redirect_enabled"); err != nil {
			return err
		}
	}

	customDomains := d.Get("cdn_frontdoor_custom_domain_ids").(*pluginsdk.Set).List()
	linkToDefaultDomain := d.Get("link_to_default_domain").(bool)
	if !linkToDefaultDomain && len(customDomains) == 0 {
		return errors.New("`link_to_default_domain` cannot be disabled when no `cdn_frontdoor_custom_domain_ids` are defined")
	}

	if len(customDomains) != 0 {
		if err := validateRoutesCustomDomainProfile(customDomains, id.ProfileName); err != nil {
			return err
		}
	}

	if d.HasChange("cache") {
		props.CacheConfiguration = expandCdnFrontdoorRouteCacheConfiguration(d.Get("cache").([]interface{}))
	}

	if d.HasChange("enabled") {
		props.EnabledState = expandCdnFrontDoorRouteEnabled(d.Get("enabled").(bool))
	}

	if d.HasChange("forwarding_protocol") {
		props.ForwardingProtocol = pointer.ToEnum[routes.ForwardingProtocol](d.Get("forwarding_protocol").(string))
	}

	if d.HasChange("https_redirect_enabled") {
		props.HTTPSRedirect = expandCdnFrontDoorRouteHttpsRedirect(httpsRedirect)
	}

	if d.HasChange("link_to_default_domain") {
		props.LinkToDefaultDomain = expandCdnFrontDoorRouteDefaultDomain(d.Get("link_to_default_domain").(bool))
	}

	if d.HasChange("cdn_frontdoor_custom_domain_ids") {
		props.CustomDomains = expandCustomDomainActivatedResourceArray(customDomains)
	}

	if d.HasChange("cdn_frontdoor_origin_group_id") {
		props.OriginGroup = expandCdnFrontDoorRouteResourceReference(d.Get("cdn_frontdoor_origin_group_id").(string))
	}

	if d.HasChange("cdn_frontdoor_origin_path") {
		props.OriginPath = pointer.ToOrNil(d.Get("cdn_frontdoor_origin_path").(string))
	}

	if d.HasChange("patterns_to_match") {
		props.PatternsToMatch = utils.ExpandStringSlice(d.Get("patterns_to_match").([]interface{}))
	}

	if d.HasChange("cdn_frontdoor_rule_set_ids") {
		props.RuleSets = expandCdnFrontdoorRouteRuleSetReferenceArray(d.Get("cdn_frontdoor_rule_set_ids").(*pluginsdk.Set).List())
	}

	if d.HasChange("supported_protocols") {
		props.SupportedProtocols = expandCdnFrontDoorRouteEndpointProtocolsArray(protocols)
	}

	if err := client.CreateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	// NOTE: These are not sent to the API, they are only here so Terraform
	// can provision/destroy the resources in the correct order.
	if originIds := d.Get("cdn_frontdoor_origin_ids").([]interface{}); len(originIds) > 0 {
		d.Set("cdn_frontdoor_origin_ids", originIds)
	}

	return resourceCdnFrontDoorRouteRead(d, meta)
}

func resourceCdnFrontDoorRouteDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := routes.ParseRouteID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandCdnFrontDoorRouteResourceReference(input string) *routes.ResourceReference {
	if len(input) == 0 {
		return nil
	}

	return &routes.ResourceReference{
		Id: pointer.To(input),
	}
}

func expandCdnFrontDoorRouteEnabled(input bool) *routes.EnabledState {
	if input {
		return pointer.To(routes.EnabledStateEnabled)
	}

	return pointer.To(routes.EnabledStateDisabled)
}

func expandCdnFrontDoorRouteHttpsRedirect(input bool) *routes.HTTPSRedirect {
	if input {
		return pointer.To(routes.HTTPSRedirectEnabled)
	}

	return pointer.To(routes.HTTPSRedirectDisabled)
}

func expandCdnFrontDoorRouteDefaultDomain(input bool) *routes.LinkToDefaultDomain {
	if input {
		return pointer.To(routes.LinkToDefaultDomainEnabled)
	}

	return pointer.To(routes.LinkToDefaultDomainDisabled)
}

func expandCdnFrontDoorRouteEndpointProtocolsArray(input []interface{}) *[]routes.AFDEndpointProtocols {
	results := make([]routes.AFDEndpointProtocols, 0)

	for _, item := range input {
		results = append(results, routes.AFDEndpointProtocols(item.(string)))
	}

	return &results
}

func expandCdnFrontdoorRouteRuleSetReferenceArray(input []interface{}) *[]routes.ResourceReference {
	// NOTE: The Frontdoor service, do not treat an empty object like an empty object
	// if it is not nil they assume it is fully defined and then end up throwing errors
	// when they attempt to get a value from one of the fields.
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]routes.ResourceReference, 0)
	for _, item := range input {
		results = append(results, routes.ResourceReference{
			Id: pointer.To(item.(string)),
		})
	}

	return &results
}

func expandCdnFrontdoorRouteCacheConfiguration(input []interface{}) *routes.AfdRouteCacheConfiguration {
	// NOTE: If this is not an explicit nil you will receive an "Unsupported QueryStringCachingBehavior type:
	// Property 'RouteV2.CacheConfiguration.QueryStringCachingBehavior' is required but it was not set" error.
	// The Frontdoor service treats empty slices as if they are fully defined unlike other services.
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	cacheConfiguration := &routes.AfdRouteCacheConfiguration{
		CompressionSettings: &routes.CompressionSettings{
			IsCompressionEnabled: pointer.To(v["compression_enabled"].(bool)),
		},
		QueryParameters:            expandStringSliceToCsvFormat(v["query_strings"].([]interface{})),
		QueryStringCachingBehavior: pointer.ToEnum[routes.AfdQueryStringCachingBehavior](v["query_string_caching_behavior"].(string)),
	}

	if contentTypes := v["content_types_to_compress"].([]interface{}); len(contentTypes) > 0 {
		cacheConfiguration.CompressionSettings.ContentTypesToCompress = utils.ExpandStringSlice(contentTypes)
	}

	return cacheConfiguration
}

func flattenCdnFrontDoorRouteRuleSetResourceArray(input *[]routes.ResourceReference) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	// Normalize these values in the configuration file we know they are valid because they were set on the
	// resource... if these are modified in the portal they will all be lowercased...
	// Issue: https://github.com/Azure/azure-sdk-for-go/issues/19378
	for _, ruleSet := range *input {
		id, err := rulesets.ParseRuleSetIDInsensitively(pointer.From(ruleSet.Id))
		if err != nil {
			return results, err
		}
		results = append(results, id.ID())
	}

	return results, nil
}

func flattenCdnFrontDoorRouteOriginGroupResourceReference(input *routes.ResourceReference) (string, error) {
	if input != nil && input.Id != nil {
		id, err := afdorigingroups.ParseOriginGroupID(*input.Id)
		if err != nil {
			return "", err
		}

		return id.ID(), nil
	}

	return "", nil
}

func flattenCdnFrontDoorRouteCustomDomainActivatedResourceArray(input *[]routes.ActivatedResourceReference) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return results, nil
	}

	// Normalize these values in the configuration file we know they are valid because they were set on the
	// resource... if these are modified in the portal they will all be lowercased...
	for _, customDomain := range *input {
		if customDomain.Id == nil {
			continue
		}
		id, err := afddomains.ParseCustomDomainIDInsensitively(*customDomain.Id)
		if err != nil {
			return nil, err
		}
		results = append(results, id.ID())
	}

	return results, nil
}

func flattenCdnFrontDoorRouteCacheConfiguration(input *routes.AfdRouteCacheConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	compressionEnabled := false
	contentTypesToCompress := make([]interface{}, 0)
	if v := input.CompressionSettings; v != nil {
		compressionEnabled = pointer.From(v.IsCompressionEnabled)
		contentTypesToCompress = utils.FlattenStringSlice(v.ContentTypesToCompress)
	}

	return []interface{}{
		map[string]interface{}{
			"compression_enabled":           compressionEnabled,
			"content_types_to_compress":     contentTypesToCompress,
			"query_string_caching_behavior": pointer.FromEnum(input.QueryStringCachingBehavior),
			"query_strings":                 flattenCsvToStringSlice(input.QueryParameters),
		},
	}
}
