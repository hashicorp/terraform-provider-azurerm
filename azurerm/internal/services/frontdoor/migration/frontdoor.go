package migration

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/frontdoor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
)

func FrontDoorV0V1Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"header_frontdoor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"load_balancer_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"enforce_backend_pools_certificate_name_check": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"backend_pools_send_receive_timeout_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"location": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"routing_rule": {
				Type:     schema.TypeList,
				MaxItems: 100,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"accepted_protocols": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 2,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"patterns_to_match": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 25,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"frontend_endpoints": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 100,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"redirect_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"custom_fragment": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"custom_host": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"custom_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"custom_query_string": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"redirect_protocol": {
										Type:     schema.TypeString,
										Required: true,
									},
									"redirect_type": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"forwarding_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"backend_pool_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"cache_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"cache_use_dynamic_compression": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"cache_query_parameter_strip_directive": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"custom_forwarding_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"forwarding_protocol": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"backend_pool_load_balancing": {
				Type:     schema.TypeList,
				MaxItems: 5000,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"sample_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"successful_samples_required": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"additional_latency_milliseconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},

			"backend_pool_health_probe": {
				Type:     schema.TypeList,
				MaxItems: 5000,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"probe_method": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"interval_in_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},

			"backend_pool": {
				Type:     schema.TypeList,
				MaxItems: 50,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backend": {
							Type:     schema.TypeList,
							MaxItems: 100,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"address": {
										Type:     schema.TypeString,
										Required: true,
									},
									"http_port": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"https_port": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"priority": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"host_header": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"health_probe_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"load_balancing_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"frontend_endpoint": {
				Type:     schema.TypeList,
				MaxItems: 100,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"session_affinity_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"session_affinity_ttl_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"custom_https_provisioning_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"web_application_firewall_policy_link_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"custom_https_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate_source": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"minimum_tls_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"provisioning_state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"provisioning_substate": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"azure_key_vault_certificate_secret_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"azure_key_vault_certificate_secret_version": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"azure_key_vault_certificate_vault_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func FrontDoorV1ToV2(rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
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
	newIdStr := newId.ID("")

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newIdStr)

	rawState["id"] = newIdStr

	return rawState, nil
}
