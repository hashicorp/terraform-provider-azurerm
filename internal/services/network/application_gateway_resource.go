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

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

// See https://github.com/Azure/azure-sdk-for-go/blob/master/services/network/mgmt/2018-04-01/network/models.go
func possibleApplicationGatewaySslCipherSuiteValues() []string {
	cipherSuites := make([]string, 0)
	for _, cipherSuite := range network.PossibleApplicationGatewaySslCipherSuiteValues() {
		cipherSuites = append(cipherSuites, string(cipherSuite))
	}
	return cipherSuites
}

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
							string(network.ApplicationGatewaySslProtocolTLSv10),
							string(network.ApplicationGatewaySslProtocolTLSv11),
							string(network.ApplicationGatewaySslProtocolTLSv12),
							string(network.ApplicationGatewaySslProtocolTLSv13),
						}, false),
					},
				},

				"policy_type": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(network.ApplicationGatewaySslPolicyTypeCustom),
						string(network.ApplicationGatewaySslPolicyTypeCustomV2),
						string(network.ApplicationGatewaySslPolicyTypePredefined),
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
						ValidateFunc: validation.StringInSlice(possibleApplicationGatewaySslCipherSuiteValues(), false),
					},
				},

				"min_protocol_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(network.ApplicationGatewaySslProtocolTLSv10),
						string(network.ApplicationGatewaySslProtocolTLSv11),
						string(network.ApplicationGatewaySslProtocolTLSv12),
						string(network.ApplicationGatewaySslProtocolTLSv13),
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
			_, err := parse.ApplicationGatewayID(id)
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

			"identity": commonschema.UserAssignedIdentityOptional(),

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
								string(network.ProtocolHTTP),
								string(network.ProtocolHTTPS),
							}, false),
						},

						"cookie_based_affinity": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.ApplicationGatewayCookieBasedAffinityEnabled),
								string(network.ApplicationGatewayCookieBasedAffinityDisabled),
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
							Default:  string(network.IPAllocationMethodDynamic),
							ValidateFunc: validation.StringInSlice([]string{
								string(network.IPAllocationMethodDynamic),
								string(network.IPAllocationMethodStatic),
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
								string(network.ProtocolHTTP),
								string(network.ProtocolHTTPS),
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
											string(network.ApplicationGatewayCustomErrorStatusCodeHTTPStatus403),
											string(network.ApplicationGatewayCustomErrorStatusCodeHTTPStatus502),
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
											string(network.IPAllocationMethodDynamic),
											string(network.IPAllocationMethodStatic),
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
								string(network.ApplicationGatewayRequestRoutingRuleTypeBasic),
								string(network.ApplicationGatewayRequestRoutingRuleTypePathBasedRouting),
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
								string(network.ApplicationGatewayRedirectTypePermanent),
								string(network.ApplicationGatewayRedirectTypeTemporary),
								string(network.ApplicationGatewayRedirectTypeFound),
								string(network.ApplicationGatewayRedirectTypeSeeOther),
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
								string(network.ApplicationGatewaySkuNameStandardSmall),
								string(network.ApplicationGatewaySkuNameStandardMedium),
								string(network.ApplicationGatewaySkuNameStandardLarge),
								string(network.ApplicationGatewaySkuNameStandardV2),
								string(network.ApplicationGatewaySkuNameWAFLarge),
								string(network.ApplicationGatewaySkuNameWAFMedium),
								string(network.ApplicationGatewaySkuNameWAFV2),
							}, false),
						},

						"tier": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.ApplicationGatewayTierStandard),
								string(network.ApplicationGatewayTierStandardV2),
								string(network.ApplicationGatewayTierWAF),
								string(network.ApplicationGatewayTierWAFV2),
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
								string(network.ProtocolHTTP),
								string(network.ProtocolHTTPS),
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
													Computed: true,
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

						"verify_client_cert_issuer_dn": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
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
								string(network.ApplicationGatewayFirewallModeDetection),
								string(network.ApplicationGatewayFirewallModePrevention),
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
											string(network.OwaspCrsExclusionEntryMatchVariableRequestArgKeys),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestArgNames),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestArgValues),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestCookieKeys),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestCookieNames),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestCookieValues),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestHeaderKeys),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestHeaderNames),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestHeaderValues),
										}, false),
									},

									"selector_match_operator": {
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											string(network.OwaspCrsExclusionEntrySelectorMatchOperatorContains),
											string(network.OwaspCrsExclusionEntrySelectorMatchOperatorEndsWith),
											string(network.OwaspCrsExclusionEntrySelectorMatchOperatorEquals),
											string(network.OwaspCrsExclusionEntrySelectorMatchOperatorEqualsAny),
											string(network.OwaspCrsExclusionEntrySelectorMatchOperatorStartsWith),
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
								string(network.ApplicationGatewayCustomErrorStatusCodeHTTPStatus403),
								string(network.ApplicationGatewayCustomErrorStatusCodeHTTPStatus502),
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

			"tags": tags.Schema(),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(applicationGatewayCustomizeDiff),
	}
}

func resourceApplicationGatewayCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ApplicationGatewaysClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Application Gateway creation.")

	id := parse.NewApplicationGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_application_gateway", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
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

	gatewayIPConfigurations, stopApplicationGateway := expandApplicationGatewayIPConfigurations(d)

	globalConfiguration := expandApplicationGatewayGlobalConfiguration(d.Get("global").([]interface{}))

	httpListeners, err := expandApplicationGatewayHTTPListeners(d, id.ID())
	if err != nil {
		return fmt.Errorf("fail to expand `http_listener`: %+v", err)
	}

	rewriteRuleSets, err := expandApplicationGatewayRewriteRuleSets(d)
	if err != nil {
		return fmt.Errorf("expanding `rewrite_rule_set`: %v", err)
	}

	gateway := network.ApplicationGateway{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
		ApplicationGatewayPropertiesFormat: &network.ApplicationGatewayPropertiesFormat{
			AutoscaleConfiguration:        expandApplicationGatewayAutoscaleConfiguration(d),
			AuthenticationCertificates:    expandApplicationGatewayAuthenticationCertificates(d.Get("authentication_certificate").([]interface{})),
			TrustedRootCertificates:       trustedRootCertificates,
			CustomErrorConfigurations:     expandApplicationGatewayCustomErrorConfigurations(d.Get("custom_error_configuration").([]interface{})),
			BackendAddressPools:           expandApplicationGatewayBackendAddressPools(d),
			BackendHTTPSettingsCollection: expandApplicationGatewayBackendHTTPSettings(d, id.ID()),
			EnableHTTP2:                   utils.Bool(enablehttp2),
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
			URLPathMaps:     urlPathMaps,
		},
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		gateway.Zones = &zones
	}

	if v, ok := d.GetOk("fips_enabled"); ok {
		gateway.ApplicationGatewayPropertiesFormat.EnableFips = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("force_firewall_policy_association"); ok {
		gateway.ApplicationGatewayPropertiesFormat.ForceFirewallPolicyAssociation = utils.Bool(v.(bool))
	}

	if _, ok := d.GetOk("identity"); ok {
		expandedIdentity, err := expandApplicationGatewayIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		gateway.Identity = expandedIdentity
	}

	// validation (todo these should probably be moved into their respective expand functions, which would then return an error?)
	for _, backendHttpSettings := range *gateway.ApplicationGatewayPropertiesFormat.BackendHTTPSettingsCollection {
		if props := backendHttpSettings.ApplicationGatewayBackendHTTPSettingsPropertiesFormat; props != nil {
			if props.HostName == nil || props.PickHostNameFromBackendAddress == nil {
				continue
			}

			if *props.HostName != "" && *props.PickHostNameFromBackendAddress {
				return fmt.Errorf("Only one of `host_name` or `pick_host_name_from_backend_address` can be set")
			}
		}
	}

	for _, probe := range *gateway.ApplicationGatewayPropertiesFormat.Probes {
		if props := probe.ApplicationGatewayProbePropertiesFormat; props != nil {
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
		gateway.ApplicationGatewayPropertiesFormat.WebApplicationFirewallConfiguration = expandApplicationGatewayWafConfig(d)
	}

	appGWSkuTier := d.Get("sku.0.tier").(string)
	wafFileUploadLimit := d.Get("waf_configuration.0.file_upload_limit_mb").(int)

	if appGWSkuTier != string(network.ApplicationGatewayTierWAFV2) && wafFileUploadLimit > 500 {
		return fmt.Errorf("Only SKU `%s` allows `file_upload_limit_mb` to exceed 500MB", network.ApplicationGatewayTierWAFV2)
	}

	if v, ok := d.GetOk("firewall_policy_id"); ok {
		id := v.(string)
		gateway.ApplicationGatewayPropertiesFormat.FirewallPolicy = &network.SubResource{
			ID: &id,
		}
	}

	if stopApplicationGateway {
		future, err := client.Stop(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("stopping %s: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for %s to stop: %+v", id, err)
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, gateway)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for create of %s: %+v", id, err)
	}

	if stopApplicationGateway {
		future, err := client.Start(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("starting %s: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for %s to start: %+v", id, err)
		}
	}

	d.SetId(id.ID())
	return resourceApplicationGatewayRead(d, meta)
}

func resourceApplicationGatewayUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ApplicationGatewaysClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApplicationGatewayID(d.Id())
	if err != nil {
		return err
	}

	applicationGateway, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if d.HasChange("tags") {
		applicationGateway.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if applicationGateway.ApplicationGatewayPropertiesFormat == nil {
		applicationGateway.ApplicationGatewayPropertiesFormat = &network.ApplicationGatewayPropertiesFormat{}
	}

	if d.HasChange("enable_http2") {
		applicationGateway.ApplicationGatewayPropertiesFormat.EnableHTTP2 = utils.Bool(d.Get("enable_http2").(bool))
	}

	if d.HasChange("trusted_root_certificate") {
		trustedRootCertificates, err := expandApplicationGatewayTrustedRootCertificates(d.Get("trusted_root_certificate").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `trusted_root_certificate`: %+v", err)
		}
		applicationGateway.ApplicationGatewayPropertiesFormat.TrustedRootCertificates = trustedRootCertificates
	}

	if d.HasChange("request_routing_rule") {
		requestRoutingRules, err := expandApplicationGatewayRequestRoutingRules(d, id.ID())
		if err != nil {
			return fmt.Errorf("expanding `request_routing_rule`: %+v", err)
		}
		applicationGateway.ApplicationGatewayPropertiesFormat.RequestRoutingRules = requestRoutingRules
	}

	if d.HasChange("url_path_map") {
		urlPathMaps, err := expandApplicationGatewayURLPathMaps(d, id.ID())
		if err != nil {
			return fmt.Errorf("expanding `url_path_map`: %+v", err)
		}

		applicationGateway.ApplicationGatewayPropertiesFormat.URLPathMaps = urlPathMaps
	}

	if d.HasChange("redirect_configuration") {
		redirectConfigurations, err := expandApplicationGatewayRedirectConfigurations(d, id.ID())
		if err != nil {
			return fmt.Errorf("expanding `redirect_configuration`: %+v", err)
		}

		applicationGateway.ApplicationGatewayPropertiesFormat.RedirectConfigurations = redirectConfigurations
	}

	if d.HasChange("ssl_certificate") {
		sslCertificates, err := expandApplicationGatewaySslCertificates(d)
		if err != nil {
			return fmt.Errorf("expanding `ssl_certificate`: %+v", err)
		}

		applicationGateway.ApplicationGatewayPropertiesFormat.SslCertificates = sslCertificates
	}

	if d.HasChange("trusted_client_certificate") {
		trustedClientCertificates, err := expandApplicationGatewayTrustedClientCertificates(d)
		if err != nil {
			return fmt.Errorf("expanding `trusted_client_certificate`: %+v", err)
		}

		applicationGateway.ApplicationGatewayPropertiesFormat.TrustedClientCertificates = trustedClientCertificates
	}

	if d.HasChange("ssl_profile") {
		applicationGateway.ApplicationGatewayPropertiesFormat.SslProfiles = expandApplicationGatewaySslProfiles(d, id.ID())
	}

	gatewayIPConfigurations, stopApplicationGateway := expandApplicationGatewayIPConfigurations(d)
	if d.HasChange("gateway_ip_configuration") {
		applicationGateway.ApplicationGatewayPropertiesFormat.GatewayIPConfigurations = gatewayIPConfigurations
	}

	if d.HasChange("global") {
		globalConfiguration := expandApplicationGatewayGlobalConfiguration(d.Get("global").([]interface{}))
		applicationGateway.ApplicationGatewayPropertiesFormat.GlobalConfiguration = globalConfiguration
	}

	if d.HasChange("http_listener") {
		httpListeners, err := expandApplicationGatewayHTTPListeners(d, id.ID())
		if err != nil {
			return fmt.Errorf("fail to expand `http_listener`: %+v", err)
		}

		applicationGateway.ApplicationGatewayPropertiesFormat.HTTPListeners = httpListeners
	}

	if d.HasChange("rewrite_rule_set") {
		rewriteRuleSets, err := expandApplicationGatewayRewriteRuleSets(d)
		if err != nil {
			return fmt.Errorf("expanding `rewrite_rule_set`: %v", err)
		}

		applicationGateway.ApplicationGatewayPropertiesFormat.RewriteRuleSets = rewriteRuleSets
	}

	if d.HasChange("autoscale_configuration") {
		applicationGateway.ApplicationGatewayPropertiesFormat.AutoscaleConfiguration = expandApplicationGatewayAutoscaleConfiguration(d)
	}

	if d.HasChange("authentication_certificate") {
		applicationGateway.ApplicationGatewayPropertiesFormat.AuthenticationCertificates = expandApplicationGatewayAuthenticationCertificates(d.Get("authentication_certificate").([]interface{}))
	}

	if d.HasChange("custom_error_configuration") {
		applicationGateway.ApplicationGatewayPropertiesFormat.CustomErrorConfigurations = expandApplicationGatewayCustomErrorConfigurations(d.Get("custom_error_configuration").([]interface{}))
	}

	if d.HasChange("backend_address_pool") {
		applicationGateway.ApplicationGatewayPropertiesFormat.BackendAddressPools = expandApplicationGatewayBackendAddressPools(d)
	}

	if d.HasChange("backend_http_settings") {
		applicationGateway.ApplicationGatewayPropertiesFormat.BackendHTTPSettingsCollection = expandApplicationGatewayBackendHTTPSettings(d, id.ID())
	}

	if d.HasChange("frontend_ip_configuration") {
		applicationGateway.ApplicationGatewayPropertiesFormat.FrontendIPConfigurations = expandApplicationGatewayFrontendIPConfigurations(d, id.ID())
	}

	if d.HasChange("frontend_port") {
		applicationGateway.ApplicationGatewayPropertiesFormat.FrontendPorts = expandApplicationGatewayFrontendPorts(d)
	}

	if d.HasChange("private_link_configuration") {
		applicationGateway.ApplicationGatewayPropertiesFormat.PrivateLinkConfigurations = expandApplicationGatewayPrivateLinkConfigurations(d)
	}

	if d.HasChange("probe") {
		applicationGateway.ApplicationGatewayPropertiesFormat.Probes = expandApplicationGatewayProbes(d)
	}

	if d.HasChange("sku") {
		applicationGateway.ApplicationGatewayPropertiesFormat.Sku = expandApplicationGatewaySku(d)
	}

	if d.HasChange("ssl_policy") {
		applicationGateway.ApplicationGatewayPropertiesFormat.SslPolicy = expandApplicationGatewaySslPolicy(d.Get("ssl_policy").([]interface{}))
	}

	if d.HasChange("zones") {
		zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
		if len(zones) > 0 {
			applicationGateway.Zones = &zones
		}
	}

	if d.HasChange("fips_enabled") {
		applicationGateway.ApplicationGatewayPropertiesFormat.EnableFips = utils.Bool(d.Get("fips_enabled").(bool))
	}

	if d.HasChange("force_firewall_policy_association") {
		applicationGateway.ApplicationGatewayPropertiesFormat.ForceFirewallPolicyAssociation = utils.Bool(d.Get("force_firewall_policy_association").(bool))
	}

	if d.HasChange("identity") {
		expandedIdentity, err := expandApplicationGatewayIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		applicationGateway.Identity = expandedIdentity
	}

	// validation (todo these should probably be moved into their respective expand functions, which would then return an error?)
	if applicationGateway.ApplicationGatewayPropertiesFormat != nil && applicationGateway.ApplicationGatewayPropertiesFormat.BackendHTTPSettingsCollection != nil {
		for _, backendHttpSettings := range *applicationGateway.ApplicationGatewayPropertiesFormat.BackendHTTPSettingsCollection {
			if props := backendHttpSettings.ApplicationGatewayBackendHTTPSettingsPropertiesFormat; props != nil {
				if props.HostName == nil || props.PickHostNameFromBackendAddress == nil {
					continue
				}

				if *props.HostName != "" && *props.PickHostNameFromBackendAddress {
					return fmt.Errorf("Only one of `host_name` or `pick_host_name_from_backend_address` can be set")
				}
			}
		}
	}

	if applicationGateway.ApplicationGatewayPropertiesFormat != nil && applicationGateway.ApplicationGatewayPropertiesFormat.Probes != nil {
		for _, probe := range *applicationGateway.ApplicationGatewayPropertiesFormat.Probes {
			if props := probe.ApplicationGatewayProbePropertiesFormat; props != nil {
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
		applicationGateway.ApplicationGatewayPropertiesFormat.WebApplicationFirewallConfiguration = expandApplicationGatewayWafConfig(d)
	}

	appGWSkuTier := d.Get("sku.0.tier").(string)
	wafFileUploadLimit := d.Get("waf_configuration.0.file_upload_limit_mb").(int)

	if appGWSkuTier != string(network.ApplicationGatewayTierWAFV2) && wafFileUploadLimit > 500 {
		return fmt.Errorf("Only SKU `%s` allows `file_upload_limit_mb` to exceed 500MB", network.ApplicationGatewayTierWAFV2)
	}

	if d.HasChange("firewall_policy_id") {
		applicationGateway.ApplicationGatewayPropertiesFormat.FirewallPolicy = &network.SubResource{
			ID: utils.String(d.Get("firewall_policy_id").(string)),
		}
	}

	if stopApplicationGateway {
		future, err := client.Stop(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("stopping %s: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for %s to stop: %+v", id, err)
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, applicationGateway)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}

	if stopApplicationGateway {
		future, err := client.Start(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("starting %s: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for %s to start: %+v", id, err)
		}
	}

	d.SetId(id.ID())
	return resourceApplicationGatewayRead(d, meta)
}

func resourceApplicationGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ApplicationGatewaysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApplicationGatewayID(d.Id())
	if err != nil {
		return err
	}

	applicationGateway, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(applicationGateway.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("location", location.NormalizeNilable(applicationGateway.Location))
	d.Set("zones", zones.FlattenUntyped(applicationGateway.Zones))

	identity, err := flattenApplicationGatewayIdentity(applicationGateway.Identity)
	if err != nil {
		return err
	}
	if err = d.Set("identity", identity); err != nil {
		return err
	}

	if props := applicationGateway.ApplicationGatewayPropertiesFormat; props != nil {
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

		urlPathMaps, err := flattenApplicationGatewayURLPathMaps(props.URLPathMaps)
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
		if props.FirewallPolicy != nil && props.FirewallPolicy.ID != nil {
			firewallPolicyId = *props.FirewallPolicy.ID
			policyId, err := webapplicationfirewallpolicies.ParseApplicationGatewayWebApplicationFirewallPolicyIDInsensitively(firewallPolicyId)
			if err == nil {
				firewallPolicyId = policyId.ID()
			}
		}
		d.Set("firewall_policy_id", firewallPolicyId)
	}

	return tags.FlattenAndSet(d, applicationGateway.Tags)
}

func resourceApplicationGatewayDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ApplicationGatewaysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ApplicationGatewayID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandApplicationGatewayIdentity(input []interface{}) (*network.ManagedServiceIdentity, error) {
	expanded, err := identity.ExpandUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := network.ManagedServiceIdentity{
		Type: network.ResourceIdentityType(string(expanded.Type)),
	}

	if expanded.Type == identity.TypeUserAssigned {
		out.UserAssignedIdentities = make(map[string]*network.ManagedServiceIdentityUserAssignedIdentitiesValue)
		for k := range expanded.IdentityIds {
			out.UserAssignedIdentities[k] = &network.ManagedServiceIdentityUserAssignedIdentitiesValue{}
		}
	}

	return &out, nil
}

func flattenApplicationGatewayIdentity(input *network.ManagedServiceIdentity) (*[]interface{}, error) {
	var transform *identity.UserAssignedMap

	if input != nil {
		transform = &identity.UserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		if input.UserAssignedIdentities != nil {
			for k, v := range input.UserAssignedIdentities {
				transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
					ClientId:    v.ClientID,
					PrincipalId: v.PrincipalID,
				}
			}
		}
	}

	return identity.FlattenUserAssignedMap(transform)
}

func expandApplicationGatewayAuthenticationCertificates(certs []interface{}) *[]network.ApplicationGatewayAuthenticationCertificate {
	results := make([]network.ApplicationGatewayAuthenticationCertificate, 0)

	for _, raw := range certs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		data := v["data"].(string)

		// data must be base64 encoded
		encodedData := utils.Base64EncodeIfNot(data)

		output := network.ApplicationGatewayAuthenticationCertificate{
			Name: utils.String(name),
			ApplicationGatewayAuthenticationCertificatePropertiesFormat: &network.ApplicationGatewayAuthenticationCertificatePropertiesFormat{
				Data: utils.String(encodedData),
			},
		}

		results = append(results, output)
	}

	return &results
}

func expandApplicationGatewayTrustedRootCertificates(certs []interface{}) (*[]network.ApplicationGatewayTrustedRootCertificate, error) {
	results := make([]network.ApplicationGatewayTrustedRootCertificate, 0)

	for _, raw := range certs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		data := v["data"].(string)
		kvsid := v["key_vault_secret_id"].(string)

		output := network.ApplicationGatewayTrustedRootCertificate{
			Name: utils.String(name),
			ApplicationGatewayTrustedRootCertificatePropertiesFormat: &network.ApplicationGatewayTrustedRootCertificatePropertiesFormat{},
		}

		switch {
		case data != "" && kvsid != "":
			return nil, fmt.Errorf("only one of `key_vault_secret_id` or `data` must be specified for the `trusted_root_certificate` block %q", name)
		case data != "":
			output.ApplicationGatewayTrustedRootCertificatePropertiesFormat.Data = utils.String(utils.Base64EncodeIfNot(data))
		case kvsid != "":
			output.ApplicationGatewayTrustedRootCertificatePropertiesFormat.KeyVaultSecretID = utils.String(kvsid)
		default:
			return nil, fmt.Errorf("either `key_vault_secret_id` or `data` must be specified for the `trusted_root_certificate` block %q", name)
		}

		results = append(results, output)
	}

	return &results, nil
}

func flattenApplicationGatewayAuthenticationCertificates(certs *[]network.ApplicationGatewayAuthenticationCertificate, d *pluginsdk.ResourceData) []interface{} {
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

		if v := cert.ID; v != nil {
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

func flattenApplicationGatewayTrustedRootCertificates(certs *[]network.ApplicationGatewayTrustedRootCertificate, d *pluginsdk.ResourceData) []interface{} {
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

		if v := cert.ID; v != nil {
			output["id"] = *v
		}

		kvsid := ""
		if props := cert.ApplicationGatewayTrustedRootCertificatePropertiesFormat; props != nil {
			if v := props.KeyVaultSecretID; v != nil {
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

func expandApplicationGatewayBackendAddressPools(d *pluginsdk.ResourceData) *[]network.ApplicationGatewayBackendAddressPool {
	vs := d.Get("backend_address_pool").(*schema.Set).List()
	results := make([]network.ApplicationGatewayBackendAddressPool, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})
		backendAddresses := make([]network.ApplicationGatewayBackendAddress, 0)

		if fqdnsConfig, ok := v["fqdns"]; ok {
			fqdns := fqdnsConfig.(*schema.Set).List()
			for _, ip := range fqdns {
				backendAddresses = append(backendAddresses, network.ApplicationGatewayBackendAddress{
					Fqdn: utils.String(ip.(string)),
				})
			}
		}

		if ipAddressesConfig, ok := v["ip_addresses"]; ok {
			ipAddresses := ipAddressesConfig.(*schema.Set).List()

			for _, ip := range ipAddresses {
				backendAddresses = append(backendAddresses, network.ApplicationGatewayBackendAddress{
					IPAddress: utils.String(ip.(string)),
				})
			}
		}

		name := v["name"].(string)
		output := network.ApplicationGatewayBackendAddressPool{
			Name: utils.String(name),
			ApplicationGatewayBackendAddressPoolPropertiesFormat: &network.ApplicationGatewayBackendAddressPoolPropertiesFormat{
				BackendAddresses: &backendAddresses,
			},
		}

		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayBackendAddressPools(input *[]network.ApplicationGatewayBackendAddressPool) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, config := range *input {
		ipAddressList := make([]interface{}, 0)
		fqdnList := make([]interface{}, 0)

		if props := config.ApplicationGatewayBackendAddressPoolPropertiesFormat; props != nil {
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

		if config.ID != nil {
			output["id"] = *config.ID
		}

		if config.Name != nil {
			output["name"] = *config.Name
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewayBackendHTTPSettings(d *pluginsdk.ResourceData, gatewayID string) *[]network.ApplicationGatewayBackendHTTPSettings {
	results := make([]network.ApplicationGatewayBackendHTTPSettings, 0)
	vs := d.Get("backend_http_settings").(*schema.Set).List()

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		path := v["path"].(string)
		port := int32(v["port"].(int))
		protocol := v["protocol"].(string)
		cookieBasedAffinity := v["cookie_based_affinity"].(string)
		pickHostNameFromBackendAddress := v["pick_host_name_from_backend_address"].(bool)
		requestTimeout := int32(v["request_timeout"].(int))

		setting := network.ApplicationGatewayBackendHTTPSettings{
			Name: &name,
			ApplicationGatewayBackendHTTPSettingsPropertiesFormat: &network.ApplicationGatewayBackendHTTPSettingsPropertiesFormat{
				CookieBasedAffinity:            network.ApplicationGatewayCookieBasedAffinity(cookieBasedAffinity),
				Path:                           utils.String(path),
				PickHostNameFromBackendAddress: utils.Bool(pickHostNameFromBackendAddress),
				Port:                           utils.Int32(port),
				Protocol:                       network.ApplicationGatewayProtocol(protocol),
				RequestTimeout:                 utils.Int32(requestTimeout),
				ConnectionDraining:             expandApplicationGatewayConnectionDraining(v),
			},
		}

		hostName := v["host_name"].(string)
		if hostName != "" {
			setting.ApplicationGatewayBackendHTTPSettingsPropertiesFormat.HostName = utils.String(hostName)
		}

		affinityCookieName := v["affinity_cookie_name"].(string)
		if affinityCookieName != "" {
			setting.AffinityCookieName = utils.String(affinityCookieName)
		}

		if v["authentication_certificate"] != nil {
			authCerts := v["authentication_certificate"].([]interface{})
			authCertSubResources := make([]network.SubResource, 0)

			for _, rawAuthCert := range authCerts {
				authCert := rawAuthCert.(map[string]interface{})
				authCertName := authCert["name"].(string)
				authCertID := fmt.Sprintf("%s/authenticationCertificates/%s", gatewayID, authCertName)
				authCertSubResource := network.SubResource{
					ID: utils.String(authCertID),
				}

				authCertSubResources = append(authCertSubResources, authCertSubResource)
			}

			setting.ApplicationGatewayBackendHTTPSettingsPropertiesFormat.AuthenticationCertificates = &authCertSubResources
		}

		if v["trusted_root_certificate_names"] != nil {
			trustedRootCertNames := v["trusted_root_certificate_names"].([]interface{})
			trustedRootCertSubResources := make([]network.SubResource, 0)

			for _, rawTrustedRootCertName := range trustedRootCertNames {
				trustedRootCertName := rawTrustedRootCertName.(string)
				trustedRootCertID := fmt.Sprintf("%s/trustedRootCertificates/%s", gatewayID, trustedRootCertName)
				trustedRootCertSubResource := network.SubResource{
					ID: utils.String(trustedRootCertID),
				}

				trustedRootCertSubResources = append(trustedRootCertSubResources, trustedRootCertSubResource)
			}

			setting.ApplicationGatewayBackendHTTPSettingsPropertiesFormat.TrustedRootCertificates = &trustedRootCertSubResources
		}

		probeName := v["probe_name"].(string)
		if probeName != "" {
			probeID := fmt.Sprintf("%s/probes/%s", gatewayID, probeName)
			setting.ApplicationGatewayBackendHTTPSettingsPropertiesFormat.Probe = &network.SubResource{
				ID: utils.String(probeID),
			}
		}

		results = append(results, setting)
	}

	return &results
}

func flattenApplicationGatewayBackendHTTPSettings(input *[]network.ApplicationGatewayBackendHTTPSettings) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.ID != nil {
			output["id"] = *v.ID
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.ApplicationGatewayBackendHTTPSettingsPropertiesFormat; props != nil {
			output["cookie_based_affinity"] = string(props.CookieBasedAffinity)

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

			output["protocol"] = string(props.Protocol)

			if timeout := props.RequestTimeout; timeout != nil {
				output["request_timeout"] = int(*timeout)
			}

			authenticationCertificates := make([]interface{}, 0)
			if certs := props.AuthenticationCertificates; certs != nil {
				for _, cert := range *certs {
					if cert.ID == nil {
						continue
					}

					certId, err := parse.AuthenticationCertificateIDInsensitively(*cert.ID)
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
					if cert.ID == nil {
						continue
					}

					certId, err := parse.TrustedRootCertificateIDInsensitively(*cert.ID)
					if err != nil {
						return nil, err
					}

					trustedRootCertificateNames = append(trustedRootCertificateNames, certId.Name)
				}
			}
			output["trusted_root_certificate_names"] = trustedRootCertificateNames

			if probe := props.Probe; probe != nil {
				if probe.ID != nil {
					id, err := parse.ProbeIDInsensitively(*probe.ID)
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

func expandApplicationGatewayConnectionDraining(d map[string]interface{}) *network.ApplicationGatewayConnectionDraining {
	connectionsRaw := d["connection_draining"].([]interface{})

	if len(connectionsRaw) == 0 {
		return nil
	}

	connectionRaw := connectionsRaw[0].(map[string]interface{})

	return &network.ApplicationGatewayConnectionDraining{
		Enabled:           utils.Bool(connectionRaw["enabled"].(bool)),
		DrainTimeoutInSec: utils.Int32(int32(connectionRaw["drain_timeout_sec"].(int))),
	}
}

func flattenApplicationGatewayConnectionDraining(input *network.ApplicationGatewayConnectionDraining) []interface{} {
	result := map[string]interface{}{}
	if input == nil {
		return []interface{}{}
	}

	if v := input.Enabled; v != nil {
		result["enabled"] = *v
	}
	if v := input.DrainTimeoutInSec; v != nil {
		result["drain_timeout_sec"] = *v
	}

	return []interface{}{result}
}

func expandApplicationGatewaySslPolicy(vs []interface{}) *network.ApplicationGatewaySslPolicy {
	policy := network.ApplicationGatewaySslPolicy{}
	disabledSSLProtocols := make([]network.ApplicationGatewaySslProtocol, 0)

	if len(vs) > 0 && vs[0] != nil {
		v := vs[0].(map[string]interface{})
		policyType := network.ApplicationGatewaySslPolicyType(v["policy_type"].(string))

		for _, policy := range v["disabled_protocols"].([]interface{}) {
			disabledSSLProtocols = append(disabledSSLProtocols, network.ApplicationGatewaySslProtocol(policy.(string)))
		}

		if policyType == network.ApplicationGatewaySslPolicyTypePredefined {
			policyName := network.ApplicationGatewaySslPolicyName(v["policy_name"].(string))
			policy = network.ApplicationGatewaySslPolicy{
				PolicyType: policyType,
				PolicyName: policyName,
			}
		} else if policyType == network.ApplicationGatewaySslPolicyTypeCustom || policyType == network.ApplicationGatewaySslPolicyTypeCustomV2 {
			minProtocolVersion := network.ApplicationGatewaySslProtocol(v["min_protocol_version"].(string))
			cipherSuites := make([]network.ApplicationGatewaySslCipherSuite, 0)

			for _, cipherSuite := range v["cipher_suites"].([]interface{}) {
				cipherSuites = append(cipherSuites, network.ApplicationGatewaySslCipherSuite(cipherSuite.(string)))
			}

			policy = network.ApplicationGatewaySslPolicy{
				PolicyType:         policyType,
				MinProtocolVersion: minProtocolVersion,
				CipherSuites:       &cipherSuites,
			}
		}
	}

	if len(disabledSSLProtocols) > 0 {
		policy = network.ApplicationGatewaySslPolicy{
			DisabledSslProtocols: &disabledSSLProtocols,
		}
	}

	return &policy
}

func flattenApplicationGatewaySslPolicy(input *network.ApplicationGatewaySslPolicy) []interface{} {
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

func expandApplicationGatewayHTTPListeners(d *pluginsdk.ResourceData, gatewayID string) (*[]network.ApplicationGatewayHTTPListener, error) {
	vs := d.Get("http_listener").(*schema.Set).List()

	results := make([]network.ApplicationGatewayHTTPListener, 0)

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

		listener := network.ApplicationGatewayHTTPListener{
			Name: utils.String(name),
			ApplicationGatewayHTTPListenerPropertiesFormat: &network.ApplicationGatewayHTTPListenerPropertiesFormat{
				FrontendIPConfiguration: &network.SubResource{
					ID: utils.String(frontendIPConfigID),
				},
				FrontendPort: &network.SubResource{
					ID: utils.String(frontendPortID),
				},
				Protocol:                    network.ApplicationGatewayProtocol(protocol),
				RequireServerNameIndication: utils.Bool(requireSNI),
				CustomErrorConfigurations:   customErrorConfigurations,
			},
		}

		host := v["host_name"].(string)
		hosts := v["host_names"].(*pluginsdk.Set).List()

		if host != "" && len(hosts) > 0 {
			return nil, fmt.Errorf("`host_name` and `host_names` cannot be specified together")
		}

		if host != "" {
			listener.ApplicationGatewayHTTPListenerPropertiesFormat.HostName = &host
		}

		if len(hosts) > 0 {
			listener.ApplicationGatewayHTTPListenerPropertiesFormat.HostNames = utils.ExpandStringSlice(hosts)
		}

		if sslCertName := v["ssl_certificate_name"].(string); sslCertName != "" {
			certID := fmt.Sprintf("%s/sslCertificates/%s", gatewayID, sslCertName)
			listener.ApplicationGatewayHTTPListenerPropertiesFormat.SslCertificate = &network.SubResource{
				ID: utils.String(certID),
			}
		}

		if firewallPolicyID != "" && len(firewallPolicyID) > 0 {
			listener.ApplicationGatewayHTTPListenerPropertiesFormat.FirewallPolicy = &network.SubResource{
				ID: utils.String(firewallPolicyID),
			}
		}

		if sslProfileName != "" && len(sslProfileName) > 0 {
			sslProfileID := fmt.Sprintf("%s/sslProfiles/%s", gatewayID, sslProfileName)
			listener.ApplicationGatewayHTTPListenerPropertiesFormat.SslProfile = &network.SubResource{
				ID: utils.String(sslProfileID),
			}
		}

		results = append(results, listener)
	}

	return &results, nil
}

func flattenApplicationGatewayHTTPListeners(input *[]network.ApplicationGatewayHTTPListener) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.ID != nil {
			output["id"] = *v.ID
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.ApplicationGatewayHTTPListenerPropertiesFormat; props != nil {
			if port := props.FrontendPort; port != nil {
				if port.ID != nil {
					portId, err := parse.FrontendPortIDInsensitively(*port.ID)
					if err != nil {
						return nil, err
					}
					output["frontend_port_name"] = portId.Name
					output["frontend_port_id"] = portId.ID()
				}
			}

			if feConfig := props.FrontendIPConfiguration; feConfig != nil {
				if feConfig.ID != nil {
					feConfigId, err := parse.FrontendIPConfigurationIDInsensitively(*feConfig.ID)
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

			output["protocol"] = string(props.Protocol)

			if cert := props.SslCertificate; cert != nil {
				if cert.ID != nil {
					certId, err := parse.SslCertificateIDInsensitively(*cert.ID)
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

			if fwp := props.FirewallPolicy; fwp != nil && fwp.ID != nil {
				policyId, err := webapplicationfirewallpolicies.ParseApplicationGatewayWebApplicationFirewallPolicyIDInsensitively(*fwp.ID)
				if err != nil {
					return nil, err
				}
				output["firewall_policy_id"] = policyId.ID()
			}

			if sslp := props.SslProfile; sslp != nil {
				if sslp.ID != nil {
					sslProfileId, err := parse.SslProfileIDInsensitively(*sslp.ID)
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

func expandApplicationGatewayIPConfigurations(d *pluginsdk.ResourceData) (*[]network.ApplicationGatewayIPConfiguration, bool) {
	vs := d.Get("gateway_ip_configuration").([]interface{})
	results := make([]network.ApplicationGatewayIPConfiguration, 0)
	stopApplicationGateway := false

	for _, configRaw := range vs {
		data := configRaw.(map[string]interface{})

		name := data["name"].(string)
		subnetID := data["subnet_id"].(string)

		output := network.ApplicationGatewayIPConfiguration{
			Name: utils.String(name),
			ApplicationGatewayIPConfigurationPropertiesFormat: &network.ApplicationGatewayIPConfigurationPropertiesFormat{
				Subnet: &network.SubResource{
					ID: utils.String(subnetID),
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

func flattenApplicationGatewayIPConfigurations(input *[]network.ApplicationGatewayIPConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.ID != nil {
			output["id"] = *v.ID
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.ApplicationGatewayIPConfigurationPropertiesFormat; props != nil {
			if subnet := props.Subnet; subnet != nil {
				if subnet.ID != nil {
					output["subnet_id"] = *subnet.ID
				}
			}
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewayGlobalConfiguration(input []interface{}) *network.ApplicationGatewayGlobalConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &network.ApplicationGatewayGlobalConfiguration{
		EnableRequestBuffering:  utils.Bool(v["request_buffering_enabled"].(bool)),
		EnableResponseBuffering: utils.Bool(v["response_buffering_enabled"].(bool)),
	}
}

func flattenApplicationGatewayGlobalConfiguration(input *network.ApplicationGatewayGlobalConfiguration) []interface{} {
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

func expandApplicationGatewayFrontendPorts(d *pluginsdk.ResourceData) *[]network.ApplicationGatewayFrontendPort {
	vs := d.Get("frontend_port").(*pluginsdk.Set).List()
	results := make([]network.ApplicationGatewayFrontendPort, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		port := int32(v["port"].(int))

		output := network.ApplicationGatewayFrontendPort{
			Name: utils.String(name),
			ApplicationGatewayFrontendPortPropertiesFormat: &network.ApplicationGatewayFrontendPortPropertiesFormat{
				Port: utils.Int32(port),
			},
		}
		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayFrontendPorts(input *[]network.ApplicationGatewayFrontendPort) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.ID != nil {
			output["id"] = *v.ID
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.ApplicationGatewayFrontendPortPropertiesFormat; props != nil {
			if props.Port != nil {
				output["port"] = int(*props.Port)
			}
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewayFrontendIPConfigurations(d *pluginsdk.ResourceData, gatewayID string) *[]network.ApplicationGatewayFrontendIPConfiguration {
	vs := d.Get("frontend_ip_configuration").([]interface{})
	results := make([]network.ApplicationGatewayFrontendIPConfiguration, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		properties := network.ApplicationGatewayFrontendIPConfigurationPropertiesFormat{}

		if val := v["subnet_id"].(string); val != "" {
			properties.Subnet = &network.SubResource{
				ID: utils.String(val),
			}
		}

		if val := v["private_ip_address_allocation"].(string); val != "" {
			properties.PrivateIPAllocationMethod = network.IPAllocationMethod(val)
		}

		if val := v["private_ip_address"].(string); val != "" {
			properties.PrivateIPAddress = utils.String(val)
		}

		if val := v["public_ip_address_id"].(string); val != "" {
			properties.PublicIPAddress = &network.SubResource{
				ID: utils.String(val),
			}
		}

		if val := v["private_link_configuration_name"].(string); val != "" {
			privateLinkConfigurationID := fmt.Sprintf("%s/privateLinkConfigurations/%s", gatewayID, val)
			properties.PrivateLinkConfiguration = &network.SubResource{
				ID: utils.String(privateLinkConfigurationID),
			}
		}

		name := v["name"].(string)
		output := network.ApplicationGatewayFrontendIPConfiguration{
			Name: utils.String(name),
			ApplicationGatewayFrontendIPConfigurationPropertiesFormat: &properties,
		}

		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayFrontendIPConfigurations(input *[]network.ApplicationGatewayFrontendIPConfiguration) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, config := range *input {
		output := make(map[string]interface{})
		if config.ID != nil {
			output["id"] = *config.ID
		}

		if config.Name != nil {
			output["name"] = *config.Name
		}

		if props := config.ApplicationGatewayFrontendIPConfigurationPropertiesFormat; props != nil {
			output["private_ip_address_allocation"] = string(props.PrivateIPAllocationMethod)

			if props.Subnet != nil && props.Subnet.ID != nil {
				output["subnet_id"] = *props.Subnet.ID
			}

			if props.PrivateIPAddress != nil {
				output["private_ip_address"] = *props.PrivateIPAddress
			}

			if props.PublicIPAddress != nil && props.PublicIPAddress.ID != nil {
				output["public_ip_address_id"] = *props.PublicIPAddress.ID
			}

			if props.PrivateLinkConfiguration != nil && props.PrivateLinkConfiguration.ID != nil {
				configurationID, err := parse.ApplicationGatewayPrivateLinkConfigurationIDInsensitively(*props.PrivateLinkConfiguration.ID)
				if err != nil {
					return nil, err
				}
				output["private_link_configuration_name"] = configurationID.PrivateLinkConfigurationName
				output["private_link_configuration_id"] = *props.PrivateLinkConfiguration.ID
			}
		}

		results = append(results, output)
	}

	return results, nil
}

func expandApplicationGatewayProbes(d *pluginsdk.ResourceData) *[]network.ApplicationGatewayProbe {
	vs := d.Get("probe").(*schema.Set).List()
	results := make([]network.ApplicationGatewayProbe, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		host := v["host"].(string)
		interval := int32(v["interval"].(int))
		minServers := int32(v["minimum_servers"].(int))
		name := v["name"].(string)
		probePath := v["path"].(string)
		protocol := v["protocol"].(string)
		port := int32(v["port"].(int))
		timeout := int32(v["timeout"].(int))
		unhealthyThreshold := int32(v["unhealthy_threshold"].(int))
		pickHostNameFromBackendHTTPSettings := v["pick_host_name_from_backend_http_settings"].(bool)

		output := network.ApplicationGatewayProbe{
			Name: utils.String(name),
			ApplicationGatewayProbePropertiesFormat: &network.ApplicationGatewayProbePropertiesFormat{
				Host:                                utils.String(host),
				Interval:                            utils.Int32(interval),
				MinServers:                          utils.Int32(minServers),
				Path:                                utils.String(probePath),
				Protocol:                            network.ApplicationGatewayProtocol(protocol),
				Timeout:                             utils.Int32(timeout),
				UnhealthyThreshold:                  utils.Int32(unhealthyThreshold),
				PickHostNameFromBackendHTTPSettings: utils.Bool(pickHostNameFromBackendHTTPSettings),
			},
		}

		matchConfigs := v["match"].([]interface{})
		if len(matchConfigs) > 0 {
			matchBody := ""
			outputMatch := &network.ApplicationGatewayProbeHealthResponseMatch{}
			if matchConfigs[0] != nil {
				match := matchConfigs[0].(map[string]interface{})
				matchBody = match["body"].(string)

				statusCodes := make([]string, 0)
				for _, statusCode := range match["status_code"].([]interface{}) {
					statusCodes = append(statusCodes, statusCode.(string))
				}
				outputMatch.StatusCodes = &statusCodes
			}
			outputMatch.Body = utils.String(matchBody)
			output.ApplicationGatewayProbePropertiesFormat.Match = outputMatch
		}

		if port != 0 {
			output.ApplicationGatewayProbePropertiesFormat.Port = utils.Int32(port)
		}

		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayProbes(input *[]network.ApplicationGatewayProbe) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.ID != nil {
			output["id"] = *v.ID
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.ApplicationGatewayProbePropertiesFormat; props != nil {
			output["protocol"] = string(props.Protocol)

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

func expandApplicationGatewayPrivateLinkConfigurations(d *pluginsdk.ResourceData) *[]network.ApplicationGatewayPrivateLinkConfiguration {
	vs := d.Get("private_link_configuration").(*pluginsdk.Set).List()
	plConfigResults := make([]network.ApplicationGatewayPrivateLinkConfiguration, 0)

	for _, rawPl := range vs {
		v := rawPl.(map[string]interface{})
		name := v["name"].(string)
		ipConfigurations := v["ip_configuration"].([]interface{})
		ipConfigurationResults := make([]network.ApplicationGatewayPrivateLinkIPConfiguration, 0)
		for _, rawIp := range ipConfigurations {
			v := rawIp.(map[string]interface{})
			name := v["name"].(string)
			subnetId := v["subnet_id"].(string)
			primary := v["primary"].(bool)
			ipConfiguration := network.ApplicationGatewayPrivateLinkIPConfiguration{
				Name: utils.String(name),
				ApplicationGatewayPrivateLinkIPConfigurationProperties: &network.ApplicationGatewayPrivateLinkIPConfigurationProperties{
					Primary: &primary,
					Subnet: &network.SubResource{
						ID: utils.String(subnetId),
					},
				},
			}
			if privateIpAddress := v["private_ip_address"].(string); privateIpAddress != "" {
				ipConfiguration.ApplicationGatewayPrivateLinkIPConfigurationProperties.PrivateIPAddress = utils.String(privateIpAddress)
			}
			if privateIpAddressAllocation := v["private_ip_address_allocation"].(string); privateIpAddressAllocation != "" {
				ipConfiguration.ApplicationGatewayPrivateLinkIPConfigurationProperties.PrivateIPAllocationMethod = network.IPAllocationMethod(privateIpAddressAllocation)
			}
			ipConfigurationResults = append(ipConfigurationResults, ipConfiguration)
		}

		configuration := network.ApplicationGatewayPrivateLinkConfiguration{
			Name: utils.String(name),
			ApplicationGatewayPrivateLinkConfigurationProperties: &network.ApplicationGatewayPrivateLinkConfigurationProperties{
				IPConfigurations: &ipConfigurationResults,
			},
		}
		plConfigResults = append(plConfigResults, configuration)
	}

	return &plConfigResults
}

func flattenApplicationGatewayPrivateEndpoints(input *[]network.ApplicationGatewayPrivateEndpointConnection) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, endpoint := range *input {
		result := map[string]interface{}{}
		if endpoint.Name != nil {
			result["name"] = *endpoint.Name
		}
		if endpoint.ID != nil {
			result["id"] = *endpoint.ID
		}
	}
	return results
}

func flattenApplicationGatewayPrivateLinkConfigurations(input *[]network.ApplicationGatewayPrivateLinkConfiguration) []interface{} {
	plConfigResults := make([]interface{}, 0)
	if input == nil {
		return plConfigResults
	}

	for _, plConfig := range *input {
		plConfigResult := map[string]interface{}{}
		if plConfig.Name != nil {
			plConfigResult["name"] = *plConfig.Name
		}
		if plConfig.ID != nil {
			plConfigResult["id"] = *plConfig.ID
		}
		ipConfigResults := make([]interface{}, 0)
		if props := plConfig.ApplicationGatewayPrivateLinkConfigurationProperties; props != nil {
			for _, ipConfig := range *props.IPConfigurations {
				ipConfigResult := map[string]interface{}{}
				if ipConfig.Name != nil {
					ipConfigResult["name"] = *ipConfig.Name
				}
				if subnet := ipConfig.Subnet; subnet != nil {
					if subnet.ID != nil {
						ipConfigResult["subnet_id"] = *subnet.ID
					}
				}
				if ipConfig.PrivateIPAddress != nil {
					ipConfigResult["private_ip_address"] = *ipConfig.PrivateIPAddress
				}
				ipConfigResult["private_ip_address_allocation"] = string(ipConfig.PrivateIPAllocationMethod)
				if ipConfig.Primary != nil {
					ipConfigResult["primary"] = *ipConfig.Primary
				}
				ipConfigResults = append(ipConfigResults, ipConfigResult)
			}
		}
		plConfigResult["ip_configuration"] = ipConfigResults
		plConfigResults = append(plConfigResults, plConfigResult)
	}
	return plConfigResults
}

func expandApplicationGatewayRequestRoutingRules(d *pluginsdk.ResourceData, gatewayID string) (*[]network.ApplicationGatewayRequestRoutingRule, error) {
	vs := d.Get("request_routing_rule").(*pluginsdk.Set).List()
	results := make([]network.ApplicationGatewayRequestRoutingRule, 0)
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
		priority := int32(v["priority"].(int))

		rule := network.ApplicationGatewayRequestRoutingRule{
			Name: utils.String(name),
			ApplicationGatewayRequestRoutingRulePropertiesFormat: &network.ApplicationGatewayRequestRoutingRulePropertiesFormat{
				RuleType: network.ApplicationGatewayRequestRoutingRuleType(ruleType),
				HTTPListener: &network.SubResource{
					ID: utils.String(httpListenerID),
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
			rule.ApplicationGatewayRequestRoutingRulePropertiesFormat.BackendAddressPool = &network.SubResource{
				ID: utils.String(backendAddressPoolID),
			}
		}

		if backendHTTPSettingsName != "" {
			backendHTTPSettingsID := fmt.Sprintf("%s/backendHttpSettingsCollection/%s", gatewayID, backendHTTPSettingsName)
			rule.ApplicationGatewayRequestRoutingRulePropertiesFormat.BackendHTTPSettings = &network.SubResource{
				ID: utils.String(backendHTTPSettingsID),
			}
		}

		if redirectConfigName != "" {
			redirectConfigID := fmt.Sprintf("%s/redirectConfigurations/%s", gatewayID, redirectConfigName)
			rule.ApplicationGatewayRequestRoutingRulePropertiesFormat.RedirectConfiguration = &network.SubResource{
				ID: utils.String(redirectConfigID),
			}
		}

		if urlPathMapName := v["url_path_map_name"].(string); urlPathMapName != "" {
			urlPathMapID := fmt.Sprintf("%s/urlPathMaps/%s", gatewayID, urlPathMapName)
			rule.ApplicationGatewayRequestRoutingRulePropertiesFormat.URLPathMap = &network.SubResource{
				ID: utils.String(urlPathMapID),
			}
		}

		if rewriteRuleSetName := v["rewrite_rule_set_name"].(string); rewriteRuleSetName != "" {
			rewriteRuleSetID := fmt.Sprintf("%s/rewriteRuleSets/%s", gatewayID, rewriteRuleSetName)
			rule.ApplicationGatewayRequestRoutingRulePropertiesFormat.RewriteRuleSet = &network.SubResource{
				ID: utils.String(rewriteRuleSetID),
			}
		}

		if priority != 0 {
			rule.ApplicationGatewayRequestRoutingRulePropertiesFormat.Priority = &priority
			priorityset = true
		}

		results = append(results, rule)
	}

	if priorityset {
		for _, rule := range results {
			if rule.ApplicationGatewayRequestRoutingRulePropertiesFormat.Priority == nil {
				return nil, fmt.Errorf("If you wish to use rule priority, you will have to specify rule-priority field values for all the existing request routing rules.")
			}
		}
	}

	return &results, nil
}

func flattenApplicationGatewayRequestRoutingRules(input *[]network.ApplicationGatewayRequestRoutingRule) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, config := range *input {
		if props := config.ApplicationGatewayRequestRoutingRulePropertiesFormat; props != nil {
			output := map[string]interface{}{
				"rule_type": string(props.RuleType),
			}

			if config.ID != nil {
				output["id"] = *config.ID
			}

			if config.Name != nil {
				output["name"] = *config.Name
			}

			if config.Priority != nil {
				output["priority"] = *config.Priority
			}

			if pool := props.BackendAddressPool; pool != nil {
				if pool.ID != nil {
					poolId, err := parse.BackendAddressPoolIDInsensitively(*pool.ID)
					if err != nil {
						return nil, err
					}
					output["backend_address_pool_name"] = poolId.Name
					output["backend_address_pool_id"] = poolId.ID()
				}
			}

			if settings := props.BackendHTTPSettings; settings != nil {
				if settings.ID != nil {
					settingsId, err := parse.BackendHttpSettingsCollectionIDInsensitively(*settings.ID)
					if err != nil {
						return nil, err
					}

					output["backend_http_settings_name"] = settingsId.BackendHttpSettingsCollectionName
					output["backend_http_settings_id"] = *settings.ID
				}
			}

			if listener := props.HTTPListener; listener != nil {
				if listener.ID != nil {
					listenerId, err := parse.HttpListenerIDInsensitively(*listener.ID)
					if err != nil {
						return nil, err
					}
					output["http_listener_id"] = listenerId.ID()
					output["http_listener_name"] = listenerId.Name
				}
			}

			if pathMap := props.URLPathMap; pathMap != nil {
				if pathMap.ID != nil {
					pathMapId, err := parse.UrlPathMapIDInsensitively(*pathMap.ID)
					if err != nil {
						return nil, err
					}
					output["url_path_map_name"] = pathMapId.Name
					output["url_path_map_id"] = pathMapId.ID()
				}
			}

			if redirect := props.RedirectConfiguration; redirect != nil {
				if redirect.ID != nil {
					redirectId, err := parse.RedirectConfigurationsIDInsensitively(*redirect.ID)
					if err != nil {
						return nil, err
					}
					output["redirect_configuration_name"] = redirectId.RedirectConfigurationName
					output["redirect_configuration_id"] = redirectId.ID()
				}
			}

			if rewrite := props.RewriteRuleSet; rewrite != nil {
				if rewrite.ID != nil {
					rewriteId, err := parse.RewriteRuleSetIDInsensitively(*rewrite.ID)
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

func expandApplicationGatewayRewriteRuleSets(d *pluginsdk.ResourceData) (*[]network.ApplicationGatewayRewriteRuleSet, error) {
	vs := d.Get("rewrite_rule_set").([]interface{})
	ruleSets := make([]network.ApplicationGatewayRewriteRuleSet, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})
		rules := make([]network.ApplicationGatewayRewriteRule, 0)

		name := v["name"].(string)

		for _, ruleConfig := range v["rewrite_rule"].([]interface{}) {
			r := ruleConfig.(map[string]interface{})
			conditions := make([]network.ApplicationGatewayRewriteRuleCondition, 0)
			requestConfigurations := make([]network.ApplicationGatewayHeaderConfiguration, 0)
			responseConfigurations := make([]network.ApplicationGatewayHeaderConfiguration, 0)
			urlConfiguration := network.ApplicationGatewayURLConfiguration{}

			rule := network.ApplicationGatewayRewriteRule{
				Name:         utils.String(r["name"].(string)),
				RuleSequence: utils.Int32(int32(r["rule_sequence"].(int))),
			}

			for _, rawCondition := range r["condition"].([]interface{}) {
				c := rawCondition.(map[string]interface{})
				condition := network.ApplicationGatewayRewriteRuleCondition{
					Variable:   utils.String(c["variable"].(string)),
					Pattern:    utils.String(c["pattern"].(string)),
					IgnoreCase: utils.Bool(c["ignore_case"].(bool)),
					Negate:     utils.Bool(c["negate"].(bool)),
				}
				conditions = append(conditions, condition)
			}
			rule.Conditions = &conditions

			for _, rawConfig := range r["request_header_configuration"].([]interface{}) {
				c := rawConfig.(map[string]interface{})
				config := network.ApplicationGatewayHeaderConfiguration{
					HeaderName:  utils.String(c["header_name"].(string)),
					HeaderValue: utils.String(c["header_value"].(string)),
				}
				requestConfigurations = append(requestConfigurations, config)
			}

			for _, rawConfig := range r["response_header_configuration"].([]interface{}) {
				c := rawConfig.(map[string]interface{})
				config := network.ApplicationGatewayHeaderConfiguration{
					HeaderName:  utils.String(c["header_name"].(string)),
					HeaderValue: utils.String(c["header_value"].(string)),
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
					urlConfiguration.ModifiedPath = utils.String(c["path"].(string))
				}
				if c["query_string"] != nil && components != "path_only" {
					urlConfiguration.ModifiedQueryString = utils.String(c["query_string"].(string))
				}
				if c["reroute"] != nil {
					urlConfiguration.Reroute = utils.Bool(c["reroute"].(bool))
				}
			}

			rule.ActionSet = &network.ApplicationGatewayRewriteRuleActionSet{
				RequestHeaderConfigurations:  &requestConfigurations,
				ResponseHeaderConfigurations: &responseConfigurations,
			}

			if len(r["url"].([]interface{})) > 0 {
				rule.ActionSet.URLConfiguration = &urlConfiguration
			}

			rules = append(rules, rule)
		}

		ruleSet := network.ApplicationGatewayRewriteRuleSet{
			Name: utils.String(name),
			ApplicationGatewayRewriteRuleSetPropertiesFormat: &network.ApplicationGatewayRewriteRuleSetPropertiesFormat{
				RewriteRules: &rules,
			},
		}

		ruleSets = append(ruleSets, ruleSet)
	}

	return &ruleSets, nil
}

func flattenApplicationGatewayRewriteRuleSets(input *[]network.ApplicationGatewayRewriteRuleSet) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, config := range *input {
		if props := config.ApplicationGatewayRewriteRuleSetPropertiesFormat; props != nil {
			output := map[string]interface{}{}

			if config.ID != nil {
				output["id"] = *config.ID
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

						if actionSet.URLConfiguration != nil {
							config := *actionSet.URLConfiguration
							components := ""
							path := ""
							if config.ModifiedPath != nil {
								path = *config.ModifiedPath
							}

							queryString := ""
							if config.ModifiedQueryString != nil {
								queryString = *config.ModifiedQueryString
							}

							if path != queryString {
								if path != "" && queryString == "" {
									components = "path_only"
								} else if queryString != "" && path == "" {
									components = "query_string_only"
								}
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

func expandApplicationGatewayRedirectConfigurations(d *pluginsdk.ResourceData, gatewayID string) (*[]network.ApplicationGatewayRedirectConfiguration, error) {
	vs := d.Get("redirect_configuration").(*pluginsdk.Set).List()
	results := make([]network.ApplicationGatewayRedirectConfiguration, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		redirectType := v["redirect_type"].(string)
		targetListenerName := v["target_listener_name"].(string)
		targetUrl := v["target_url"].(string)
		includePath := v["include_path"].(bool)
		includeQueryString := v["include_query_string"].(bool)

		output := network.ApplicationGatewayRedirectConfiguration{
			Name: utils.String(name),
			ApplicationGatewayRedirectConfigurationPropertiesFormat: &network.ApplicationGatewayRedirectConfigurationPropertiesFormat{
				RedirectType:       network.ApplicationGatewayRedirectType(redirectType),
				IncludeQueryString: utils.Bool(includeQueryString),
				IncludePath:        utils.Bool(includePath),
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
			output.ApplicationGatewayRedirectConfigurationPropertiesFormat.TargetListener = &network.SubResource{
				ID: utils.String(targetListenerID),
			}
		}

		if targetUrl != "" {
			output.ApplicationGatewayRedirectConfigurationPropertiesFormat.TargetURL = utils.String(targetUrl)
		}

		results = append(results, output)
	}

	return &results, nil
}

func flattenApplicationGatewayRedirectConfigurations(input *[]network.ApplicationGatewayRedirectConfiguration) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, config := range *input {
		if props := config.ApplicationGatewayRedirectConfigurationPropertiesFormat; props != nil {
			output := map[string]interface{}{
				"redirect_type": string(props.RedirectType),
			}

			if config.ID != nil {
				output["id"] = *config.ID
			}

			if config.Name != nil {
				output["name"] = *config.Name
			}

			if listener := props.TargetListener; listener != nil {
				if listener.ID != nil {
					listenerId, err := parse.HttpListenerIDInsensitively(*listener.ID)
					if err != nil {
						return nil, err
					}
					output["target_listener_name"] = listenerId.Name
					output["target_listener_id"] = listenerId.ID()
				}
			}

			if config.TargetURL != nil {
				output["target_url"] = *config.TargetURL
			}

			if config.IncludePath != nil {
				output["include_path"] = *config.IncludePath
			}

			if config.IncludeQueryString != nil {
				output["include_query_string"] = *config.IncludeQueryString
			}

			results = append(results, output)
		}
	}

	return results, nil
}

func expandApplicationGatewayAutoscaleConfiguration(d *pluginsdk.ResourceData) *network.ApplicationGatewayAutoscaleConfiguration {
	vs := d.Get("autoscale_configuration").([]interface{})
	if len(vs) == 0 {
		return nil
	}
	v := vs[0].(map[string]interface{})

	minCapacity := int32(v["min_capacity"].(int))
	maxCapacity := int32(v["max_capacity"].(int))

	configuration := network.ApplicationGatewayAutoscaleConfiguration{
		MinCapacity: utils.Int32(minCapacity),
	}

	if maxCapacity != 0 {
		configuration.MaxCapacity = utils.Int32(maxCapacity)
	}

	return &configuration
}

func flattenApplicationGatewayAutoscaleConfiguration(input *network.ApplicationGatewayAutoscaleConfiguration) []interface{} {
	result := make(map[string]interface{})
	if input == nil {
		return []interface{}{}
	}
	if v := input.MinCapacity; v != nil {
		result["min_capacity"] = *v
	}
	if input.MaxCapacity != nil {
		result["max_capacity"] = *input.MaxCapacity
	}

	return []interface{}{result}
}

func expandApplicationGatewaySku(d *pluginsdk.ResourceData) *network.ApplicationGatewaySku {
	vs := d.Get("sku").([]interface{})
	v := vs[0].(map[string]interface{})

	name := v["name"].(string)
	tier := v["tier"].(string)
	capacity := int32(v["capacity"].(int))

	sku := network.ApplicationGatewaySku{
		Name: network.ApplicationGatewaySkuName(name),
		Tier: network.ApplicationGatewayTier(tier),
	}

	if capacity != 0 {
		sku.Capacity = utils.Int32(capacity)
	}

	return &sku
}

func flattenApplicationGatewaySku(input *network.ApplicationGatewaySku) []interface{} {
	result := make(map[string]interface{})

	result["name"] = string(input.Name)
	result["tier"] = string(input.Tier)
	if input.Capacity != nil {
		result["capacity"] = int(*input.Capacity)
	}

	return []interface{}{result}
}

func expandApplicationGatewaySslCertificates(d *pluginsdk.ResourceData) (*[]network.ApplicationGatewaySslCertificate, error) {
	vs := d.Get("ssl_certificate").(*schema.Set).List()
	results := make([]network.ApplicationGatewaySslCertificate, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		data := v["data"].(string)
		password := v["password"].(string)
		kvsid := v["key_vault_secret_id"].(string)
		cert := v["public_cert_data"].(string)

		output := network.ApplicationGatewaySslCertificate{
			Name: utils.String(name),
			ApplicationGatewaySslCertificatePropertiesFormat: &network.ApplicationGatewaySslCertificatePropertiesFormat{},
		}

		// nolint gocritic
		if data != "" && kvsid != "" {
			return nil, fmt.Errorf("only one of `key_vault_secret_id` or `data` must be specified for the `ssl_certificate` block %q", name)
		} else if data != "" {
			// data must be base64 encoded
			output.ApplicationGatewaySslCertificatePropertiesFormat.Data = utils.String(utils.Base64EncodeIfNot(data))

			output.ApplicationGatewaySslCertificatePropertiesFormat.Password = utils.String(password)
		} else if kvsid != "" {
			if password != "" {
				return nil, fmt.Errorf("only one of `key_vault_secret_id` or `password` must be specified for the `ssl_certificate` block %q", name)
			}

			output.ApplicationGatewaySslCertificatePropertiesFormat.KeyVaultSecretID = utils.String(kvsid)
		} else if cert != "" {
			output.ApplicationGatewaySslCertificatePropertiesFormat.PublicCertData = utils.String(cert)
		} else {
			return nil, fmt.Errorf("either `key_vault_secret_id` or `data` must be specified for the `ssl_certificate` block %q", name)
		}

		results = append(results, output)
	}

	return &results, nil
}

func flattenApplicationGatewaySslCertificates(input *[]network.ApplicationGatewaySslCertificate, d *pluginsdk.ResourceData) []interface{} {
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

		if v.ID != nil {
			output["id"] = *v.ID
		}

		output["name"] = name

		if props := v.ApplicationGatewaySslCertificatePropertiesFormat; props != nil {
			if data := props.PublicCertData; data != nil {
				output["public_cert_data"] = *data
			}

			if kvsid := props.KeyVaultSecretID; kvsid != nil {
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

func expandApplicationGatewayTrustedClientCertificates(d *pluginsdk.ResourceData) (*[]network.ApplicationGatewayTrustedClientCertificate, error) {
	vs := d.Get("trusted_client_certificate").([]interface{})
	results := make([]network.ApplicationGatewayTrustedClientCertificate, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		data := v["data"].(string)

		output := network.ApplicationGatewayTrustedClientCertificate{
			Name: utils.String(name),
			ApplicationGatewayTrustedClientCertificatePropertiesFormat: &network.ApplicationGatewayTrustedClientCertificatePropertiesFormat{},
		}

		// nolint gocritic
		if data != "" {
			// data must be base64 encoded
			output.ApplicationGatewayTrustedClientCertificatePropertiesFormat.Data = utils.String(utils.Base64EncodeIfNot(data))
		} else {
			return nil, fmt.Errorf("`data` must be specified for the `trusted_client_certificate` block %q", name)
		}

		results = append(results, output)
	}

	return &results, nil
}

func flattenApplicationGatewayTrustedClientCertificates(input *[]network.ApplicationGatewayTrustedClientCertificate) []interface{} {
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

		if v.ID != nil {
			output["id"] = *v.ID
		}

		output["name"] = name

		if props := v.ApplicationGatewayTrustedClientCertificatePropertiesFormat; props != nil {
			if data := props.Data; data != nil {
				output["data"] = *data
			}
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewaySslProfiles(d *pluginsdk.ResourceData, gatewayID string) *[]network.ApplicationGatewaySslProfile {
	vs := d.Get("ssl_profile").([]interface{})
	results := make([]network.ApplicationGatewaySslProfile, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		verifyClientCertIssuerDn := v["verify_client_cert_issuer_dn"].(bool)

		output := network.ApplicationGatewaySslProfile{
			Name: utils.String(name),
			ApplicationGatewaySslProfilePropertiesFormat: &network.ApplicationGatewaySslProfilePropertiesFormat{
				ClientAuthConfiguration: &network.ApplicationGatewayClientAuthConfiguration{VerifyClientCertIssuerDN: utils.Bool(verifyClientCertIssuerDn)},
			},
		}

		if v["trusted_client_certificate_names"] != nil {
			clientCerts := v["trusted_client_certificate_names"].([]interface{})
			clientCertSubResources := make([]network.SubResource, 0)

			for _, rawClientCert := range clientCerts {
				clientCertName := rawClientCert
				clientCertID := fmt.Sprintf("%s/trustedClientCertificates/%s", gatewayID, clientCertName)
				clientCertSubResource := network.SubResource{
					ID: utils.String(clientCertID),
				}
				clientCertSubResources = append(clientCertSubResources, clientCertSubResource)
			}
			output.ApplicationGatewaySslProfilePropertiesFormat.TrustedClientCertificates = &clientCertSubResources
		}

		sslPolicy := v["ssl_policy"].([]interface{})
		if len(sslPolicy) > 0 {
			output.ApplicationGatewaySslProfilePropertiesFormat.SslPolicy = expandApplicationGatewaySslPolicy(sslPolicy)
		} else {
			output.ApplicationGatewaySslProfilePropertiesFormat.SslPolicy = nil
		}
		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewaySslProfiles(input *[]network.ApplicationGatewaySslProfile) ([]interface{}, error) {
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

		if v.ID != nil {
			output["id"] = *v.ID
		}

		output["name"] = name
		output["verify_client_cert_issuer_dn"] = *v.ClientAuthConfiguration.VerifyClientCertIssuerDN

		output["ssl_policy"] = flattenApplicationGatewaySslPolicy(v.SslPolicy)

		if props := v.ApplicationGatewaySslProfilePropertiesFormat; props != nil {
			trustedClientCertificateNames := make([]interface{}, 0)
			if certs := props.TrustedClientCertificates; certs != nil {
				for _, cert := range *certs {
					if cert.ID == nil {
						continue
					}

					certId, err := parse.TrustedClientCertificateIDInsensitively(*cert.ID)
					if err != nil {
						return nil, err
					}

					trustedClientCertificateNames = append(trustedClientCertificateNames, certId.Name)
				}
			}
			output["trusted_client_certificate_names"] = trustedClientCertificateNames
		}

		results = append(results, output)
	}

	return results, nil
}

func expandApplicationGatewayURLPathMaps(d *pluginsdk.ResourceData, gatewayID string) (*[]network.ApplicationGatewayURLPathMap, error) {
	vs := d.Get("url_path_map").([]interface{})
	results := make([]network.ApplicationGatewayURLPathMap, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)

		pathRules := make([]network.ApplicationGatewayPathRule, 0)
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

			rule := network.ApplicationGatewayPathRule{
				Name: utils.String(ruleName),
				ApplicationGatewayPathRulePropertiesFormat: &network.ApplicationGatewayPathRulePropertiesFormat{
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
				rule.ApplicationGatewayPathRulePropertiesFormat.BackendAddressPool = &network.SubResource{
					ID: utils.String(backendAddressPoolID),
				}
			}

			if backendHTTPSettingsName != "" {
				backendHTTPSettingsID := fmt.Sprintf("%s/backendHttpSettingsCollection/%s", gatewayID, backendHTTPSettingsName)
				rule.ApplicationGatewayPathRulePropertiesFormat.BackendHTTPSettings = &network.SubResource{
					ID: utils.String(backendHTTPSettingsID),
				}
			}

			if redirectConfigurationName != "" {
				redirectConfigurationID := fmt.Sprintf("%s/redirectConfigurations/%s", gatewayID, redirectConfigurationName)
				rule.ApplicationGatewayPathRulePropertiesFormat.RedirectConfiguration = &network.SubResource{
					ID: utils.String(redirectConfigurationID),
				}
			}

			if rewriteRuleSetName := ruleConfigMap["rewrite_rule_set_name"].(string); rewriteRuleSetName != "" {
				rewriteRuleSetID := fmt.Sprintf("%s/rewriteRuleSets/%s", gatewayID, rewriteRuleSetName)
				rule.ApplicationGatewayPathRulePropertiesFormat.RewriteRuleSet = &network.SubResource{
					ID: utils.String(rewriteRuleSetID),
				}
			}

			if firewallPolicyID != "" && len(firewallPolicyID) > 0 {
				rule.ApplicationGatewayPathRulePropertiesFormat.FirewallPolicy = &network.SubResource{
					ID: utils.String(firewallPolicyID),
				}
			}

			pathRules = append(pathRules, rule)
		}

		output := network.ApplicationGatewayURLPathMap{
			Name: utils.String(name),
			ApplicationGatewayURLPathMapPropertiesFormat: &network.ApplicationGatewayURLPathMapPropertiesFormat{
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
			output.ApplicationGatewayURLPathMapPropertiesFormat.DefaultBackendAddressPool = &network.SubResource{
				ID: utils.String(defaultBackendAddressPoolID),
			}
		}

		if defaultBackendHTTPSettingsName != "" {
			defaultBackendHTTPSettingsID := fmt.Sprintf("%s/backendHttpSettingsCollection/%s", gatewayID, defaultBackendHTTPSettingsName)
			output.ApplicationGatewayURLPathMapPropertiesFormat.DefaultBackendHTTPSettings = &network.SubResource{
				ID: utils.String(defaultBackendHTTPSettingsID),
			}
		}

		if defaultRedirectConfigurationName != "" {
			defaultRedirectConfigurationID := fmt.Sprintf("%s/redirectConfigurations/%s", gatewayID, defaultRedirectConfigurationName)
			output.ApplicationGatewayURLPathMapPropertiesFormat.DefaultRedirectConfiguration = &network.SubResource{
				ID: utils.String(defaultRedirectConfigurationID),
			}
		}

		if defaultRewriteRuleSetName := v["default_rewrite_rule_set_name"].(string); defaultRewriteRuleSetName != "" {
			defaultRewriteRuleSetID := fmt.Sprintf("%s/rewriteRuleSets/%s", gatewayID, defaultRewriteRuleSetName)
			output.ApplicationGatewayURLPathMapPropertiesFormat.DefaultRewriteRuleSet = &network.SubResource{
				ID: utils.String(defaultRewriteRuleSetID),
			}
		}

		results = append(results, output)
	}

	return &results, nil
}

func flattenApplicationGatewayURLPathMaps(input *[]network.ApplicationGatewayURLPathMap) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil {
		return results, nil
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		if v.ID != nil {
			output["id"] = *v.ID
		}

		if v.Name != nil {
			output["name"] = *v.Name
		}

		if props := v.ApplicationGatewayURLPathMapPropertiesFormat; props != nil {
			if backendPool := props.DefaultBackendAddressPool; backendPool != nil && backendPool.ID != nil {
				poolId, err := parse.BackendAddressPoolIDInsensitively(*backendPool.ID)
				if err != nil {
					return nil, err
				}
				output["default_backend_address_pool_name"] = poolId.Name
				output["default_backend_address_pool_id"] = poolId.ID()
			}

			if settings := props.DefaultBackendHTTPSettings; settings != nil && settings.ID != nil {
				settingsId, err := parse.BackendHttpSettingsCollectionIDInsensitively(*settings.ID)
				if err != nil {
					return nil, err
				}
				output["default_backend_http_settings_name"] = settingsId.BackendHttpSettingsCollectionName
				output["default_backend_http_settings_id"] = settingsId.ID()
			}

			if redirect := props.DefaultRedirectConfiguration; redirect != nil && redirect.ID != nil {
				redirectId, err := parse.RedirectConfigurationsIDInsensitively(*redirect.ID)
				if err != nil {
					return nil, err
				}
				output["default_redirect_configuration_name"] = redirectId.RedirectConfigurationName
				output["default_redirect_configuration_id"] = redirectId.ID()
			}

			if rewrite := props.DefaultRewriteRuleSet; rewrite != nil && rewrite.ID != nil {
				rewriteId, err := parse.RewriteRuleSetIDInsensitively(*rewrite.ID)
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

					if rule.ID != nil {
						ruleOutput["id"] = *rule.ID
					}

					if rule.Name != nil {
						ruleOutput["name"] = *rule.Name
					}

					if ruleProps := rule.ApplicationGatewayPathRulePropertiesFormat; ruleProps != nil {
						if pool := ruleProps.BackendAddressPool; pool != nil && pool.ID != nil {
							poolId, err := parse.BackendAddressPoolIDInsensitively(*pool.ID)
							if err != nil {
								return nil, err
							}
							ruleOutput["backend_address_pool_name"] = poolId.Name
							ruleOutput["backend_address_pool_id"] = poolId.ID()
						}

						if backend := ruleProps.BackendHTTPSettings; backend != nil && backend.ID != nil {
							backendId, err := parse.BackendHttpSettingsCollectionIDInsensitively(*backend.ID)
							if err != nil {
								return nil, err
							}
							ruleOutput["backend_http_settings_name"] = backendId.BackendHttpSettingsCollectionName
							ruleOutput["backend_http_settings_id"] = backendId.ID()
						}

						if redirect := ruleProps.RedirectConfiguration; redirect != nil && redirect.ID != nil {
							redirectId, err := parse.RedirectConfigurationsIDInsensitively(*redirect.ID)
							if err != nil {
								return nil, err
							}
							ruleOutput["redirect_configuration_name"] = redirectId.RedirectConfigurationName
							ruleOutput["redirect_configuration_id"] = redirectId.ID()
						}

						if rewrite := ruleProps.RewriteRuleSet; rewrite != nil && rewrite.ID != nil {
							rewriteId, err := parse.RewriteRuleSetIDInsensitively(*rewrite.ID)
							if err != nil {
								return nil, err
							}
							ruleOutput["rewrite_rule_set_name"] = rewriteId.Name
							ruleOutput["rewrite_rule_set_id"] = rewriteId.ID()
						}

						if fwp := ruleProps.FirewallPolicy; fwp != nil && fwp.ID != nil {
							policyId, err := webapplicationfirewallpolicies.ParseApplicationGatewayWebApplicationFirewallPolicyIDInsensitively(*fwp.ID)
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

func expandApplicationGatewayWafConfig(d *pluginsdk.ResourceData) *network.ApplicationGatewayWebApplicationFirewallConfiguration {
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

	return &network.ApplicationGatewayWebApplicationFirewallConfiguration{
		Enabled:                utils.Bool(enabled),
		FirewallMode:           network.ApplicationGatewayFirewallMode(mode),
		RuleSetType:            utils.String(ruleSetType),
		RuleSetVersion:         utils.String(ruleSetVersion),
		FileUploadLimitInMb:    utils.Int32(int32(fileUploadLimitInMb)),
		RequestBodyCheck:       utils.Bool(requestBodyCheck),
		MaxRequestBodySizeInKb: utils.Int32(int32(maxRequestBodySizeInKb)),
		DisabledRuleGroups:     expandApplicationGatewayFirewallDisabledRuleGroup(v["disabled_rule_group"].([]interface{})),
		Exclusions:             expandApplicationGatewayFirewallExclusion(v["exclusion"].([]interface{})),
	}
}

func flattenApplicationGatewayWafConfig(input *network.ApplicationGatewayWebApplicationFirewallConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	output := make(map[string]interface{})

	if input.Enabled != nil {
		output["enabled"] = *input.Enabled
	}

	output["firewall_mode"] = string(input.FirewallMode)

	if input.RuleSetType != nil {
		output["rule_set_type"] = *input.RuleSetType
	}

	if input.RuleSetVersion != nil {
		output["rule_set_version"] = *input.RuleSetVersion
	}

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

func expandApplicationGatewayFirewallDisabledRuleGroup(d []interface{}) *[]network.ApplicationGatewayFirewallDisabledRuleGroup {
	if len(d) == 0 {
		return nil
	}

	disabledRuleGroups := make([]network.ApplicationGatewayFirewallDisabledRuleGroup, 0)
	for _, disabledRuleGroup := range d {
		disabledRuleGroupMap := disabledRuleGroup.(map[string]interface{})

		ruleGroupName := disabledRuleGroupMap["rule_group_name"].(string)

		ruleGroup := network.ApplicationGatewayFirewallDisabledRuleGroup{
			RuleGroupName: utils.String(ruleGroupName),
		}

		rules := make([]int32, 0)
		for _, rule := range disabledRuleGroupMap["rules"].([]interface{}) {
			rules = append(rules, int32(rule.(int)))
		}

		if len(rules) > 0 {
			ruleGroup.Rules = &rules
		}

		disabledRuleGroups = append(disabledRuleGroups, ruleGroup)
	}
	return &disabledRuleGroups
}

func flattenApplicationGateWayDisabledRuleGroups(input *[]network.ApplicationGatewayFirewallDisabledRuleGroup) []interface{} {
	ruleGroups := make([]interface{}, 0)
	for _, ruleGroup := range *input {
		ruleGroupOutput := map[string]interface{}{}

		if ruleGroup.RuleGroupName != nil {
			ruleGroupOutput["rule_group_name"] = *ruleGroup.RuleGroupName
		}

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

func expandApplicationGatewayFirewallExclusion(d []interface{}) *[]network.ApplicationGatewayFirewallExclusion {
	if len(d) == 0 {
		return nil
	}

	exclusions := make([]network.ApplicationGatewayFirewallExclusion, 0)
	for _, exclusion := range d {
		exclusionMap := exclusion.(map[string]interface{})

		matchVariable := exclusionMap["match_variable"].(string)
		selectorMatchOperator := exclusionMap["selector_match_operator"].(string)
		selector := exclusionMap["selector"].(string)

		exclusionList := network.ApplicationGatewayFirewallExclusion{
			MatchVariable:         utils.String(matchVariable),
			SelectorMatchOperator: utils.String(selectorMatchOperator),
			Selector:              utils.String(selector),
		}

		exclusions = append(exclusions, exclusionList)
	}

	return &exclusions
}

func flattenApplicationGatewayFirewallExclusion(input *[]network.ApplicationGatewayFirewallExclusion) []interface{} {
	exclusionLists := make([]interface{}, 0)
	for _, exclusionList := range *input {
		exclusionListOutput := map[string]interface{}{}

		if exclusionList.MatchVariable != nil {
			exclusionListOutput["match_variable"] = *exclusionList.MatchVariable
		}

		if exclusionList.SelectorMatchOperator != nil {
			exclusionListOutput["selector_match_operator"] = *exclusionList.SelectorMatchOperator
		}

		if exclusionList.Selector != nil {
			exclusionListOutput["selector"] = *exclusionList.Selector
		}
		exclusionLists = append(exclusionLists, exclusionListOutput)
	}
	return exclusionLists
}

func expandApplicationGatewayCustomErrorConfigurations(vs []interface{}) *[]network.ApplicationGatewayCustomError {
	results := make([]network.ApplicationGatewayCustomError, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})
		statusCode := v["status_code"].(string)
		customErrorPageUrl := v["custom_error_page_url"].(string)

		output := network.ApplicationGatewayCustomError{
			StatusCode:         network.ApplicationGatewayCustomErrorStatusCode(statusCode),
			CustomErrorPageURL: utils.String(customErrorPageUrl),
		}
		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayCustomErrorConfigurations(input *[]network.ApplicationGatewayCustomError) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		output := map[string]interface{}{}

		output["status_code"] = string(v.StatusCode)

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

func applicationGatewayCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
	_, hasAutoscaleConfig := d.GetOk("autoscale_configuration.0")
	capacity, hasCapacity := d.GetOk("sku.0.capacity")
	tier := d.Get("sku.0.tier").(string)

	if !hasAutoscaleConfig && !hasCapacity {
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
		if (strings.EqualFold(tier, string(network.ApplicationGatewayTierStandard)) || strings.EqualFold(tier, string(network.ApplicationGatewayTierWAF))) && (capacity.(int) < 1 || capacity.(int) > 32) {
			return fmt.Errorf("The value '%d' exceeds the maximum capacity allowed for a %q V1 SKU, the %q SKU must have a capacity value between 1 and 32", capacity, tier, tier)
		}

		if (strings.EqualFold(tier, string(network.ApplicationGatewayTierStandardV2)) || strings.EqualFold(tier, string(network.ApplicationGatewayTierWAFV2))) && (capacity.(int) < 1 || capacity.(int) > 125) {
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
