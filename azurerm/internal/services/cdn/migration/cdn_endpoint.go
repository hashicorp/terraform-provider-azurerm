package migration

import (
	"context"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = CdnEndpointV0ToV1{}

type CdnEndpointV0ToV1 struct{}

func (CdnEndpointV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

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
			Computed: true,
		},

		"querystring_caching_behaviour": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "IgnoreQueryString",
		},

		"content_types_to_compress": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true,
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
			Computed: true,
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
		},

		"host_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		//lintignore:XS003
		"global_delivery_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cache_expiration_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"behavior": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"duration": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},

					"cache_key_query_string_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"behavior": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"parameters": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},

					"modify_request_header_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"action": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"value": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},

					"modify_response_header_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"action": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"value": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},

					"url_redirect_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"redirect_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"protocol": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "MatchRequest",
								},

								"hostname": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"path": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"query_string": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"fragment": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},

					"url_rewrite_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"source_pattern": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"destination": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"preserve_unmatched_path": {
									Type:     pluginsdk.TypeBool,
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
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"order": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"cookies_condition": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"selector": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"negate_condition": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"match_values": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"transforms": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"http_version_condition": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "Equal",
								},

								"negate_condition": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"match_values": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"device_condition": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "Equal",
								},

								"negate_condition": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"match_values": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"post_arg_condition": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"selector": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"negate_condition": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"match_values": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"transforms": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"query_string_condition": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"negate_condition": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"match_values": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"transforms": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"remote_address_condition": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"negate_condition": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"match_values": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"request_body_condition": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"negate_condition": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"match_values": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"transforms": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"request_header_condition": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"selector": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"negate_condition": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"match_values": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"transforms": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"request_method_condition": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "Equal",
								},

								"negate_condition": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"match_values": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"request_scheme_condition": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "Equal",
								},

								"negate_condition": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"match_values": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"request_uri_condition": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"negate_condition": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"match_values": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"transforms": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"url_file_extension_condition": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"negate_condition": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"match_values": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"transforms": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"url_file_name_condition": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"negate_condition": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"match_values": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"transforms": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"url_path_condition": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"negate_condition": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"match_values": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"transforms": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"cache_expiration_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"behavior": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"duration": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},

					"cache_key_query_string_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"behavior": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"parameters": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},

					"modify_request_header_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"action": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"value": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},

					"modify_response_header_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"action": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"value": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},

					"url_redirect_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"redirect_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"protocol": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "MatchRequest",
								},

								"hostname": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"path": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"query_string": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"fragment": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},

					"url_rewrite_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"source_pattern": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"destination": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"preserve_unmatched_path": {
									Type:     pluginsdk.TypeBool,
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
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (CdnEndpointV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
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

		newId := parse.NewEndpointID(oldParsedId.SubscriptionID, resourceGroup, profileName, name)
		newIdStr := newId.ID()

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newIdStr)

		rawState["id"] = newIdStr

		return rawState, nil
	}
}
