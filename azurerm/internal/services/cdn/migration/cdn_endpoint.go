package migration

import (
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
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

			"location": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				StateFunc:        location.StateFunc,
				DiffSuppressFunc: location.DiffSuppressFunc,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

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
							Type:             schema.TypeString,
							Required:         true,
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
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

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
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"behavior": {
										Type:     schema.TypeString,
										Required: true,
									},

									"duration": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"cache_key_query_string_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"behavior": {
										Type:     schema.TypeString,
										Required: true,
									},

									"parameters": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"modify_request_header_action": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:     schema.TypeString,
										Required: true,
									},

									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"modify_response_header_action": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:     schema.TypeString,
										Required: true,
									},

									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"url_redirect_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"redirect_type": {
										Type:     schema.TypeString,
										Required: true,
									},

									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  string(cdn.MatchRequest),
									},

									"hostname": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"path": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"query_string": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"fragment": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"url_rewrite_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_pattern": {
										Type:     schema.TypeString,
										Required: true,
									},

									"destination": {
										Type:     schema.TypeString,
										Required: true,
									},

									"preserve_unmatched_path": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
								},
							},
						},
					},
				},
			},

			"delivery_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"order": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"cookies_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"selector": {
										Type:     schema.TypeString,
										Required: true,
									},

									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},

									"negate_condition": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"match_values": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},

									"transforms": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"http_version_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "Equal",
									},

									"negate_condition": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"match_values": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"device_condition": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "Equal",
									},

									"negate_condition": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"match_values": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"post_arg_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"selector": {
										Type:     schema.TypeString,
										Required: true,
									},

									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},

									"negate_condition": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"match_values": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},

									"transforms": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"query_string_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},

									"negate_condition": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"match_values": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},

									"transforms": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"remote_address_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},

									"negate_condition": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"match_values": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"request_body_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},

									"negate_condition": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"match_values": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},

									"transforms": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"request_header_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"selector": {
										Type:     schema.TypeString,
										Required: true,
									},

									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},

									"negate_condition": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"match_values": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},

									"transforms": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"request_method_condition": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "Equal",
									},

									"negate_condition": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"match_values": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"request_scheme_condition": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "Equal",
									},

									"negate_condition": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"match_values": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"request_uri_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},

									"negate_condition": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"match_values": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},

									"transforms": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"url_file_extension_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},

									"negate_condition": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"match_values": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},

									"transforms": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"url_file_name_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},

									"negate_condition": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"match_values": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},

									"transforms": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"url_path_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},

									"negate_condition": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},

									"match_values": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},

									"transforms": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"cache_expiration_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"behavior": {
										Type:     schema.TypeString,
										Required: true,
									},

									"duration": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"cache_key_query_string_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"behavior": {
										Type:     schema.TypeString,
										Required: true,
									},

									"parameters": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"modify_request_header_action": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:     schema.TypeString,
										Required: true,
									},

									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"modify_response_header_action": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:     schema.TypeString,
										Required: true,
									},

									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"url_redirect_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"redirect_type": {
										Type:     schema.TypeString,
										Required: true,
									},

									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  string(cdn.MatchRequest),
									},

									"hostname": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"path": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"query_string": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"fragment": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"url_rewrite_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_pattern": {
										Type:     schema.TypeString,
										Required: true,
									},

									"destination": {
										Type:     schema.TypeString,
										Required: true,
									},

									"preserve_unmatched_path": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
								},
							},
						},
					},
				},
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
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

	newId := parse.NewCdnEndpointID(parse.NewCdnProfileID(oldParsedId.SubscriptionID, resourceGroup, profileName), name)
	newIdStr := newId.ID("")

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newIdStr)

	rawState["id"] = newIdStr

	return rawState, nil
}
