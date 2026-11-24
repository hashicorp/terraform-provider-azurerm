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
			Description:  "The name of the Front Door service.",
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
			Description: "A friendly name for the Front Door service.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"load_balancer_enabled": {
			Description: "Should the Front Door Load Balancer be enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"routing_rule": {
			Type:     pluginsdk.TypeList,
			MaxItems: 500,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Description: "The ID of the Azure Front Door Backend Routing Rule.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"name": {
						Description:  "The name of the Routing Rule.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.BackendPoolRoutingRuleName,
					},
					"enabled": {
						Description: "`Enable` or `Disable` use of this Backend Routing Rule.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
					"accepted_protocols": {
						Description: "Protocol schemes to match for the Backend Routing Rule.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						MaxItems:    2,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(frontdoors.PossibleValuesForFrontDoorProtocol(), false),
						},
					},
					"patterns_to_match": {
						Description: "The route patterns for the Backend Routing Rule.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						MaxItems:    25,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"frontend_endpoints": {
						Description: "The names of the `frontend_endpoint` blocks within this resource to associate with this `routing_rule`.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						MaxItems:    500,
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
									Description: "The destination fragment in the portion of URL after '#'. Set this to add a fragment to the redirect URL.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
								},
								"custom_host": {
									Description: "Set this to change the URL for the redirection.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
								},
								"custom_path": {
									Description: "The path to retain as per the incoming request, or update in the URL for the redirection.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
								},
								"custom_query_string": {
									Description: "Replace any existing query string from the incoming request URL.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
								},
								"redirect_protocol": {
									Description:  "Protocol to use when redirecting.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(frontdoors.PossibleValuesForFrontDoorRedirectProtocol(), false),
								},
								"redirect_type": {
									Description:  "Status code for the redirect.",
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
									Description:  "The name of the Backend Pool to forward the incoming traffic to.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validate.BackendPoolRoutingRuleName,
								},
								"cache_enabled": {
									Description: "Whether to enable caching.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     false,
								},
								"cache_use_dynamic_compression": {
									Description: "Whether to use dynamic compression when caching.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     false,
								},
								"cache_query_parameter_strip_directive": {
									Description:  "The cache behaviour in relation to query string parameters.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(frontdoors.FrontDoorQueryStripAll),
									ValidateFunc: validation.StringInSlice(frontdoors.PossibleValuesForFrontDoorQuery(), false),
								},
								"cache_query_parameters": {
									Description: "A list of query parameters. Works only in combination with `cache_query_parameter_strip_directive` set to `StripAllExcept` or `StripOnly`.",
									Type:        pluginsdk.TypeList,
									Optional:    true,
									MaxItems:    25,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"cache_duration": {
									Description:  "The minimum caching duration (in ISO8601 notation e.g. `P1DT2H` for 1 day and 2 hours). Needs to be greater than 0 and smaller than 365 days. `cache_duration` works only in combination with `cache_enabled` set to `true`.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.ISO8601DurationBetween("PT1S", "P365D"),
								},
								"custom_forwarding_path": {
									Description: "Path to use when constructing the request to forward to the backend. This functions as a URL Rewrite.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
								},
								"forwarding_protocol": {
									Description:  "Protocol to use when redirecting.",
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
						Description: "The ID of the Azure Front Door Backend Load Balancer.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"name": {
						Description:  "The name of the Load Balancer.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.BackendPoolRoutingRuleName,
					},
					"sample_size": {
						Description: "The number of samples to consider for load balancing decisions.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						Default:     4,
					},
					"successful_samples_required": {
						Description: "The number of samples within the sample period that must succeed.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						Default:     2,
					},
					"additional_latency_milliseconds": {
						Description: "The additional latency in milliseconds for probes to fall into the lowest latency bucket.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						Default:     0,
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
						Description: "The ID of the Azure Front Door Backend Health Probe.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"name": {
						Description:  "The name of the Health Probe.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.BackendPoolRoutingRuleName,
					},
					"enabled": {
						Description: "Is this health probe enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
					"path": {
						Description: "The path to use for the Health Probe.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     "/",
					},
					"protocol": {
						Description:  "Protocol scheme to use for the Health Probe.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(frontdoors.FrontDoorProtocolHTTP),
						ValidateFunc: validation.StringInSlice(frontdoors.PossibleValuesForFrontDoorProtocol(), false),
					},
					"probe_method": {
						Description:  "The HTTP method the health probe uses when querying the backend pool instances.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(frontdoors.FrontDoorHealthProbeMethodGET),
						ValidateFunc: validation.StringInSlice(frontdoors.PossibleValuesForFrontDoorHealthProbeMethod(), false),
					},
					"interval_in_seconds": {
						Description: "The number of seconds between each Health Probe.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						Default:     120,
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
									Description: "Whether the backend is enabled or not.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     true,
								},
								"address": {
									Description: "Location of the backend (IP address or FQDN).",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
								"http_port": {
									Description:  "The HTTP TCP port number.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 65535),
								},
								"https_port": {
									Description:  "The HTTPS TCP port number.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 65535),
								},
								"weight": {
									Description:  "Weight of this endpoint for load balancing purposes.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      50,
									ValidateFunc: validation.IntBetween(1, 1000),
								},
								"priority": {
									Description:  "Priority to use for load balancing. Higher priorities will not be used for load balancing if any lower priority backend is healthy.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      1,
									ValidateFunc: validation.IntBetween(1, 5),
								},
								"host_header": {
									Description: "The value to use as the host header sent to the backend.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
							},
						},
					},
					"id": {
						Description: "The ID of the Azure Front Door Backend Pool.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"name": {
						Description:  "The name of the Backend Pool.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.BackendPoolRoutingRuleName,
					},
					"health_probe_name": {
						Description: "The name of the `backend_pool_health_probe` block within this resource to use for this `Backend Pool`.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"load_balancing_name": {
						Description: "The name of the `backend_pool_load_balancing` block within this resource to use for this `Backend Pool`.",
						Type:        pluginsdk.TypeString,
						Required:    true,
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
						Description: "Enforce certificate name check on `HTTPS` requests to all backend pools, this setting will have no effect on `HTTP` requests.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},

					"backend_pools_send_receive_timeout_seconds": {
						Description:  "The send and receive timeout on forwarding request to the backend. When the timeout is reached, the request fails and returns.",
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
						Description: "The ID of the Azure Front Door Frontend Endpoint.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"name": {
						Description:  "The name of the `frontend_endpoint`.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.BackendPoolRoutingRuleName,
					},
					"host_name": {
						Description: "The host name of the `frontend_endpoint`. In order to use a name.azurefd.net domain, the name value must match the Front Door name.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"session_affinity_enabled": {
						Description: "Whether to allow session affinity on this host.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"session_affinity_ttl_seconds": {
						Description: "The TTL to use in seconds for session affinity, if applicable.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						Default:     0,
					},
					"web_application_firewall_policy_link_id": {
						Description:  "The Web Application Firewall Policy `ID` for each host.",
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
			Description: "A map/dictionary of Backend Pool Health Probe Names (key) to the Backend Pool Health Probe ID (value).",
			Type:        pluginsdk.TypeMap,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"backend_pool_load_balancing_settings": {
			Description: "A map/dictionary of Backend Pool Load Balancing Setting Names (key) to the Backend Pool Load Balancing Setting ID (value).",
			Type:        pluginsdk.TypeMap,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"backend_pools": {
			Description: "A map/dictionary of Backend Pool Names (key) to the Backend Pool ID (value).",
			Type:        pluginsdk.TypeMap,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"frontend_endpoints": {
			Description: "A map/dictionary of Frontend Endpoint Names (key) to the Frontend Endpoint ID (value).",
			Type:        pluginsdk.TypeMap,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"routing_rules": {
			Description: "A map/dictionary of Routing Rule Names (key) to the Routing Rule ID (value).",
			Type:        pluginsdk.TypeMap,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.Tags(),
	}
}
