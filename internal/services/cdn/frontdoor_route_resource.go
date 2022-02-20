package cdn

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/routes"
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
			_, err := routes.ParseRouteID(id)
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
				ValidateFunc: afdendpoints.ValidateAfdEndpointID,
			},

			"cache_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"query_strings": {
							Type:     pluginsdk.TypeList,
							Optional: true,

							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
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

						"enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"deployment_status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"endpoint_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"forwarding_protocol": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(routes.ForwardingProtocolHttpsOnly),
				ValidateFunc: validation.StringInSlice([]string{
					string(routes.ForwardingProtocolHttpOnly),
					string(routes.ForwardingProtocolHttpsOnly),
					string(routes.ForwardingProtocolMatchRequest),
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

			"origin_group_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
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

			"provisioning_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
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

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(routes.AFDEndpointProtocolsHttp),
						string(routes.AFDEndpointProtocolsHttps),
					}, false),
				},
			},
		},
	}
}

func resourceFrontdoorRouteCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRoutesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	afdEndpointId, err := afdendpoints.ParseAfdEndpointID(d.Get("frontdoor_endpoint_id").(string))
	if err != nil {
		return err
	}

	id := routes.NewRouteID(afdEndpointId.SubscriptionId, afdEndpointId.ResourceGroupName, afdEndpointId.ProfileName, afdEndpointId.EndpointName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_frontdoor_route", id.ID())
		}
	}

	forwardingProtocolValue := routes.ForwardingProtocol(d.Get("forwarding_protocol").(string))

	props := routes.Route{
		Properties: &routes.RouteProperties{
			CacheConfiguration:  expandRouteAfdRouteCacheConfiguration(d.Get("cache_configuration").([]interface{})),
			CustomDomains:       expandRouteActivatedResourceReferenceArray(d.Get("custom_domains").([]interface{})),
			EnabledState:        ConvertBoolToRoutesEnabledState(d.Get("enabled").(bool)),
			ForwardingProtocol:  &forwardingProtocolValue,
			HttpsRedirect:       ConvertBoolToRouteHttpsRedirect(d.Get("https_redirect").(bool)),
			LinkToDefaultDomain: ConvertBoolToRouteLinkToDefaultDomain(d.Get("link_to_default_domain").(bool)),
			OriginGroup:         *expandRouteResourceReference(d.Get("origin_group_id").(string)),
			PatternsToMatch:     utils.ExpandStringSlice(d.Get("patterns_to_match").([]interface{})),
			RuleSets:            expandRouteResourceReferenceArray(d.Get("rule_set_ids").([]interface{})),
			SupportedProtocols:  expandRouteAFDEndpointProtocolsArray(d.Get("supported_protocols").([]interface{})),
		},
	}

	if originPath := d.Get("origin_path").(string); originPath != "" {
		props.Properties.OriginPath = utils.String(originPath)
	}

	if err := client.CreateThenPoll(ctx, id, props); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceFrontdoorRouteRead(d, meta)
}

func resourceFrontdoorRouteRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRoutesClient
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

	d.Set("name", id.RouteName)

	d.Set("frontdoor_endpoint_id", afdendpoints.NewAfdEndpointID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.EndpointName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("deployment_status", props.DeploymentStatus)
			d.Set("enabled", ConvertRoutesEnabledStateToBool(props.EnabledState))
			d.Set("forwarding_protocol", props.ForwardingProtocol)
			d.Set("https_redirect", ConvertRouteHttpsRedirectToBool(props.HttpsRedirect))
			d.Set("link_to_default_domain", ConvertRouteLinkToDefaultDomainToBool(props.LinkToDefaultDomain))
			d.Set("origin_path", props.OriginPath)
			d.Set("patterns_to_match", props.PatternsToMatch)
			d.Set("provisioning_state", props.ProvisioningState)

			// BUG: Endpoint name is not being returned by the API
			d.Set("endpoint_name", id.EndpointName)

			if err := d.Set("cache_configuration", flattenFrontdoorRouteCacheConfiguration(props.CacheConfiguration)); err != nil {
				return fmt.Errorf("setting `cache_configuration`: %+v", err)
			}

			if err := d.Set("custom_domains", flattenRouteActivatedResourceReferenceArray(props.CustomDomains)); err != nil {
				return fmt.Errorf("setting `custom_domains`: %+v", err)
			}

			if err := d.Set("origin_group_id", flattenRouteResourceReference(&props.OriginGroup)); err != nil {
				return fmt.Errorf("setting `origin_group_id`: %+v", err)
			}

			if err := d.Set("rule_set_ids", flattenRouteResourceReferenceArry(props.RuleSets)); err != nil {
				return fmt.Errorf("setting `rule_set_ids`: %+v", err)
			}

			if err := d.Set("supported_protocols", flattenRouteAFDEndpointProtocolsArray(props.SupportedProtocols)); err != nil {
				return fmt.Errorf("setting `supported_protocols`: %+v", err)
			}
		}
	}
	return nil
}

func resourceFrontdoorRouteUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRoutesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := routes.ParseRouteID(d.Id())
	if err != nil {
		return err
	}

	forwardingProtocolValue := routes.ForwardingProtocol(d.Get("forwarding_protocol").(string))

	props := routes.RouteUpdateParameters{
		Properties: &routes.RouteUpdatePropertiesParameters{
			CacheConfiguration:  expandRouteAfdRouteCacheConfiguration(d.Get("cache_configuration").([]interface{})),
			CustomDomains:       expandRouteActivatedResourceReferenceArray(d.Get("custom_domains").([]interface{})),
			EnabledState:        ConvertBoolToRoutesEnabledState(d.Get("enabled").(bool)),
			ForwardingProtocol:  &forwardingProtocolValue,
			HttpsRedirect:       ConvertBoolToRouteHttpsRedirect(d.Get("https_redirect").(bool)),
			LinkToDefaultDomain: ConvertBoolToRouteLinkToDefaultDomain(d.Get("link_to_default_domain").(bool)),
			OriginGroup:         expandRouteResourceReference(d.Get("origin_group_id").(string)),
			PatternsToMatch:     utils.ExpandStringSlice(d.Get("patterns_to_match").([]interface{})),
			RuleSets:            expandRouteResourceReferenceArray(d.Get("rule_set_ids").([]interface{})),
			SupportedProtocols:  expandRouteAFDEndpointProtocolsArray(d.Get("supported_protocols").([]interface{})),
		},
	}

	if originPath := d.Get("origin_path").(string); originPath != "" {
		props.Properties.OriginPath = &originPath
	}

	if err := client.UpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceFrontdoorRouteRead(d, meta)
}

func resourceFrontdoorRouteDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRoutesClient
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

func expandRouteResourceReference(input string) *routes.ResourceReference {
	if len(input) == 0 || input == "" {
		return nil
	}

	return &routes.ResourceReference{
		Id: utils.String(input),
	}
}

func expandRouteResourceReferenceArray(input []interface{}) *[]routes.ResourceReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]routes.ResourceReference, 0)

	for _, item := range input {
		results = append(results, routes.ResourceReference{
			Id: utils.String(item.(string)),
		})
	}

	return &results
}

func expandRouteAFDEndpointProtocolsArray(input []interface{}) *[]routes.AFDEndpointProtocols {
	results := make([]routes.AFDEndpointProtocols, 0)

	for _, item := range input {
		results = append(results, routes.AFDEndpointProtocols(item.(string)))
	}

	return &results
}

func expandRouteAfdRouteCacheConfiguration(input []interface{}) *routes.AfdRouteCacheConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	queryStringCachingBehaviorValue := routes.AfdQueryStringCachingBehavior(v["query_string_caching_behavior"].(string))
	return &routes.AfdRouteCacheConfiguration{
		QueryParameters:            expandFrontdoorRouteQueryParameters(v["query_strings"].([]interface{})),
		QueryStringCachingBehavior: &queryStringCachingBehaviorValue,
	}
}

func expandFrontdoorRouteQueryParameters(input []interface{}) *string {
	if len(input) == 0 {
		return nil
	}

	v := utils.ExpandStringSlice(input)
	csv := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(*v)), ","), "[]")

	return &csv
}

func expandRouteActivatedResourceReferenceArray(input []interface{}) *[]routes.ActivatedResourceReference {
	results := make([]routes.ActivatedResourceReference, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, routes.ActivatedResourceReference{
			Id: utils.String(v["id"].(string)),
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

func flattenRouteActivatedResourceReferenceArray(inputs *[]routes.ActivatedResourceReference) []interface{} {
	results := make([]interface{}, 0)
	if inputs == nil {
		return results
	}

	for _, input := range *inputs {
		result := make(map[string]interface{})

		if input.Id != nil {
			result["id"] = *input.Id
		}

		if input.IsActive != nil {
			result["enabled"] = *input.IsActive
		}
		results = append(results, result)
	}

	return results
}

func flattenRouteResourceReference(input *routes.ResourceReference) string {
	result := ""
	if input == nil {
		return result
	}

	if input.Id != nil {
		result = *input.Id
	}

	return result
}

func flattenRouteResourceReferenceArry(input *[]routes.ResourceReference) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if item.Id != nil {
			results = append(results, *item.Id)
		}
	}

	return results
}

func flattenRouteAFDEndpointProtocolsArray(input *[]routes.AFDEndpointProtocols) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, item)
	}

	return results
}

func flattenFrontdoorRouteCacheConfiguration(input *routes.AfdRouteCacheConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	if input.QueryParameters != nil {
		result["query_strings"] = flattenFrontdoorRouteQueryParameters(input.QueryParameters)
	}

	if input.QueryStringCachingBehavior != nil {
		result["query_string_caching_behavior"] = *input.QueryStringCachingBehavior
	}

	return append(results, result)
}
