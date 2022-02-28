package cdn

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceFrontdoorRoute() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontdoorRouteCreate,
		Read:   resourceFrontdoorRouteRead,
		Update: resourceFrontdoorRouteUpdate,
		Delete: resourceFrontdoorRouteDelete,

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
			},

			"frontdoor_endpoint_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorEndpointID,
			},

			"frontdoor_origin_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorOriginGroupID,
			},

			"frontdoor_origin_ids": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,

				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.FrontdoorOriginID,
				},
			},

			// NOTE: AfdRouteCacheConfiguration to disable caching, do not provide block in API call.
			"cache_configuration": {
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
							Default:  string(track1.AfdQueryStringCachingBehaviorIgnoreQueryString),
							ValidateFunc: validation.StringInSlice([]string{
								string(track1.AfdQueryStringCachingBehaviorIgnoreQueryString),
								string(track1.AfdQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings),
								string(track1.AfdQueryStringCachingBehaviorIncludeSpecifiedQueryStrings),
								string(track1.AfdQueryStringCachingBehaviorUseQueryString),
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
								ValidateFunc: ValidateContentTypes,
							},
						},
					},
				},
			},

			"custom_domains": {
				Type:     pluginsdk.TypeList,
				Optional: true,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"id": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"active": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
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
				Default:  string(track1.ForwardingProtocolHTTPSOnly),
				ValidateFunc: validation.StringInSlice([]string{
					string(track1.ForwardingProtocolHTTPOnly),
					string(track1.ForwardingProtocolHTTPSOnly),
					string(track1.ForwardingProtocolMatchRequest),
				}, false),
			},

			"https_redirect": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"link_to_default_domain": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"origin_path": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"patterns_to_match": {
				Type:     pluginsdk.TypeList,
				Optional: true,

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"rule_set_ids": {
				Type:     pluginsdk.TypeList,
				Optional: true,

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"supported_protocols": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 2,

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(track1.AFDEndpointProtocolsHTTP),
						string(track1.AFDEndpointProtocolsHTTPS),
					}, false),
				},
			},

			"frontdoor_endpoint_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFrontdoorRouteCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRoutesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	afdEndpointId, err := parse.FrontdoorEndpointID(d.Get("frontdoor_endpoint_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFrontdoorRouteID(afdEndpointId.SubscriptionId, afdEndpointId.ResourceGroup, afdEndpointId.ProfileName, afdEndpointId.AfdEndpointName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_frontdoor_route", id.ID())
		}
	}

	props := track1.Route{
		RouteProperties: &track1.RouteProperties{
			CustomDomains:       expandRouteActivatedResourceReferenceArray(d.Get("custom_domains").([]interface{})),
			EnabledState:        ConvertBoolToEnabledState(d.Get("enabled").(bool)),
			ForwardingProtocol:  track1.ForwardingProtocol(d.Get("forwarding_protocol").(string)),
			HTTPSRedirect:       ConvertBoolToRouteHttpsRedirect(d.Get("https_redirect").(bool)),
			LinkToDefaultDomain: ConvertBoolToRouteLinkToDefaultDomain(d.Get("link_to_default_domain").(bool)),
			OriginGroup:         expandResourceReference(d.Get("frontdoor_origin_group_id").(string)),
			PatternsToMatch:     utils.ExpandStringSlice(d.Get("patterns_to_match").([]interface{})),
			RuleSets:            expandRouteResourceReferenceArray(d.Get("rule_set_ids").([]interface{})),
			SupportedProtocols:  expandRouteAFDEndpointProtocolsArray(d.Get("supported_protocols").([]interface{})),
		},
	}

	if cacheConfiguration := expandRouteAfdRouteCacheConfiguration(d.Get("cache_configuration").([]interface{})); cacheConfiguration != nil {
		props.RouteProperties.CacheConfiguration = cacheConfiguration
	}

	if originPath := d.Get("origin_path").(string); originPath != "" {
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
	d.Set("frontdoor_origin_ids", utils.ExpandStringSlice(d.Get("frontdoor_origin_ids").([]interface{})))

	return resourceFrontdoorRouteRead(d, meta)
}

func resourceFrontdoorRouteRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRoutesClient
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

	d.Set("name", id.RouteName)

	d.Set("frontdoor_endpoint_id", parse.NewFrontdoorEndpointID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.AfdEndpointName).ID())

	if props := resp.RouteProperties; props != nil {
		d.Set("enabled", ConvertEnabledStateToBool(&props.EnabledState))
		d.Set("forwarding_protocol", props.ForwardingProtocol)
		d.Set("https_redirect", ConvertRouteHttpsRedirectToBool(&props.HTTPSRedirect))
		d.Set("link_to_default_domain", ConvertRouteLinkToDefaultDomainToBool(&props.LinkToDefaultDomain))
		d.Set("origin_path", props.OriginPath)
		d.Set("patterns_to_match", props.PatternsToMatch)

		// BUG: Endpoint name is not being returned by the API
		d.Set("frontdoor_endpoint_name", id.AfdEndpointName)

		if err := d.Set("cache_configuration", flattenFrontdoorRouteCacheConfiguration(props.CacheConfiguration)); err != nil {
			return fmt.Errorf("setting `cache_configuration`: %+v", err)
		}

		if err := d.Set("custom_domains", flattenRouteActivatedResourceReferenceArray(props.CustomDomains)); err != nil {
			return fmt.Errorf("setting `custom_domains`: %+v", err)
		}

		if err := d.Set("frontdoor_origin_group_id", flattenResourceReference(props.OriginGroup)); err != nil {
			return fmt.Errorf("setting `frontdoor_origin_group_id`: %+v", err)
		}

		if err := d.Set("rule_set_ids", flattenRouteResourceReferenceArry(props.RuleSets)); err != nil {
			return fmt.Errorf("setting `rule_set_ids`: %+v", err)
		}

		if err := d.Set("supported_protocols", flattenRouteAFDEndpointProtocolsArray(props.SupportedProtocols)); err != nil {
			return fmt.Errorf("setting `supported_protocols`: %+v", err)
		}
	}

	return nil
}

func resourceFrontdoorRouteUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRoutesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorRouteID(d.Id())
	if err != nil {
		return err
	}

	props := track1.RouteUpdateParameters{
		RouteUpdatePropertiesParameters: &track1.RouteUpdatePropertiesParameters{
			CustomDomains:       expandRouteActivatedResourceReferenceArray(d.Get("custom_domains").([]interface{})),
			EnabledState:        ConvertBoolToEnabledState(d.Get("enabled").(bool)),
			ForwardingProtocol:  track1.ForwardingProtocol(d.Get("forwarding_protocol").(string)),
			HTTPSRedirect:       ConvertBoolToRouteHttpsRedirect(d.Get("https_redirect").(bool)),
			LinkToDefaultDomain: ConvertBoolToRouteLinkToDefaultDomain(d.Get("link_to_default_domain").(bool)),
			OriginGroup:         expandResourceReference(d.Get("frontdoor_origin_group_id").(string)),
			PatternsToMatch:     utils.ExpandStringSlice(d.Get("patterns_to_match").([]interface{})),
			RuleSets:            expandRouteResourceReferenceArray(d.Get("rule_set_ids").([]interface{})),
			SupportedProtocols:  expandRouteAFDEndpointProtocolsArray(d.Get("supported_protocols").([]interface{})),
		},
	}

	if cacheConfiguration := expandRouteAfdRouteCacheConfiguration(d.Get("cache_configuration").([]interface{})); cacheConfiguration != nil {
		props.RouteUpdatePropertiesParameters.CacheConfiguration = cacheConfiguration
	}

	if originPath := d.Get("origin_path").(string); originPath != "" {
		props.RouteUpdatePropertiesParameters.OriginPath = &originPath
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName, props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceFrontdoorRouteRead(d, meta)
}

func resourceFrontdoorRouteDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRoutesClient
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

func expandRouteResourceReferenceArray(input []interface{}) *[]track1.ResourceReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]track1.ResourceReference, 0)

	for _, item := range input {
		results = append(results, track1.ResourceReference{
			ID: utils.String(item.(string)),
		})
	}

	return &results
}

func expandRouteAFDEndpointProtocolsArray(input []interface{}) *[]track1.AFDEndpointProtocols {
	results := make([]track1.AFDEndpointProtocols, 0)

	for _, item := range input {
		results = append(results, track1.AFDEndpointProtocols(item.(string)))
	}

	return &results
}

func expandRouteAfdRouteCacheConfiguration(input []interface{}) *track1.AfdRouteCacheConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	queryStringCachingBehaviorValue := track1.AfdQueryStringCachingBehavior(v["query_string_caching_behavior"].(string))
	comprssionEnabled := v["compression_enabled"].(bool)

	cacheConfiguration := &track1.AfdRouteCacheConfiguration{
		QueryParameters:            expandFrontdoorRouteQueryParameters(v["query_strings"].([]interface{})),
		QueryStringCachingBehavior: queryStringCachingBehaviorValue,
	}

	compressionSettings := &track1.CompressionSettings{}
	compressionSettings.IsCompressionEnabled = utils.Bool(comprssionEnabled)

	if contentTypes := v["content_types_to_compress"].([]interface{}); len(contentTypes) > 0 {
		compressionSettings.ContentTypesToCompress = utils.ExpandStringSlice(contentTypes)
	}

	cacheConfiguration.CompressionSettings = compressionSettings

	return cacheConfiguration
}

func expandFrontdoorRouteQueryParameters(input []interface{}) *string {
	if len(input) == 0 {
		return nil
	}

	v := utils.ExpandStringSlice(input)
	csv := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(*v)), ","), "[]")

	return &csv
}

func expandRouteActivatedResourceReferenceArray(input []interface{}) *[]track1.ActivatedResourceReference {
	results := make([]track1.ActivatedResourceReference, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, track1.ActivatedResourceReference{
			ID: utils.String(v["id"].(string)),
		})
	}
	return &results
}

func flattenFrontdoorRouteQueryParameters(input *string) []interface{} {
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

func flattenRouteActivatedResourceReferenceArray(inputs *[]track1.ActivatedResourceReference) []interface{} {
	results := make([]interface{}, 0)
	if inputs == nil {
		return results
	}

	for _, input := range *inputs {
		result := make(map[string]interface{})

		if input.ID != nil {
			result["id"] = *input.ID
		}

		if input.IsActive != nil {
			result["active"] = *input.IsActive
		}
		results = append(results, result)
	}

	return results
}

func flattenRouteResourceReferenceArry(input *[]track1.ResourceReference) []interface{} {
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

func flattenRouteAFDEndpointProtocolsArray(input *[]track1.AFDEndpointProtocols) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, item)
	}

	return results
}

func flattenFrontdoorRouteCacheConfiguration(input *track1.AfdRouteCacheConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	if input.QueryParameters != nil {
		result["query_strings"] = flattenFrontdoorRouteQueryParameters(input.QueryParameters)
	}

	if input.QueryStringCachingBehavior != "" {
		result["query_string_caching_behavior"] = input.QueryStringCachingBehavior
	}

	if input.CompressionSettings != nil {
		compressionSettings := input.CompressionSettings.(map[string]interface{})
		compressionEnabled := compressionSettings["isCompressionEnabled"].(bool)
		contentTypesToCompress := compressionSettings["contentTypesToCompress"].([]interface{})

		contentTypes := make([]string, 0)
		for _, contentType := range contentTypesToCompress {
			contentTypes = append(contentTypes, contentType.(string))
		}

		result["compression_enabled"] = compressionEnabled
		result["content_types_to_compress"] = utils.FlattenStringSlice(&contentTypes)
	}

	return append(results, result)
}
