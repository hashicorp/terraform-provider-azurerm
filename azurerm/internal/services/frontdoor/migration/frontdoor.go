package migration

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = FrontDoorUpgradeV0ToV1{}

type FrontDoorUpgradeV0ToV1 struct{}

func (FrontDoorUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return frontDoorSchemaForV0AndV1()
}

func (FrontDoorUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// this resource was set to "schema version 1" unintentionally.. so we're adding
		// a "fake" upgrade here to account for it
		return rawState, nil
	}
}

var _ pluginsdk.StateUpgrade = FrontDoorUpgradeV1ToV2{}

type FrontDoorUpgradeV1ToV2 struct{}

func (FrontDoorUpgradeV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return frontDoorSchemaForV0AndV1()
}

func (FrontDoorUpgradeV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/frontdoors/{frontDoorName}
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/frontDoors/{frontDoorName}
		oldId := rawState["id"].(string)
		oldParsedId, err := azure.ParseAzureResourceID(oldId)
		if err != nil {
			return rawState, err
		}

		resourceGroup := oldParsedId.ResourceGroup
		frontDoorName := ""
		for key, value := range oldParsedId.Path {
			if strings.EqualFold(key, "frontDoors") {
				frontDoorName = value
				break
			}
		}

		if frontDoorName == "" {
			return rawState, fmt.Errorf("couldn't find the `frontDoors` segment in the old resource id %q", oldId)
		}

		newId := parse.NewFrontDoorID(oldParsedId.SubscriptionID, resourceGroup, frontDoorName)
		newIdStr := newId.ID()

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newIdStr)

		rawState["id"] = newIdStr

		return rawState, nil
	}
}

func frontDoorSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"cname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"header_frontdoor_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"friendly_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"load_balancer_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"enforce_backend_pools_certificate_name_check": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"backend_pools_send_receive_timeout_seconds": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"routing_rule": {
			Type:     pluginsdk.TypeList,
			MaxItems: 100,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"accepted_protocols": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 2,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"patterns_to_match": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 25,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"frontend_endpoints": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 100,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"redirect_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"custom_fragment": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"custom_host": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"custom_path": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"custom_query_string": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"redirect_protocol": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"redirect_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
					"forwarding_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"backend_pool_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"cache_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},
								"cache_use_dynamic_compression": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},
								"cache_query_parameter_strip_directive": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"custom_forwarding_path": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"forwarding_protocol": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},

		"backend_pool_load_balancing": {
			Type:     pluginsdk.TypeList,
			MaxItems: 5000,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"sample_size": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
					"successful_samples_required": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
					"additional_latency_milliseconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"backend_pool_health_probe": {
			Type:     pluginsdk.TypeList,
			MaxItems: 5000,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"path": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"protocol": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"probe_method": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"interval_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"backend_pool": {
			Type:     pluginsdk.TypeList,
			MaxItems: 50,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"backend": {
						Type:     pluginsdk.TypeList,
						MaxItems: 100,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  true,
								},
								"address": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"http_port": {
									Type:     pluginsdk.TypeInt,
									Required: true,
								},
								"https_port": {
									Type:     pluginsdk.TypeInt,
									Required: true,
								},
								"weight": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
								"priority": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
								"host_header": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"health_probe_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"load_balancing_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"frontend_endpoint": {
			Type:     pluginsdk.TypeList,
			MaxItems: 100,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"host_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"session_affinity_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"session_affinity_ttl_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
					"custom_https_provisioning_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true,
					},
					"web_application_firewall_policy_link_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					//lintignore:XS003
					"custom_https_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"certificate_source": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"minimum_tls_version": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"provisioning_state": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"provisioning_substate": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"azure_key_vault_certificate_secret_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"azure_key_vault_certificate_secret_version": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"azure_key_vault_certificate_vault_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
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
