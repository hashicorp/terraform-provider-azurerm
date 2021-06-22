package network

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	keyVaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	msiParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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

func resourceApplicationGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApplicationGatewayCreateUpdate,
		Read:   resourceApplicationGatewayRead,
		Update: resourceApplicationGatewayCreateUpdate,
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

			"location": azure.SchemaLocation(),

			"zones": azure.SchemaZones(),

			"resource_group_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
			"identity": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(network.ResourceIdentityTypeUserAssigned),
							ValidateFunc: validation.StringInSlice([]string{
								string(network.ResourceIdentityTypeUserAssigned),
							}, false),
						},
						"identity_ids": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.NoZeroValues,
							},
						},
					},
				},
			},

			// Required
			"backend_address_pool": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"fqdns": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MinItems: 1,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.NoZeroValues,
							},
						},

						"ip_addresses": {
							Type:     pluginsdk.TypeList,
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
			},

			"backend_http_settings": {
				Type:     pluginsdk.TypeList,
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
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.ProtocolHTTP),
								string(network.ProtocolHTTPS),
							}, true),
						},

						"cookie_based_affinity": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.ApplicationGatewayCookieBasedAffinityEnabled),
								string(network.ApplicationGatewayCookieBasedAffinityDisabled),
							}, true),
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
							Computed: true,
						},

						"private_ip_address": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
						},

						"public_ip_address_id": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
						},

						"private_ip_address_allocation": {
							Type:             pluginsdk.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.IPAllocationMethodDynamic),
								string(network.IPAllocationMethodStatic),
							}, true),
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
							ValidateFunc: azure.ValidateResourceID,
						},

						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"http_listener": {
				Type:     pluginsdk.TypeList,
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
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.ProtocolHTTP),
								string(network.ProtocolHTTPS),
							}, true),
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
							ValidateFunc: azure.ValidateResourceID,
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
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.ApplicationGatewaySkuNameStandardSmall),
								string(network.ApplicationGatewaySkuNameStandardMedium),
								string(network.ApplicationGatewaySkuNameStandardLarge),
								string(network.ApplicationGatewaySkuNameStandardV2),
								string(network.ApplicationGatewaySkuNameWAFLarge),
								string(network.ApplicationGatewaySkuNameWAFMedium),
								string(network.ApplicationGatewaySkuNameWAFV2),
							}, true),
						},

						"tier": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.ApplicationGatewayTierStandard),
								string(network.ApplicationGatewayTierStandardV2),
								string(network.ApplicationGatewayTierWAF),
								string(network.ApplicationGatewayTierWAFV2),
							}, true),
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
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							Sensitive:    true,
						},

						// TODO required soft delete on the keyvault
						/*"key_vault_secret_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateKeyVaultChildId,
						},*/

						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			//lintignore:XS003
			"ssl_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
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
								}, false),
							},
						},

						"policy_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.ApplicationGatewaySslPolicyTypeCustom),
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
							}, false),
						},
					},
				},
			},

			"enable_http2": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"probe": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"protocol": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.ProtocolHTTP),
								string(network.ProtocolHTTPS),
							}, true),
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

						//lintignore:XS003
						"match": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"body": {
										Type:     pluginsdk.TypeString,
										Optional: true,
									},

									"status_code": {
										Type:     pluginsdk.TypeList,
										Optional: true,
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

			"ssl_certificate": {
				// TODO: should this become a Set?
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
							Optional:  true,
							Sensitive: true,
							StateFunc: base64EncodedStateFunc,
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
										ValidateFunc: networkValidate.ApplicationGatewayWebApplicationFirewallPolicyID,
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
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.ApplicationGatewayFirewallModeDetection),
								string(network.ApplicationGatewayFirewallModePrevention),
							}, true),
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
											string(network.OwaspCrsExclusionEntryMatchVariableRequestArgNames),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestCookieNames),
											string(network.OwaspCrsExclusionEntryMatchVariableRequestHeaderNames),
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
				ValidateFunc: azure.ValidateResourceID,
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

func resourceApplicationGatewayCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
	trustedRootCertificates := expandApplicationGatewayTrustedRootCertificates(d.Get("trusted_root_certificate").([]interface{}))

	requestRoutingRules, err := expandApplicationGatewayRequestRoutingRules(d, id.ID())
	if err != nil {
		return fmt.Errorf("Error expanding `request_routing_rule`: %+v", err)
	}

	urlPathMaps, err := expandApplicationGatewayURLPathMaps(d, id.ID())
	if err != nil {
		return fmt.Errorf("Error expanding `url_path_map`: %+v", err)
	}

	redirectConfigurations, err := expandApplicationGatewayRedirectConfigurations(d, id.ID())
	if err != nil {
		return fmt.Errorf("Error expanding `redirect_configuration`: %+v", err)
	}

	sslCertificates, err := expandApplicationGatewaySslCertificates(d)
	if err != nil {
		return fmt.Errorf("Error expanding `ssl_certificate`: %+v", err)
	}

	gatewayIPConfigurations, stopApplicationGateway := expandApplicationGatewayIPConfigurations(d)

	httpListeners, err := expandApplicationGatewayHTTPListeners(d, id.ID())
	if err != nil {
		return fmt.Errorf("fail to expand `http_listener`: %+v", err)
	}

	rewriteRuleSets, err := expandApplicationGatewayRewriteRuleSets(d)
	if err != nil {
		return fmt.Errorf("error expanding `rewrite_rule_set`: %v", err)
	}

	gateway := network.ApplicationGateway{
		Location: utils.String(location),
		Zones:    azure.ExpandZones(d.Get("zones").([]interface{})),

		Tags: tags.Expand(t),
		ApplicationGatewayPropertiesFormat: &network.ApplicationGatewayPropertiesFormat{
			AutoscaleConfiguration:        expandApplicationGatewayAutoscaleConfiguration(d),
			AuthenticationCertificates:    expandApplicationGatewayAuthenticationCertificates(d.Get("authentication_certificate").([]interface{})),
			TrustedRootCertificates:       trustedRootCertificates,
			CustomErrorConfigurations:     expandApplicationGatewayCustomErrorConfigurations(d.Get("custom_error_configuration").([]interface{})),
			BackendAddressPools:           expandApplicationGatewayBackendAddressPools(d),
			BackendHTTPSettingsCollection: expandApplicationGatewayBackendHTTPSettings(d, id.ID()),
			EnableHTTP2:                   utils.Bool(enablehttp2),
			FrontendIPConfigurations:      expandApplicationGatewayFrontendIPConfigurations(d),
			FrontendPorts:                 expandApplicationGatewayFrontendPorts(d),
			GatewayIPConfigurations:       gatewayIPConfigurations,
			HTTPListeners:                 httpListeners,
			Probes:                        expandApplicationGatewayProbes(d),
			RequestRoutingRules:           requestRoutingRules,
			RedirectConfigurations:        redirectConfigurations,
			Sku:                           expandApplicationGatewaySku(d),
			SslCertificates:               sslCertificates,
			SslPolicy:                     expandApplicationGatewaySslPolicy(d),

			RewriteRuleSets: rewriteRuleSets,
			URLPathMaps:     urlPathMaps,
		},
	}

	if _, ok := d.GetOk("identity"); ok {
		gateway.Identity = expandAzureRmApplicationGatewayIdentity(d)
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
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for create/update of %s: %+v", id, err)
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

	d.Set("name", applicationGateway.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := applicationGateway.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("zones", applicationGateway.Zones)

	identity, err := flattenRmApplicationGatewayIdentity(applicationGateway.Identity)
	if err != nil {
		return err
	}
	if err = d.Set("identity", identity); err != nil {
		return err
	}

	if props := applicationGateway.ApplicationGatewayPropertiesFormat; props != nil {
		if err = d.Set("authentication_certificate", flattenApplicationGatewayAuthenticationCertificates(props.AuthenticationCertificates, d)); err != nil {
			return fmt.Errorf("Error setting `authentication_certificate`: %+v", err)
		}

		if err = d.Set("trusted_root_certificate", flattenApplicationGatewayTrustedRootCertificates(props.TrustedRootCertificates, d)); err != nil {
			return fmt.Errorf("Error setting `trusted_root_certificate`: %+v", err)
		}

		if setErr := d.Set("backend_address_pool", flattenApplicationGatewayBackendAddressPools(props.BackendAddressPools)); setErr != nil {
			return fmt.Errorf("Error setting `backend_address_pool`: %+v", setErr)
		}

		backendHttpSettings, err := flattenApplicationGatewayBackendHTTPSettings(props.BackendHTTPSettingsCollection)
		if err != nil {
			return fmt.Errorf("Error flattening `backend_http_settings`: %+v", err)
		}
		if setErr := d.Set("backend_http_settings", backendHttpSettings); setErr != nil {
			return fmt.Errorf("Error setting `backend_http_settings`: %+v", setErr)
		}

		if setErr := d.Set("ssl_policy", flattenApplicationGatewaySslPolicy(props.SslPolicy)); setErr != nil {
			return fmt.Errorf("Error setting `ssl_policy`: %+v", setErr)
		}

		d.Set("enable_http2", props.EnableHTTP2)

		httpListeners, err := flattenApplicationGatewayHTTPListeners(props.HTTPListeners)
		if err != nil {
			return fmt.Errorf("Error flattening `http_listener`: %+v", err)
		}
		if setErr := d.Set("http_listener", httpListeners); setErr != nil {
			return fmt.Errorf("Error setting `http_listener`: %+v", setErr)
		}

		if setErr := d.Set("frontend_port", flattenApplicationGatewayFrontendPorts(props.FrontendPorts)); setErr != nil {
			return fmt.Errorf("Error setting `frontend_port`: %+v", setErr)
		}

		if setErr := d.Set("frontend_ip_configuration", flattenApplicationGatewayFrontendIPConfigurations(props.FrontendIPConfigurations)); setErr != nil {
			return fmt.Errorf("Error setting `frontend_ip_configuration`: %+v", setErr)
		}

		if setErr := d.Set("gateway_ip_configuration", flattenApplicationGatewayIPConfigurations(props.GatewayIPConfigurations)); setErr != nil {
			return fmt.Errorf("Error setting `gateway_ip_configuration`: %+v", setErr)
		}

		if setErr := d.Set("probe", flattenApplicationGatewayProbes(props.Probes)); setErr != nil {
			return fmt.Errorf("Error setting `probe`: %+v", setErr)
		}

		requestRoutingRules, err := flattenApplicationGatewayRequestRoutingRules(props.RequestRoutingRules)
		if err != nil {
			return fmt.Errorf("Error flattening `request_routing_rule`: %+v", err)
		}
		if setErr := d.Set("request_routing_rule", requestRoutingRules); setErr != nil {
			return fmt.Errorf("Error setting `request_routing_rule`: %+v", setErr)
		}

		redirectConfigurations, err := flattenApplicationGatewayRedirectConfigurations(props.RedirectConfigurations)
		if err != nil {
			return fmt.Errorf("Error flattening `redirect configuration`: %+v", err)
		}
		if setErr := d.Set("redirect_configuration", redirectConfigurations); setErr != nil {
			return fmt.Errorf("Error setting `redirect configuration`: %+v", setErr)
		}

		rewriteRuleSets := flattenApplicationGatewayRewriteRuleSets(props.RewriteRuleSets)
		if setErr := d.Set("rewrite_rule_set", rewriteRuleSets); setErr != nil {
			return fmt.Errorf("Error setting `rewrite_rule_set`: %+v", setErr)
		}

		if setErr := d.Set("sku", flattenApplicationGatewaySku(props.Sku)); setErr != nil {
			return fmt.Errorf("Error setting `sku`: %+v", setErr)
		}

		if setErr := d.Set("autoscale_configuration", flattenApplicationGatewayAutoscaleConfiguration(props.AutoscaleConfiguration)); setErr != nil {
			return fmt.Errorf("Error setting `autoscale_configuration`: %+v", setErr)
		}

		if setErr := d.Set("ssl_certificate", flattenApplicationGatewaySslCertificates(props.SslCertificates, d)); setErr != nil {
			return fmt.Errorf("Error setting `ssl_certificate`: %+v", setErr)
		}

		if setErr := d.Set("custom_error_configuration", flattenApplicationGatewayCustomErrorConfigurations(props.CustomErrorConfigurations)); setErr != nil {
			return fmt.Errorf("Error setting `custom_error_configuration`: %+v", setErr)
		}

		urlPathMaps, err := flattenApplicationGatewayURLPathMaps(props.URLPathMaps)
		if err != nil {
			return fmt.Errorf("Error flattening `url_path_map`: %+v", err)
		}
		if setErr := d.Set("url_path_map", urlPathMaps); setErr != nil {
			return fmt.Errorf("Error setting `url_path_map`: %+v", setErr)
		}

		if setErr := d.Set("waf_configuration", flattenApplicationGatewayWafConfig(props.WebApplicationFirewallConfiguration)); setErr != nil {
			return fmt.Errorf("Error setting `waf_configuration`: %+v", setErr)
		}

		if props.FirewallPolicy != nil {
			d.Set("firewall_policy_id", props.FirewallPolicy.ID)
		}
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

func expandAzureRmApplicationGatewayIdentity(d *pluginsdk.ResourceData) *network.ManagedServiceIdentity {
	v := d.Get("identity")
	identities := v.([]interface{})
	identity := identities[0].(map[string]interface{})
	identityType := network.ResourceIdentityType(identity["type"].(string))

	identityIds := make(map[string]*network.ManagedServiceIdentityUserAssignedIdentitiesValue)
	for _, id := range identity["identity_ids"].([]interface{}) {
		identityIds[id.(string)] = &network.ManagedServiceIdentityUserAssignedIdentitiesValue{}
	}

	appGatewayIdentity := network.ManagedServiceIdentity{
		Type: identityType,
	}

	if identityType == network.ResourceIdentityTypeUserAssigned {
		appGatewayIdentity.UserAssignedIdentities = identityIds
	}

	return &appGatewayIdentity
}

func flattenRmApplicationGatewayIdentity(identity *network.ManagedServiceIdentity) ([]interface{}, error) {
	if identity == nil {
		return make([]interface{}, 0), nil
	}

	result := make(map[string]interface{})
	result["type"] = string(identity.Type)
	if result["type"] == "userAssigned" {
		result["type"] = "UserAssigned"
	}

	identityIds := make([]string, 0)
	if identity.UserAssignedIdentities != nil {
		for key := range identity.UserAssignedIdentities {
			parsedId, err := msiParse.UserAssignedIdentityID(key)
			if err != nil {
				return nil, err
			}
			identityIds = append(identityIds, parsedId.ID())
		}
	}
	result["identity_ids"] = identityIds

	return []interface{}{result}, nil
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

func expandApplicationGatewayTrustedRootCertificates(certs []interface{}) *[]network.ApplicationGatewayTrustedRootCertificate {
	results := make([]network.ApplicationGatewayTrustedRootCertificate, 0)

	for _, raw := range certs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		data := v["data"].(string)

		output := network.ApplicationGatewayTrustedRootCertificate{
			Name: utils.String(name),
			ApplicationGatewayTrustedRootCertificatePropertiesFormat: &network.ApplicationGatewayTrustedRootCertificatePropertiesFormat{},
		}

		if data != "" {
			output.ApplicationGatewayTrustedRootCertificatePropertiesFormat.Data = utils.String(utils.Base64EncodeIfNot(data))
		}

		results = append(results, output)
	}

	return &results
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

		/*kvsid := ""
		if props := cert.ApplicationGatewayTrustedRootCertificatePropertiesFormat; props != nil {
			if v := props.KeyVaultSecretID; v != nil {
				kvsid = *v
				output["key_vault_secret_id"] = *v
			}
		}*/

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
	vs := d.Get("backend_address_pool").([]interface{})
	results := make([]network.ApplicationGatewayBackendAddressPool, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})
		backendAddresses := make([]network.ApplicationGatewayBackendAddress, 0)

		for _, ip := range v["fqdns"].([]interface{}) {
			backendAddresses = append(backendAddresses, network.ApplicationGatewayBackendAddress{
				Fqdn: utils.String(ip.(string)),
			})
		}
		for _, ip := range v["ip_addresses"].([]interface{}) {
			backendAddresses = append(backendAddresses, network.ApplicationGatewayBackendAddress{
				IPAddress: utils.String(ip.(string)),
			})
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
	vs := d.Get("backend_http_settings").([]interface{})

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

					certId, err := azure.ParseAzureResourceID(*cert.ID)
					if err != nil {
						return nil, err
					}

					name := certId.Path["authenticationCertificates"]
					certificate := map[string]interface{}{
						"id":   *cert.ID,
						"name": name,
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

					certId, err := azure.ParseAzureResourceID(*cert.ID)
					if err != nil {
						return nil, err
					}

					certName := certId.Path["trustedRootCertificates"]
					trustedRootCertificateNames = append(trustedRootCertificateNames, certName)
				}
			}
			output["trusted_root_certificate_names"] = trustedRootCertificateNames

			if probe := props.Probe; probe != nil {
				if probe.ID != nil {
					id, err := azure.ParseAzureResourceID(*probe.ID)
					if err != nil {
						return results, err
					}

					output["probe_name"] = id.Path["probes"]
					output["probe_id"] = *probe.ID
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

func expandApplicationGatewaySslPolicy(d *pluginsdk.ResourceData) *network.ApplicationGatewaySslPolicy {
	policy := network.ApplicationGatewaySslPolicy{}
	disabledSSLPolicies := make([]network.ApplicationGatewaySslProtocol, 0)

	vs := d.Get("ssl_policy").([]interface{})

	if len(vs) > 0 && vs[0] != nil {
		v := vs[0].(map[string]interface{})
		policyType := network.ApplicationGatewaySslPolicyType(v["policy_type"].(string))

		for _, policy := range v["disabled_protocols"].([]interface{}) {
			disabledSSLPolicies = append(disabledSSLPolicies, network.ApplicationGatewaySslProtocol(policy.(string)))
		}

		if policyType == network.ApplicationGatewaySslPolicyTypePredefined {
			policyName := network.ApplicationGatewaySslPolicyName(v["policy_name"].(string))
			policy = network.ApplicationGatewaySslPolicy{
				PolicyType: policyType,
				PolicyName: policyName,
			}
		} else if policyType == network.ApplicationGatewaySslPolicyTypeCustom {
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

	if len(disabledSSLPolicies) > 0 {
		policy = network.ApplicationGatewaySslPolicy{
			DisabledSslProtocols: &disabledSSLPolicies,
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
	vs := d.Get("http_listener").([]interface{})
	results := make([]network.ApplicationGatewayHTTPListener, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		frontendIPConfigName := v["frontend_ip_configuration_name"].(string)
		frontendPortName := v["frontend_port_name"].(string)
		protocol := v["protocol"].(string)
		requireSNI := v["require_sni"].(bool)

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
					portId, err := azure.ParseAzureResourceID(*port.ID)
					if err != nil {
						return nil, err
					}
					portName := portId.Path["frontendPorts"]
					output["frontend_port_name"] = portName
					output["frontend_port_id"] = *port.ID
				}
			}

			if feConfig := props.FrontendIPConfiguration; feConfig != nil {
				if feConfig.ID != nil {
					feConfigId, err := azure.ParseAzureResourceID(*feConfig.ID)
					if err != nil {
						return nil, err
					}
					frontendName := feConfigId.Path["frontendIPConfigurations"]
					output["frontend_ip_configuration_name"] = frontendName
					output["frontend_ip_configuration_id"] = *feConfig.ID
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
					certId, err := azure.ParseAzureResourceID(*cert.ID)
					if err != nil {
						return nil, err
					}
					sslCertName := certId.Path["sslCertificates"]

					output["ssl_certificate_name"] = sslCertName
					output["ssl_certificate_id"] = *cert.ID
				}
			}

			if sni := props.RequireServerNameIndication; sni != nil {
				output["require_sni"] = *sni
			}

			if fwp := props.FirewallPolicy; fwp != nil && fwp.ID != nil {
				output["firewall_policy_id"] = *fwp.ID
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

func expandApplicationGatewayFrontendIPConfigurations(d *pluginsdk.ResourceData) *[]network.ApplicationGatewayFrontendIPConfiguration {
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

		name := v["name"].(string)
		output := network.ApplicationGatewayFrontendIPConfiguration{
			Name: utils.String(name),
			ApplicationGatewayFrontendIPConfigurationPropertiesFormat: &properties,
		}

		results = append(results, output)
	}

	return &results
}

func flattenApplicationGatewayFrontendIPConfigurations(input *[]network.ApplicationGatewayFrontendIPConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
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
		}

		results = append(results, output)
	}

	return results
}

func expandApplicationGatewayProbes(d *pluginsdk.ResourceData) *[]network.ApplicationGatewayProbe {
	vs := d.Get("probe").([]interface{})
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

func expandApplicationGatewayRequestRoutingRules(d *pluginsdk.ResourceData, gatewayID string) (*[]network.ApplicationGatewayRequestRoutingRule, error) {
	vs := d.Get("request_routing_rule").(*pluginsdk.Set).List()
	results := make([]network.ApplicationGatewayRequestRoutingRule, 0)

	for _, raw := range vs {
		v := raw.(map[string]interface{})

		name := v["name"].(string)
		ruleType := v["rule_type"].(string)
		httpListenerName := v["http_listener_name"].(string)
		httpListenerID := fmt.Sprintf("%s/httpListeners/%s", gatewayID, httpListenerName)
		backendAddressPoolName := v["backend_address_pool_name"].(string)
		backendHTTPSettingsName := v["backend_http_settings_name"].(string)
		redirectConfigName := v["redirect_configuration_name"].(string)

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

		results = append(results, rule)
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

			if pool := props.BackendAddressPool; pool != nil {
				if pool.ID != nil {
					poolId, err := azure.ParseAzureResourceID(*pool.ID)
					if err != nil {
						return nil, err
					}
					backendAddressPoolName := poolId.Path["backendAddressPools"]
					output["backend_address_pool_name"] = backendAddressPoolName
					output["backend_address_pool_id"] = *pool.ID
				}
			}

			if settings := props.BackendHTTPSettings; settings != nil {
				if settings.ID != nil {
					settingsId, err := azure.ParseAzureResourceID(*settings.ID)
					if err != nil {
						return nil, err
					}
					backendHTTPSettingsName := settingsId.Path["backendHttpSettingsCollection"]
					output["backend_http_settings_name"] = backendHTTPSettingsName
					output["backend_http_settings_id"] = *settings.ID
				}
			}

			if listener := props.HTTPListener; listener != nil {
				if listener.ID != nil {
					listenerId, err := azure.ParseAzureResourceID(*listener.ID)
					if err != nil {
						return nil, err
					}
					httpListenerName := listenerId.Path["httpListeners"]
					output["http_listener_id"] = *listener.ID
					output["http_listener_name"] = httpListenerName
				}
			}

			if pathMap := props.URLPathMap; pathMap != nil {
				if pathMap.ID != nil {
					pathMapId, err := azure.ParseAzureResourceID(*pathMap.ID)
					if err != nil {
						return nil, err
					}
					urlPathMapName := pathMapId.Path["urlPathMaps"]
					output["url_path_map_name"] = urlPathMapName
					output["url_path_map_id"] = *pathMap.ID
				}
			}

			if redirect := props.RedirectConfiguration; redirect != nil {
				if redirect.ID != nil {
					redirectId, err := azure.ParseAzureResourceID(*redirect.ID)
					if err != nil {
						return nil, err
					}
					redirectName := redirectId.Path["redirectConfigurations"]
					output["redirect_configuration_name"] = redirectName
					output["redirect_configuration_id"] = *redirect.ID
				}
			}

			if rewrite := props.RewriteRuleSet; rewrite != nil {
				if rewrite.ID != nil {
					rewriteId, err := azure.ParseAzureResourceID(*rewrite.ID)
					if err != nil {
						return nil, err
					}
					rewriteName := rewriteId.Path["rewriteRuleSets"]
					output["rewrite_rule_set_name"] = rewriteName
					output["rewrite_rule_set_id"] = *rewrite.ID
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
				if c["path"] != nil {
					urlConfiguration.ModifiedPath = utils.String(c["path"].(string))
				}
				if c["query_string"] != nil {
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
							urlConfig := map[string]interface{}{}

							if config.ModifiedPath != nil {
								urlConfig["path"] = *config.ModifiedPath
							}

							if config.ModifiedQueryString != nil {
								urlConfig["query_string"] = *config.ModifiedQueryString
							}

							if config.Reroute != nil {
								urlConfig["reroute"] = *config.Reroute
							}
							urlConfigs = append(urlConfigs, urlConfig)
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
					listenerId, err := azure.ParseAzureResourceID(*listener.ID)
					if err != nil {
						return nil, err
					}
					targetListenerName := listenerId.Path["httpListeners"]
					output["target_listener_name"] = targetListenerName
					output["target_listener_id"] = *listener.ID
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
	vs := d.Get("ssl_certificate").([]interface{})
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
			existingVals := existing.([]interface{})
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
				rulePaths = append(rulePaths, rulePath.(string))
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
				poolId, err := azure.ParseAzureResourceID(*backendPool.ID)
				if err != nil {
					return nil, err
				}
				backendAddressPoolName := poolId.Path["backendAddressPools"]
				output["default_backend_address_pool_name"] = backendAddressPoolName
				output["default_backend_address_pool_id"] = *backendPool.ID
			}

			if settings := props.DefaultBackendHTTPSettings; settings != nil && settings.ID != nil {
				settingsId, err := azure.ParseAzureResourceID(*settings.ID)
				if err != nil {
					return nil, err
				}
				backendHTTPSettingsName := settingsId.Path["backendHttpSettingsCollection"]
				output["default_backend_http_settings_name"] = backendHTTPSettingsName
				output["default_backend_http_settings_id"] = *settings.ID
			}

			if redirect := props.DefaultRedirectConfiguration; redirect != nil && redirect.ID != nil {
				settingsId, err := azure.ParseAzureResourceID(*redirect.ID)
				if err != nil {
					return nil, err
				}
				redirectConfigurationName := settingsId.Path["redirectConfigurations"]
				output["default_redirect_configuration_name"] = redirectConfigurationName
				output["default_redirect_configuration_id"] = *redirect.ID
			}

			if rewrite := props.DefaultRewriteRuleSet; rewrite != nil && rewrite.ID != nil {
				settingsId, err := azure.ParseAzureResourceID(*rewrite.ID)
				if err != nil {
					return nil, err
				}
				defaultRewriteRuleSetName := settingsId.Path["rewriteRuleSets"]
				output["default_rewrite_rule_set_name"] = defaultRewriteRuleSetName
				output["default_rewrite_rule_set_id"] = *rewrite.ID
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
							poolId, err := azure.ParseAzureResourceID(*pool.ID)
							if err != nil {
								return nil, err
							}
							backendAddressPoolName2 := poolId.Path["backendAddressPools"]
							ruleOutput["backend_address_pool_name"] = backendAddressPoolName2
							ruleOutput["backend_address_pool_id"] = *pool.ID
						}

						if backend := ruleProps.BackendHTTPSettings; backend != nil && backend.ID != nil {
							backendId, err := azure.ParseAzureResourceID(*backend.ID)
							if err != nil {
								return nil, err
							}
							backendHTTPSettingsName2 := backendId.Path["backendHttpSettingsCollection"]
							ruleOutput["backend_http_settings_name"] = backendHTTPSettingsName2
							ruleOutput["backend_http_settings_id"] = *backend.ID
						}

						if redirect := ruleProps.RedirectConfiguration; redirect != nil && redirect.ID != nil {
							redirectId, err := azure.ParseAzureResourceID(*redirect.ID)
							if err != nil {
								return nil, err
							}
							redirectConfigurationName2 := redirectId.Path["redirectConfigurations"]
							ruleOutput["redirect_configuration_name"] = redirectConfigurationName2
							ruleOutput["redirect_configuration_id"] = *redirect.ID
						}

						if rewrite := ruleProps.RewriteRuleSet; rewrite != nil && rewrite.ID != nil {
							rewriteId, err := azure.ParseAzureResourceID(*rewrite.ID)
							if err != nil {
								return nil, err
							}
							rewriteRuleSet := rewriteId.Path["rewriteRuleSets"]
							ruleOutput["rewrite_rule_set_name"] = rewriteRuleSet
							ruleOutput["rewrite_rule_set_id"] = *rewrite.ID
						}

						if fwp := ruleProps.FirewallPolicy; fwp != nil && fwp.ID != nil {
							ruleOutput["firewall_policy_id"] = *fwp.ID
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

func applicationGatewayCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) error {
	_, hasAutoscaleConfig := d.GetOk("autoscale_configuration.0")
	capacity, hasCapacity := d.GetOk("sku.0.capacity")
	tier := d.Get("sku.0.tier").(string)

	if !hasAutoscaleConfig && !hasCapacity {
		return fmt.Errorf("The Application Gateway must specify either `capacity` or `autoscale_configuration` for the selected SKU tier %q", tier)
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
