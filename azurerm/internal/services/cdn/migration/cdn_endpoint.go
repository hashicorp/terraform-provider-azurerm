package migration

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/deliveryruleactions"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/deliveryruleconditions"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/parse"
)

func CdnEndpointV0Schema() *schema.Resource {
	return &schema.Resource{
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

			// This is a point-in-time copy paste of the return value of `endpointGlobalDeliveryRule()` used in V1 schema
			"global_delivery_rule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cache_expiration_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     deliveryruleactions.CacheExpiration(),
						},

						"cache_key_query_string_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     deliveryruleactions.CacheKeyQueryString(),
						},

						"modify_request_header_action": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     deliveryruleactions.ModifyRequestHeader(),
						},

						"modify_response_header_action": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     deliveryruleactions.ModifyResponseHeader(),
						},

						"url_redirect_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     deliveryruleactions.URLRedirect(),
						},

						"url_rewrite_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     deliveryruleactions.URLRewrite(),
						},
					},
				},
			},

			// This is a point-in-time copy paste of the return value of `endpointDeliveryRule()` used in V1 schema
			"delivery_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.EndpointDeliveryRuleName(),
						},

						"order": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},

						"cookies_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     deliveryruleconditions.Cookies(),
						},

						"http_version_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     deliveryruleconditions.HTTPVersion(),
						},

						"device_condition": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     deliveryruleconditions.Device(),
						},

						"post_arg_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     deliveryruleconditions.PostArg(),
						},

						"query_string_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     deliveryruleconditions.QueryString(),
						},

						"remote_address_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     deliveryruleconditions.RemoteAddress(),
						},

						"request_body_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     deliveryruleconditions.RequestBody(),
						},

						"request_header_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     deliveryruleconditions.RequestHeader(),
						},

						"request_method_condition": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     deliveryruleconditions.RequestMethod(),
						},

						"request_scheme_condition": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     deliveryruleconditions.RequestScheme(),
						},

						"request_uri_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     deliveryruleconditions.RequestURI(),
						},

						"url_file_extension_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     deliveryruleconditions.URLFileExtension(),
						},

						"url_file_name_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     deliveryruleconditions.URLFileName(),
						},

						"url_path_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     deliveryruleconditions.URLPath(),
						},

						"cache_expiration_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     deliveryruleactions.CacheExpiration(),
						},

						"cache_key_query_string_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     deliveryruleactions.CacheKeyQueryString(),
						},

						"modify_request_header_action": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     deliveryruleactions.ModifyRequestHeader(),
						},

						"modify_response_header_action": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     deliveryruleactions.ModifyResponseHeader(),
						},

						"url_redirect_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     deliveryruleactions.URLRedirect(),
						},

						"url_rewrite_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     deliveryruleactions.URLRewrite(),
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func CdnEndpointV0ToV1(rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	// old
	// 	/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/endpoints/{endpointName}
	// new:
	// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/endpoints/{endpointName}
	// summary:
	// resourcegroups -> resourceGroups
	oldId := rawState["id"].(string)
	oldParsedId, err := azure.ParseAzureResourceID(oldId)
	if err != nil {
		return rawState, err
	}

	resourceGroup := oldParsedId.ResourceGroup
	profileName, err := oldParsedId.PopSegment("profiles")
	if err != nil {
		return rawState, err
	}
	name, err := oldParsedId.PopSegment("endpoints")
	if err != nil {
		return rawState, err
	}

	newId := parse.NewCdnEndpointID(parse.NewCdnProfileID(resourceGroup, profileName), name)
	newIdStr := newId.ID(oldParsedId.SubscriptionID)

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newIdStr)

	rawState["id"] = newIdStr

	return rawState, nil
}
