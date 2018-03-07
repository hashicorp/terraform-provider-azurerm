package azurerm

import (
	"bytes"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2017-04-02/cdn"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
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

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

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
						},

						"host_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"http_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  80,
						},

						"https_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  443,
						},
					},
				},
				Set: resourceArmCdnEndpointOriginHash,
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

			"host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmCdnEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cdnEndpointsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM CDN EndPoint creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	profileName := d.Get("profile_name").(string)
	httpAllowed := d.Get("is_http_allowed").(bool)
	httpsAllowed := d.Get("is_https_allowed").(bool)
	compressionEnabled := d.Get("is_compression_enabled").(bool)
	cachingBehaviour := d.Get("querystring_caching_behaviour").(string)
	originHostHeader := d.Get("origin_host_header").(string)
	originPath := d.Get("origin_path").(string)
	tags := d.Get("tags").(map[string]interface{})

	endpoint := cdn.Endpoint{
		Location: &location,
		EndpointProperties: &cdn.EndpointProperties{
			IsHTTPAllowed:              &httpAllowed,
			IsHTTPSAllowed:             &httpsAllowed,
			IsCompressionEnabled:       &compressionEnabled,
			QueryStringCachingBehavior: cdn.QueryStringCachingBehavior(cachingBehaviour),
			OriginHostHeader:           utils.String(originHostHeader),
			OriginPath:                 utils.String(originPath),
			//GeoFilters: []cdn.GeoFilter{{
			//  Action: cdn.Allow || cdn.Block
			//  CountryCodes: []string{}
			//  RelativePath: ""
			//}}
			//ProbePath: ""
		},
		Tags: expandTags(tags),
	}

	origins, err := expandAzureRmCdnEndpointOrigins(d)
	if err != nil {
		return fmt.Errorf("Error Building list of CDN Endpoint Origins: %s", err)
	}
	if len(origins) > 0 {
		endpoint.EndpointProperties.Origins = &origins
	}

	if v, ok := d.GetOk("content_types_to_compress"); ok {
		contentTypes := expandArmCdnEndpointContentTypesToCompress(v)
		endpoint.EndpointProperties.ContentTypesToCompress = &contentTypes
	}

	future, err := client.Create(ctx, resourceGroup, profileName, name, endpoint)
	if err != nil {
		return fmt.Errorf("Error creating CDN Endpoint %q (Profile %q / Resource Group %q): %+v", name, profileName, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
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
	client := meta.(*ArmClient).cdnEndpointsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	profileName := d.Get("profile_name").(string)
	httpAllowed := d.Get("is_http_allowed").(bool)
	httpsAllowed := d.Get("is_https_allowed").(bool)
	compressionEnabled := d.Get("is_compression_enabled").(bool)
	cachingBehaviour := d.Get("querystring_caching_behaviour").(string)
	hostHeader := d.Get("origin_host_header").(string)
	originPath := d.Get("origin_path").(string)
	tags := d.Get("tags").(map[string]interface{})

	endpoint := cdn.EndpointUpdateParameters{
		EndpointPropertiesUpdateParameters: &cdn.EndpointPropertiesUpdateParameters{
			IsHTTPAllowed:              &httpAllowed,
			IsHTTPSAllowed:             &httpsAllowed,
			IsCompressionEnabled:       &compressionEnabled,
			QueryStringCachingBehavior: cdn.QueryStringCachingBehavior(cachingBehaviour),
			OriginHostHeader:           &hostHeader,
			OriginPath:                 &originPath,
		},
		Tags: expandTags(tags),
	}

	if d.HasChange("content_types_to_compress") {
		v := d.Get("content_types_to_compress")
		contentTypes := expandArmCdnEndpointContentTypesToCompress(v)
		endpoint.EndpointPropertiesUpdateParameters.ContentTypesToCompress = &contentTypes
	}

	future, err := client.Update(ctx, resGroup, profileName, name, endpoint)
	if err != nil {
		return fmt.Errorf("Error updating CDN Endpoint %q (Profile %q / Resource Group %q): %s", name, profileName, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the CDN Endpoint %q (Profile %q / Resource Group %q) to finish updating: %+v", name, profileName, resGroup, err)
	}

	return resourceArmCdnEndpointRead(d, meta)
}

func resourceArmCdnEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cdnEndpointsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.EndpointProperties; props != nil {
		d.Set("host_name", props.HostName)
		d.Set("is_compression_enabled", props.IsCompressionEnabled)
		d.Set("is_http_allowed", props.IsHTTPAllowed)
		d.Set("is_https_allowed", props.IsHTTPSAllowed)
		d.Set("querystring_caching_behaviour", props.QueryStringCachingBehavior)
		d.Set("origin_host_header", props.OriginHostHeader)
		d.Set("origin_path", props.OriginPath)

		contentTypes := flattenAzureRMCdnEndpointContentTypes(props.ContentTypesToCompress)
		if err := d.Set("content_types_to_compress", contentTypes); err != nil {
			return fmt.Errorf("Error flattening `content_types_to_compress`: %+v", err)
		}
		origins := flattenAzureRMCdnEndpointOrigin(props.Origins)
		if err := d.Set("origin", origins); err != nil {
			return fmt.Errorf("Error flattening `origin`: %+v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmCdnEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).cdnEndpointsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error waiting for CDN Endpoint %q (Profile %q / Resource Group %q) to be deleted: %+v", name, profileName, resourceGroup, err)
	}

	return nil
}
func expandArmCdnEndpointContentTypesToCompress(v interface{}) []string {
	var contentTypes []string
	inputContentTypes := v.(*schema.Set).List()
	for _, ct := range inputContentTypes {
		str := ct.(string)
		contentTypes = append(contentTypes, str)
	}
}

func resourceArmCdnEndpointOriginHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["host_name"].(string)))

	// TODO: can we remove this; and do we need a state migration?

	return hashcode.String(buf.String())
}

func expandAzureRmCdnEndpointOrigins(d *schema.ResourceData) ([]cdn.DeepCreatedOrigin, error) {
	configs := d.Get("origin").(*schema.Set).List()
	origins := make([]cdn.DeepCreatedOrigin, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		hostName := data["host_name"].(string)
		properties := cdn.DeepCreatedOriginProperties{
			HostName: utils.String(hostName),
		}

		if v, ok := data["https_port"]; ok {
			port := v.(int)
			properties.HTTPSPort = utils.Int32(int32(port))
		}

		if v, ok := data["http_port"]; ok {
			port := v.(int)
			properties.HTTPPort = utils.Int32(int32(port))
		}

		name := data["name"].(string)

		origin := cdn.DeepCreatedOrigin{
			Name: utils.String(name),
			DeepCreatedOriginProperties: &properties,
		}

		origins = append(origins, origin)
	}

	return origins, nil
}

func flattenAzureRMCdnEndpointOrigin(input *[]cdn.DeepCreatedOrigin) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	if list := input; list != nil {
		for _, i := range *list {
			l := map[string]interface{}{}

			if name := i.Name; name != nil {
				l["name"] = *name
			}

			if props := i.DeepCreatedOriginProperties; props != nil {
				if hostName := props.HostName; hostName != nil {
					l["host_name"] = *hostName
				}
				if port := props.HTTPPort; port != nil {
					l["http_port"] = int(*port)
				}
				if port := props.HTTPSPort; port != nil {
					l["https_port"] = int(*port)
				}
			}

			result = append(result, l)
		}
	}

	return result
}

func flattenAzureRMCdnEndpointContentTypes(list *[]string) []interface{} {
	vs := make([]interface{}, 0, len(*list))
	for _, v := range *list {
		vs = append(vs, v)
	}
	return vs
}
