// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/applicationgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func base64EncodedStateFunc(v interface{}) string {
	switch s := v.(type) {
	case string:
		return utils.Base64EncodeIfNot(s)
	default:
		return ""
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
						Type: pluginsdk.TypeString,
						ValidateFunc: validation.StringInSlice([]string{
							string(applicationgateways.ApplicationGatewaySslProtocolTLSvOneZero),
							string(applicationgateways.ApplicationGatewaySslProtocolTLSvOneOne),
							string(applicationgateways.ApplicationGatewaySslProtocolTLSvOneTwo),
							string(applicationgateways.ApplicationGatewaySslProtocolTLSvOneThree),
						}, false),
					},
				},

				"policy_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(applicationgateways.ApplicationGatewaySslPolicyTypeCustom),
						string(applicationgateways.ApplicationGatewaySslPolicyTypeCustomVTwo),
						string(applicationgateways.ApplicationGatewaySslPolicyTypePredefined),
					}, false),
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
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(applicationgateways.ApplicationGatewaySslProtocolTLSvOneZero),
						string(applicationgateways.ApplicationGatewaySslProtocolTLSvOneOne),
						string(applicationgateways.ApplicationGatewaySslProtocolTLSvOneTwo),
						string(applicationgateways.ApplicationGatewaySslProtocolTLSvOneThree),
					}, false),
				},
			},
		},
	}
}

func resourceApplicationGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApplicationGatewayCreate,
		Read:   resourceApplicationGatewayRead,
		Update: resourceApplicationGatewayUpdate,
		Delete: resourceApplicationGatewayDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := applicationgateways.ParseApplicationGatewayID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
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
								ValidateFunc: validate.IPv4Address,
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
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"path": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"port": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validate.PortNumber,
						},

						"protocol": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(applicationgateways.ApplicationGatewayProtocolHTTP),
								string(applicationgateways.ApplicationGatewayProtocolHTTPS),
							}, false),
						},

						"cookie_based_affinity": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(applicationgateways.ApplicationGatewayCookieBasedAffinityEnabled),
								string(applicationgateways.ApplicationGatewayCookieBasedAffinityDisabled),
							}, false),
						},

						"affinity_cookie_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"host_name": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"pick_host_name_from_backend_address": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"request_timeout": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      30,
							ValidateFunc: validation.IntBetween(1, 86400),
						},

						"authentication_certificate": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
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

						"trusted_root_certificate_names": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
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

						"probe_name": {
							Type:     pluginsdk.TypeString,
							Optional: true,
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
							Computed: true,
						},

						"public_ip_address_id": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"private_ip_address_allocation": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(applicationgateways.IPAllocationMethodDynamic),
							ValidateFunc: validation.StringInSlice([]string{
								string(applicationgateways.IPAllocationMethodDynamic),
								string(applicationgateways.IPAllocationMethodStatic),
							}, false),
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
							ValidateFunc: validate.PortNumber,
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
				Type:     pluginsdk.TypeSet,
				Required: true,
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
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(applicationgateways.ApplicationGatewayCustomErrorStatusCodeHTTPStatusFourZeroThree),
											string(applicationgateways.ApplicationGatewayCustomErrorStatusCodeHTTPStatusFiveZeroTwo),
										}, false),
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
										Computed: true,
									},

									"private_ip_address_allocation": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(applicationgateways.IPAllocationMethodDynamic),
											string(applicationgateways.IPAllocationMethodStatic),
										}, false),
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
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"rule_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(applicationgateways.ApplicationGatewayRequestRoutingRuleTypeBasic),
								string(applicationgateways.ApplicationGatewayRequestRoutingRuleTypePathBasedRouting),
							}, false),
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
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(applicationgateways.ApplicationGatewayRedirectTypePermanent),
								string(applicationgateways.ApplicationGatewayRedirectTypeTemporary),
								string(applicationgateways.ApplicationGatewayRedirectTypeFound),
								string(applicationgateways.ApplicationGatewayRedirectTypeSeeOther),
							}, false),
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
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(applicationgateways.ApplicationGatewaySkuNameBasic),
								string(applicationgateways.ApplicationGatewaySkuNameStandardSmall),
								string(applicationgateways.ApplicationGatewaySkuNameStandardMedium),
								string(applicationgateways.ApplicationGatewaySkuNameStandardLarge),
								string(applicationgateways.ApplicationGatewaySkuNameStandardVTwo),
								string(applicationgateways.ApplicationGatewaySkuNameWAFLarge),
								string(applicationgateways.ApplicationGatewaySkuNameWAFMedium),
								string(applicationgateways.ApplicationGatewaySkuNameWAFVTwo),
							}, false),
						},

						"tier": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(applicationgateways.ApplicationGatewayTierBasic),
								string(applicationgateways.ApplicationGatewayTierStandard),
								string(applicationgateways.ApplicationGatewayTierStandardVTwo),
								string(applicationgateways.ApplicationGatewayTierWAF),
								string(applicationgateways.ApplicationGatewayTierWAFVTwo),
							}, false),
						},

						"capacity": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
						},
					},
				},
			},

			// Optional
			"authentication_certificate": {
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
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							Sensitive:    true,
						},

						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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
							ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
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

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_http2": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
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
						"name": {
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

						"path": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"host": {
							Type:     pluginsdk.TypeString,
							Optional: true,
						},

						"interval": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},

						"timeout": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},

						"unhealthy_threshold": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},

						"port": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validate.PortNumber,
						},

						"pick_host_name_from_backend_http_settings": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"minimum_servers": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Default:  0,
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
							ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
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

						"trusted_client_certificate_names": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						// TODO: replace cert by certificate in 4.0
						"verify_client_cert_issuer_dn": {
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
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(applicationgateways.ApplicationGatewayFirewallModeDetection),
								string(applicationgateways.ApplicationGatewayFirewallModePrevention),
							}, false),
						},

						"rule_set_type": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "OWASP",
							ValidateFunc: networkValidate.ValidateWebApplicationFirewallPolicyRuleSetType,
						},

						"rule_set_version": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: networkValidate.ValidateWebApplicationFirewallPolicyRuleSetVersion,
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
										ValidateFunc: networkValidate.ValidateWebApplicationFirewallPolicyRuleGroupName,
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
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestArgKeys),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestArgNames),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestArgValues),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestCookieKeys),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestCookieNames),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestCookieValues),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestHeaderKeys),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestHeaderNames),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariableRequestHeaderValues),
										}, false),
									},

									"selector_match_operator": {
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntrySelectorMatchOperatorContains),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntrySelectorMatchOperatorEndsWith),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntrySelectorMatchOperatorEquals),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntrySelectorMatchOperatorEqualsAny),
											string(webapplicationfirewallpolicies.OwaspCrsExclusionEntrySelectorMatchOperatorStartsWith),
										}, false),
										Optional: true,
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
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(applicationgateways.ApplicationGatewayCustomErrorStatusCodeHTTPStatusFourZeroThree),
								string(applicationgateways.ApplicationGatewayCustomErrorStatusCodeHTTPStatusFiveZeroTwo),
							}, false),
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
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(applicationGatewayCustomizeDiff),
	}
}

func resourceApplicationGatewayCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ApplicationGateways
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := applicationgateways.NewApplicationGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_application_gateway", id.ID())
	}

	enablehttp2 := d.Get("enable_http2").(bool)
	t := d.Get("tags").(map[string]interface{})

	// Gateway ID is needed to link sub-resources together in expand functions
	trustedRootCertificates, err := expandApplicationGatewayTrustedRootCertificates(d.Get("trusted_root_certificate").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `trusted_root_certificate`: %+v", err)
	}

	requestRoutingRules, err := expandApplicationGatewayRequestRoutingRules(d, id.ID())
	if err != nil {
		return fmt.Errorf("expanding `request_routing_rule`: %+v", err)
	}

	urlPathMaps, err := expandApplicationGatewayURLPathMaps(d, id.ID())
	if err != nil {
		return fmt.Errorf("expanding `url_path_map`: %+v", err)
	}

	redirectConfigurations, err := expandApplicationGatewayRedirectConfigurations(d, id.ID())
	if err != nil {
		return fmt.Errorf("expanding `redirect_configuration`: %+v", err)
	}

	sslCertificates, err := expandApplicationGatewaySslCertificates(d)
	if err != nil {
		return fmt.Errorf("expanding `ssl_certificate`: %+v", err)
	}

	trustedClientCertificates, err := expandApplicationGatewayTrustedClientCertificates(d)
	if err != nil {
		return fmt.Errorf("expanding `trusted_client_certificate`: %+v", err)
	}

	sslProfiles := expandApplicationGatewaySslProfiles(d, id.ID())

	gatewayIPConfigurations, _ := expandApplicationGatewayIPConfigurations(d)

	globalConfiguration := expandApplicationGatewayGlobalConfiguration(d.Get("global").([]interface{}))

	httpListeners, err := expandApplicationGatewayHTTPListeners(d, id.ID())
	if err != nil {
		return fmt.Errorf("fail to expand `http_listener`: %+v", err)
	}

	rewriteRuleSets, err := expandApplicationGatewayRewriteRuleSets(d)
	if err != nil {
		return fmt.Errorf("expanding `rewrite_rule_set`: %v", err)
	}

	gateway := applicationgateways.ApplicationGateway{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(t),
		Properties: &applicationgateways.ApplicationGatewayPropertiesFormat{
			AutoscaleConfiguration:        expandApplicationGatewayAutoscaleConfiguration(d),
			AuthenticationCertificates:    expandApplicationGatewayAuthenticationCertificates(d.Get("authentication_certificate").([]interface{})),
			TrustedRootCertificates:       trustedRootCertificates,
			CustomErrorConfigurations:     expandApplicationGatewayCustomErrorConfigurations(d.Get("custom_error_configuration").([]interface{})),
			BackendAddressPools:           expandApplicationGatewayBackendAddressPools(d),
			BackendHTTPSettingsCollection: expandApplicationGatewayBackendHTTPSettings(d, id.ID()),
			EnableHTTP2:                   pointer.To(enablehttp2),
			FrontendIPConfigurations:      expandApplicationGatewayFrontendIPConfigurations(d, id.ID()),
			FrontendPorts:                 expandApplicationGatewayFrontendPorts(d),
			GatewayIPConfigurations:       gatewayIPConfigurations,
			GlobalConfiguration:           globalConfiguration,
			HTTPListeners:                 httpListeners,
			PrivateLinkConfigurations:     expandApplicationGatewayPrivateLinkConfigurations(d),
			Probes:                        expandApplicationGatewayProbes(d),
			RequestRoutingRules:           requestRoutingRules,
			RedirectConfigurations:        redirectConfigurations,
			Sku:                           expandApplicationGatewaySku(d),
			SslCertificates:               sslCertificates,
			TrustedClientCertificates:     trustedClientCertificates,
			SslProfiles:                   sslProfiles,
			SslPolicy:                     expandApplicationGatewaySslPolicy(d.Get("ssl_policy").([]interface{})),

			RewriteRuleSets: rewriteRuleSets,
			UrlPathMaps:     urlPathMaps,
		},
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		gateway.Zones = &zones
	}

	if v, ok := d.GetOk("fips_enabled"); ok {
		gateway.Properties.EnableFips = pointer.To(v.(bool))
	}

	if v, ok := d.GetOk("force_firewall_policy_association"); ok {
		gateway.Properties.ForceFirewallPolicyAssociation = pointer.To(v.(bool))
	}

	if _, ok := d.GetOk("identity"); ok {
		expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		gateway.Identity = expandedIdentity
	}

	// validation (todo these should probably be moved into their respective expand functions, which would then return an error?)
	for _, backendHttpSettings := range *gateway.Properties.BackendHTTPSettingsCollection {
		if props := backendHttpSettings.Properties; props != nil {
			if props.HostName == nil || props.PickHostNameFromBackendAddress == nil {
				continue
			}

			if *props.HostName != "" && *props.PickHostNameFromBackendAddress {
				return fmt.Errorf("Only one of `host_name` or `pick_host_name_from_backend_address` can be set")
			}
		}
	}

	for _, probe := range *gateway.Properties.Probes {
		if props := probe.Properties; props != nil {
			if props.Host == nil || props.PickHostNameFromBackendHTTPSettings == nil {
				continue
			}

			if *props.Host == "" && !*props.PickHostNameFromBackendHTTPSettings {
				return fmt.Errorf("One of `host` or `pick_host_name_from_backend_http_settings` must be set")
			}

			if *props.Host != "" && *props.PickHostNameFromBackendHTTPSettings {
				return fmt.Errorf("Only one of `host` or `pick_host_name_from_backend_http_settings` can be set")
			}
		}
	}

	if _, ok := d.GetOk("waf_configuration"); ok {
		gateway.Properties.WebApplicationFirewallConfiguration = expandApplicationGatewayWafConfig(d)
	}

	appGWSkuTier := d.Get("sku.0.tier").(string)
	wafFileUploadLimit := d.Get("waf_configuration.0.file_upload_limit_mb").(int)

	if appGWSkuTier != string(applicationgateways.ApplicationGatewayTierWAFVTwo) && wafFileUploadLimit > 500 {
		return fmt.Errorf("Only SKU `%s` allows `file_upload_limit_mb` to exceed 500MB", applicationgateways.ApplicationGatewayTierWAFVTwo)
	}

	if v, ok := d.GetOk("firewall_policy_id"); ok {
		id := v.(string)
		gateway.Properties.FirewallPolicy = &applicationgateways.SubResource{
			Id: &id,
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, gateway); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceApplicationGatewayRead(d, meta)
}

func resourceApplicationGatewayUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ApplicationGateways
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := applicationgateways.ParseApplicationGatewayID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if payload.Properties == nil {
		payload.Properties = &applicationgateways.ApplicationGatewayPropertiesFormat{}
	}

	if d.HasChange("enable_http2") {
		payload.Properties.EnableHTTP2 = pointer.To(d.Get("enable_http2").(bool))
	}

	if d.HasChange("trusted_root_certificate") {
		trustedRootCertificates, err := expandApplicationGatewayTrustedRootCertificates(d.Get("trusted_root_certificate").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `trusted_root_certificate`: %+v", err)
		}
		payload.Properties.TrustedRootCertificates = trustedRootCertificates
	}

	if d.HasChange("request_routing_rule") {
		requestRoutingRules, err := expandApplicationGatewayRequestRoutingRules(d, id.ID())
		if err != nil {
			return fmt.Errorf("expanding `request_routing_rule`: %+v", err)
		}
		payload.Properties.RequestRoutingRules = requestRoutingRules
	}

	if d.HasChange("url_path_map") {
		urlPathMaps, err := expandApplicationGatewayURLPathMaps(d, id.ID())
		if err != nil {
			return fmt.Errorf("expanding `url_path_map`: %+v", err)
		}

		payload.Properties.UrlPathMaps = urlPathMaps
	}

	if d.HasChange("redirect_configuration") {
		redirectConfigurations, err := expandApplicationGatewayRedirectConfigurations(d, id.ID())
		if err != nil {
			return fmt.Errorf("expanding `redirect_configuration`: %+v", err)
		}

		payload.Properties.RedirectConfigurations = redirectConfigurations
	}

	if d.HasChange("ssl_certificate") {
		sslCertificates, err := expandApplicationGatewaySslCertificates(d)
		if err != nil {
			return fmt.Errorf("expanding `ssl_certificate`: %+v", err)
		}

		payload.Properties.SslCertificates = sslCertificates
	}

	if d.HasChange("trusted_client_certificate") {
		trustedClientCertificates, err := expandApplicationGatewayTrustedClientCertificates(d)
		if err != nil {
			return fmt.Errorf("expanding `trusted_client_certificate`: %+v", err)
		}

		payload.Properties.TrustedClientCertificates = trustedClientCertificates
	}

	if d.HasChange("ssl_profile") {
		payload.Properties.SslProfiles = expandApplicationGatewaySslProfiles(d, id.ID())
	}

	gatewayIPConfigurations, stopApplicationGateway := expandApplicationGatewayIPConfigurations(d)
	if d.HasChange("gateway_ip_configuration") {
		payload.Properties.GatewayIPConfigurations = gatewayIPConfigurations
	}

	if d.HasChange("global") {
		globalConfiguration := expandApplicationGatewayGlobalConfiguration(d.Get("global").([]interface{}))
		payload.Properties.GlobalConfiguration = globalConfiguration
	}

	if d.HasChange("http_listener") {
		httpListeners, err := expandApplicationGatewayHTTPListeners(d, id.ID())
		if err != nil {
			return fmt.Errorf("fail to expand `http_listener`: %+v", err)
		}

		payload.Properties.HTTPListeners = httpListeners
	}

	if d.HasChange("rewrite_rule_set") {
		rewriteRuleSets, err := expandApplicationGatewayRewriteRuleSets(d)
		if err != nil {
			return fmt.Errorf("expanding `rewrite_rule_set`: %v", err)
		}

		payload.Properties.RewriteRuleSets = rewriteRuleSets
	}

	if d.HasChange("autoscale_configuration") {
		payload.Properties.AutoscaleConfiguration = expandApplicationGatewayAutoscaleConfiguration(d)
	}

	if d.HasChange("authentication_certificate") {
		payload.Properties.AuthenticationCertificates = expandApplicationGatewayAuthenticationCertificates(d.Get("authentication_certificate").([]interface{}))
	}

	if d.HasChange("custom_error_configuration") {
		payload.Properties.CustomErrorConfigurations = expandApplicationGatewayCustomErrorConfigurations(d.Get("custom_error_configuration").([]interface{}))
	}

	if d.HasChange("backend_address_pool") {
		payload.Properties.BackendAddressPools = expandApplicationGatewayBackendAddressPools(d)
	}

	if d.HasChange("backend_http_settings") {
		payload.Properties.BackendHTTPSettingsCollection = expandApplicationGatewayBackendHTTPSettings(d, id.ID())
	}

	if d.HasChange("frontend_ip_configuration") {
		payload.Properties.FrontendIPConfigurations = expandApplicationGatewayFrontendIPConfigurations(d, id.ID())
	}

	if d.HasChange("frontend_port") {
		payload.Properties.FrontendPorts = expandApplicationGatewayFrontendPorts(d)
	}

	if d.HasChange("private_link_configuration") {
		payload.Properties.PrivateLinkConfigurations = expandApplicationGatewayPrivateLinkConfigurations(d)
	}

	if d.HasChange("probe") {
		payload.Properties.Probes = expandApplicationGatewayProbes(d)
	}

	if d.HasChange("sku") {
		payload.Properties.Sku = expandApplicationGatewaySku(d)
	}

	if d.HasChange("ssl_policy") {
		payload.Properties.SslPolicy = expandApplicationGatewaySslPolicy(d.Get("ssl_policy").([]interface{}))
	}

	if d.HasChange("zones") {
		zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
		if len(zones) > 0 {
			payload.Zones = &zones
		}
	}

	if d.HasChange("fips_enabled") {
		payload.Properties.EnableFips = pointer.To(d.Get("fips_enabled").(bool))
	}

	if d.HasChange("force_firewall_policy_association") {
		payload.Properties.ForceFirewallPolicyAssociation = pointer.To(d.Get("force_firewall_policy_association").(bool))
	}

	if _, ok := d.GetOk("identity"); ok {
		expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		payload.Identity = expandedIdentity
	}

	// validation (todo these should probably be moved into their respective expand functions, which would then return an error?)
	if payload.Properties != nil && payload.Properties.BackendHTTPSettingsCollection != nil {
		for _, backendHttpSettings := range *payload.Properties.BackendHTTPSettingsCollection {
			if props := backendHttpSettings.Properties; props != nil {
				if props.HostName == nil || props.PickHostNameFromBackendAddress == nil {
					continue
				}

				if *props.HostName != "" && *props.PickHostNameFromBackendAddress {
					return fmt.Errorf("Only one of `host_name` or `pick_host_name_from_backend_address` can be set")
				}
			}
		}
	}

	if payload.Properties != nil && payload.Properties.Probes != nil {
		for _, probe := range *payload.Properties.Probes {
			if props := probe.Properties; props != nil {
				if props.Host == nil || props.PickHostNameFromBackendHTTPSettings == nil {
					continue
				}

				if *props.Host == "" && !*props.PickHostNameFromBackendHTTPSettings {
					return fmt.Errorf("One of `host` or `pick_host_name_from_backend_http_settings` must be set")
				}

				if *props.Host != "" && *props.PickHostNameFromBackendHTTPSettings {
					return fmt.Errorf("Only one of `host` or `pick_host_name_from_backend_http_settings` can be set")
				}
			}
		}
	}

	if d.HasChange("waf_configuration") {
		payload.Properties.WebApplicationFirewallConfiguration = expandApplicationGatewayWafConfig(d)
	}

	appGWSkuTier := d.Get("sku.0.tier").(string)
	wafFileUploadLimit := d.Get("waf_configuration.0.file_upload_limit_mb").(int)

	if appGWSkuTier != string(applicationgateways.ApplicationGatewayTierWAFVTwo) && wafFileUploadLimit > 500 {
		return fmt.Errorf("Only SKU `%s` allows `file_upload_limit_mb` to exceed 500MB", applicationgateways.ApplicationGatewayTierWAFVTwo)
	}

	if d.HasChange("firewall_policy_id") {
		if d.Get("firewall_policy_id").(string) != "" {
			payload.Properties.FirewallPolicy = &applicationgateways.SubResource{
				Id: pointer.To(d.Get("firewall_policy_id").(string)),
			}
		} else {
			payload.Properties.FirewallPolicy = nil
		}
	}

	if stopApplicationGateway {
		if err := client.StopThenPoll(ctx, *id); err != nil {
			return fmt.Errorf("stopping %s: %+v", id, err)
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if stopApplicationGateway {
		if err := client.StartThenPoll(ctx, *id); err != nil {
			return fmt.Errorf("starting %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())
	return resourceApplicationGatewayRead(d, meta)
}

func resourceApplicationGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ApplicationGateways
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := applicationgateways.ParseApplicationGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ApplicationGatewayName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("zones", zones.FlattenUntyped(model.Zones))

		identity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return err
		}
		if err = d.Set("identity", identity); err != nil {
			return err
		}

		if props := model.Properties; props != nil {
			if err = d.Set("authentication_certificate", flattenApplicationGatewayAuthenticationCertificates(props.AuthenticationCertificates, d)); err != nil {
				return fmt.Errorf("setting `authentication_certificate`: %+v", err)
			}

			if err = d.Set("trusted_root_certificate", flattenApplicationGatewayTrustedRootCertificates(props.TrustedRootCertificates, d)); err != nil {
				return fmt.Errorf("setting `trusted_root_certificate`: %+v", err)
			}

			if setErr := d.Set("backend_address_pool", flattenApplicationGatewayBackendAddressPools(props.BackendAddressPools)); setErr != nil {
				return fmt.Errorf("setting `backend_address_pool`: %+v", setErr)
			}

			backendHttpSettings, err := flattenApplicationGatewayBackendHTTPSettings(props.BackendHTTPSettingsCollection)
			if err != nil {
				return fmt.Errorf("flattening `backend_http_settings`: %+v", err)
			}
			if setErr := d.Set("backend_http_settings", backendHttpSettings); setErr != nil {
				return fmt.Errorf("setting `backend_http_settings`: %+v", setErr)
			}

			if setErr := d.Set("ssl_policy", flattenApplicationGatewaySslPolicy(props.SslPolicy)); setErr != nil {
				return fmt.Errorf("setting `ssl_policy`: %+v", setErr)
			}

			d.Set("enable_http2", props.EnableHTTP2)
			d.Set("fips_enabled", props.EnableFips)
			d.Set("force_firewall_policy_association", props.ForceFirewallPolicyAssociation)

			httpListeners, err := flattenApplicationGatewayHTTPListeners(props.HTTPListeners)
			if err != nil {
				return fmt.Errorf("flattening `http_listener`: %+v", err)
			}
			if setErr := d.Set("http_listener", httpListeners); setErr != nil {
				return fmt.Errorf("setting `http_listener`: %+v", setErr)
			}

			if setErr := d.Set("frontend_port", flattenApplicationGatewayFrontendPorts(props.FrontendPorts)); setErr != nil {
				return fmt.Errorf("setting `frontend_port`: %+v", setErr)
			}

			frontendIPConfigurations, err := flattenApplicationGatewayFrontendIPConfigurations(props.FrontendIPConfigurations)
			if err != nil {
				return fmt.Errorf("flattening `frontend IP configuration`: %+v", err)
			}
			if setErr := d.Set("frontend_ip_configuration", frontendIPConfigurations); setErr != nil {
				return fmt.Errorf("setting `frontend_ip_configuration`: %+v", setErr)
			}

			if setErr := d.Set("gateway_ip_configuration", flattenApplicationGatewayIPConfigurations(props.GatewayIPConfigurations)); setErr != nil {
				return fmt.Errorf("setting `gateway_ip_configuration`: %+v", setErr)
			}

			if setErr := d.Set("global", flattenApplicationGatewayGlobalConfiguration(props.GlobalConfiguration)); setErr != nil {
				return fmt.Errorf("setting `global`: %+v", setErr)
			}

			if setErr := d.Set("private_endpoint_connection", flattenApplicationGatewayPrivateEndpoints(props.PrivateEndpointConnections)); setErr != nil {
				return fmt.Errorf("setting `private_endpoint_connection`: %+v", setErr)
			}

			if setErr := d.Set("private_link_configuration", flattenApplicationGatewayPrivateLinkConfigurations(props.PrivateLinkConfigurations)); setErr != nil {
				return fmt.Errorf("setting `private_link_configuration`: %+v", setErr)
			}

			if setErr := d.Set("probe", flattenApplicationGatewayProbes(props.Probes)); setErr != nil {
				return fmt.Errorf("setting `probe`: %+v", setErr)
			}

			requestRoutingRules, err := flattenApplicationGatewayRequestRoutingRules(props.RequestRoutingRules)
			if err != nil {
				return fmt.Errorf("flattening `request_routing_rule`: %+v", err)
			}
			if setErr := d.Set("request_routing_rule", requestRoutingRules); setErr != nil {
				return fmt.Errorf("setting `request_routing_rule`: %+v", setErr)
			}

			redirectConfigurations, err := flattenApplicationGatewayRedirectConfigurations(props.RedirectConfigurations)
			if err != nil {
				return fmt.Errorf("flattening `redirect configuration`: %+v", err)
			}
			if setErr := d.Set("redirect_configuration", redirectConfigurations); setErr != nil {
				return fmt.Errorf("setting `redirect_configuration`: %+v", setErr)
			}

			rewriteRuleSets := flattenApplicationGatewayRewriteRuleSets(props.RewriteRuleSets)
			if setErr := d.Set("rewrite_rule_set", rewriteRuleSets); setErr != nil {
				return fmt.Errorf("setting `rewrite_rule_set`: %+v", setErr)
			}

			if setErr := d.Set("sku", flattenApplicationGatewaySku(props.Sku)); setErr != nil {
				return fmt.Errorf("setting `sku`: %+v", setErr)
			}

			if setErr := d.Set("autoscale_configuration", flattenApplicationGatewayAutoscaleConfiguration(props.AutoscaleConfiguration)); setErr != nil {
				return fmt.Errorf("setting `autoscale_configuration`: %+v", setErr)
			}

			if setErr := d.Set("ssl_certificate", flattenApplicationGatewaySslCertificates(props.SslCertificates, d)); setErr != nil {
				return fmt.Errorf("setting `ssl_certificate`: %+v", setErr)
			}

			if setErr := d.Set("trusted_client_certificate", flattenApplicationGatewayTrustedClientCertificates(props.TrustedClientCertificates)); setErr != nil {
				return fmt.Errorf("setting `trusted_client_certificate`: %+v", setErr)
			}

			sslProfiles, err := flattenApplicationGatewaySslProfiles(props.SslProfiles)
			if err != nil {
				return fmt.Errorf("flattening `ssl_profile`: %+v", err)
			}
			if setErr := d.Set("ssl_profile", sslProfiles); setErr != nil {
				return fmt.Errorf("setting `ssl_profile`: %+v", setErr)
			}

			if setErr := d.Set("custom_error_configuration", flattenApplicationGatewayCustomErrorConfigurations(props.CustomErrorConfigurations)); setErr != nil {
				return fmt.Errorf("setting `custom_error_configuration`: %+v", setErr)
			}

			urlPathMaps, err := flattenApplicationGatewayURLPathMaps(props.UrlPathMaps)
			if err != nil {
				return fmt.Errorf("flattening `url_path_map`: %+v", err)
			}
			if setErr := d.Set("url_path_map", urlPathMaps); setErr != nil {
				return fmt.Errorf("setting `url_path_map`: %+v", setErr)
			}

			if setErr := d.Set("waf_configuration", flattenApplicationGatewayWafConfig(props.WebApplicationFirewallConfiguration)); setErr != nil {
				return fmt.Errorf("setting `waf_configuration`: %+v", setErr)
			}

			firewallPolicyId := ""
			if props.FirewallPolicy != nil && props.FirewallPolicy.Id != nil {
				firewallPolicyId = *props.FirewallPolicy.Id
				policyId, err := webapplicationfirewallpolicies.ParseApplicationGatewayWebApplicationFirewallPolicyIDInsensitively(firewallPolicyId)
				if err == nil {
					firewallPolicyId = policyId.ID()
				}
			}
			d.Set("firewall_policy_id", firewallPolicyId)
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceApplicationGatewayDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ApplicationGateways
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := applicationgateways.ParseApplicationGatewayID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandApplicationGatewayAuthenticationCertificates(certs []interface{}) *[]applicationgateways.ApplicationGatewayAuthenticationCertificate {
	results := make([]applicationgateways.ApplicationGatewayAuthenticationCertificate, 0)

	for _, raw := range certs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		data := v["data"].(string)

		// data must be base64 encoded
		encodedData := utils.Base64EncodeIfNot(data)

		output := applicationgateways.ApplicationGatewayAuthenticationCertificate{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayAuthenticationCertificatePropertiesFormat{
				Data: pointer.To(encodedData),
			},
		}

		results = append(results, output)
	}

	return &results
}

func expandApplicationGatewayTrustedRootCertificates(certs []interface{}) (*[]applicationgateways.ApplicationGatewayTrustedRootCertificate, error) {
	results := make([]applicationgateways.ApplicationGatewayTrustedRootCertificate, 0)

	for _, raw := range certs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		data := v["data"].(string)
		kvsid := v["key_vault_secret_id"].(string)

		output := applicationgateways.ApplicationGatewayTrustedRootCertificate{
			Name:       pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayTrustedRootCertificatePropertiesFormat{},
		}

		switch {
		case data != "" && kvsid != "":
			return nil, fmt.Errorf("only one of `key_vault_secret_id` or `data` must be specified for the `trusted_root_certificate` block %q", name)
		case data != "":
			output.Properties.Data = pointer.To(utils.Base64EncodeIfNot(data))
		case kvsid != "":
			output.Properties.KeyVaultSecretId = pointer.To(kvsid)
		default:
			return nil, fmt.Errorf("either `key_vault_secret_id` or `data` must be specified for the `trusted_root_certificate` block %q", name)
		}

		results = append(results, output)
	}

	return &results, nil
}

func flattenApplicationGatewayAuthenticationCertificates(certs *[]applicationgateways.ApplicationGatewayAuthenticationCertificate, d *pluginsdk.ResourceData) []interface{} {
	results := make([]interface{}, 0)
	if certs == nil {
		return results
	}

	// since the certificate data isn't returned lets load any existing data
	nameToDataMap := map[string]string{}
	if existing, ok := d.GetOk("authentication_certificate"); ok && existing != nil {
		for _, c := range existing.([]interface{}) {
			b := c.(map[string]interface{})
			nameToDataMap[b["name"].(string)] = b["data"].(string)
		}
	}

	for _, cert := range *certs {
		output := map[string]interface{}{}

		if v := cert.Id; v != nil {
			output["id"] = *v
		}

		if v := cert.Name; v != nil {
			output["name"] = *v

			// we have a name, so try and look up the old data to pass it along
			if data, ok := nameToDataMap[*v]; ok && data != "" {
				output["data"] = data
			}
		}

		results = append(results, output)
	}

	return results
}

func flattenApplicationGatewayTrustedRootCertificates(certs *[]applicationgateways.ApplicationGatewayTrustedRootCertificate, d *pluginsdk.ResourceData) []interface{} {
	results := make([]interface{}, 0)
	if certs == nil {
		return results
	}

	// since the certificate data isn't returned lets load any existing data
	nameToDataMap := map[string]string{}
	if existing, ok := d.GetOk("trusted_root_certificate"); ok && existing != nil {
		for _, c := range existing.([]interface{}) {
			b := c.(map[string]interface{})
			nameToDataMap[b["name"].(string)] = b["data"].(string)
		}
	}

	for _, cert := range *certs {
		output := map[string]interface{}{}

		if v := cert.Id; v != nil {
			output["id"] = *v
		}

		kvsid := ""
		if props := cert.Properties; props != nil {
			if v := props.KeyVaultSecretId; v != nil {
				kvsid = *v
			}
		}
		output["key_vault_secret_id"] = kvsid

		if v := cert.Name; v != nil {
			output["name"] = *v

			// if theres no key vauld ID and we have a name, so try and look up the old data to pass it along
			if data, ok := nameToDataMap[*v]; ok && data != "" {
				output["data"] = data
			}
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewayBackendAddressPools(d *pluginsdk.ResourceData) *[]applicationgateways.ApplicationGatewayBackendAddressPool {
	vs := d.Get("backend_address_pool").(*schema.Set).List()
	results := make([]applicationgateways.ApplicationGatewayBackendAddressPool, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})
		backendAddresses := make([]applicationgateways.ApplicationGatewayBackendAddress, 0)

		if fqdnsConfig, ok := v["fqdns"]; ok {
			fqdns := fqdnsConfig.(*schema.Set).List()
			for _, ip := range fqdns {
				backendAddresses = append(backendAddresses, applicationgateways.ApplicationGatewayBackendAddress{
					Fqdn: pointer.To(ip.(string)),
				})
			}
		}

		if ipAddressesConfig, ok := v["ip_addresses"]; ok {
			ipAddresses := ipAddressesConfig.(*schema.Set).List()

			for _, ip := range ipAddresses {
				backendAddresses = append(backendAddresses, applicationgateways.ApplicationGatewayBackendAddress{
					IPAddress: pointer.To(ip.(string)),
				})
			}
		}

		name := v["name"].(string)
		output := applicationgateways.ApplicationGatewayBackendAddressPool{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayBackendAddressPoolPropertiesFormat{
				BackendAddresses: &backendAddresses,
			},
		}

		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayBackendAddressPools(input *[]applicationgateways.ApplicationGatewayBackendAddressPool) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, config := range *input {
		ipAddressList := make([]interface{}, 0)
		fqdnList := make([]interface{}, 0)

		if props := config.Properties; props != nil {
			if props.BackendAddresses != nil {
				for _, address := range *props.BackendAddresses {
					if address.IPAddress != nil {
						ipAddressList = append(ipAddressList, *address.IPAddress)
					} else if address.Fqdn != nil {
						fqdnList = append(fqdnList, *address.Fqdn)
					}
				}
			}
		}

		output := map[string]interface{}{
			"fqdns":        fqdnList,
			"ip_addresses": ipAddressList,
		}

		if config.Id != nil {
			output["id"] = *config.Id
		}

		if config.Name != nil {
			output["name"] = *config.Name
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewayBackendHTTPSettings(d *pluginsdk.ResourceData, gatewayID string) *[]applicationgateways.ApplicationGatewayBackendHTTPSettings {
	results := make([]applicationgateways.ApplicationGatewayBackendHTTPSettings, 0)
	vs := d.Get("backend_http_settings").(*schema.Set).List()

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		path := v["path"].(string)
		port := int64(v["port"].(int))
		protocol := v["protocol"].(string)
		cookieBasedAffinity := v["cookie_based_affinity"].(string)
		pickHostNameFromBackendAddress := v["pick_host_name_from_backend_address"].(bool)
		requestTimeout := int64(v["request_timeout"].(int))

		setting := applicationgateways.ApplicationGatewayBackendHTTPSettings{
			Name: &name,
			Properties: &applicationgateways.ApplicationGatewayBackendHTTPSettingsPropertiesFormat{
				CookieBasedAffinity:            pointer.To(applicationgateways.ApplicationGatewayCookieBasedAffinity(cookieBasedAffinity)),
				Path:                           pointer.To(path),
				PickHostNameFromBackendAddress: pointer.To(pickHostNameFromBackendAddress),
				Port:                           pointer.To(port),
				Protocol:                       pointer.To(applicationgateways.ApplicationGatewayProtocol(protocol)),
				RequestTimeout:                 pointer.To(requestTimeout),
				ConnectionDraining:             expandApplicationGatewayConnectionDraining(v),
			},
		}

		hostName := v["host_name"].(string)
		if hostName != "" {
			setting.Properties.HostName = pointer.To(hostName)
		}

		affinityCookieName := v["affinity_cookie_name"].(string)
		if affinityCookieName != "" {
			setting.Properties.AffinityCookieName = pointer.To(affinityCookieName)
		}

		if v["authentication_certificate"] != nil {
			authCerts := v["authentication_certificate"].([]interface{})
			authCertSubResources := make([]applicationgateways.SubResource, 0)

			for _, rawAuthCert := range authCerts {
				authCert := rawAuthCert.(map[string]interface{})
				authCertName := authCert["name"].(string)
				authCertID := fmt.Sprintf("%s/authenticationCertificates/%s", gatewayID, authCertName)
				authCertSubResource := applicationgateways.SubResource{
					Id: pointer.To(authCertID),
				}

				authCertSubResources = append(authCertSubResources, authCertSubResource)
			}

			setting.Properties.AuthenticationCertificates = &authCertSubResources
		}

		if v["trusted_root_certificate_names"] != nil {
			trustedRootCertNames := v["trusted_root_certificate_names"].([]interface{})
			trustedRootCertSubResources := make([]applicationgateways.SubResource, 0)

			for _, rawTrustedRootCertName := range trustedRootCertNames {
				trustedRootCertName := rawTrustedRootCertName.(string)
				trustedRootCertID := fmt.Sprintf("%s/trustedRootCertificates/%s", gatewayID, trustedRootCertName)
				trustedRootCertSubResource := applicationgateways.SubResource{
					Id: pointer.To(trustedRootCertID),
				}

				trustedRootCertSubResources = append(trustedRootCertSubResources, trustedRootCertSubResource)
			}

			setting.Properties.TrustedRootCertificates = &trustedRootCertSubResources
		}

		probeName := v["probe_name"].(string)
		if probeName != "" {
			probeID := fmt.Sprintf("%s/probes/%s", gatewayID, probeName)
			setting.Properties.Probe = &applicationgateways.SubResource{
				Id: pointer.To(probeID),
			}
		}

		results = append(results, setting)
	}

	return &results
}

func flattenApplicationGatewayBackendHTTPSettings(input *[]applicationgateways.ApplicationGatewayBackendHTTPSettings) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.Id != nil {
			output["id"] = *v.Id
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.Properties; props != nil {
			output["cookie_based_affinity"] = props.CookieBasedAffinity

			if affinityCookieName := props.AffinityCookieName; affinityCookieName != nil {
				output["affinity_cookie_name"] = affinityCookieName
			}

			if path := props.Path; path != nil {
				output["path"] = *path
			}
			output["connection_draining"] = flattenApplicationGatewayConnectionDraining(props.ConnectionDraining)

			if port := props.Port; port != nil {
				output["port"] = int(*port)
			}

			if hostName := props.HostName; hostName != nil {
				output["host_name"] = *hostName
			}

			if pickHostNameFromBackendAddress := props.PickHostNameFromBackendAddress; pickHostNameFromBackendAddress != nil {
				output["pick_host_name_from_backend_address"] = *pickHostNameFromBackendAddress
			}

			output["protocol"] = props.Protocol

			if timeout := props.RequestTimeout; timeout != nil {
				output["request_timeout"] = int(*timeout)
			}

			authenticationCertificates := make([]interface{}, 0)
			if certs := props.AuthenticationCertificates; certs != nil {
				for _, cert := range *certs {
					if cert.Id == nil {
						continue
					}

					certId, err := parse.AuthenticationCertificateIDInsensitively(*cert.Id)
					if err != nil {
						return nil, err
					}

					certificate := map[string]interface{}{
						"id":   certId.ID(),
						"name": certId.Name,
					}
					authenticationCertificates = append(authenticationCertificates, certificate)
				}
			}
			output["authentication_certificate"] = authenticationCertificates

			trustedRootCertificateNames := make([]interface{}, 0)
			if certs := props.TrustedRootCertificates; certs != nil {
				for _, cert := range *certs {
					if cert.Id == nil {
						continue
					}

					certId, err := parse.TrustedRootCertificateIDInsensitively(*cert.Id)
					if err != nil {
						return nil, err
					}

					trustedRootCertificateNames = append(trustedRootCertificateNames, certId.Name)
				}
			}
			output["trusted_root_certificate_names"] = trustedRootCertificateNames

			if probe := props.Probe; probe != nil {
				if probe.Id != nil {
					id, err := parse.ProbeIDInsensitively(*probe.Id)
					if err != nil {
						return results, err
					}

					output["probe_name"] = id.Name
					output["probe_id"] = id.ID()
				}
			}
		}

		results = append(results, output)
	}

	return results, nil
}

func expandApplicationGatewayConnectionDraining(d map[string]interface{}) *applicationgateways.ApplicationGatewayConnectionDraining {
	connectionsRaw := d["connection_draining"].([]interface{})

	if len(connectionsRaw) == 0 {
		return nil
	}

	connectionRaw := connectionsRaw[0].(map[string]interface{})

	return &applicationgateways.ApplicationGatewayConnectionDraining{
		Enabled:           connectionRaw["enabled"].(bool),
		DrainTimeoutInSec: int64(connectionRaw["drain_timeout_sec"].(int)),
	}
}

func flattenApplicationGatewayConnectionDraining(input *applicationgateways.ApplicationGatewayConnectionDraining) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{map[string]interface{}{
		"enabled":           input.Enabled,
		"drain_timeout_sec": input.DrainTimeoutInSec,
	}}
}

func expandApplicationGatewaySslPolicy(vs []interface{}) *applicationgateways.ApplicationGatewaySslPolicy {
	policy := applicationgateways.ApplicationGatewaySslPolicy{}
	disabledSSLProtocols := make([]applicationgateways.ApplicationGatewaySslProtocol, 0)

	if len(vs) > 0 && vs[0] != nil {
		v := vs[0].(map[string]interface{})
		policyType := applicationgateways.ApplicationGatewaySslPolicyType(v["policy_type"].(string))

		for _, policy := range v["disabled_protocols"].([]interface{}) {
			disabledSSLProtocols = append(disabledSSLProtocols, applicationgateways.ApplicationGatewaySslProtocol(policy.(string)))
		}

		if policyType == applicationgateways.ApplicationGatewaySslPolicyTypePredefined {
			policyName := applicationgateways.ApplicationGatewaySslPolicyName(v["policy_name"].(string))
			policy = applicationgateways.ApplicationGatewaySslPolicy{
				PolicyType: pointer.To(policyType),
				PolicyName: pointer.To(policyName),
			}
		} else if policyType == applicationgateways.ApplicationGatewaySslPolicyTypeCustom || policyType == applicationgateways.ApplicationGatewaySslPolicyTypeCustomVTwo {
			minProtocolVersion := applicationgateways.ApplicationGatewaySslProtocol(v["min_protocol_version"].(string))
			cipherSuites := make([]applicationgateways.ApplicationGatewaySslCipherSuite, 0)

			for _, cipherSuite := range v["cipher_suites"].([]interface{}) {
				cipherSuites = append(cipherSuites, applicationgateways.ApplicationGatewaySslCipherSuite(cipherSuite.(string)))
			}

			policy = applicationgateways.ApplicationGatewaySslPolicy{
				PolicyType:         pointer.To(policyType),
				MinProtocolVersion: pointer.To(minProtocolVersion),
				CipherSuites:       &cipherSuites,
			}
		}
	}

	if len(disabledSSLProtocols) > 0 {
		policy = applicationgateways.ApplicationGatewaySslPolicy{
			DisabledSslProtocols: &disabledSSLProtocols,
		}
	}

	return &policy
}

func flattenApplicationGatewaySslPolicy(input *applicationgateways.ApplicationGatewaySslPolicy) []interface{} {
	results := make([]interface{}, 0)

	if input == nil {
		return results
	}

	output := map[string]interface{}{}
	output["policy_name"] = input.PolicyName
	output["policy_type"] = input.PolicyType
	output["min_protocol_version"] = input.MinProtocolVersion

	cipherSuites := make([]interface{}, 0)
	if input.CipherSuites != nil {
		for _, v := range *input.CipherSuites {
			cipherSuites = append(cipherSuites, string(v))
		}
	}
	output["cipher_suites"] = cipherSuites

	disabledSslProtocols := make([]interface{}, 0)
	if input.DisabledSslProtocols != nil {
		for _, v := range *input.DisabledSslProtocols {
			disabledSslProtocols = append(disabledSslProtocols, string(v))
		}
	}
	output["disabled_protocols"] = disabledSslProtocols

	results = append(results, output)
	return results
}

func expandApplicationGatewayHTTPListeners(d *pluginsdk.ResourceData, gatewayID string) (*[]applicationgateways.ApplicationGatewayHTTPListener, error) {
	vs := d.Get("http_listener").(*schema.Set).List()

	results := make([]applicationgateways.ApplicationGatewayHTTPListener, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		frontendIPConfigName := v["frontend_ip_configuration_name"].(string)
		frontendPortName := v["frontend_port_name"].(string)
		protocol := v["protocol"].(string)
		requireSNI := v["require_sni"].(bool)
		sslProfileName := v["ssl_profile_name"].(string)

		frontendIPConfigID := fmt.Sprintf("%s/frontendIPConfigurations/%s", gatewayID, frontendIPConfigName)
		frontendPortID := fmt.Sprintf("%s/frontendPorts/%s", gatewayID, frontendPortName)
		firewallPolicyID := v["firewall_policy_id"].(string)

		customErrorConfigurations := expandApplicationGatewayCustomErrorConfigurations(v["custom_error_configuration"].([]interface{}))

		listener := applicationgateways.ApplicationGatewayHTTPListener{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayHTTPListenerPropertiesFormat{
				FrontendIPConfiguration: &applicationgateways.SubResource{
					Id: pointer.To(frontendIPConfigID),
				},
				FrontendPort: &applicationgateways.SubResource{
					Id: pointer.To(frontendPortID),
				},
				Protocol:                    pointer.To(applicationgateways.ApplicationGatewayProtocol(protocol)),
				RequireServerNameIndication: pointer.To(requireSNI),
				CustomErrorConfigurations:   customErrorConfigurations,
			},
		}

		host := v["host_name"].(string)
		hosts := v["host_names"].(*pluginsdk.Set).List()

		if host != "" && len(hosts) > 0 {
			return nil, fmt.Errorf("`host_name` and `host_names` cannot be specified together")
		}

		if host != "" {
			listener.Properties.HostName = &host
		}

		if len(hosts) > 0 {
			listener.Properties.HostNames = utils.ExpandStringSlice(hosts)
		}

		if sslCertName := v["ssl_certificate_name"].(string); sslCertName != "" {
			certID := fmt.Sprintf("%s/sslCertificates/%s", gatewayID, sslCertName)
			listener.Properties.SslCertificate = &applicationgateways.SubResource{
				Id: pointer.To(certID),
			}
		}

		if firewallPolicyID != "" && len(firewallPolicyID) > 0 {
			listener.Properties.FirewallPolicy = &applicationgateways.SubResource{
				Id: pointer.To(firewallPolicyID),
			}
		}

		if sslProfileName != "" && len(sslProfileName) > 0 {
			sslProfileID := fmt.Sprintf("%s/sslProfiles/%s", gatewayID, sslProfileName)
			listener.Properties.SslProfile = &applicationgateways.SubResource{
				Id: pointer.To(sslProfileID),
			}
		}

		results = append(results, listener)
	}

	return &results, nil
}

func flattenApplicationGatewayHTTPListeners(input *[]applicationgateways.ApplicationGatewayHTTPListener) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.Id != nil {
			output["id"] = *v.Id
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.Properties; props != nil {
			if port := props.FrontendPort; port != nil {
				if port.Id != nil {
					portId, err := parse.FrontendPortIDInsensitively(*port.Id)
					if err != nil {
						return nil, err
					}
					output["frontend_port_name"] = portId.Name
					output["frontend_port_id"] = portId.ID()
				}
			}

			if feConfig := props.FrontendIPConfiguration; feConfig != nil {
				if feConfig.Id != nil {
					feConfigId, err := parse.FrontendIPConfigurationIDInsensitively(*feConfig.Id)
					if err != nil {
						return nil, err
					}
					output["frontend_ip_configuration_name"] = feConfigId.Name
					output["frontend_ip_configuration_id"] = feConfigId.ID()
				}
			}

			if hostname := props.HostName; hostname != nil {
				output["host_name"] = *hostname
			}

			if hostnames := props.HostNames; hostnames != nil {
				output["host_names"] = utils.FlattenStringSlice(hostnames)
			}

			output["protocol"] = props.Protocol

			if cert := props.SslCertificate; cert != nil {
				if cert.Id != nil {
					certId, err := parse.SslCertificateIDInsensitively(*cert.Id)
					if err != nil {
						return nil, err
					}

					output["ssl_certificate_name"] = certId.Name
					output["ssl_certificate_id"] = certId.ID()
				}
			}

			if sni := props.RequireServerNameIndication; sni != nil {
				output["require_sni"] = *sni
			}

			if fwp := props.FirewallPolicy; fwp != nil && fwp.Id != nil {
				policyId, err := webapplicationfirewallpolicies.ParseApplicationGatewayWebApplicationFirewallPolicyIDInsensitively(*fwp.Id)
				if err != nil {
					return nil, err
				}
				output["firewall_policy_id"] = policyId.ID()
			}

			if sslp := props.SslProfile; sslp != nil {
				if sslp.Id != nil {
					sslProfileId, err := parse.SslProfileIDInsensitively(*sslp.Id)
					if err != nil {
						return nil, err
					}

					output["ssl_profile_name"] = sslProfileId.Name
					output["ssl_profile_id"] = sslProfileId.ID()
				}
			}

			output["custom_error_configuration"] = flattenApplicationGatewayCustomErrorConfigurations(props.CustomErrorConfigurations)
		}

		results = append(results, output)
	}

	return results, nil
}

func expandApplicationGatewayIPConfigurations(d *pluginsdk.ResourceData) (*[]applicationgateways.ApplicationGatewayIPConfiguration, bool) {
	vs := d.Get("gateway_ip_configuration").([]interface{})
	results := make([]applicationgateways.ApplicationGatewayIPConfiguration, 0)
	stopApplicationGateway := false

	for _, configRaw := range vs {
		data := configRaw.(map[string]interface{})

		name := data["name"].(string)
		subnetID := data["subnet_id"].(string)

		output := applicationgateways.ApplicationGatewayIPConfiguration{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayIPConfigurationPropertiesFormat{
				Subnet: &applicationgateways.SubResource{
					Id: pointer.To(subnetID),
				},
			},
		}
		results = append(results, output)
	}

	if d.HasChange("gateway_ip_configuration") {
		oldRaw, newRaw := d.GetChange("gateway_ip_configuration")
		oldVS := oldRaw.([]interface{})
		newVS := newRaw.([]interface{})

		// If we're creating the application gateway return the current gateway ip configuration.
		if len(oldVS) == 0 {
			return &results, false
		}

		// The application gateway needs to be stopped if a gateway ip configuration is added or removed
		if len(oldVS) != len(newVS) {
			return &results, true
		}

		for i, configRaw := range newVS {
			newData := configRaw.(map[string]interface{})
			oldData := oldVS[i].(map[string]interface{})

			newSubnetID := newData["subnet_id"].(string)
			oldSubnetID := oldData["subnet_id"].(string)
			// The application gateway needs to be shutdown if the subnet ids don't match
			if newSubnetID != oldSubnetID {
				stopApplicationGateway = true
			}
		}
	}

	return &results, stopApplicationGateway
}

func flattenApplicationGatewayIPConfigurations(input *[]applicationgateways.ApplicationGatewayIPConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.Id != nil {
			output["id"] = *v.Id
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.Properties; props != nil {
			if subnet := props.Subnet; subnet != nil {
				if subnet.Id != nil {
					output["subnet_id"] = *subnet.Id
				}
			}
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewayGlobalConfiguration(input []interface{}) *applicationgateways.ApplicationGatewayGlobalConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &applicationgateways.ApplicationGatewayGlobalConfiguration{
		EnableRequestBuffering:  pointer.To(v["request_buffering_enabled"].(bool)),
		EnableResponseBuffering: pointer.To(v["response_buffering_enabled"].(bool)),
	}
}

func flattenApplicationGatewayGlobalConfiguration(input *applicationgateways.ApplicationGatewayGlobalConfiguration) []interface{} {
	if input == nil {
		return nil
	}

	output := make(map[string]interface{})

	if input.EnableRequestBuffering != nil {
		output["request_buffering_enabled"] = *input.EnableRequestBuffering
	}

	if input.EnableResponseBuffering != nil {
		output["response_buffering_enabled"] = *input.EnableResponseBuffering
	}

	return []interface{}{output}
}

func expandApplicationGatewayFrontendPorts(d *pluginsdk.ResourceData) *[]applicationgateways.ApplicationGatewayFrontendPort {
	vs := d.Get("frontend_port").(*pluginsdk.Set).List()
	results := make([]applicationgateways.ApplicationGatewayFrontendPort, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		port := int64(v["port"].(int))

		output := applicationgateways.ApplicationGatewayFrontendPort{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayFrontendPortPropertiesFormat{
				Port: pointer.To(port),
			},
		}
		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayFrontendPorts(input *[]applicationgateways.ApplicationGatewayFrontendPort) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.Id != nil {
			output["id"] = *v.Id
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.Properties; props != nil {
			if props.Port != nil {
				output["port"] = int(*props.Port)
			}
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewayFrontendIPConfigurations(d *pluginsdk.ResourceData, gatewayID string) *[]applicationgateways.ApplicationGatewayFrontendIPConfiguration {
	vs := d.Get("frontend_ip_configuration").([]interface{})
	results := make([]applicationgateways.ApplicationGatewayFrontendIPConfiguration, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		properties := applicationgateways.ApplicationGatewayFrontendIPConfigurationPropertiesFormat{}

		if val := v["subnet_id"].(string); val != "" {
			properties.Subnet = &applicationgateways.SubResource{
				Id: pointer.To(val),
			}
		}

		if val := v["private_ip_address_allocation"].(string); val != "" {
			properties.PrivateIPAllocationMethod = pointer.To(applicationgateways.IPAllocationMethod(val))
		}

		if val := v["private_ip_address"].(string); val != "" {
			properties.PrivateIPAddress = pointer.To(val)
		}

		if val := v["public_ip_address_id"].(string); val != "" {
			properties.PublicIPAddress = &applicationgateways.SubResource{
				Id: pointer.To(val),
			}
		}

		if val := v["private_link_configuration_name"].(string); val != "" {
			privateLinkConfigurationID := fmt.Sprintf("%s/privateLinkConfigurations/%s", gatewayID, val)
			properties.PrivateLinkConfiguration = &applicationgateways.SubResource{
				Id: pointer.To(privateLinkConfigurationID),
			}
		}

		name := v["name"].(string)
		output := applicationgateways.ApplicationGatewayFrontendIPConfiguration{
			Name:       pointer.To(name),
			Properties: &properties,
		}

		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayFrontendIPConfigurations(input *[]applicationgateways.ApplicationGatewayFrontendIPConfiguration) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, config := range *input {
		output := make(map[string]interface{})
		if config.Id != nil {
			output["id"] = *config.Id
		}

		if config.Name != nil {
			output["name"] = *config.Name
		}

		if props := config.Properties; props != nil {
			output["private_ip_address_allocation"] = string(pointer.From(props.PrivateIPAllocationMethod))

			if props.Subnet != nil && props.Subnet.Id != nil {
				output["subnet_id"] = *props.Subnet.Id
			}

			if props.PrivateIPAddress != nil {
				output["private_ip_address"] = *props.PrivateIPAddress
			}

			if props.PublicIPAddress != nil && props.PublicIPAddress.Id != nil {
				output["public_ip_address_id"] = *props.PublicIPAddress.Id
			}

			if props.PrivateLinkConfiguration != nil && props.PrivateLinkConfiguration.Id != nil {
				configurationID, err := parse.ApplicationGatewayPrivateLinkConfigurationIDInsensitively(*props.PrivateLinkConfiguration.Id)
				if err != nil {
					return nil, err
				}
				output["private_link_configuration_name"] = configurationID.PrivateLinkConfigurationName
				output["private_link_configuration_id"] = *props.PrivateLinkConfiguration.Id
			}
		}

		results = append(results, output)
	}

	return results, nil
}

func expandApplicationGatewayProbes(d *pluginsdk.ResourceData) *[]applicationgateways.ApplicationGatewayProbe {
	vs := d.Get("probe").(*schema.Set).List()
	results := make([]applicationgateways.ApplicationGatewayProbe, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		host := v["host"].(string)
		interval := int64(v["interval"].(int))
		minServers := int64(v["minimum_servers"].(int))
		name := v["name"].(string)
		probePath := v["path"].(string)
		protocol := v["protocol"].(string)
		port := int64(v["port"].(int))
		timeout := int64(v["timeout"].(int))
		unhealthyThreshold := int64(v["unhealthy_threshold"].(int))
		pickHostNameFromBackendHTTPSettings := v["pick_host_name_from_backend_http_settings"].(bool)

		output := applicationgateways.ApplicationGatewayProbe{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayProbePropertiesFormat{
				Host:                                pointer.To(host),
				Interval:                            pointer.To(interval),
				MinServers:                          pointer.To(minServers),
				Path:                                pointer.To(probePath),
				Protocol:                            pointer.To(applicationgateways.ApplicationGatewayProtocol(protocol)),
				Timeout:                             pointer.To(timeout),
				UnhealthyThreshold:                  pointer.To(unhealthyThreshold),
				PickHostNameFromBackendHTTPSettings: pointer.To(pickHostNameFromBackendHTTPSettings),
			},
		}

		matchConfigs := v["match"].([]interface{})
		if len(matchConfigs) > 0 {
			matchBody := ""
			outputMatch := &applicationgateways.ApplicationGatewayProbeHealthResponseMatch{}
			if matchConfigs[0] != nil {
				match := matchConfigs[0].(map[string]interface{})
				matchBody = match["body"].(string)

				statusCodes := make([]string, 0)
				for _, statusCode := range match["status_code"].([]interface{}) {
					statusCodes = append(statusCodes, statusCode.(string))
				}
				outputMatch.StatusCodes = &statusCodes
			}
			outputMatch.Body = pointer.To(matchBody)
			output.Properties.Match = outputMatch
		}

		if port != 0 {
			output.Properties.Port = pointer.To(port)
		}

		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayProbes(input *[]applicationgateways.ApplicationGatewayProbe) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.Id != nil {
			output["id"] = *v.Id
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.Properties; props != nil {
			output["protocol"] = string(pointer.From(props.Protocol))

			if host := props.Host; host != nil {
				output["host"] = *host
			}

			if path := props.Path; path != nil {
				output["path"] = *path
			}

			if interval := props.Interval; interval != nil {
				output["interval"] = int(*interval)
			}

			if timeout := props.Timeout; timeout != nil {
				output["timeout"] = int(*timeout)
			}

			if threshold := props.UnhealthyThreshold; threshold != nil {
				output["unhealthy_threshold"] = int(*threshold)
			}

			port := 0
			if props.Port != nil {
				port = int(*props.Port)
			}
			output["port"] = port

			if pickHostNameFromBackendHTTPSettings := props.PickHostNameFromBackendHTTPSettings; pickHostNameFromBackendHTTPSettings != nil {
				output["pick_host_name_from_backend_http_settings"] = *pickHostNameFromBackendHTTPSettings
			}

			if minServers := props.MinServers; minServers != nil {
				output["minimum_servers"] = int(*minServers)
			}

			matches := make([]interface{}, 0)
			if match := props.Match; match != nil {
				matchConfig := map[string]interface{}{}
				if body := match.Body; body != nil {
					matchConfig["body"] = *body
				}

				statusCodes := make([]interface{}, 0)
				if match.StatusCodes != nil {
					for _, status := range *match.StatusCodes {
						statusCodes = append(statusCodes, status)
					}
				}
				matchConfig["status_code"] = statusCodes
				matches = append(matches, matchConfig)
			}
			output["match"] = matches
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewayPrivateLinkConfigurations(d *pluginsdk.ResourceData) *[]applicationgateways.ApplicationGatewayPrivateLinkConfiguration {
	vs := d.Get("private_link_configuration").(*pluginsdk.Set).List()
	plConfigResults := make([]applicationgateways.ApplicationGatewayPrivateLinkConfiguration, 0)

	for _, rawPl := range vs {
		v := rawPl.(map[string]interface{})
		name := v["name"].(string)
		ipConfigurations := v["ip_configuration"].([]interface{})
		ipConfigurationResults := make([]applicationgateways.ApplicationGatewayPrivateLinkIPConfiguration, 0)
		for _, rawIp := range ipConfigurations {
			v := rawIp.(map[string]interface{})
			name := v["name"].(string)
			subnetId := v["subnet_id"].(string)
			primary := v["primary"].(bool)
			ipConfiguration := applicationgateways.ApplicationGatewayPrivateLinkIPConfiguration{
				Name: pointer.To(name),
				Properties: &applicationgateways.ApplicationGatewayPrivateLinkIPConfigurationProperties{
					Primary: &primary,
					Subnet: &applicationgateways.SubResource{
						Id: pointer.To(subnetId),
					},
				},
			}
			if privateIpAddress := v["private_ip_address"].(string); privateIpAddress != "" {
				ipConfiguration.Properties.PrivateIPAddress = pointer.To(privateIpAddress)
			}
			if privateIpAddressAllocation := v["private_ip_address_allocation"].(string); privateIpAddressAllocation != "" {
				ipConfiguration.Properties.PrivateIPAllocationMethod = pointer.To(applicationgateways.IPAllocationMethod(privateIpAddressAllocation))
			}
			ipConfigurationResults = append(ipConfigurationResults, ipConfiguration)
		}

		configuration := applicationgateways.ApplicationGatewayPrivateLinkConfiguration{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayPrivateLinkConfigurationProperties{
				IPConfigurations: &ipConfigurationResults,
			},
		}
		plConfigResults = append(plConfigResults, configuration)
	}

	return &plConfigResults
}

func flattenApplicationGatewayPrivateEndpoints(input *[]applicationgateways.ApplicationGatewayPrivateEndpointConnection) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, endpoint := range *input {
		result := map[string]interface{}{}
		if endpoint.Name != nil {
			result["name"] = *endpoint.Name
		}
		if endpoint.Id != nil {
			result["id"] = *endpoint.Id
		}
	}
	return results
}

func flattenApplicationGatewayPrivateLinkConfigurations(input *[]applicationgateways.ApplicationGatewayPrivateLinkConfiguration) []interface{} {
	plConfigResults := make([]interface{}, 0)
	if input == nil {
		return plConfigResults
	}

	for _, plConfig := range *input {
		plConfigResult := map[string]interface{}{}
		if plConfig.Name != nil {
			plConfigResult["name"] = *plConfig.Name
		}
		if plConfig.Id != nil {
			plConfigResult["id"] = *plConfig.Id
		}
		ipConfigResults := make([]interface{}, 0)
		if props := plConfig.Properties; props != nil {
			for _, ipConfig := range *props.IPConfigurations {
				ipConfigResult := map[string]interface{}{}
				if ipConfig.Name != nil {
					ipConfigResult["name"] = *ipConfig.Name
				}
				if ipConfigProps := ipConfig.Properties; ipConfigProps != nil {
					if ipConfigProps.Subnet != nil {
						ipConfigResult["subnet_id"] = *ipConfigProps.Subnet.Id
					}
					if ipConfigProps.PrivateIPAddress != nil {
						ipConfigResult["private_ip_address"] = *ipConfigProps.PrivateIPAddress
					}
					ipConfigResult["private_ip_address_allocation"] = string(pointer.From(ipConfigProps.PrivateIPAllocationMethod))
					if ipConfigProps.Primary != nil {
						ipConfigResult["primary"] = *ipConfigProps.Primary
					}
					ipConfigResults = append(ipConfigResults, ipConfigResult)
				}
			}
		}
		plConfigResult["ip_configuration"] = ipConfigResults
		plConfigResults = append(plConfigResults, plConfigResult)
	}
	return plConfigResults
}

func expandApplicationGatewayRequestRoutingRules(d *pluginsdk.ResourceData, gatewayID string) (*[]applicationgateways.ApplicationGatewayRequestRoutingRule, error) {
	vs := d.Get("request_routing_rule").(*pluginsdk.Set).List()
	results := make([]applicationgateways.ApplicationGatewayRequestRoutingRule, 0)
	priorityset := false

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		ruleType := v["rule_type"].(string)
		httpListenerName := v["http_listener_name"].(string)
		httpListenerID := fmt.Sprintf("%s/httpListeners/%s", gatewayID, httpListenerName)
		backendAddressPoolName := v["backend_address_pool_name"].(string)
		backendHTTPSettingsName := v["backend_http_settings_name"].(string)
		redirectConfigName := v["redirect_configuration_name"].(string)
		priority := int64(v["priority"].(int))

		rule := applicationgateways.ApplicationGatewayRequestRoutingRule{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayRequestRoutingRulePropertiesFormat{
				RuleType: pointer.To(applicationgateways.ApplicationGatewayRequestRoutingRuleType(ruleType)),
				HTTPListener: &applicationgateways.SubResource{
					Id: pointer.To(httpListenerID),
				},
			},
		}

		if backendAddressPoolName != "" && redirectConfigName != "" {
			return nil, fmt.Errorf("Conflict between `backend_address_pool_name` and `redirect_configuration_name` (back-end pool not applicable when redirection specified)")
		}

		if backendHTTPSettingsName != "" && redirectConfigName != "" {
			return nil, fmt.Errorf("Conflict between `backend_http_settings_name` and `redirect_configuration_name` (back-end settings not applicable when redirection specified)")
		}

		if backendAddressPoolName != "" {
			backendAddressPoolID := fmt.Sprintf("%s/backendAddressPools/%s", gatewayID, backendAddressPoolName)
			rule.Properties.BackendAddressPool = &applicationgateways.SubResource{
				Id: pointer.To(backendAddressPoolID),
			}
		}

		if backendHTTPSettingsName != "" {
			backendHTTPSettingsID := fmt.Sprintf("%s/backendHttpSettingsCollection/%s", gatewayID, backendHTTPSettingsName)
			rule.Properties.BackendHTTPSettings = &applicationgateways.SubResource{
				Id: pointer.To(backendHTTPSettingsID),
			}
		}

		if redirectConfigName != "" {
			redirectConfigID := fmt.Sprintf("%s/redirectConfigurations/%s", gatewayID, redirectConfigName)
			rule.Properties.RedirectConfiguration = &applicationgateways.SubResource{
				Id: pointer.To(redirectConfigID),
			}
		}

		if urlPathMapName := v["url_path_map_name"].(string); urlPathMapName != "" {
			urlPathMapID := fmt.Sprintf("%s/urlPathMaps/%s", gatewayID, urlPathMapName)
			rule.Properties.UrlPathMap = &applicationgateways.SubResource{
				Id: pointer.To(urlPathMapID),
			}
		}

		if rewriteRuleSetName := v["rewrite_rule_set_name"].(string); rewriteRuleSetName != "" {
			rewriteRuleSetID := fmt.Sprintf("%s/rewriteRuleSets/%s", gatewayID, rewriteRuleSetName)
			rule.Properties.RewriteRuleSet = &applicationgateways.SubResource{
				Id: pointer.To(rewriteRuleSetID),
			}
		}

		if priority != 0 {
			rule.Properties.Priority = &priority
			priorityset = true
		}

		results = append(results, rule)
	}

	if priorityset {
		for _, rule := range results {
			if rule.Properties.Priority == nil {
				return nil, fmt.Errorf("If you wish to use rule priority, you will have to specify rule-priority field values for all the existing request routing rules.")
			}
		}
	}

	return &results, nil
}

func flattenApplicationGatewayRequestRoutingRules(input *[]applicationgateways.ApplicationGatewayRequestRoutingRule) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, config := range *input {
		if props := config.Properties; props != nil {
			output := map[string]interface{}{
				"rule_type": string(pointer.From(props.RuleType)),
			}

			if config.Id != nil {
				output["id"] = *config.Id
			}

			if config.Name != nil {
				output["name"] = *config.Name
			}

			if props.Priority != nil {
				output["priority"] = *props.Priority
			}

			if pool := props.BackendAddressPool; pool != nil {
				if pool.Id != nil {
					poolId, err := parse.BackendAddressPoolIDInsensitively(*pool.Id)
					if err != nil {
						return nil, err
					}
					output["backend_address_pool_name"] = poolId.Name
					output["backend_address_pool_id"] = poolId.ID()
				}
			}

			if settings := props.BackendHTTPSettings; settings != nil {
				if settings.Id != nil {
					settingsId, err := parse.BackendHttpSettingsCollectionIDInsensitively(*settings.Id)
					if err != nil {
						return nil, err
					}

					output["backend_http_settings_name"] = settingsId.BackendHttpSettingsCollectionName
					output["backend_http_settings_id"] = *settings.Id
				}
			}

			if listener := props.HTTPListener; listener != nil {
				if listener.Id != nil {
					listenerId, err := parse.HttpListenerIDInsensitively(*listener.Id)
					if err != nil {
						return nil, err
					}
					output["http_listener_id"] = listenerId.ID()
					output["http_listener_name"] = listenerId.Name
				}
			}

			if pathMap := props.UrlPathMap; pathMap != nil {
				if pathMap.Id != nil {
					pathMapId, err := parse.UrlPathMapIDInsensitively(*pathMap.Id)
					if err != nil {
						return nil, err
					}
					output["url_path_map_name"] = pathMapId.Name
					output["url_path_map_id"] = pathMapId.ID()
				}
			}

			if redirect := props.RedirectConfiguration; redirect != nil {
				if redirect.Id != nil {
					redirectId, err := parse.RedirectConfigurationsIDInsensitively(*redirect.Id)
					if err != nil {
						return nil, err
					}
					output["redirect_configuration_name"] = redirectId.RedirectConfigurationName
					output["redirect_configuration_id"] = redirectId.ID()
				}
			}

			if rewrite := props.RewriteRuleSet; rewrite != nil {
				if rewrite.Id != nil {
					rewriteId, err := parse.RewriteRuleSetIDInsensitively(*rewrite.Id)
					if err != nil {
						return nil, err
					}
					output["rewrite_rule_set_name"] = rewriteId.Name
					output["rewrite_rule_set_id"] = rewriteId.ID()
				}
			}

			results = append(results, output)
		}
	}

	return results, nil
}

func expandApplicationGatewayRewriteRuleSets(d *pluginsdk.ResourceData) (*[]applicationgateways.ApplicationGatewayRewriteRuleSet, error) {
	vs := d.Get("rewrite_rule_set").([]interface{})
	ruleSets := make([]applicationgateways.ApplicationGatewayRewriteRuleSet, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})
		rules := make([]applicationgateways.ApplicationGatewayRewriteRule, 0)

		name := v["name"].(string)

		for _, ruleConfig := range v["rewrite_rule"].([]interface{}) {
			r := ruleConfig.(map[string]interface{})
			conditions := make([]applicationgateways.ApplicationGatewayRewriteRuleCondition, 0)
			requestConfigurations := make([]applicationgateways.ApplicationGatewayHeaderConfiguration, 0)
			responseConfigurations := make([]applicationgateways.ApplicationGatewayHeaderConfiguration, 0)
			urlConfiguration := applicationgateways.ApplicationGatewayURLConfiguration{}

			rule := applicationgateways.ApplicationGatewayRewriteRule{
				Name:         pointer.To(r["name"].(string)),
				RuleSequence: pointer.To(int64(r["rule_sequence"].(int))),
			}

			for _, rawCondition := range r["condition"].([]interface{}) {
				c := rawCondition.(map[string]interface{})
				condition := applicationgateways.ApplicationGatewayRewriteRuleCondition{
					Variable:   pointer.To(c["variable"].(string)),
					Pattern:    pointer.To(c["pattern"].(string)),
					IgnoreCase: pointer.To(c["ignore_case"].(bool)),
					Negate:     pointer.To(c["negate"].(bool)),
				}
				conditions = append(conditions, condition)
			}
			rule.Conditions = &conditions

			for _, rawConfig := range r["request_header_configuration"].([]interface{}) {
				c := rawConfig.(map[string]interface{})
				config := applicationgateways.ApplicationGatewayHeaderConfiguration{
					HeaderName:  pointer.To(c["header_name"].(string)),
					HeaderValue: pointer.To(c["header_value"].(string)),
				}
				requestConfigurations = append(requestConfigurations, config)
			}

			for _, rawConfig := range r["response_header_configuration"].([]interface{}) {
				c := rawConfig.(map[string]interface{})
				config := applicationgateways.ApplicationGatewayHeaderConfiguration{
					HeaderName:  pointer.To(c["header_name"].(string)),
					HeaderValue: pointer.To(c["header_value"].(string)),
				}
				responseConfigurations = append(responseConfigurations, config)
			}

			for _, rawConfig := range r["url"].([]interface{}) {
				c := rawConfig.(map[string]interface{})
				if c["path"] == nil && c["query_string"] == nil {
					return nil, fmt.Errorf("At least one of `path` or `query_string` must be set")
				}
				components := ""
				if c["components"] != nil {
					components = c["components"].(string)
				}
				if c["path"] != nil && components != "query_string_only" {
					urlConfiguration.ModifiedPath = pointer.To(c["path"].(string))
				}
				if c["query_string"] != nil && components != "path_only" {
					urlConfiguration.ModifiedQueryString = pointer.To(c["query_string"].(string))
				}
				if c["reroute"] != nil {
					urlConfiguration.Reroute = pointer.To(c["reroute"].(bool))
				}
			}

			rule.ActionSet = &applicationgateways.ApplicationGatewayRewriteRuleActionSet{
				RequestHeaderConfigurations:  &requestConfigurations,
				ResponseHeaderConfigurations: &responseConfigurations,
			}

			if len(r["url"].([]interface{})) > 0 {
				rule.ActionSet.UrlConfiguration = &urlConfiguration
			}

			rules = append(rules, rule)
		}

		ruleSet := applicationgateways.ApplicationGatewayRewriteRuleSet{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayRewriteRuleSetPropertiesFormat{
				RewriteRules: &rules,
			},
		}

		ruleSets = append(ruleSets, ruleSet)
	}

	return &ruleSets, nil
}

func flattenApplicationGatewayRewriteRuleSets(input *[]applicationgateways.ApplicationGatewayRewriteRuleSet) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, config := range *input {
		if props := config.Properties; props != nil {
			output := map[string]interface{}{}

			if config.Id != nil {
				output["id"] = *config.Id
			}

			if config.Name != nil {
				output["name"] = *config.Name
			}

			if rulesConfig := props.RewriteRules; rulesConfig != nil {
				rules := make([]interface{}, 0)
				for _, rule := range *rulesConfig {
					ruleOutput := map[string]interface{}{}

					if rule.Name != nil {
						ruleOutput["name"] = *rule.Name
					}

					if rule.RuleSequence != nil {
						ruleOutput["rule_sequence"] = *rule.RuleSequence
					}

					conditions := make([]interface{}, 0)
					if rule.Conditions != nil {
						for _, config := range *rule.Conditions {
							condition := map[string]interface{}{}

							if config.Variable != nil {
								condition["variable"] = *config.Variable
							}

							if config.Pattern != nil {
								condition["pattern"] = *config.Pattern
							}

							if config.IgnoreCase != nil {
								condition["ignore_case"] = *config.IgnoreCase
							}

							if config.Negate != nil {
								condition["negate"] = *config.Negate
							}

							conditions = append(conditions, condition)
						}
					}
					ruleOutput["condition"] = conditions

					requestConfigs := make([]interface{}, 0)
					responseConfigs := make([]interface{}, 0)
					urlConfigs := make([]interface{}, 0)

					if rule.ActionSet != nil {
						actionSet := *rule.ActionSet

						if actionSet.RequestHeaderConfigurations != nil {
							for _, config := range *actionSet.RequestHeaderConfigurations {
								requestConfig := map[string]interface{}{}

								if config.HeaderName != nil {
									requestConfig["header_name"] = *config.HeaderName
								}

								if config.HeaderValue != nil {
									requestConfig["header_value"] = *config.HeaderValue
								}

								requestConfigs = append(requestConfigs, requestConfig)
							}
						}

						if actionSet.ResponseHeaderConfigurations != nil {
							for _, config := range *actionSet.ResponseHeaderConfigurations {
								responseConfig := map[string]interface{}{}

								if config.HeaderName != nil {
									responseConfig["header_name"] = *config.HeaderName
								}

								if config.HeaderValue != nil {
									responseConfig["header_value"] = *config.HeaderValue
								}

								responseConfigs = append(responseConfigs, responseConfig)
							}
						}

						if actionSet.UrlConfiguration != nil {
							config := *actionSet.UrlConfiguration
							components := ""

							path := ""
							if config.ModifiedPath != nil {
								path = *config.ModifiedPath
							}

							queryString := ""
							if config.ModifiedQueryString != nil {
								queryString = *config.ModifiedQueryString
							}

							// `components` doesn't exist in the API - it appears to be purely a UI state in the Portal
							// as such we should consider removing this field in the future.
							if path == queryString {
								// used to represent `both`
								components = ""
							}
							if config.ModifiedQueryString != nil && config.ModifiedPath == nil {
								components = "query_string_only"
							}
							if config.ModifiedQueryString == nil && config.ModifiedPath != nil {
								components = "path_only"
							}

							reroute := false
							if config.Reroute != nil {
								reroute = *config.Reroute
							}

							urlConfigs = append(urlConfigs, map[string]interface{}{
								"components":   components,
								"query_string": queryString,
								"path":         path,
								"reroute":      reroute,
							})
						}
					}
					ruleOutput["request_header_configuration"] = requestConfigs
					ruleOutput["response_header_configuration"] = responseConfigs
					ruleOutput["url"] = urlConfigs

					rules = append(rules, ruleOutput)
				}
				output["rewrite_rule"] = rules
			}
			results = append(results, output)
		}
	}

	return results
}

func expandApplicationGatewayRedirectConfigurations(d *pluginsdk.ResourceData, gatewayID string) (*[]applicationgateways.ApplicationGatewayRedirectConfiguration, error) {
	vs := d.Get("redirect_configuration").(*pluginsdk.Set).List()
	results := make([]applicationgateways.ApplicationGatewayRedirectConfiguration, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		redirectType := v["redirect_type"].(string)
		targetListenerName := v["target_listener_name"].(string)
		targetUrl := v["target_url"].(string)
		includePath := v["include_path"].(bool)
		includeQueryString := v["include_query_string"].(bool)

		output := applicationgateways.ApplicationGatewayRedirectConfiguration{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayRedirectConfigurationPropertiesFormat{
				RedirectType:       pointer.To(applicationgateways.ApplicationGatewayRedirectType(redirectType)),
				IncludeQueryString: pointer.To(includeQueryString),
				IncludePath:        pointer.To(includePath),
			},
		}

		if targetListenerName == "" && targetUrl == "" {
			return nil, fmt.Errorf("Set either `target_listener_name` or `target_url`")
		}

		if targetListenerName != "" && targetUrl != "" {
			return nil, fmt.Errorf("Conflict between `target_listener_name` and `target_url` (redirection is either to URL or target listener)")
		}

		if targetListenerName != "" {
			targetListenerID := fmt.Sprintf("%s/httpListeners/%s", gatewayID, targetListenerName)
			output.Properties.TargetListener = &applicationgateways.SubResource{
				Id: pointer.To(targetListenerID),
			}
		}

		if targetUrl != "" {
			output.Properties.TargetURL = pointer.To(targetUrl)
		}

		results = append(results, output)
	}

	return &results, nil
}

func flattenApplicationGatewayRedirectConfigurations(input *[]applicationgateways.ApplicationGatewayRedirectConfiguration) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, config := range *input {
		if props := config.Properties; props != nil {
			output := map[string]interface{}{
				"redirect_type": string(pointer.From(props.RedirectType)),
			}

			if config.Id != nil {
				output["id"] = *config.Id
			}

			if config.Name != nil {
				output["name"] = *config.Name
			}

			if listener := props.TargetListener; listener != nil {
				if listener.Id != nil {
					listenerId, err := parse.HttpListenerIDInsensitively(*listener.Id)
					if err != nil {
						return nil, err
					}
					output["target_listener_name"] = listenerId.Name
					output["target_listener_id"] = listenerId.ID()
				}
			}

			if props.TargetURL != nil {
				output["target_url"] = *props.TargetURL
			}

			if props.IncludePath != nil {
				output["include_path"] = *props.IncludePath
			}

			if props.IncludeQueryString != nil {
				output["include_query_string"] = *props.IncludeQueryString
			}

			results = append(results, output)
		}
	}

	return results, nil
}

func expandApplicationGatewayAutoscaleConfiguration(d *pluginsdk.ResourceData) *applicationgateways.ApplicationGatewayAutoscaleConfiguration {
	vs := d.Get("autoscale_configuration").([]interface{})
	if len(vs) == 0 {
		return nil
	}
	v := vs[0].(map[string]interface{})

	minCapacity := int64(v["min_capacity"].(int))
	maxCapacity := int64(v["max_capacity"].(int))

	configuration := applicationgateways.ApplicationGatewayAutoscaleConfiguration{
		MinCapacity: minCapacity,
	}

	if maxCapacity != 0 {
		configuration.MaxCapacity = pointer.To(maxCapacity)
	}

	return &configuration
}

func flattenApplicationGatewayAutoscaleConfiguration(input *applicationgateways.ApplicationGatewayAutoscaleConfiguration) []interface{} {
	result := make(map[string]interface{})
	if input == nil {
		return []interface{}{}
	}

	result["min_capacity"] = input.MinCapacity
	if input.MaxCapacity != nil {
		result["max_capacity"] = *input.MaxCapacity
	}

	return []interface{}{result}
}

func expandApplicationGatewaySku(d *pluginsdk.ResourceData) *applicationgateways.ApplicationGatewaySku {
	vs := d.Get("sku").([]interface{})
	v := vs[0].(map[string]interface{})

	name := v["name"].(string)
	tier := v["tier"].(string)
	capacity := int64(v["capacity"].(int))

	sku := applicationgateways.ApplicationGatewaySku{
		Name: pointer.To(applicationgateways.ApplicationGatewaySkuName(name)),
		Tier: pointer.To(applicationgateways.ApplicationGatewayTier(tier)),
	}

	if capacity != 0 {
		sku.Capacity = pointer.To(capacity)
	}

	return &sku
}

func flattenApplicationGatewaySku(input *applicationgateways.ApplicationGatewaySku) []interface{} {
	result := make(map[string]interface{})

	result["name"] = string(pointer.From(input.Name))
	result["tier"] = string(pointer.From(input.Tier))
	if input.Capacity != nil {
		result["capacity"] = int(*input.Capacity)
	}

	return []interface{}{result}
}

func expandApplicationGatewaySslCertificates(d *pluginsdk.ResourceData) (*[]applicationgateways.ApplicationGatewaySslCertificate, error) {
	vs := d.Get("ssl_certificate").(*schema.Set).List()
	results := make([]applicationgateways.ApplicationGatewaySslCertificate, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		data := v["data"].(string)
		password := v["password"].(string)
		kvsid := v["key_vault_secret_id"].(string)
		cert := v["public_cert_data"].(string)

		output := applicationgateways.ApplicationGatewaySslCertificate{
			Name:       pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewaySslCertificatePropertiesFormat{},
		}

		// nolint gocritic
		if data != "" && kvsid != "" {
			return nil, fmt.Errorf("only one of `key_vault_secret_id` or `data` must be specified for the `ssl_certificate` block %q", name)
		} else if data != "" {
			// data must be base64 encoded
			output.Properties.Data = pointer.To(utils.Base64EncodeIfNot(data))

			output.Properties.Password = pointer.To(password)
		} else if kvsid != "" {
			if password != "" {
				return nil, fmt.Errorf("only one of `key_vault_secret_id` or `password` must be specified for the `ssl_certificate` block %q", name)
			}

			output.Properties.KeyVaultSecretId = pointer.To(kvsid)
		} else if cert != "" {
			output.Properties.PublicCertData = pointer.To(cert)
		} else {
			return nil, fmt.Errorf("either `key_vault_secret_id` or `data` must be specified for the `ssl_certificate` block %q", name)
		}

		results = append(results, output)
	}

	return &results, nil
}

func flattenApplicationGatewaySslCertificates(input *[]applicationgateways.ApplicationGatewaySslCertificate, d *pluginsdk.ResourceData) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]interface{}{}
		if v.Name == nil {
			continue
		}

		name := *v.Name

		if v.Id != nil {
			output["id"] = *v.Id
		}

		output["name"] = name

		if props := v.Properties; props != nil {
			if data := props.PublicCertData; data != nil {
				output["public_cert_data"] = *data
			}

			if kvsid := props.KeyVaultSecretId; kvsid != nil {
				output["key_vault_secret_id"] = *kvsid
			}
		}

		// since the certificate data isn't returned we have to load it from the same index
		if existing, ok := d.GetOk("ssl_certificate"); ok && existing != nil {
			existingVals := existing.(*schema.Set).List()

			for _, existingVal := range existingVals {
				existingCerts := existingVal.(map[string]interface{})
				existingName := existingCerts["name"].(string)

				if name == existingName {
					if data := existingCerts["data"]; data != nil {
						v := utils.Base64EncodeIfNot(data.(string))
						output["data"] = v
					}

					if password := existingCerts["password"]; password != nil {
						output["password"] = password.(string)
					}
				}
			}
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewayTrustedClientCertificates(d *pluginsdk.ResourceData) (*[]applicationgateways.ApplicationGatewayTrustedClientCertificate, error) {
	vs := d.Get("trusted_client_certificate").([]interface{})
	results := make([]applicationgateways.ApplicationGatewayTrustedClientCertificate, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		data := v["data"].(string)

		output := applicationgateways.ApplicationGatewayTrustedClientCertificate{
			Name:       pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayTrustedClientCertificatePropertiesFormat{},
		}

		// nolint gocritic
		if data != "" {
			// data must be base64 encoded
			output.Properties.Data = pointer.To(utils.Base64EncodeIfNot(data))
		} else {
			return nil, fmt.Errorf("`data` must be specified for the `trusted_client_certificate` block %q", name)
		}

		results = append(results, output)
	}

	return &results, nil
}

func flattenApplicationGatewayTrustedClientCertificates(input *[]applicationgateways.ApplicationGatewayTrustedClientCertificate) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]interface{}{}
		if v.Name == nil {
			continue
		}

		name := *v.Name

		if v.Id != nil {
			output["id"] = *v.Id
		}

		output["name"] = name

		if props := v.Properties; props != nil {
			if data := props.Data; data != nil {
				output["data"] = *data
			}
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewaySslProfiles(d *pluginsdk.ResourceData, gatewayID string) *[]applicationgateways.ApplicationGatewaySslProfile {
	vs := d.Get("ssl_profile").([]interface{})
	results := make([]applicationgateways.ApplicationGatewaySslProfile, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		verifyClientCertIssuerDn := v["verify_client_cert_issuer_dn"].(bool)
		verifyClientCertificateRevocation := applicationgateways.ApplicationGatewayClientRevocationOptionsNone
		if v["verify_client_certificate_revocation"].(string) != "" {
			verifyClientCertificateRevocation = applicationgateways.ApplicationGatewayClientRevocationOptions(v["verify_client_certificate_revocation"].(string))
		}

		output := applicationgateways.ApplicationGatewaySslProfile{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewaySslProfilePropertiesFormat{
				ClientAuthConfiguration: &applicationgateways.ApplicationGatewayClientAuthConfiguration{
					VerifyClientCertIssuerDN: pointer.To(verifyClientCertIssuerDn),
					VerifyClientRevocation:   pointer.To(verifyClientCertificateRevocation),
				},
			},
		}

		if v["trusted_client_certificate_names"] != nil {
			clientCerts := v["trusted_client_certificate_names"].([]interface{})
			clientCertSubResources := make([]applicationgateways.SubResource, 0)

			for _, rawClientCert := range clientCerts {
				clientCertName := rawClientCert
				clientCertID := fmt.Sprintf("%s/trustedClientCertificates/%s", gatewayID, clientCertName)
				clientCertSubResource := applicationgateways.SubResource{
					Id: pointer.To(clientCertID),
				}
				clientCertSubResources = append(clientCertSubResources, clientCertSubResource)
			}
			output.Properties.TrustedClientCertificates = &clientCertSubResources
		}

		sslPolicy := v["ssl_policy"].([]interface{})
		if len(sslPolicy) > 0 {
			output.Properties.SslPolicy = expandApplicationGatewaySslPolicy(sslPolicy)
		} else {
			output.Properties.SslPolicy = nil
		}
		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewaySslProfiles(input *[]applicationgateways.ApplicationGatewaySslProfile) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, v := range *input {
		output := map[string]interface{}{}
		if v.Name == nil {
			continue
		}

		name := *v.Name

		if v.Id != nil {
			output["id"] = *v.Id
		}

		output["name"] = name
		output["ssl_policy"] = flattenApplicationGatewaySslPolicy(v.Properties.SslPolicy)

		verifyClientCertIssuerDn := false
		verifyClientCertificateRevocation := ""

		if props := v.Properties; props != nil {
			if props.ClientAuthConfiguration != nil {
				verifyClientCertIssuerDn = pointer.From(props.ClientAuthConfiguration.VerifyClientCertIssuerDN)
				if *props.ClientAuthConfiguration.VerifyClientRevocation != applicationgateways.ApplicationGatewayClientRevocationOptionsNone {
					verifyClientCertificateRevocation = string(pointer.From(props.ClientAuthConfiguration.VerifyClientRevocation))
				}
			}

			trustedClientCertificateNames := make([]interface{}, 0)
			if certs := props.TrustedClientCertificates; certs != nil {
				for _, cert := range *certs {
					if cert.Id == nil {
						continue
					}

					certId, err := parse.TrustedClientCertificateIDInsensitively(*cert.Id)
					if err != nil {
						return nil, err
					}

					trustedClientCertificateNames = append(trustedClientCertificateNames, certId.Name)
				}
			}
			output["trusted_client_certificate_names"] = trustedClientCertificateNames
			output["verify_client_cert_issuer_dn"] = verifyClientCertIssuerDn
			output["verify_client_certificate_revocation"] = verifyClientCertificateRevocation
		}

		results = append(results, output)
	}

	return results, nil
}

func expandApplicationGatewayURLPathMaps(d *pluginsdk.ResourceData, gatewayID string) (*[]applicationgateways.ApplicationGatewayURLPathMap, error) {
	vs := d.Get("url_path_map").([]interface{})
	results := make([]applicationgateways.ApplicationGatewayURLPathMap, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)

		pathRules := make([]applicationgateways.ApplicationGatewayPathRule, 0)
		for _, ruleConfig := range v["path_rule"].([]interface{}) {
			ruleConfigMap := ruleConfig.(map[string]interface{})

			ruleName := ruleConfigMap["name"].(string)
			backendAddressPoolName := ruleConfigMap["backend_address_pool_name"].(string)
			backendHTTPSettingsName := ruleConfigMap["backend_http_settings_name"].(string)
			redirectConfigurationName := ruleConfigMap["redirect_configuration_name"].(string)
			firewallPolicyID := ruleConfigMap["firewall_policy_id"].(string)

			rulePaths := make([]string, 0)
			for _, rulePath := range ruleConfigMap["paths"].([]interface{}) {
				p, ok := rulePath.(string)
				if ok {
					rulePaths = append(rulePaths, p)
				}
			}

			rule := applicationgateways.ApplicationGatewayPathRule{
				Name: pointer.To(ruleName),
				Properties: &applicationgateways.ApplicationGatewayPathRulePropertiesFormat{
					Paths: &rulePaths,
				},
			}

			if backendAddressPoolName != "" && redirectConfigurationName != "" {
				return nil, fmt.Errorf("Conflict between `backend_address_pool_name` and `redirect_configuration_name` (back-end pool not applicable when redirection specified)")
			}

			if backendHTTPSettingsName != "" && redirectConfigurationName != "" {
				return nil, fmt.Errorf("Conflict between `backend_http_settings_name` and `redirect_configuration_name` (back-end settings not applicable when redirection specified)")
			}

			if backendAddressPoolName != "" {
				backendAddressPoolID := fmt.Sprintf("%s/backendAddressPools/%s", gatewayID, backendAddressPoolName)
				rule.Properties.BackendAddressPool = &applicationgateways.SubResource{
					Id: pointer.To(backendAddressPoolID),
				}
			}

			if backendHTTPSettingsName != "" {
				backendHTTPSettingsID := fmt.Sprintf("%s/backendHttpSettingsCollection/%s", gatewayID, backendHTTPSettingsName)
				rule.Properties.BackendHTTPSettings = &applicationgateways.SubResource{
					Id: pointer.To(backendHTTPSettingsID),
				}
			}

			if redirectConfigurationName != "" {
				redirectConfigurationID := fmt.Sprintf("%s/redirectConfigurations/%s", gatewayID, redirectConfigurationName)
				rule.Properties.RedirectConfiguration = &applicationgateways.SubResource{
					Id: pointer.To(redirectConfigurationID),
				}
			}

			if rewriteRuleSetName := ruleConfigMap["rewrite_rule_set_name"].(string); rewriteRuleSetName != "" {
				rewriteRuleSetID := fmt.Sprintf("%s/rewriteRuleSets/%s", gatewayID, rewriteRuleSetName)
				rule.Properties.RewriteRuleSet = &applicationgateways.SubResource{
					Id: pointer.To(rewriteRuleSetID),
				}
			}

			if firewallPolicyID != "" && len(firewallPolicyID) > 0 {
				rule.Properties.FirewallPolicy = &applicationgateways.SubResource{
					Id: pointer.To(firewallPolicyID),
				}
			}

			pathRules = append(pathRules, rule)
		}

		output := applicationgateways.ApplicationGatewayURLPathMap{
			Name: pointer.To(name),
			Properties: &applicationgateways.ApplicationGatewayURLPathMapPropertiesFormat{
				PathRules: &pathRules,
			},
		}

		defaultBackendAddressPoolName := v["default_backend_address_pool_name"].(string)
		defaultBackendHTTPSettingsName := v["default_backend_http_settings_name"].(string)
		defaultRedirectConfigurationName := v["default_redirect_configuration_name"].(string)

		if defaultBackendAddressPoolName == "" && defaultBackendHTTPSettingsName == "" && defaultRedirectConfigurationName == "" {
			return nil, fmt.Errorf("both the `default_backend_address_pool_name` and `default_backend_http_settings_name` or `default_redirect_configuration_name` must be specified")
		}

		if defaultBackendAddressPoolName != "" && defaultRedirectConfigurationName != "" {
			return nil, fmt.Errorf("Conflict between `default_backend_address_pool_name` and `default_redirect_configuration_name` (back-end pool not applicable when redirection specified)")
		}

		if defaultBackendHTTPSettingsName != "" && defaultRedirectConfigurationName != "" {
			return nil, fmt.Errorf("Conflict between `default_backend_http_settings_name` and `default_redirect_configuration_name` (back-end settings not applicable when redirection specified)")
		}

		if defaultBackendAddressPoolName != "" {
			defaultBackendAddressPoolID := fmt.Sprintf("%s/backendAddressPools/%s", gatewayID, defaultBackendAddressPoolName)
			output.Properties.DefaultBackendAddressPool = &applicationgateways.SubResource{
				Id: pointer.To(defaultBackendAddressPoolID),
			}
		}

		if defaultBackendHTTPSettingsName != "" {
			defaultBackendHTTPSettingsID := fmt.Sprintf("%s/backendHttpSettingsCollection/%s", gatewayID, defaultBackendHTTPSettingsName)
			output.Properties.DefaultBackendHTTPSettings = &applicationgateways.SubResource{
				Id: pointer.To(defaultBackendHTTPSettingsID),
			}
		}

		if defaultRedirectConfigurationName != "" {
			defaultRedirectConfigurationID := fmt.Sprintf("%s/redirectConfigurations/%s", gatewayID, defaultRedirectConfigurationName)
			output.Properties.DefaultRedirectConfiguration = &applicationgateways.SubResource{
				Id: pointer.To(defaultRedirectConfigurationID),
			}
		}

		if defaultRewriteRuleSetName := v["default_rewrite_rule_set_name"].(string); defaultRewriteRuleSetName != "" {
			defaultRewriteRuleSetID := fmt.Sprintf("%s/rewriteRuleSets/%s", gatewayID, defaultRewriteRuleSetName)
			output.Properties.DefaultRewriteRuleSet = &applicationgateways.SubResource{
				Id: pointer.To(defaultRewriteRuleSetID),
			}
		}

		results = append(results, output)
	}

	return &results, nil
}

func flattenApplicationGatewayURLPathMaps(input *[]applicationgateways.ApplicationGatewayURLPathMap) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.Id != nil {
			output["id"] = *v.Id
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.Properties; props != nil {
			if backendPool := props.DefaultBackendAddressPool; backendPool != nil && backendPool.Id != nil {
				poolId, err := parse.BackendAddressPoolIDInsensitively(*backendPool.Id)
				if err != nil {
					return nil, err
				}
				output["default_backend_address_pool_name"] = poolId.Name
				output["default_backend_address_pool_id"] = poolId.ID()
			}

			if settings := props.DefaultBackendHTTPSettings; settings != nil && settings.Id != nil {
				settingsId, err := parse.BackendHttpSettingsCollectionIDInsensitively(*settings.Id)
				if err != nil {
					return nil, err
				}
				output["default_backend_http_settings_name"] = settingsId.BackendHttpSettingsCollectionName
				output["default_backend_http_settings_id"] = settingsId.ID()
			}

			if redirect := props.DefaultRedirectConfiguration; redirect != nil && redirect.Id != nil {
				redirectId, err := parse.RedirectConfigurationsIDInsensitively(*redirect.Id)
				if err != nil {
					return nil, err
				}
				output["default_redirect_configuration_name"] = redirectId.RedirectConfigurationName
				output["default_redirect_configuration_id"] = redirectId.ID()
			}

			if rewrite := props.DefaultRewriteRuleSet; rewrite != nil && rewrite.Id != nil {
				rewriteId, err := parse.RewriteRuleSetIDInsensitively(*rewrite.Id)
				if err != nil {
					return nil, err
				}
				output["default_rewrite_rule_set_name"] = rewriteId.Name
				output["default_rewrite_rule_set_id"] = rewriteId.ID()
			}

			pathRules := make([]interface{}, 0)
			if rules := props.PathRules; rules != nil {
				for _, rule := range *rules {
					ruleOutput := map[string]interface{}{}

					if rule.Id != nil {
						ruleOutput["id"] = *rule.Id
					}

					if rule.Name != nil {
						ruleOutput["name"] = *rule.Name
					}

					if ruleProps := rule.Properties; ruleProps != nil {
						if pool := ruleProps.BackendAddressPool; pool != nil && pool.Id != nil {
							poolId, err := parse.BackendAddressPoolIDInsensitively(*pool.Id)
							if err != nil {
								return nil, err
							}
							ruleOutput["backend_address_pool_name"] = poolId.Name
							ruleOutput["backend_address_pool_id"] = poolId.ID()
						}

						if backend := ruleProps.BackendHTTPSettings; backend != nil && backend.Id != nil {
							backendId, err := parse.BackendHttpSettingsCollectionIDInsensitively(*backend.Id)
							if err != nil {
								return nil, err
							}
							ruleOutput["backend_http_settings_name"] = backendId.BackendHttpSettingsCollectionName
							ruleOutput["backend_http_settings_id"] = backendId.ID()
						}

						if redirect := ruleProps.RedirectConfiguration; redirect != nil && redirect.Id != nil {
							redirectId, err := parse.RedirectConfigurationsIDInsensitively(*redirect.Id)
							if err != nil {
								return nil, err
							}
							ruleOutput["redirect_configuration_name"] = redirectId.RedirectConfigurationName
							ruleOutput["redirect_configuration_id"] = redirectId.ID()
						}

						if rewrite := ruleProps.RewriteRuleSet; rewrite != nil && rewrite.Id != nil {
							rewriteId, err := parse.RewriteRuleSetIDInsensitively(*rewrite.Id)
							if err != nil {
								return nil, err
							}
							ruleOutput["rewrite_rule_set_name"] = rewriteId.Name
							ruleOutput["rewrite_rule_set_id"] = rewriteId.ID()
						}

						if fwp := ruleProps.FirewallPolicy; fwp != nil && fwp.Id != nil {
							policyId, err := webapplicationfirewallpolicies.ParseApplicationGatewayWebApplicationFirewallPolicyIDInsensitively(*fwp.Id)
							if err != nil {
								return nil, err
							}
							ruleOutput["firewall_policy_id"] = policyId.ID()
						}

						pathOutputs := make([]interface{}, 0)
						if paths := ruleProps.Paths; paths != nil {
							for _, rulePath := range *paths {
								pathOutputs = append(pathOutputs, rulePath)
							}
						}
						ruleOutput["paths"] = pathOutputs
					}

					pathRules = append(pathRules, ruleOutput)
				}
			}
			output["path_rule"] = pathRules
		}

		results = append(results, output)
	}

	return results, nil
}

func expandApplicationGatewayWafConfig(d *pluginsdk.ResourceData) *applicationgateways.ApplicationGatewayWebApplicationFirewallConfiguration {
	vs := d.Get("waf_configuration").([]interface{})
	if len(vs) == 0 || vs[0] == nil {
		return nil
	}
	v := vs[0].(map[string]interface{})

	enabled := v["enabled"].(bool)
	mode := v["firewall_mode"].(string)
	ruleSetType := v["rule_set_type"].(string)
	ruleSetVersion := v["rule_set_version"].(string)
	fileUploadLimitInMb := v["file_upload_limit_mb"].(int)
	requestBodyCheck := v["request_body_check"].(bool)
	maxRequestBodySizeInKb := v["max_request_body_size_kb"].(int)

	return &applicationgateways.ApplicationGatewayWebApplicationFirewallConfiguration{
		Enabled:                enabled,
		FirewallMode:           applicationgateways.ApplicationGatewayFirewallMode(mode),
		RuleSetType:            ruleSetType,
		RuleSetVersion:         ruleSetVersion,
		FileUploadLimitInMb:    pointer.To(int64(fileUploadLimitInMb)),
		RequestBodyCheck:       pointer.To(requestBodyCheck),
		MaxRequestBodySizeInKb: pointer.To(int64(maxRequestBodySizeInKb)),
		DisabledRuleGroups:     expandApplicationGatewayFirewallDisabledRuleGroup(v["disabled_rule_group"].([]interface{})),
		Exclusions:             expandApplicationGatewayFirewallExclusion(v["exclusion"].([]interface{})),
	}
}

func flattenApplicationGatewayWafConfig(input *applicationgateways.ApplicationGatewayWebApplicationFirewallConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	output := make(map[string]interface{})

	output["enabled"] = input.Enabled
	output["firewall_mode"] = string(input.FirewallMode)
	output["rule_set_type"] = input.RuleSetType
	output["rule_set_version"] = input.RuleSetVersion

	if input.DisabledRuleGroups != nil {
		output["disabled_rule_group"] = flattenApplicationGateWayDisabledRuleGroups(input.DisabledRuleGroups)
	}

	if input.FileUploadLimitInMb != nil {
		output["file_upload_limit_mb"] = int(*input.FileUploadLimitInMb)
	}

	if input.RequestBodyCheck != nil {
		output["request_body_check"] = *input.RequestBodyCheck
	}

	if input.MaxRequestBodySizeInKb != nil {
		output["max_request_body_size_kb"] = int(*input.MaxRequestBodySizeInKb)
	}

	if input.Exclusions != nil {
		output["exclusion"] = flattenApplicationGatewayFirewallExclusion(input.Exclusions)
	}
	results = append(results, output)

	return results
}

func expandApplicationGatewayFirewallDisabledRuleGroup(d []interface{}) *[]applicationgateways.ApplicationGatewayFirewallDisabledRuleGroup {
	if len(d) == 0 {
		return nil
	}

	disabledRuleGroups := make([]applicationgateways.ApplicationGatewayFirewallDisabledRuleGroup, 0)
	for _, disabledRuleGroup := range d {
		disabledRuleGroupMap := disabledRuleGroup.(map[string]interface{})

		ruleGroupName := disabledRuleGroupMap["rule_group_name"].(string)

		ruleGroup := applicationgateways.ApplicationGatewayFirewallDisabledRuleGroup{
			RuleGroupName: ruleGroupName,
		}

		rules := make([]int64, 0)
		for _, rule := range disabledRuleGroupMap["rules"].([]interface{}) {
			rules = append(rules, int64(rule.(int)))
		}

		if len(rules) > 0 {
			ruleGroup.Rules = &rules
		}

		disabledRuleGroups = append(disabledRuleGroups, ruleGroup)
	}
	return &disabledRuleGroups
}

func flattenApplicationGateWayDisabledRuleGroups(input *[]applicationgateways.ApplicationGatewayFirewallDisabledRuleGroup) []interface{} {
	ruleGroups := make([]interface{}, 0)
	for _, ruleGroup := range *input {
		ruleGroupOutput := map[string]interface{}{}

		ruleGroupOutput["rule_group_name"] = ruleGroup.RuleGroupName

		ruleOutputs := make([]interface{}, 0)
		if rules := ruleGroup.Rules; rules != nil {
			for _, rule := range *rules {
				ruleOutputs = append(ruleOutputs, rule)
			}
		}
		ruleGroupOutput["rules"] = ruleOutputs

		ruleGroups = append(ruleGroups, ruleGroupOutput)
	}
	return ruleGroups
}

func expandApplicationGatewayFirewallExclusion(d []interface{}) *[]applicationgateways.ApplicationGatewayFirewallExclusion {
	if len(d) == 0 {
		return nil
	}

	exclusions := make([]applicationgateways.ApplicationGatewayFirewallExclusion, 0)
	for _, exclusion := range d {
		exclusionMap := exclusion.(map[string]interface{})

		matchVariable := exclusionMap["match_variable"].(string)
		selectorMatchOperator := exclusionMap["selector_match_operator"].(string)
		selector := exclusionMap["selector"].(string)

		exclusionList := applicationgateways.ApplicationGatewayFirewallExclusion{
			MatchVariable:         matchVariable,
			SelectorMatchOperator: selectorMatchOperator,
			Selector:              selector,
		}

		exclusions = append(exclusions, exclusionList)
	}

	return &exclusions
}

func flattenApplicationGatewayFirewallExclusion(input *[]applicationgateways.ApplicationGatewayFirewallExclusion) []interface{} {
	exclusionLists := make([]interface{}, 0)
	for _, exclusionList := range *input {
		exclusionListOutput := map[string]interface{}{}
		exclusionListOutput["match_variable"] = exclusionList.MatchVariable
		exclusionListOutput["selector_match_operator"] = exclusionList.SelectorMatchOperator
		exclusionListOutput["selector"] = exclusionList.Selector
		exclusionLists = append(exclusionLists, exclusionListOutput)
	}
	return exclusionLists
}

func expandApplicationGatewayCustomErrorConfigurations(vs []interface{}) *[]applicationgateways.ApplicationGatewayCustomError {
	results := make([]applicationgateways.ApplicationGatewayCustomError, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})
		statusCode := v["status_code"].(string)
		customErrorPageUrl := v["custom_error_page_url"].(string)

		output := applicationgateways.ApplicationGatewayCustomError{
			StatusCode:         pointer.To(applicationgateways.ApplicationGatewayCustomErrorStatusCode(statusCode)),
			CustomErrorPageURL: pointer.To(customErrorPageUrl),
		}
		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayCustomErrorConfigurations(input *[]applicationgateways.ApplicationGatewayCustomError) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		output["status_code"] = v.StatusCode

		if v.CustomErrorPageURL != nil {
			output["custom_error_page_url"] = *v.CustomErrorPageURL
		}

		results = append(results, output)
	}

	return results
}

func checkSslPolicy(sslPolicy []interface{}) error {
	if len(sslPolicy) > 0 && sslPolicy[0] != nil {
		v := sslPolicy[0].(map[string]interface{})
		disabledProtocols := v["disabled_protocols"].([]interface{})
		policyType := v["policy_type"].(string)
		if len(disabledProtocols) > 0 && policyType != "" {
			return fmt.Errorf("setting disabled_protocols is not allowed when policy_type is defined")
		}
	}
	return nil
}

func checkBasicSkuFeatures(d *pluginsdk.ResourceDiff) error {
	_, hasAutoscaleConfig := d.GetOk("autoscale_configuration.0")
	if hasAutoscaleConfig {
		return fmt.Errorf("The Application Gateway does not support `autoscale_configuration` blocks for the selected SKU tier %q", applicationgateways.ApplicationGatewaySkuNameBasic)
	}

	capacity, hasCapacityConfig := d.GetOk("sku.0.capacity")
	if hasCapacityConfig {
		if capacity.(int) > 2 || capacity.(int) < 1 {
			return fmt.Errorf("`capacity` value %q for the selected SKU tier %q is invalid. Value must be between [1-2]", capacity, applicationgateways.ApplicationGatewaySkuNameBasic)
		}
	} else {
		return fmt.Errorf("The Application Gateway must specify a `capacity` value between [1-2] for the selected SKU tier %q", applicationgateways.ApplicationGatewaySkuNameBasic)
	}

	_, hasMtlsConfig := d.GetOk("trusted_client_certificate")
	if hasMtlsConfig {
		return fmt.Errorf("The Application Gateway does not support `trusted_client_certificate` blocks for the selected SKU tier %q", applicationgateways.ApplicationGatewaySkuNameBasic)
	}

	_, hasRewriteRuleSetConfig := d.GetOk("rewrite_rule_set")
	if hasRewriteRuleSetConfig {
		return fmt.Errorf("The Application Gateway does not support `rewrite_rule_set` blocks for the selected SKU tier %q", applicationgateways.ApplicationGatewaySkuNameBasic)
	}

	return nil
}

func applicationGatewayCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
	_, hasAutoscaleConfig := d.GetOk("autoscale_configuration.0")
	capacity, hasCapacity := d.GetOk("sku.0.capacity")
	tier := d.Get("sku.0.tier").(string)

	if tier == string(applicationgateways.ApplicationGatewaySkuNameBasic) {
		err := checkBasicSkuFeatures(d)
		if err != nil {
			return err
		}
	} else if !hasAutoscaleConfig && !hasCapacity {
		return fmt.Errorf("The Application Gateway must specify either `capacity` or `autoscale_configuration` for the selected SKU tier %q", tier)
	}

	sslPolicy := d.Get("ssl_policy").([]interface{})
	if err := checkSslPolicy(sslPolicy); err != nil {
		return err
	}

	sslProfiles := d.Get("ssl_profile").([]interface{})
	if len(sslProfiles) > 0 {
		for _, profile := range sslProfiles {
			if profile == nil {
				continue
			}
			v := profile.(map[string]interface{})
			if policy, ok := v["ssl_policy"]; ok && policy != nil {
				if err := checkSslPolicy(policy.([]interface{})); err != nil {
					return err
				}
			}
		}
	}

	if hasCapacity {
		if (strings.EqualFold(tier, string(applicationgateways.ApplicationGatewayTierStandard)) || strings.EqualFold(tier, string(applicationgateways.ApplicationGatewayTierWAF))) && (capacity.(int) < 1 || capacity.(int) > 32) {
			return fmt.Errorf("The value '%d' exceeds the maximum capacity allowed for a %q V1 SKU, the %q SKU must have a capacity value between 1 and 32", capacity, tier, tier)
		}

		if (strings.EqualFold(tier, string(applicationgateways.ApplicationGatewayTierStandardVTwo)) || strings.EqualFold(tier, string(applicationgateways.ApplicationGatewayTierWAFVTwo))) && (capacity.(int) < 1 || capacity.(int) > 125) {
			return fmt.Errorf("The value '%d' exceeds the maximum capacity allowed for a %q V2 SKU, the %q SKU must have a capacity value between 1 and 125", capacity, tier, tier)
		}
	}

	return nil
}

func applicationGatewayHttpListnerHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(m["name"].(string))
		buf.WriteString(m["frontend_ip_configuration_name"].(string))
		buf.WriteString(m["frontend_port_name"].(string))
		buf.WriteString(m["protocol"].(string))
		if v, ok := m["host_name"]; ok {
			buf.WriteString(v.(string))
		}
		if hostNames, ok := m["host_names"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", hostNames.(*pluginsdk.Set).List()))
		}
		if v, ok := m["ssl_certificate_name"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["require_sni"]; ok {
			buf.WriteString(fmt.Sprintf("%t", v.(bool)))
		}
		if v, ok := m["firewall_policy_id"]; ok {
			buf.WriteString(strings.ToLower(v.(string)))
		}
		if v, ok := m["ssl_profile_name"]; ok {
			buf.WriteString(v.(string))
		}
		if customErrorConfiguration, ok := m["custom_error_configuration"].([]interface{}); ok {
			for _, customErrorAttrs := range customErrorConfiguration {
				customError := customErrorAttrs.(map[string]interface{})
				if statusCode, ok := customError["status_code"]; ok {
					buf.WriteString(statusCode.(string))
				}
				if pageUrl, ok := customError["custom_error_page_url"]; ok {
					buf.WriteString(fmt.Sprintf(pageUrl.(string)))
				}
			}
		}
	}

	return pluginsdk.HashString(buf.String())
}

func applicationGatewayBackendSettingsHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(m["name"].(string))
		buf.WriteString(fmt.Sprintf("%d", m["port"].(int)))
		buf.WriteString(m["protocol"].(string))
		buf.WriteString(m["cookie_based_affinity"].(string))

		if v, ok := m["path"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["affinity_cookie_name"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["host_name"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["probe_name"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["pick_host_name_from_backend_address"]; ok {
			buf.WriteString(fmt.Sprintf("%t", v.(bool)))
		}
		if v, ok := m["request_timeout"]; ok {
			buf.WriteString(fmt.Sprintf("%d", v.(int)))
		}
		if authCert, ok := m["authentication_certificate"].([]interface{}); ok {
			for _, ac := range authCert {
				config := ac.(map[string]interface{})
				buf.WriteString(config["name"].(string))
			}
		}
		if connectionDraining, ok := m["connection_draining"].([]interface{}); ok {
			for _, ac := range connectionDraining {
				config := ac.(map[string]interface{})
				buf.WriteString(fmt.Sprintf("%t", config["enabled"].(bool)))
				buf.WriteString(fmt.Sprintf("%d", config["drain_timeout_sec"].(int)))
			}
		}
		if trustedRootCertificateNames, ok := m["trusted_root_certificate_names"]; ok {
			buf.WriteString(fmt.Sprintf("%s", trustedRootCertificateNames.([]interface{})))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func applicationGatewaySSLCertificate(v interface{}) int {
	var buf bytes.Buffer
	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(m["name"].(string))

		if v, ok := m["data"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["password"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["key_vault_secret_id"]; ok {
			buf.WriteString(strings.ToLower(v.(string)))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func applicationGatewayBackendAddressPool(v interface{}) int {
	var buf bytes.Buffer
	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(m["name"].(string))

		if fqdns, ok := m["fqdns"]; ok {
			buf.WriteString(fmt.Sprintf("%s", fqdns.(*pluginsdk.Set).List()))
		}
		if ips, ok := m["ip_addresses"]; ok {
			buf.WriteString(fmt.Sprintf("%s", ips.(*pluginsdk.Set).List()))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func applicationGatewayProbeHash(v interface{}) int {
	var buf bytes.Buffer
	if m, ok := v.(map[string]interface{}); ok {
		buf.WriteString(m["name"].(string))
		buf.WriteString(m["protocol"].(string))
		buf.WriteString(m["path"].(string))
		buf.WriteString(fmt.Sprintf("%d", m["interval"].(int)))
		buf.WriteString(fmt.Sprintf("%d", m["timeout"].(int)))
		buf.WriteString(fmt.Sprintf("%d", m["unhealthy_threshold"].(int)))

		if v, ok := m["host"]; ok {
			buf.WriteString(v.(string))
		}
		if v, ok := m["port"]; ok {
			buf.WriteString(fmt.Sprintf("%d", v.(int)))
		}
		if v, ok := m["pick_host_name_from_backend_http_settings"]; ok {
			buf.WriteString(fmt.Sprintf("%t", v.(bool)))
		}
		if v, ok := m["minimum_servers"]; ok {
			buf.WriteString(fmt.Sprintf("%d", v.(int)))
		}
		if match, ok := m["match"]; ok {
			if attrs := match.([]interface{}); len(attrs) == 1 {
				attr := attrs[0].(map[string]interface{})
				if attr["body"].(string) != "" || len(attr["status_code"].([]interface{})) != 0 {
					buf.WriteString(fmt.Sprintf("%s-%+v", attr["body"].(string), attr["status_code"].([]interface{})))
				}
			}
		}
	}

	return pluginsdk.HashString(buf.String())
}
