package cdn

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/migration"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCdnEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCdnEndpointCreate,
		Read:   resourceArmCdnEndpointRead,
		Update: resourceArmCdnEndpointUpdate,
		Delete: resourceArmCdnEndpointDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.CdnEndpointID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"profile_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"origin_host_header": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"is_http_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"is_https_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"origin": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"host_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"http_port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  80,
						},

						"https_port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  443,
						},
					},
				},
			},

			"origin_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"querystring_caching_behaviour": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(cdn.IgnoreQueryString),
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.BypassCaching),
					string(cdn.IgnoreQueryString),
					string(cdn.NotSet),
					string(cdn.UseQueryString),
				}, false),
			},

			"content_types_to_compress": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"is_compression_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"probe_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"geo_filter": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"relative_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"action": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.Allow),
								string(cdn.Block),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"country_codes": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"optimization_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.DynamicSiteAcceleration),
					string(cdn.GeneralMediaStreaming),
					string(cdn.GeneralWebDelivery),
					string(cdn.LargeFileDownload),
					string(cdn.VideoOnDemandMediaStreaming),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"global_delivery_rule": endpointGlobalDeliveryRule(),

			"delivery_rule": endpointDeliveryRule(),

			"tags": tags.Schema(),
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    migration.CdnEndpointV0Schema().CoreConfigSchema().ImpliedType(),
				Upgrade: migration.CdnEndpointV0ToV1,
				Version: 0,
			},
		},
	}
}

func resourceArmCdnEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	endpointsClient := meta.(*clients.Client).Cdn.EndpointsClient
	profilesClient := meta.(*clients.Client).Cdn.ProfilesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM CDN EndPoint creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	profileName := d.Get("profile_name").(string)

	existing, err := endpointsClient.Get(ctx, resourceGroup, profileName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing CDN Endpoint %q (Profile %q / Resource Group %q): %s", name, profileName, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_cdn_endpoint", *existing.ID)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	httpAllowed := d.Get("is_http_allowed").(bool)
	httpsAllowed := d.Get("is_https_allowed").(bool)
	cachingBehaviour := d.Get("querystring_caching_behaviour").(string)
	originHostHeader := d.Get("origin_host_header").(string)
	originPath := d.Get("origin_path").(string)
	probePath := d.Get("probe_path").(string)
	optimizationType := d.Get("optimization_type").(string)
	contentTypes := expandArmCdnEndpointContentTypesToCompress(d)
	geoFilters := expandCdnEndpointGeoFilters(d)
	t := d.Get("tags").(map[string]interface{})

	endpoint := cdn.Endpoint{
		Location: &location,
		EndpointProperties: &cdn.EndpointProperties{
			ContentTypesToCompress:     &contentTypes,
			GeoFilters:                 geoFilters,
			IsHTTPAllowed:              &httpAllowed,
			IsHTTPSAllowed:             &httpsAllowed,
			QueryStringCachingBehavior: cdn.QueryStringCachingBehavior(cachingBehaviour),
			OriginHostHeader:           utils.String(originHostHeader),
		},
		Tags: tags.Expand(t),
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

	profile, err := profilesClient.Get(ctx, resourceGroup, profileName)
	if err != nil {
		return fmt.Errorf("Error creating CDN Endpoint %q while getting CDN Profile (Profile %q / Resource Group %q): %+v", name, profileName, resourceGroup, err)
	}

	if profile.Sku != nil {
		globalDeliveryRulesRaw := d.Get("global_delivery_rule").([]interface{})
		deliveryRulesRaw := d.Get("delivery_rule").([]interface{})
		deliveryPolicy, err := expandArmCdnEndpointDeliveryPolicy(globalDeliveryRulesRaw, deliveryRulesRaw)
		if err != nil {
			return fmt.Errorf("Error expanding `global_delivery_rule` or `delivery_rule`: %s", err)
		}

		if profile.Sku.Name != cdn.StandardMicrosoft && len(*deliveryPolicy.Rules) > 0 {
			return fmt.Errorf("`global_delivery_policy` and `delivery_rule` are only allowed when `Standard_Microsoft` sku is used. Profile sku:  %s", profile.Sku.Name)
		}

		endpoint.EndpointProperties.DeliveryPolicy = deliveryPolicy
	}

	future, err := endpointsClient.Create(ctx, resourceGroup, profileName, name, endpoint)
	if err != nil {
		return fmt.Errorf("Error creating CDN Endpoint %q (Profile %q / Resource Group %q): %+v", name, profileName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, endpointsClient.Client); err != nil {
		return fmt.Errorf("Error waiting for CDN Endpoint %q (Profile %q / Resource Group %q) to finish creating: %+v", name, profileName, resourceGroup, err)
	}

	read, err := endpointsClient.Get(ctx, resourceGroup, profileName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving CDN Endpoint %q (Profile %q / Resource Group %q): %+v", name, profileName, resourceGroup, err)
	}

	id, err := parse.CdnEndpointID(*read.ID)
	if err != nil {
		return err
	}

	d.SetId(id.ID(subscriptionId))

	return resourceArmCdnEndpointRead(d, meta)
}

func resourceArmCdnEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	endpointsClient := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CdnEndpointID(d.Id())
	if err != nil {
		return err
	}

	httpAllowed := d.Get("is_http_allowed").(bool)
	httpsAllowed := d.Get("is_https_allowed").(bool)
	cachingBehaviour := d.Get("querystring_caching_behaviour").(string)
	hostHeader := d.Get("origin_host_header").(string)
	originPath := d.Get("origin_path").(string)
	probePath := d.Get("probe_path").(string)
	optimizationType := d.Get("optimization_type").(string)
	contentTypes := expandArmCdnEndpointContentTypesToCompress(d)
	geoFilters := expandCdnEndpointGeoFilters(d)
	t := d.Get("tags").(map[string]interface{})

	endpoint := cdn.EndpointUpdateParameters{
		EndpointPropertiesUpdateParameters: &cdn.EndpointPropertiesUpdateParameters{
			ContentTypesToCompress:     &contentTypes,
			GeoFilters:                 geoFilters,
			IsHTTPAllowed:              utils.Bool(httpAllowed),
			IsHTTPSAllowed:             utils.Bool(httpsAllowed),
			QueryStringCachingBehavior: cdn.QueryStringCachingBehavior(cachingBehaviour),
			OriginHostHeader:           utils.String(hostHeader),
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("is_compression_enabled"); ok {
		endpoint.EndpointPropertiesUpdateParameters.IsCompressionEnabled = utils.Bool(v.(bool))
	}
	if optimizationType != "" {
		endpoint.EndpointPropertiesUpdateParameters.OptimizationType = cdn.OptimizationType(optimizationType)
	}
	if originPath != "" {
		endpoint.EndpointPropertiesUpdateParameters.OriginPath = utils.String(originPath)
	}
	if probePath != "" {
		endpoint.EndpointPropertiesUpdateParameters.ProbePath = utils.String(probePath)
	}

	profilesClient := meta.(*clients.Client).Cdn.ProfilesClient
	profileGetCtx, profileGetCancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer profileGetCancel()

	profile, err := profilesClient.Get(profileGetCtx, id.ResourceGroup, id.ProfileName)
	if err != nil {
		return fmt.Errorf("Error creating CDN Endpoint %q while getting CDN Profile (Profile %q / Resource Group %q): %+v", id.Name, id.ProfileName, id.ResourceGroup, err)
	}

	if profile.Sku != nil {
		globalDeliveryRulesRaw := d.Get("global_delivery_rule").([]interface{})
		deliveryRulesRaw := d.Get("delivery_rule").([]interface{})
		deliveryPolicy, err := expandArmCdnEndpointDeliveryPolicy(globalDeliveryRulesRaw, deliveryRulesRaw)
		if err != nil {
			return fmt.Errorf("Error expanding `global_delivery_rule` or `delivery_rule`: %s", err)
		}

		if profile.Sku.Name != cdn.StandardMicrosoft && len(*deliveryPolicy.Rules) > 0 {
			return fmt.Errorf("`global_delivery_policy` and `delivery_rule` are only allowed when `Standard_Microsoft` sku is used. Profile sku:  %s", profile.Sku.Name)
		}

		endpoint.EndpointPropertiesUpdateParameters.DeliveryPolicy = deliveryPolicy
	}

	future, err := endpointsClient.Update(ctx, id.ResourceGroup, id.ProfileName, id.Name, endpoint)
	if err != nil {
		return fmt.Errorf("Error updating CDN Endpoint %q (Profile %q / Resource Group %q): %s", id.Name, id.ProfileName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, endpointsClient.Client); err != nil {
		return fmt.Errorf("Error waiting for the CDN Endpoint %q (Profile %q / Resource Group %q) to finish updating: %+v", id.Name, id.ProfileName, id.ResourceGroup, err)
	}

	return resourceArmCdnEndpointRead(d, meta)
}

func resourceArmCdnEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CdnEndpointID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Retrieving CDN Endpoint %q (Profile %q / Resource Group %q)", id.Name, id.ProfileName, id.ResourceGroup)

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure CDN Endpoint %q (Profile %q / Resource Group %q): %+v", id.Name, id.ProfileName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("profile_name", id.ProfileName)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.EndpointProperties; props != nil {
		d.Set("host_name", props.HostName)
		d.Set("is_http_allowed", props.IsHTTPAllowed)
		d.Set("is_https_allowed", props.IsHTTPSAllowed)
		d.Set("querystring_caching_behaviour", props.QueryStringCachingBehavior)
		d.Set("origin_host_header", props.OriginHostHeader)
		d.Set("origin_path", props.OriginPath)
		d.Set("probe_path", props.ProbePath)
		d.Set("optimization_type", string(props.OptimizationType))

		if is_compression_enabled := props.IsCompressionEnabled; is_compression_enabled != nil {
			d.Set("is_compression_enabled", *is_compression_enabled)
		}

		contentTypes := flattenAzureRMCdnEndpointContentTypes(props.ContentTypesToCompress)
		if err := d.Set("content_types_to_compress", contentTypes); err != nil {
			return fmt.Errorf("Error setting `content_types_to_compress`: %+v", err)
		}

		geoFilters := flattenCdnEndpointGeoFilters(props.GeoFilters)
		if err := d.Set("geo_filter", geoFilters); err != nil {
			return fmt.Errorf("Error setting `geo_filter`: %+v", err)
		}

		origins := flattenAzureRMCdnEndpointOrigin(props.Origins)
		if err := d.Set("origin", origins); err != nil {
			return fmt.Errorf("Error setting `origin`: %+v", err)
		}

		flattenedDeliveryPolicies, err := flattenArmCdnEndpointDeliveryPolicy(props.DeliveryPolicy)
		if err != nil {
			return err
		}
		if err := d.Set("global_delivery_rule", flattenedDeliveryPolicies.globalDeliveryRules); err != nil {
			return fmt.Errorf("Error setting `global_delivery_rule`: %+v", err)
		}
		if err := d.Set("delivery_rule", flattenedDeliveryPolicies.deliveryRules); err != nil {
			return fmt.Errorf("Error setting `delivery_rule`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmCdnEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CdnEndpointID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting CDN Endpoint %q (Profile %q / Resource Group %q): %+v", id.Name, id.ProfileName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error waiting for CDN Endpoint %q (Profile %q / Resource Group %q) to be deleted: %+v", id.Name, id.ProfileName, id.ResourceGroup, err)
	}

	return nil
}

func expandCdnEndpointGeoFilters(d *schema.ResourceData) *[]cdn.GeoFilter {
	filters := make([]cdn.GeoFilter, 0)

	inputFilters := d.Get("geo_filter").([]interface{})
	for _, v := range inputFilters {
		input := v.(map[string]interface{})
		action := input["action"].(string)
		relativePath := input["relative_path"].(string)

		inputCountryCodes := input["country_codes"].([]interface{})
		countryCodes := make([]string, 0)

		for _, v := range inputCountryCodes {
			countryCode := v.(string)
			countryCodes = append(countryCodes, countryCode)
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

func expandArmCdnEndpointContentTypesToCompress(d *schema.ResourceData) []string {
	results := make([]string, 0)
	input := d.Get("content_types_to_compress").(*schema.Set).List()

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

func expandAzureRmCdnEndpointOrigins(d *schema.ResourceData) []cdn.DeepCreatedOrigin {
	configs := d.Get("origin").(*schema.Set).List()
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
			httpPort := 0
			httpsPort := 0
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

	if len(globalRulesRaw) > 0 {
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

func flattenArmCdnEndpointDeliveryPolicy(input *cdn.EndpointPropertiesUpdateParametersDeliveryPolicy) (*flattenedEndpointDeliveryPolicies, error) {
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
