package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2017-10-12/cdn"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCdnEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCdnEndpointCreate,
		Read:   resourceArmCdnEndpointRead,
		Update: resourceArmCdnEndpointUpdate,
		Delete: resourceArmCdnEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
				Computed: true,
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
				Default:  false,
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

			"tags": tags.Schema(),
		},
	}
}

func resourceArmCdnEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cdn.EndpointsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM CDN EndPoint creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	profileName := d.Get("profile_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, profileName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing CDN Endpoint %q (Profile %q / Resource Group %q): %s", name, profileName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_cdn_endpoint", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	httpAllowed := d.Get("is_http_allowed").(bool)
	httpsAllowed := d.Get("is_https_allowed").(bool)
	compressionEnabled := d.Get("is_compression_enabled").(bool)
	cachingBehaviour := d.Get("querystring_caching_behaviour").(string)
	originHostHeader := d.Get("origin_host_header").(string)
	originPath := d.Get("origin_path").(string)
	probePath := d.Get("probe_path").(string)
	optimizationType := d.Get("optimization_type").(string)
	contentTypes := expandArmCdnEndpointContentTypesToCompress(d)
	t := d.Get("tags").(map[string]interface{})

	geoFilters, err := expandArmCdnEndpointGeoFilters(d)
	if err != nil {
		return fmt.Errorf("Error expanding `geo_filter`: %s", err)
	}

	endpoint := cdn.Endpoint{
		Location: &location,
		EndpointProperties: &cdn.EndpointProperties{
			ContentTypesToCompress:     &contentTypes,
			GeoFilters:                 geoFilters,
			IsHTTPAllowed:              &httpAllowed,
			IsHTTPSAllowed:             &httpsAllowed,
			IsCompressionEnabled:       &compressionEnabled,
			QueryStringCachingBehavior: cdn.QueryStringCachingBehavior(cachingBehaviour),
			OriginHostHeader:           utils.String(originHostHeader),
		},
		Tags: tags.Expand(t),
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

	origins, err := expandAzureRmCdnEndpointOrigins(d)
	if err != nil {
		return fmt.Errorf("Error Building list of CDN Endpoint Origins: %s", err)
	}
	if len(origins) > 0 {
		endpoint.EndpointProperties.Origins = &origins
	}

	future, err := client.Create(ctx, resourceGroup, profileName, name, endpoint)
	if err != nil {
		return fmt.Errorf("Error creating CDN Endpoint %q (Profile %q / Resource Group %q): %+v", name, profileName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for CDN Endpoint %q (Profile %q / Resource Group %q) to finish creating: %+v", name, profileName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, profileName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving CDN Endpoint %q (Profile %q / Resource Group %q): %+v", name, profileName, resourceGroup, err)
	}

	d.SetId(*read.ID)

	return resourceArmCdnEndpointRead(d, meta)
}

func resourceArmCdnEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	endpointsClient := meta.(*ArmClient).cdn.EndpointsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	profileName := d.Get("profile_name").(string)
	httpAllowed := d.Get("is_http_allowed").(bool)
	httpsAllowed := d.Get("is_https_allowed").(bool)
	compressionEnabled := d.Get("is_compression_enabled").(bool)
	cachingBehaviour := d.Get("querystring_caching_behaviour").(string)
	hostHeader := d.Get("origin_host_header").(string)
	originPath := d.Get("origin_path").(string)
	probePath := d.Get("probe_path").(string)
	optimizationType := d.Get("optimization_type").(string)
	contentTypes := expandArmCdnEndpointContentTypesToCompress(d)
	t := d.Get("tags").(map[string]interface{})

	geoFilters, err := expandArmCdnEndpointGeoFilters(d)
	if err != nil {
		return fmt.Errorf("Error expanding `geo_filter`: %s", err)
	}

	endpoint := cdn.EndpointUpdateParameters{
		EndpointPropertiesUpdateParameters: &cdn.EndpointPropertiesUpdateParameters{
			ContentTypesToCompress:     &contentTypes,
			GeoFilters:                 geoFilters,
			IsHTTPAllowed:              utils.Bool(httpAllowed),
			IsHTTPSAllowed:             utils.Bool(httpsAllowed),
			IsCompressionEnabled:       utils.Bool(compressionEnabled),
			QueryStringCachingBehavior: cdn.QueryStringCachingBehavior(cachingBehaviour),
			OriginHostHeader:           utils.String(hostHeader),
		},
		Tags: tags.Expand(t),
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

	future, err := endpointsClient.Update(ctx, resourceGroup, profileName, name, endpoint)
	if err != nil {
		return fmt.Errorf("Error updating CDN Endpoint %q (Profile %q / Resource Group %q): %s", name, profileName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, endpointsClient.Client); err != nil {
		return fmt.Errorf("Error waiting for the CDN Endpoint %q (Profile %q / Resource Group %q) to finish updating: %+v", name, profileName, resourceGroup, err)
	}

	return resourceArmCdnEndpointRead(d, meta)
}

func resourceArmCdnEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cdn.EndpointsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["endpoints"]
	profileName := id.Path["profiles"]
	if profileName == "" {
		profileName = id.Path["Profiles"]
	}
	log.Printf("[INFO] Retrieving CDN Endpoint %q (Profile %q / Resource Group %q)", name, profileName, resourceGroup)
	resp, err := client.Get(ctx, resourceGroup, profileName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure CDN Endpoint %q (Profile %q / Resource Group %q): %+v", name, profileName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("profile_name", profileName)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.EndpointProperties; props != nil {
		d.Set("host_name", props.HostName)
		d.Set("is_compression_enabled", props.IsCompressionEnabled)
		d.Set("is_http_allowed", props.IsHTTPAllowed)
		d.Set("is_https_allowed", props.IsHTTPSAllowed)
		d.Set("querystring_caching_behaviour", props.QueryStringCachingBehavior)
		d.Set("origin_host_header", props.OriginHostHeader)
		d.Set("origin_path", props.OriginPath)
		d.Set("probe_path", props.ProbePath)
		d.Set("optimization_type", string(props.OptimizationType))

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
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmCdnEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cdn.EndpointsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	profileName := id.Path["profiles"]
	if profileName == "" {
		profileName = id.Path["Profiles"]
	}
	name := id.Path["endpoints"]

	future, err := client.Delete(ctx, resourceGroup, profileName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting CDN Endpoint %q (Profile %q / Resource Group %q): %+v", name, profileName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error waiting for CDN Endpoint %q (Profile %q / Resource Group %q) to be deleted: %+v", name, profileName, resourceGroup, err)
	}

	return nil
}

func expandArmCdnEndpointGeoFilters(d *schema.ResourceData) (*[]cdn.GeoFilter, error) {
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

	return &filters, nil
}

func flattenCdnEndpointGeoFilters(input *[]cdn.GeoFilter) []interface{} {
	results := make([]interface{}, 0)

	if filters := input; filters != nil {
		for _, filter := range *filters {
			output := make(map[string]interface{})

			output["action"] = string(filter.Action)
			if path := filter.RelativePath; path != nil {
				output["relative_path"] = *path
			}

			outputCodes := make([]interface{}, 0)
			if codes := filter.CountryCodes; codes != nil {
				for _, code := range *codes {
					outputCodes = append(outputCodes, code)
				}
			}
			output["country_codes"] = outputCodes

			results = append(results, output)
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

func expandAzureRmCdnEndpointOrigins(d *schema.ResourceData) ([]cdn.DeepCreatedOrigin, error) {
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

	return origins, nil
}

func flattenAzureRMCdnEndpointOrigin(input *[]cdn.DeepCreatedOrigin) []interface{} {
	results := make([]interface{}, 0)

	if list := input; list != nil {
		for _, i := range *list {
			output := map[string]interface{}{}

			if name := i.Name; name != nil {
				output["name"] = *name
			}

			if props := i.DeepCreatedOriginProperties; props != nil {
				if hostName := props.HostName; hostName != nil {
					output["host_name"] = *hostName
				}
				if port := props.HTTPPort; port != nil {
					output["http_port"] = int(*port)
				}
				if port := props.HTTPSPort; port != nil {
					output["https_port"] = int(*port)
				}
			}

			results = append(results, output)
		}
	}

	return results
}
