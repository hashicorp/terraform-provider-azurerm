// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-04-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceFrontDoorSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorName,
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
			Default:  true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"routing_rule": {
			Type:     pluginsdk.TypeList,
			MaxItems: 500,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.BackendPoolRoutingRuleName,
					},
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"accepted_protocols": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 2,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(frontdoors.PossibleValuesForFrontDoorProtocol(), false),
						},
					},
					"patterns_to_match": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 25,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"frontend_endpoints": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 500,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
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
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(frontdoors.PossibleValuesForFrontDoorRedirectProtocol(), false),
								},
								"redirect_type": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(frontdoors.PossibleValuesForFrontDoorRedirectType(), false),
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
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validate.BackendPoolRoutingRuleName,
								},
								"cache_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
								"cache_use_dynamic_compression": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
								"cache_query_parameter_strip_directive": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(frontdoors.FrontDoorQueryStripAll),
									ValidateFunc: validation.StringInSlice(frontdoors.PossibleValuesForFrontDoorQuery(), false),
								},
								"cache_query_parameters": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 25,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"cache_duration": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.ISO8601DurationBetween("PT1S", "P365D"),
								},
								"custom_forwarding_path": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"forwarding_protocol": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(frontdoors.FrontDoorForwardingProtocolHTTPSOnly),
									ValidateFunc: validation.StringInSlice(frontdoors.PossibleValuesForFrontDoorForwardingProtocol(), false),
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.BackendPoolRoutingRuleName,
					},
					"sample_size": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  4,
					},
					"successful_samples_required": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  2,
					},
					"additional_latency_milliseconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  0,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.BackendPoolRoutingRuleName,
					},
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"path": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "/",
					},
					"protocol": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(frontdoors.FrontDoorProtocolHTTP),
						ValidateFunc: validation.StringInSlice(frontdoors.PossibleValuesForFrontDoorProtocol(), false),
					},
					"probe_method": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(frontdoors.FrontDoorHealthProbeMethodGET),
						ValidateFunc: validation.StringInSlice(frontdoors.PossibleValuesForFrontDoorHealthProbeMethod(), false),
					},
					"interval_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  120,
					},
				},
			},
		},

		"backend_pool": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"backend": {
						Type:     pluginsdk.TypeList,
						MaxItems: 500,
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
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 65535),
								},
								"https_port": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 65535),
								},
								"weight": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      50,
									ValidateFunc: validation.IntBetween(1, 1000),
								},
								"priority": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      1,
									ValidateFunc: validation.IntBetween(1, 5),
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.BackendPoolRoutingRuleName,
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

		"backend_pool_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enforce_backend_pools_certificate_name_check": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"backend_pools_send_receive_timeout_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      60,
						ValidateFunc: validation.IntBetween(0, 240),
					},
				},
			},
		},

		"frontend_endpoint": {
			Type:     pluginsdk.TypeList,
			MaxItems: 500,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.BackendPoolRoutingRuleName,
					},
					"host_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"session_affinity_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"session_affinity_ttl_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  0,
					},
					"web_application_firewall_policy_link_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.AsGeneratedID(webapplicationfirewallpolicies.ParseFrontDoorWebApplicationFirewallPolicyIDInsensitively),
					},
				},
			},
		},

		// Computed values
		"explicit_resource_order": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"backend_pool_ids": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"frontend_endpoint_ids": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"routing_rule_ids": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"backend_pool_load_balancing_ids": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"backend_pool_health_probe_ids": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"backend_pool_health_probes": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"backend_pool_load_balancing_settings": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"backend_pools": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"frontend_endpoints": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"routing_rules": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.Tags(),
	}
}
