// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnEndpointCreate,
		Read:   resourceCdnEndpointRead,
		Update: resourceCdnEndpointUpdate,
		Delete: resourceCdnEndpointDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.CdnEndpointV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.EndpointID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"profile_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"origin_host_header": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"is_http_allowed": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"is_https_allowed": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"origin": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
						},

						"host_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
						},

						"http_port": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  80,
						},

						"https_port": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  443,
						},
					},
				},
			},

			"origin_path": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: !features.FourPointOhBeta(),
			},

			"querystring_caching_behaviour": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(cdn.QueryStringCachingBehaviorIgnoreQueryString),
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.QueryStringCachingBehaviorBypassCaching),
					string(cdn.QueryStringCachingBehaviorIgnoreQueryString),
					string(cdn.QueryStringCachingBehaviorNotSet),
					string(cdn.QueryStringCachingBehaviorUseQueryString),
				}, false),
			},

			"content_types_to_compress": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: !features.FourPointOhBeta(),
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
				Set: pluginsdk.HashString,
			},

			"is_compression_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"probe_path": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: !features.FourPointOhBeta(),
			},

			"geo_filter": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"relative_path": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.ActionTypeAllow),
								string(cdn.ActionTypeBlock),
							}, false),
						},
						"country_codes": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"optimization_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.OptimizationTypeDynamicSiteAcceleration),
					string(cdn.OptimizationTypeGeneralMediaStreaming),
					string(cdn.OptimizationTypeGeneralWebDelivery),
					string(cdn.OptimizationTypeLargeFileDownload),
					string(cdn.OptimizationTypeVideoOnDemandMediaStreaming),
				}, false),
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"global_delivery_rule": endpointGlobalDeliveryRule(),

			"delivery_rule": endpointDeliveryRule(),

			"tags": tags.Schema(),
		},
	}
}

func resourceCdnEndpointCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	endpointsClient := meta.(*clients.Client).Cdn.EndpointsClient
	profilesClient := meta.(*clients.Client).Cdn.ProfilesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM CDN EndPoint creation.")

	id := parse.NewEndpointID(subscriptionId, d.Get("resource_group_name").(string), d.Get("profile_name").(string), d.Get("name").(string))
	existing, err := endpointsClient.Get(ctx, id.ResourceGroup, id.ProfileName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cdn_endpoint", id.ID())
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	httpAllowed := d.Get("is_http_allowed").(bool)
	httpsAllowed := d.Get("is_https_allowed").(bool)
	cachingBehaviour := d.Get("querystring_caching_behaviour").(string)
	originPath := d.Get("origin_path").(string)
	probePath := d.Get("probe_path").(string)
	optimizationType := d.Get("optimization_type").(string)
	t := d.Get("tags").(map[string]interface{})

	endpoint := cdn.Endpoint{
		Location: &location,
		EndpointProperties: &cdn.EndpointProperties{
			IsHTTPAllowed:              &httpAllowed,
			IsHTTPSAllowed:             &httpsAllowed,
			QueryStringCachingBehavior: cdn.QueryStringCachingBehavior(cachingBehaviour),
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("origin_host_header"); ok {
		endpoint.EndpointProperties.OriginHostHeader = utils.String(v.(string))
	}

	if _, ok := d.GetOk("content_types_to_compress"); ok {
		contentTypes := expandArmCdnEndpointContentTypesToCompress(d)
		endpoint.EndpointProperties.ContentTypesToCompress = &contentTypes
	}

	if _, ok := d.GetOk("geo_filter"); ok {
		geoFilters := expandCdnEndpointGeoFilters(d)
		endpoint.EndpointProperties.GeoFilters = geoFilters
	}

	if v, ok := d.GetOk("is_compression_enabled"); ok {
		endpoint.EndpointProperties.IsCompressionEnabled = utils.Bool(v.(bool))
	}

	if optimizationType != "" {
		endpoint.EndpointProperties.OptimizationType = cdn.OptimizationType(optimizationType)
	}

	if originPath != "" {
		endpoint.EndpointProperties.OriginPath = utils.String(originPath)
	}

	if probePath != "" {
		endpoint.EndpointProperties.ProbePath = utils.String(probePath)
	}

	origins := expandAzureRmCdnEndpointOrigins(d)
	if len(origins) > 0 {
		endpoint.EndpointProperties.Origins = &origins
	}

	profile, err := profilesClient.Get(ctx, id.ResourceGroup, id.ProfileName)
	if err != nil {
		return fmt.Errorf("retrieving parent CDN Profile for %s: %+v", id, err)
	}

	if profile.Sku != nil {
		globalDeliveryRulesRaw := d.Get("global_delivery_rule").([]interface{})
		deliveryRulesRaw := d.Get("delivery_rule").([]interface{})
		deliveryPolicy, err := expandArmCdnEndpointDeliveryPolicy(globalDeliveryRulesRaw, deliveryRulesRaw)
		if err != nil {
			return fmt.Errorf("expanding `global_delivery_rule` or `delivery_rule`: %s", err)
		}

		if profile.Sku.Name != cdn.SkuNameStandardMicrosoft && len(*deliveryPolicy.Rules) > 0 {
			return fmt.Errorf("`global_delivery_rule` and `delivery_rule` are only allowed when `Standard_Microsoft` sku is used. Profile sku:  %s", profile.Sku.Name)
		}

		if profile.Sku.Name == cdn.SkuNameStandardMicrosoft {
			endpoint.EndpointProperties.DeliveryPolicy = deliveryPolicy
		}
	}

	future, err := endpointsClient.Create(ctx, id.ResourceGroup, id.ProfileName, id.Name, endpoint)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, endpointsClient.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnEndpointRead(d, meta)
}

func resourceCdnEndpointUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	endpointsClient := meta.(*clients.Client).Cdn.EndpointsClient
	profilesClient := meta.(*clients.Client).Cdn.ProfilesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM CDN EndPoint update.")

	id, err := parse.EndpointID(d.Id())
	if err != nil {
		return err
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	httpAllowed := d.Get("is_http_allowed").(bool)
	httpsAllowed := d.Get("is_https_allowed").(bool)
	cachingBehaviour := d.Get("querystring_caching_behaviour").(string)
	originPath := d.Get("origin_path").(string)
	probePath := d.Get("probe_path").(string)
	optimizationType := d.Get("optimization_type").(string)
	t := d.Get("tags").(map[string]interface{})

	// NOTE: "Only tags can be updated after creating an endpoint." So only
	// call 'PATCH' if the only thing that has changed are the tags, else
	// call the 'PUT' instead. https://learn.microsoft.com/rest/api/cdn/endpoints/update?tabs=HTTP
	// see issue #22326 for more details.
	updateTypePATCH := true

	if d.HasChanges("is_http_allowed", "is_https_allowed", "querystring_caching_behaviour", "origin_path",
		"probe_path", "optimization_type", "origin_host_header", "content_types_to_compress", "geo_filter",
		"is_compression_enabled", "probe_path", "geo_filter", "optimization_type", "global_delivery_rule",
		"delivery_rule") {
		updateTypePATCH = false
	}

	if updateTypePATCH {
		log.Printf("[INFO] No changes detected using PATCH for Azure ARM CDN EndPoint update.")

		if !d.HasChange("tags") {
			log.Printf("[INFO] 'tags' did not change, skipping Azure ARM CDN EndPoint update.")
			return resourceCdnEndpointRead(d, meta)
		}

		endpoint := cdn.EndpointUpdateParameters{
			EndpointPropertiesUpdateParameters: &cdn.EndpointPropertiesUpdateParameters{},
			Tags:                               tags.Expand(t),
		}

		future, err := endpointsClient.Update(ctx, id.ResourceGroup, id.ProfileName, id.Name, endpoint)
		if err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}

		if err = future.WaitForCompletionRef(ctx, endpointsClient.Client); err != nil {
			return fmt.Errorf("waiting for update of %s: %+v", *id, err)
		}
	} else {
		log.Printf("[INFO] One or more fields have changed using PUT for Azure ARM CDN EndPoint update.")

		endpoint := cdn.Endpoint{
			Location: &location,
			EndpointProperties: &cdn.EndpointProperties{
				IsHTTPAllowed:              &httpAllowed,
				IsHTTPSAllowed:             &httpsAllowed,
				QueryStringCachingBehavior: cdn.QueryStringCachingBehavior(cachingBehaviour),
			},
			Tags: tags.Expand(t),
		}

		if v, ok := d.GetOk("origin_host_header"); ok {
			endpoint.EndpointProperties.OriginHostHeader = utils.String(v.(string))
		}

		if _, ok := d.GetOk("content_types_to_compress"); ok {
			contentTypes := expandArmCdnEndpointContentTypesToCompress(d)
			endpoint.EndpointProperties.ContentTypesToCompress = &contentTypes
		}

		if _, ok := d.GetOk("geo_filter"); ok {
			geoFilters := expandCdnEndpointGeoFilters(d)
			endpoint.EndpointProperties.GeoFilters = geoFilters
		}

		if v, ok := d.GetOk("is_compression_enabled"); ok {
			endpoint.EndpointProperties.IsCompressionEnabled = utils.Bool(v.(bool))
		}

		if optimizationType != "" {
			endpoint.EndpointProperties.OptimizationType = cdn.OptimizationType(optimizationType)
		}

		if originPath != "" {
			endpoint.EndpointProperties.OriginPath = utils.String(originPath)
		}

		if probePath != "" {
			endpoint.EndpointProperties.ProbePath = utils.String(probePath)
		}

		origins := expandAzureRmCdnEndpointOrigins(d)
		if len(origins) > 0 {
			endpoint.EndpointProperties.Origins = &origins
		}

		profile, err := profilesClient.Get(ctx, id.ResourceGroup, id.ProfileName)
		if err != nil {
			return fmt.Errorf("retrieving parent CDN Profile for %s: %+v", id, err)
		}

		if profile.Sku != nil {
			globalDeliveryRulesRaw := d.Get("global_delivery_rule").([]interface{})
			deliveryRulesRaw := d.Get("delivery_rule").([]interface{})
			deliveryPolicy, err := expandArmCdnEndpointDeliveryPolicy(globalDeliveryRulesRaw, deliveryRulesRaw)
			if err != nil {
				return fmt.Errorf("expanding `global_delivery_rule` or `delivery_rule`: %s", err)
			}

			if profile.Sku.Name != cdn.SkuNameStandardMicrosoft && len(*deliveryPolicy.Rules) > 0 {
				return fmt.Errorf("`global_delivery_rule` and `delivery_rule` are only allowed when `Standard_Microsoft` sku is used. Profile sku:  %s", profile.Sku.Name)
			}

			if profile.Sku.Name == cdn.SkuNameStandardMicrosoft {
				endpoint.EndpointProperties.DeliveryPolicy = deliveryPolicy
			}
		}

		future, err := endpointsClient.Create(ctx, id.ResourceGroup, id.ProfileName, id.Name, endpoint)
		if err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, endpointsClient.Client); err != nil {
			return fmt.Errorf("waiting for update of %s: %+v", id, err)
		}
	}

	return resourceCdnEndpointRead(d, meta)
}

func resourceCdnEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("profile_name", id.ProfileName)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.EndpointProperties; props != nil {
		d.Set("fqdn", props.HostName)
		d.Set("is_http_allowed", props.IsHTTPAllowed)
		d.Set("is_https_allowed", props.IsHTTPSAllowed)
		d.Set("querystring_caching_behaviour", props.QueryStringCachingBehavior)
		d.Set("origin_host_header", props.OriginHostHeader)
		d.Set("origin_path", props.OriginPath)
		d.Set("probe_path", props.ProbePath)
		d.Set("optimization_type", string(props.OptimizationType))

		compressionEnabled := false
		if v := props.IsCompressionEnabled; v != nil {
			compressionEnabled = *v
		}
		d.Set("is_compression_enabled", compressionEnabled)

		contentTypes := flattenAzureRMCdnEndpointContentTypes(props.ContentTypesToCompress)
		if err := d.Set("content_types_to_compress", contentTypes); err != nil {
			return fmt.Errorf("setting `content_types_to_compress`: %+v", err)
		}

		geoFilters := flattenCdnEndpointGeoFilters(props.GeoFilters)
		if err := d.Set("geo_filter", geoFilters); err != nil {
			return fmt.Errorf("setting `geo_filter`: %+v", err)
		}

		origins := flattenAzureRMCdnEndpointOrigin(props.Origins)
		if err := d.Set("origin", origins); err != nil {
			return fmt.Errorf("setting `origin`: %+v", err)
		}

		flattenedDeliveryPolicies, err := flattenEndpointDeliveryPolicy(props.DeliveryPolicy)
		if err != nil {
			return err
		}
		if err := d.Set("global_delivery_rule", flattenedDeliveryPolicies.globalDeliveryRules); err != nil {
			return fmt.Errorf("setting `global_delivery_rule`: %+v", err)
		}
		if err := d.Set("delivery_rule", flattenedDeliveryPolicies.deliveryRules); err != nil {
			return fmt.Errorf("setting `delivery_rule`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceCdnEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EndpointID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandCdnEndpointGeoFilters(d *pluginsdk.ResourceData) *[]cdn.GeoFilter {
	filters := make([]cdn.GeoFilter, 0)

	inputFilters := d.Get("geo_filter").([]interface{})
	for _, v := range inputFilters {
		input := v.(map[string]interface{})
		action := input["action"].(string)
		relativePath := input["relative_path"].(string)

		inputCountryCodes := input["country_codes"].([]interface{})
		countryCodes := make([]string, 0)

		for _, v := range inputCountryCodes {
			if v != nil {
				countryCode := v.(string)
				countryCodes = append(countryCodes, countryCode)
			}
		}

		filter := cdn.GeoFilter{
			Action:       cdn.GeoFilterActions(action),
			RelativePath: utils.String(relativePath),
			CountryCodes: &countryCodes,
		}
		filters = append(filters, filter)
	}

	return &filters
}

func flattenCdnEndpointGeoFilters(input *[]cdn.GeoFilter) []interface{} {
	results := make([]interface{}, 0)

	if filters := input; filters != nil {
		for _, filter := range *filters {
			relativePath := ""
			if filter.RelativePath != nil {
				relativePath = *filter.RelativePath
			}

			outputCodes := make([]interface{}, 0)
			if codes := filter.CountryCodes; codes != nil {
				for _, code := range *codes {
					outputCodes = append(outputCodes, code)
				}
			}

			results = append(results, map[string]interface{}{
				"action":        string(filter.Action),
				"country_codes": outputCodes,
				"relative_path": relativePath,
			})
		}
	}

	return results
}

func expandArmCdnEndpointContentTypesToCompress(d *pluginsdk.ResourceData) []string {
	results := make([]string, 0)
	input := d.Get("content_types_to_compress").(*pluginsdk.Set).List()

	for _, v := range input {
		contentType := v.(string)
		results = append(results, contentType)
	}

	return results
}

func flattenAzureRMCdnEndpointContentTypes(input *[]string) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		for _, v := range *input {
			output = append(output, v)
		}
	}

	return output
}

func expandAzureRmCdnEndpointOrigins(d *pluginsdk.ResourceData) []cdn.DeepCreatedOrigin {
	configs := d.Get("origin").(*pluginsdk.Set).List()
	origins := make([]cdn.DeepCreatedOrigin, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		name := data["name"].(string)
		hostName := data["host_name"].(string)

		origin := cdn.DeepCreatedOrigin{
			Name: utils.String(name),
			DeepCreatedOriginProperties: &cdn.DeepCreatedOriginProperties{
				HostName: utils.String(hostName),
			},
		}

		if v, ok := data["https_port"]; ok {
			port := v.(int)
			origin.DeepCreatedOriginProperties.HTTPSPort = utils.Int32(int32(port))
		}

		if v, ok := data["http_port"]; ok {
			port := v.(int)
			origin.DeepCreatedOriginProperties.HTTPPort = utils.Int32(int32(port))
		}

		origins = append(origins, origin)
	}

	return origins
}

func flattenAzureRMCdnEndpointOrigin(input *[]cdn.DeepCreatedOrigin) []interface{} {
	results := make([]interface{}, 0)

	if list := input; list != nil {
		for _, i := range *list {
			name := ""
			if i.Name != nil {
				name = *i.Name
			}

			hostName := ""
			httpPort := 80
			httpsPort := 443
			if props := i.DeepCreatedOriginProperties; props != nil {
				if props.HostName != nil {
					hostName = *props.HostName
				}
				if port := props.HTTPPort; port != nil {
					httpPort = int(*port)
				}
				if port := props.HTTPSPort; port != nil {
					httpsPort = int(*port)
				}
			}

			results = append(results, map[string]interface{}{
				"name":       name,
				"host_name":  hostName,
				"http_port":  httpPort,
				"https_port": httpsPort,
			})
		}
	}

	return results
}

func expandArmCdnEndpointDeliveryPolicy(globalRulesRaw []interface{}, deliveryRulesRaw []interface{}) (*cdn.EndpointPropertiesUpdateParametersDeliveryPolicy, error) {
	deliveryRules := make([]cdn.DeliveryRule, 0)
	deliveryPolicy := cdn.EndpointPropertiesUpdateParametersDeliveryPolicy{
		Description: utils.String(""),
		Rules:       &deliveryRules,
	}

	if len(globalRulesRaw) > 0 && globalRulesRaw[0] != nil {
		ruleRaw := globalRulesRaw[0].(map[string]interface{})
		rule, err := expandArmCdnEndpointGlobalDeliveryRule(ruleRaw)
		if err != nil {
			return nil, err
		}
		deliveryRules = append(deliveryRules, *rule)
	}

	for _, ruleV := range deliveryRulesRaw {
		ruleRaw := ruleV.(map[string]interface{})
		rule, err := expandArmCdnEndpointDeliveryRule(ruleRaw)
		if err != nil {
			return nil, err
		}
		deliveryRules = append(deliveryRules, *rule)
	}

	return &deliveryPolicy, nil
}

type flattenedEndpointDeliveryPolicies struct {
	globalDeliveryRules []interface{}
	deliveryRules       []interface{}
}

func flattenEndpointDeliveryPolicy(input *cdn.EndpointPropertiesUpdateParametersDeliveryPolicy) (*flattenedEndpointDeliveryPolicies, error) {
	output := flattenedEndpointDeliveryPolicies{
		globalDeliveryRules: make([]interface{}, 0),
		deliveryRules:       make([]interface{}, 0),
	}
	if input == nil || input.Rules == nil {
		return &output, nil
	}

	for _, rule := range *input.Rules {
		if rule.Order == nil {
			continue
		}

		if int(*rule.Order) == 0 {
			flattenedRule, err := flattenArmCdnEndpointGlobalDeliveryRule(rule)
			if err != nil {
				return nil, err
			}

			output.globalDeliveryRules = append(output.globalDeliveryRules, flattenedRule)
			continue
		}

		flattenedRule, err := flattenArmCdnEndpointDeliveryRule(rule)
		if err != nil {
			return nil, err
		}

		output.deliveryRules = append(output.deliveryRules, flattenedRule)
	}

	return &output, nil
}
