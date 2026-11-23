// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_gateway

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/applicationgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appgateway/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func applicationGatewaySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": commonschema.Location(),

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		// lintignore:S016,S023
		"backend_address_pool": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"fqdns": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.NoZeroValues,
						},
					},

					"ip_addresses": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsIPv4Address,
						},
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
			Set: applicationGatewayBackendAddressPool,
		},

		// lintignore:S016,S017,S023
		"backend_http_settings": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			MinItems:     1,
			AtLeastOneOf: []string{"backend_http_settings", "backend"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"cookie_based_affinity": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(applicationgateways.PossibleValuesForApplicationGatewayCookieBasedAffinity(), false),
					},

					"port": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IsPortNumber,
					},

					"protocol": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(applicationgateways.ApplicationGatewayProtocolHTTP),
							string(applicationgateways.ApplicationGatewayProtocolHTTPS),
						}, false),
					},

					"affinity_cookie_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"certificate_chain_validation_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"connection_draining": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},

								"drain_timeout_sec": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 3600),
								},
							},
						},
					},

					"dedicated_backend_connection_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"host_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringStartsWithOneOf("/"),
					},

					"pick_host_name_from_backend_address": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"probe_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"request_timeout": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      30,
						ValidateFunc: validation.IntBetween(1, 86400),
					},

					"sni_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"sni_validation_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"trusted_root_certificate_names": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"probe_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
			Set: applicationGatewayBackendSettingsHash,
		},

		"backend": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			AtLeastOneOf: []string{"backend_http_settings", "backend"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ApplicationGatewayName,
					},

					"port": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IsPortNumber,
					},

					"protocol": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(applicationgateways.ApplicationGatewayProtocolTcp),
							string(applicationgateways.ApplicationGatewayProtocolTls),
						}, false),
					},

					"client_ip_preservation_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"host_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"probe_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.ApplicationGatewayName,
					},

					"timeout_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      30,
						ValidateFunc: validation.IntBetween(1, 86400),
					},

					"trusted_root_certificate_names": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.ApplicationGatewayName,
						},
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"probe_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"frontend_ip_configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"private_ip_address": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// Note: O+C because Azure assigns a private IP from the subnet when not specified
						Computed: true,
					},

					"public_ip_address_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"private_ip_address_allocation": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(applicationgateways.IPAllocationMethodDynamic),
						ValidateFunc: validation.StringInSlice(applicationgateways.PossibleValuesForIPAllocationMethod(), false),
					},

					"private_link_configuration_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"private_link_configuration_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"frontend_port": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"port": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IsPortNumber,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"gateway_ip_configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 2,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"global": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"request_buffering_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"response_buffering_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},

		// lintignore:S016,S023
		"http_listener": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			AtLeastOneOf: []string{"http_listener", "listener"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"frontend_ip_configuration_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"frontend_port_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"protocol": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(applicationgateways.ApplicationGatewayProtocolHTTP),
							string(applicationgateways.ApplicationGatewayProtocolHTTPS),
						}, false),
					},

					"host_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"host_names": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"ssl_certificate_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"require_sni": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"frontend_ip_configuration_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"frontend_port_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"ssl_certificate_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"ssl_profile_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"custom_error_configuration": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"status_code": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(applicationgateways.PossibleValuesForApplicationGatewayCustomErrorStatusCode(), false),
								},

								"custom_error_page_url": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"firewall_policy_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: webapplicationfirewallpolicies.ValidateApplicationGatewayWebApplicationFirewallPolicyID,
					},

					"ssl_profile_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
			Set: applicationGatewayHttpListnerHash,
		},

		"listener": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			AtLeastOneOf: []string{"http_listener", "listener"},
			Set:          applicationGatewayListenerHash,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ApplicationGatewayName,
					},

					"frontend_ip_configuration_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ApplicationGatewayName,
					},

					"frontend_port_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ApplicationGatewayName,
					},

					"protocol": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(applicationgateways.ApplicationGatewayProtocolTcp),
							string(applicationgateways.ApplicationGatewayProtocolTls),
						}, false),
					},

					"host_names": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"ssl_certificate_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.ApplicationGatewayName,
					},

					"ssl_profile_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.ApplicationGatewayName,
					},

					"frontend_ip_configuration_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"frontend_port_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"ssl_certificate_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"ssl_profile_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"fips_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"private_endpoint_connection": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"private_link_configuration": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"ip_configuration": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"subnet_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: commonids.ValidateSubnetID,
								},

								"private_ip_address": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									// Note: O+C because Azure assigns a private IP from the subnet when not specified
									Computed: true,
								},

								"private_ip_address_allocation": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(applicationgateways.PossibleValuesForIPAllocationMethod(), false),
								},

								"primary": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
							},
						},
					},
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"request_routing_rule": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			MinItems:     1,
			AtLeastOneOf: []string{"request_routing_rule", "routing_rule"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"rule_type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(applicationgateways.PossibleValuesForApplicationGatewayRequestRoutingRuleType(), false),
					},

					"http_listener_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"backend_address_pool_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"backend_http_settings_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"url_path_map_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"redirect_configuration_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"rewrite_rule_set_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"priority": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 20000),
					},

					"backend_address_pool_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"backend_http_settings_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"http_listener_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"url_path_map_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"redirect_configuration_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"rewrite_rule_set_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"routing_rule": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			AtLeastOneOf: []string{"request_routing_rule", "routing_rule"},
			Set:          applicationGatewayRoutingRuleHash,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ApplicationGatewayName,
					},

					"backend_address_pool_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ApplicationGatewayName,
					},

					"backend_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ApplicationGatewayName,
					},

					"listener_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ApplicationGatewayName,
					},

					"priority": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 20000),
					},

					"backend_address_pool_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"backend_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"listener_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"redirect_configuration": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"redirect_type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(applicationgateways.PossibleValuesForApplicationGatewayRedirectType(), false),
					},

					"target_listener_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"target_url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"include_path": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"include_query_string": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"target_listener_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
		"autoscale_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"min_capacity": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(0, 100),
					},
					"max_capacity": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(2, 125),
					},
				},
			},
		},
		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(applicationgateways.PossibleValuesForApplicationGatewaySkuName(), false),
					},

					"tier": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(applicationgateways.PossibleValuesForApplicationGatewayTier(), false),
					},

					"capacity": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"trusted_root_certificate": {
			Type:     pluginsdk.TypeList, // todo this should probably be a map
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"data": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Sensitive:    true,
					},

					"key_vault_secret_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret),
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		// lintignore:XS003
		"ssl_policy": sslProfileSchema(true),

		"http2_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"force_firewall_policy_association": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		// lintignore:S016,S023
		"probe": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"interval": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 86400),
					},

					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ApplicationGatewayName,
					},

					"protocol": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(applicationgateways.PossibleValuesForApplicationGatewayProtocol(), false),
					},

					"timeout": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 86400),
					},

					"unhealthy_threshold": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 20),
					},

					"host": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					// lintignore:XS003
					"match": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"body": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"status_code": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"minimum_servers": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  0,
					},

					"path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringStartsWithOneOf("/"),
					},

					"pick_host_name_from_backend_http_settings": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"port": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IsPortNumber,
					},

					"proxy_protocol_header_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
			Set: applicationGatewayProbeHash,
		},

		"rewrite_rule_set": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"rewrite_rule": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"rule_sequence": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 1000),
								},

								"condition": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"variable": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"pattern": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"ignore_case": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												Default:  false,
											},
											"negate": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												Default:  false,
											},
										},
									},
								},

								"request_header_configuration": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"header_name": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"header_value": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
										},
									},
								},

								"response_header_configuration": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"header_name": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
											"header_value": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},
										},
									},
								},

								"url": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"path": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},
											"query_string": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},

											"components": {
												Type:     pluginsdk.TypeString,
												Optional: true,
												ValidateFunc: validation.StringInSlice([]string{
													"path_only",
													"query_string_only",
												}, false),
											},

											"reroute": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												Default:  false,
											},
										},
									},
								},
							},
						},
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		// lintignore:S016,S023
		"ssl_certificate": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"data": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						StateFunc:    base64EncodedStateFunc,
						ValidateFunc: validation.StringIsBase64,
					},

					"password": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
					},

					"key_vault_secret_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeSecret),
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"public_cert_data": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
			Set: applicationGatewaySSLCertificate,
		},

		"trusted_client_certificate": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"data": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
						StateFunc: base64EncodedStateFunc,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"ssl_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"client_authentication_mode": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(applicationgateways.ApplicationGatewayClientAuthVerificationModesStrict),
						ValidateFunc: validation.StringInSlice(
							applicationgateways.PossibleValuesForApplicationGatewayClientAuthVerificationModes(),
							false,
						),
					},

					"trusted_client_certificate_names": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"verify_client_certificate_issuer_dn": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"verify_client_certificate_revocation": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(applicationgateways.ApplicationGatewayClientRevocationOptionsOCSP),
						}, false),
					},

					// lintignore:XS003
					"ssl_policy": sslProfileSchema(false),

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"url_path_map": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"default_backend_address_pool_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"default_backend_http_settings_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"default_redirect_configuration_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"default_rewrite_rule_set_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"path_rule": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"paths": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"backend_address_pool_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"backend_http_settings_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"redirect_configuration_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"rewrite_rule_set_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"backend_address_pool_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"backend_http_settings_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"redirect_configuration_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"rewrite_rule_set_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"firewall_policy_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: webapplicationfirewallpolicies.ValidateApplicationGatewayWebApplicationFirewallPolicyID,
								},
							},
						},
					},

					"default_backend_address_pool_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"default_backend_http_settings_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"default_redirect_configuration_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"default_rewrite_rule_set_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"waf_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"firewall_mode": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(applicationgateways.PossibleValuesForApplicationGatewayFirewallMode(), false),
					},

					"rule_set_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "OWASP",
						ValidateFunc: validate.ValidateWebApplicationFirewallPolicyRuleSetType,
					},

					"rule_set_version": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ValidateWebApplicationFirewallPolicyRuleSetVersion,
					},
					"file_upload_limit_mb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 750),
						Default:      100,
					},
					"request_body_check": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"max_request_body_size_kb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 128),
						Default:      128,
					},
					"disabled_rule_group": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"rule_group_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validate.ValidateWebApplicationFirewallPolicyRuleGroupName,
								},

								"rules": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeInt,
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
					},
					"exclusion": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"match_variable": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForOwaspCrsExclusionEntryMatchVariable(), false),
								},

								"selector_match_operator": {
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.StringInSlice(webapplicationfirewallpolicies.PossibleValuesForOwaspCrsExclusionEntrySelectorMatchOperator(), false),
									Optional:     true,
								},
								"selector": {
									ValidateFunc: validation.StringIsNotEmpty,
									Type:         pluginsdk.TypeString,
									Optional:     true,
								},
							},
						},
					},
				},
			},
		},

		"firewall_policy_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: webapplicationfirewallpolicies.ValidateApplicationGatewayWebApplicationFirewallPolicyID,
		},

		"custom_error_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"status_code": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(applicationgateways.PossibleValuesForApplicationGatewayCustomErrorStatusCode(), false),
					},

					"custom_error_page_url": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func sslProfileSchema(computed bool) *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Computed: computed,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"disabled_protocols": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice(applicationgateways.PossibleValuesForApplicationGatewaySslProtocol(), false),
					},
				},

				"policy_type": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(applicationgateways.PossibleValuesForApplicationGatewaySslPolicyType(), false),
				},

				"policy_name": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},

				"cipher_suites": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice(applicationgateways.PossibleValuesForApplicationGatewaySslCipherSuite(), false),
					},
				},

				"min_protocol_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(applicationgateways.PossibleValuesForApplicationGatewaySslProtocol(), false),
				},
			},
		},
	}
}
