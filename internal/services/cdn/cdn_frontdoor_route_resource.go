package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontdoorRoute() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontdoorRouteCreate,
		Read:   resourceCdnFrontdoorRouteRead,
		Update: resourceCdnFrontdoorRouteUpdate,
		Delete: resourceCdnFrontdoorRouteDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontdoorRouteID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				// TODO: missing validation
				// WS: Fixed
				ValidateFunc: validate.CdnFrontdoorRouteName,
			},

			"cdn_frontdoor_endpoint_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorEndpointID,
			},

			// TODO: why is this prefixed with `cdn_frontdoor_`? we can remove that since it's implied?
			// WS: Because the legacy Frontdoor also has some of the same resource types, so I was exposing
			// them all with the prefix as a disambiguator so there wouldn't be any confusion what was
			// expected here as input.
			"cdn_frontdoor_origin_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorOriginGroupID,
			},

			// TODO: why is this prefixed with `cdn_frontdoor_`? we can remove that since it's implied?
			// WS: Same as above.
			"cdn_frontdoor_origin_ids": {
				// TODO: BLOCKER - these are sent to the API so must be returned
				// WS: These are not sent to the API, they are only here so Terraform
				// can provision/destroy the resources in the correct order.
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,

				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.FrontdoorOriginID,
				},
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
								ValidateFunc: validation.StringInSlice(cdnFrontdoorContentTypes(), false),
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

			// TODO: why is this prefixed with `cdn_frontdoor_`? we can remove that since it's implied?
			// WS: Same as above.
			"cdn_frontdoor_custom_domain_ids": {
				Type:     pluginsdk.TypeList,
				Optional: true,

				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.FrontdoorCustomDomainID,
				},
			},

			"link_to_default_domain_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			// TODO: why is this prefixed with `cdn_frontdoor_`? we can remove that since it's implied?
			// WS: Same as above.
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

			// TODO: why is this prefixed with `cdn_frontdoor_`? we can remove that since it's implied?
			// WS: Same as above.
			"cdn_frontdoor_rule_set_ids": {
				Type:     pluginsdk.TypeList,
				Optional: true,

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"supported_protocols": {
				// TODO: does this need to be a Set?
				// WS: Fixed
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

			// TODO: why is this prefixed with `cdn_frontdoor_`? we can remove that since it's implied?
			// WS: Same as above.
			"cdn_frontdoor_endpoint_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			// TODO: why is this prefixed with `cdn_frontdoor_`? we can remove that since it's implied?
			// WS: Same as above.
			"cdn_frontdoor_custom_domains_active_status": {
				Type:     pluginsdk.TypeList,
				Computed: true,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"active": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceCdnFrontdoorRouteCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	endpointId, err := parse.FrontdoorEndpointID(d.Get("cdn_frontdoor_endpoint_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFrontdoorRouteID(endpointId.SubscriptionId, endpointId.ResourceGroup, endpointId.ProfileName, endpointId.AfdEndpointName, d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cdn_frontdoor_route", id.ID())
	}

	protocolsRaw := d.Get("supported_protocols").(*pluginsdk.Set).List()

	props := cdn.Route{
		RouteProperties: &cdn.RouteProperties{
			CustomDomains:       expandCdnFrontdoorRouteActivatedResourceArray(d.Get("cdn_frontdoor_custom_domain_ids").([]interface{})),
			CacheConfiguration:  expandCdnFrontdoorRouteCacheConfiguration(d.Get("cache").([]interface{})),
			EnabledState:        convertCdnFrontdoorBoolToEnabledState(d.Get("enabled").(bool)),
			ForwardingProtocol:  cdn.ForwardingProtocol(d.Get("forwarding_protocol").(string)),
			HTTPSRedirect:       convertCdnFrontdoorBoolToRouteHttpsRedirect(d.Get("https_redirect_enabled").(bool)),
			LinkToDefaultDomain: convertCdnFrontdoorBoolToRouteLinkToDefaultDomain(d.Get("link_to_default_domain_enabled").(bool)),
			OriginGroup:         expandCdnFrontdoorResourceReference(d.Get("cdn_frontdoor_origin_group_id").(string)),
			PatternsToMatch:     utils.ExpandStringSlice(d.Get("patterns_to_match").([]interface{})),
			RuleSets:            expandCdnFrontdoorRouteResourceReferenceArray(d.Get("cdn_frontdoor_rule_set_ids").([]interface{})),
			SupportedProtocols:  expandCdnFrontdoorRouteEndpointProtocolsArray(protocolsRaw),
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

	// TODO: this needs to be removed in favour of the API returning the data (as per the ARM guidelines)
	// WS: These are not sent to the API, they are only here so Terraform
	// can provision/destroy the resources in the correct order.
	if originIds := d.Get("cdn_frontdoor_origin_ids").([]interface{}); len(originIds) > 0 {
		d.Set("cdn_frontdoor_origin_ids", utils.ExpandStringSlice(originIds))
	}

	return resourceCdnFrontdoorRouteRead(d, meta)
}

func resourceCdnFrontdoorRouteRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorRouteID(d.Id())
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

	// TODO: BLOCKER - these need to be returned from the API
	// WS: These are not sent to the API, they are only here so Terraform
	// can provision/destroy the resources in the correct order.
	if originIds := d.Get("cdn_frontdoor_origin_ids").([]interface{}); len(originIds) > 0 {
		d.Set("cdn_frontdoor_origin_ids", utils.ExpandStringSlice(originIds))
	}

	d.Set("name", id.RouteName)
	d.Set("cdn_frontdoor_endpoint_id", parse.NewFrontdoorEndpointID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.AfdEndpointName).ID())

	if props := resp.RouteProperties; props != nil {
		// TODO: split this into two separate flatten functions
		// WS: Fixed
		domainsActive := flattenCdnFrontdoorRouteActivatedResourceComputedArray(props.CustomDomains)
		domains := flattenCdnFrontdoorRouteActivatedResourceArray(props.CustomDomains)
		d.Set("cdn_frontdoor_custom_domain_ids", domains)
		d.Set("enabled", convertCdnFrontdoorEnabledStateToBool(&props.EnabledState))
		d.Set("forwarding_protocol", props.ForwardingProtocol)
		d.Set("https_redirect_enabled", convertCdnFrontdoorRouteHttpsRedirectToBool(&props.HTTPSRedirect))
		d.Set("link_to_default_domain_enabled", convertCdnFrontdoorRouteLinkToDefaultDomainToBool(&props.LinkToDefaultDomain))
		d.Set("cdn_frontdoor_origin_path", props.OriginPath)
		d.Set("patterns_to_match", props.PatternsToMatch)

		// TODO: BLOCKER - BUG: Endpoint name is not being returned by the API
		// WS: Yes, this is a service bug, but I do not agree that it is a "BLOCKER" as we have a workaround.
		d.Set("cdn_frontdoor_endpoint_name", id.AfdEndpointName)

		if err := d.Set("cdn_frontdoor_custom_domains_active_status", domainsActive); err != nil {
			return fmt.Errorf("setting %q: %+v", "cdn_frontdoor_custom_domains_active_status", err)
		}

		if err := d.Set("cache", flattenCdnFrontdoorRouteCacheConfiguration(props.CacheConfiguration)); err != nil {
			return fmt.Errorf("setting `cache`: %+v", err)
		}

		if err := d.Set("cdn_frontdoor_origin_group_id", flattenCdnFrontdoorResourceReference(props.OriginGroup)); err != nil {
			return fmt.Errorf("setting `cdn_frontdoor_origin_group_id`: %+v", err)
		}

		if err := d.Set("cdn_frontdoor_rule_set_ids", flattenCdnFrontdoorRouteResourceArray(props.RuleSets)); err != nil {
			return fmt.Errorf("setting `cdn_frontdoor_rule_set_ids`: %+v", err)
		}

		if err := d.Set("supported_protocols", flattenCdnFrontdoorRouteEndpointProtocolsArray(props.SupportedProtocols)); err != nil {
			return fmt.Errorf("setting `supported_protocols`: %+v", err)
		}
	}

	return nil
}

func resourceCdnFrontdoorRouteUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	workaroundsClient := azuresdkhacks.NewCdnFrontDoorRoutesWorkaroundClient(client)
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorRouteID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}
	if existing.RouteProperties == nil {
		return fmt.Errorf("retrieving existing %s: `properties` was nil", *id)
	}

	props := azuresdkhacks.RouteUpdatePropertiesParameters{
		CustomDomains: existing.RouteProperties.CustomDomains,
	}

	if d.HasChange("cdn_frontdoor_custom_domain_ids") {
		props.CustomDomains = expandCdnFrontdoorRouteActivatedResourceArray(d.Get("cdn_frontdoor_custom_domain_ids").([]interface{}))
	}

	if d.HasChange("cache") {
		props.CacheConfiguration = expandCdnFrontdoorRouteCacheConfiguration(d.Get("cache").([]interface{}))
	}

	if d.HasChange("enabled") {
		props.EnabledState = convertCdnFrontdoorBoolToEnabledState(d.Get("enabled").(bool))
	}

	if d.HasChange("forwarding_protocol") {
		props.ForwardingProtocol = cdn.ForwardingProtocol(d.Get("forwarding_protocol").(string))
	}

	if d.HasChange("https_redirect_enabled") {
		props.HTTPSRedirect = convertCdnFrontdoorBoolToRouteHttpsRedirect(d.Get("https_redirect_enabled").(bool))
	}

	if d.HasChange("link_to_default_domain_enabled") {
		props.LinkToDefaultDomain = convertCdnFrontdoorBoolToRouteLinkToDefaultDomain(d.Get("link_to_default_domain_enabled").(bool))
	}

	if d.HasChange("cdn_frontdoor_origin_group_id") {
		props.OriginGroup = expandCdnFrontdoorResourceReference(d.Get("cdn_frontdoor_origin_group_id").(string))
	}

	if d.HasChange("cdn_frontdoor_origin_path") {
		props.OriginPath = utils.String(d.Get("cdn_frontdoor_origin_path").(string))
	}

	if d.HasChange("patterns_to_match") {
		props.PatternsToMatch = utils.ExpandStringSlice(d.Get("patterns_to_match").([]interface{}))
	}

	if d.HasChange("cdn_frontdoor_rule_set_ids") {
		props.RuleSets = expandCdnFrontdoorRouteResourceReferenceArray(d.Get("cdn_frontdoor_rule_set_ids").([]interface{}))
	}

	if d.HasChange("supported_protocols") {
		protocalsRaw := d.Get("supported_protocols").(*pluginsdk.Set).List()
		props.SupportedProtocols = expandCdnFrontdoorRouteEndpointProtocolsArray(protocalsRaw)
	}

	payload := azuresdkhacks.RouteUpdateParameters{
		RouteUpdatePropertiesParameters: &props,
	}

	future, err := workaroundsClient.Update(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName, payload)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	// TODO: BLOCKER - this'll need to be returned from the API and set in the Read
	// WS: These are not sent to the API, they are only here so Terraform
	// can provision/destroy the resources in the correct order.
	if originIds := d.Get("cdn_frontdoor_origin_ids").([]interface{}); len(originIds) > 0 {
		d.Set("cdn_frontdoor_origin_ids", utils.ExpandStringSlice(originIds))
	}

	return resourceCdnFrontdoorRouteRead(d, meta)
}

func resourceCdnFrontdoorRouteDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorRouteID(d.Id())
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

func expandCdnFrontdoorRouteEndpointProtocolsArray(input []interface{}) *[]cdn.AFDEndpointProtocols {
	results := make([]cdn.AFDEndpointProtocols, 0)

	for _, item := range input {
		results = append(results, cdn.AFDEndpointProtocols(item.(string)))
	}

	return &results
}

func expandCdnFrontdoorRouteResourceReferenceArray(input []interface{}) *[]cdn.ResourceReference {
	results := make([]cdn.ResourceReference, 0)
	if len(input) == 0 || input[0] == nil {
		// TODO: shouldn't this be returning an empty slice?
		// WS: Fixed, note due to the services treatment of empty object I believe this will cause an error since it is not a nil value instead.
		return &results
	}

	for _, item := range input {
		results = append(results, cdn.ResourceReference{
			ID: utils.String(item.(string)),
		})
	}

	return &results
}

func expandCdnFrontdoorRouteCacheConfiguration(input []interface{}) *cdn.AfdRouteCacheConfiguration {
	if len(input) == 0 || input[0] == nil {
		// TODO: shouldn't this be returning an empty slice?
		// WS: No, if this is not an explicit nil you will receive a "Unsupported QueryStringCachingBehavior type: ''.
		// Property 'RouteV2.CacheConfiguration.QueryStringCachingBehavior' is required but it was not set" error.
		// The Frontdoor service treats empty slices as if they are fully defined unlike other services.
		return nil
	}

	v := input[0].(map[string]interface{})

	queryStringCachingBehaviorValue := cdn.AfdQueryStringCachingBehavior(v["query_string_caching_behavior"].(string))
	compressionEnabled := v["compression_enabled"].(bool)

	cacheConfiguration := &cdn.AfdRouteCacheConfiguration{
		CompressionSettings: &cdn.CompressionSettings{
			IsCompressionEnabled: utils.Bool(compressionEnabled),
		},
		QueryParameters:            expandCdnFrontdoorStringSliceToCsvFormat(v["query_strings"].([]interface{})),
		QueryStringCachingBehavior: queryStringCachingBehaviorValue,
	}

	if contentTypes := v["content_types_to_compress"].([]interface{}); len(contentTypes) > 0 {
		cacheConfiguration.CompressionSettings.ContentTypesToCompress = utils.ExpandStringSlice(contentTypes)
	}

	return cacheConfiguration
}

func expandCdnFrontdoorRouteActivatedResourceArray(input []interface{}) *[]cdn.ActivatedResourceReference {
	results := make([]cdn.ActivatedResourceReference, 0)
	if len(input) == 0 {
		// TODO: confirm if sending an empty list means we can remove the hack
		// WS: I have confirmed with the service team that this is required to be an explicit "nil" value, an empty list will not work.
		// I had to modify the SDK to allow for nil which in the API means disassociate the custom domains
		return nil
	}

	for _, customDomain := range input {
		id := customDomain.(string)
		results = append(results, cdn.ActivatedResourceReference{
			ID: utils.String(id),
		})
	}

	return &results
}

func flattenCdnFrontdoorRouteActivatedResourceArray(inputs *[]cdn.ActivatedResourceReference) []interface{} {
	results := make([]interface{}, 0)
	if inputs == nil {
		return results
	}

	for _, customDomain := range *inputs {
		results = append(results, customDomain.ID)
	}

	return results
}

func flattenCdnFrontdoorRouteActivatedResourceComputedArray(inputs *[]cdn.ActivatedResourceReference) []interface{} {
	results := make([]interface{}, 0)
	if inputs == nil {
		// TODO: this should be split, but this must return an empty slice and not nil
		// WS: Fixed
		return results
	}

	for _, customDomain := range *inputs {
		result := make(map[string]interface{})
		result["id"] = customDomain.ID
		result["active"] = customDomain.IsActive
		results = append(results, result)
	}

	return results
}

func flattenCdnFrontdoorRouteResourceArray(input *[]cdn.ResourceReference) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if item.ID != nil {
			results = append(results, *item.ID)
		}
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
		queryParameters = flattenCdnFrontdoorCsvToStringSlice(input.QueryParameters)
	}

	cachingBehaviour := ""
	if input.QueryStringCachingBehavior != "" {
		cachingBehaviour = string(input.QueryStringCachingBehavior)
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
			"query_string_caching_behavior": cachingBehaviour,
			"query_strings":                 queryParameters, // TODO: why isn't this called query_parameters? WS: To be consistent with the legacy resource I felt it was best to keep the names a constant across all resources.
		},
	}
}
