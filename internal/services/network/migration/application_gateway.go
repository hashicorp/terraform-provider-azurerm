package migration

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/applicationgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ApplicationGatewayV0ToV1{}

type ApplicationGatewayV0ToV1 struct{}

func (ApplicationGatewayV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"authentication_certificate": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"data": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"autoscale_configuration": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"max_capacity": {
					Optional: true,
					Type:     pluginsdk.TypeInt,
				},
				"min_capacity": {
					Required: true,
					Type:     pluginsdk.TypeInt,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"backend_address_pool": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"fqdns": {
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Optional: true,
					Type:     pluginsdk.TypeSet,
				},
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"ip_addresses": {
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Optional: true,
					Type:     pluginsdk.TypeSet,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Required: true,
			Set:      applicationGatewayBackendAddressPool,
			Type:     pluginsdk.TypeSet,
		},
		"backend_http_settings": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"affinity_cookie_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"authentication_certificate": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"id": {
							Computed: true,
							Type:     pluginsdk.TypeString,
						},
						"name": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"connection_draining": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"drain_timeout_sec": {
							Required: true,
							Type:     pluginsdk.TypeInt,
						},
						"enabled": {
							Required: true,
							Type:     pluginsdk.TypeBool,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"cookie_based_affinity": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"host_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"path": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"pick_host_name_from_backend_address": {
					Optional: true,
					Type:     pluginsdk.TypeBool,
				},
				"port": {
					Required: true,
					Type:     pluginsdk.TypeInt,
				},
				"probe_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"probe_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"protocol": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"request_timeout": {
					Optional: true,
					Type:     pluginsdk.TypeInt,
				},
				"trusted_root_certificate_names": {
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
			}},
			Required: true,
			Set:      applicationGatewayBackendSettingsHash,
			Type:     pluginsdk.TypeSet,
		},
		"custom_error_configuration": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"custom_error_page_url": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"status_code": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"enable_http2": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"fips_enabled": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"firewall_policy_id": {
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"force_firewall_policy_association": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"frontend_ip_configuration": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"private_ip_address": {
					Computed: true,
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"private_ip_address_allocation": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"private_link_configuration_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"private_link_configuration_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"public_ip_address_id": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"subnet_id": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Required: true,
			Type:     pluginsdk.TypeList,
		},
		"frontend_port": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"port": {
					Required: true,
					Type:     pluginsdk.TypeInt,
				},
			}},
			Required: true,
			Type:     pluginsdk.TypeSet,
		},
		"gateway_ip_configuration": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"subnet_id": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Required: true,
			Type:     pluginsdk.TypeList,
		},
		"global": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"request_buffering_enabled": {
					Required: true,
					Type:     pluginsdk.TypeBool,
				},
				"response_buffering_enabled": {
					Required: true,
					Type:     pluginsdk.TypeBool,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"http_listener": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"custom_error_configuration": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"custom_error_page_url": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"id": {
							Computed: true,
							Type:     pluginsdk.TypeString,
						},
						"status_code": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"firewall_policy_id": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"frontend_ip_configuration_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"frontend_ip_configuration_name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"frontend_port_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"frontend_port_name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"host_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"host_names": {
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Optional: true,
					Type:     pluginsdk.TypeSet,
				},
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"protocol": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"require_sni": {
					Optional: true,
					Type:     pluginsdk.TypeBool,
				},
				"ssl_certificate_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"ssl_certificate_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"ssl_profile_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"ssl_profile_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Required: true,
			Set:      applicationGatewayHttpListenerHash,
			Type:     pluginsdk.TypeSet,
		},
		"identity": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"identity_ids": {
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Required: true,
					Type:     pluginsdk.TypeSet,
				},
				"type": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"location": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"private_endpoint_connection": {
			Computed: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"name": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Type: pluginsdk.TypeSet,
		},
		"private_link_configuration": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"ip_configuration": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"name": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"primary": {
							Required: true,
							Type:     pluginsdk.TypeBool,
						},
						"private_ip_address": {
							Computed: true,
							Optional: true,
							Type:     pluginsdk.TypeString,
						},
						"private_ip_address_allocation": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"subnet_id": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
					}},
					Required: true,
					Type:     pluginsdk.TypeList,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeSet,
		},
		"probe": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"host": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"interval": {
					Required: true,
					Type:     pluginsdk.TypeInt,
				},
				"match": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"body": {
							Optional: true,
							Type:     pluginsdk.TypeString,
						},
						"status_code": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Required: true,
							Type:     pluginsdk.TypeList,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"minimum_servers": {
					Optional: true,
					Type:     pluginsdk.TypeInt,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"path": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"pick_host_name_from_backend_http_settings": {
					Optional: true,
					Type:     pluginsdk.TypeBool,
				},
				"port": {
					Optional: true,
					Type:     pluginsdk.TypeInt,
				},
				"protocol": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"timeout": {
					Required: true,
					Type:     pluginsdk.TypeInt,
				},
				"unhealthy_threshold": {
					Required: true,
					Type:     pluginsdk.TypeInt,
				},
			}},
			Optional: true,
			Set:      applicationGatewayProbeHash,
			Type:     pluginsdk.TypeSet,
		},
		"redirect_configuration": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"include_path": {
					Optional: true,
					Type:     pluginsdk.TypeBool,
				},
				"include_query_string": {
					Optional: true,
					Type:     pluginsdk.TypeBool,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"redirect_type": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"target_listener_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"target_listener_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"target_url": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeSet,
		},
		"request_routing_rule": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"backend_address_pool_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"backend_address_pool_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"backend_http_settings_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"backend_http_settings_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"http_listener_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"http_listener_name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"priority": {
					Optional: true,
					Type:     pluginsdk.TypeInt,
				},
				"redirect_configuration_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"redirect_configuration_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"rewrite_rule_set_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"rewrite_rule_set_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"rule_type": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"url_path_map_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"url_path_map_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Required: true,
			Type:     pluginsdk.TypeSet,
		},
		"resource_group_name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"rewrite_rule_set": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"rewrite_rule": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"condition": {
							Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
								"ignore_case": {
									Optional: true,
									Type:     pluginsdk.TypeBool,
								},
								"negate": {
									Optional: true,
									Type:     pluginsdk.TypeBool,
								},
								"pattern": {
									Required: true,
									Type:     pluginsdk.TypeString,
								},
								"variable": {
									Required: true,
									Type:     pluginsdk.TypeString,
								},
							}},
							Optional: true,
							Type:     pluginsdk.TypeList,
						},
						"name": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"request_header_configuration": {
							Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
								"header_name": {
									Required: true,
									Type:     pluginsdk.TypeString,
								},
								"header_value": {
									Required: true,
									Type:     pluginsdk.TypeString,
								},
							}},
							Optional: true,
							Type:     pluginsdk.TypeList,
						},
						"response_header_configuration": {
							Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
								"header_name": {
									Required: true,
									Type:     pluginsdk.TypeString,
								},
								"header_value": {
									Required: true,
									Type:     pluginsdk.TypeString,
								},
							}},
							Optional: true,
							Type:     pluginsdk.TypeList,
						},
						"rule_sequence": {
							Required: true,
							Type:     pluginsdk.TypeInt,
						},
						"url": {
							Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
								"components": {
									Computed: true,
									Optional: true,
									Type:     pluginsdk.TypeString,
								},
								"path": {
									Optional: true,
									Type:     pluginsdk.TypeString,
								},
								"query_string": {
									Optional: true,
									Type:     pluginsdk.TypeString,
								},
								"reroute": {
									Optional: true,
									Type:     pluginsdk.TypeBool,
								},
							}},
							Optional: true,
							Type:     pluginsdk.TypeList,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"sku": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"capacity": {
					Optional: true,
					Type:     pluginsdk.TypeInt,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"tier": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Required: true,
			Type:     pluginsdk.TypeList,
		},
		"ssl_certificate": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"data": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"key_vault_secret_id": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"password": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"public_cert_data": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Set:      applicationGatewaySSLCertificate,
			Type:     pluginsdk.TypeSet,
		},
		"ssl_policy": {
			Computed: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"cipher_suites": {
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"disabled_protocols": {
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"min_protocol_version": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"policy_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"policy_type": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"ssl_profile": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"ssl_policy": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"cipher_suites": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Optional: true,
							Type:     pluginsdk.TypeList,
						},
						"disabled_protocols": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Optional: true,
							Type:     pluginsdk.TypeList,
						},
						"min_protocol_version": {
							Optional: true,
							Type:     pluginsdk.TypeString,
						},
						"policy_name": {
							Optional: true,
							Type:     pluginsdk.TypeString,
						},
						"policy_type": {
							Optional: true,
							Type:     pluginsdk.TypeString,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"trusted_client_certificate_names": {
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"verify_client_cert_issuer_dn": {
					Optional: true,
					Type:     pluginsdk.TypeBool,
				},
				"verify_client_certificate_revocation": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"tags": {
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Optional: true,
			Type:     pluginsdk.TypeMap,
		},
		"trusted_client_certificate": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"data": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"trusted_root_certificate": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"data": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"key_vault_secret_id": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"url_path_map": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"default_backend_address_pool_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"default_backend_address_pool_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"default_backend_http_settings_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"default_backend_http_settings_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"default_redirect_configuration_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"default_redirect_configuration_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"default_rewrite_rule_set_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"default_rewrite_rule_set_name": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"path_rule": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"backend_address_pool_id": {
							Computed: true,
							Type:     pluginsdk.TypeString,
						},
						"backend_address_pool_name": {
							Optional: true,
							Type:     pluginsdk.TypeString,
						},
						"backend_http_settings_id": {
							Computed: true,
							Type:     pluginsdk.TypeString,
						},
						"backend_http_settings_name": {
							Optional: true,
							Type:     pluginsdk.TypeString,
						},
						"firewall_policy_id": {
							Optional: true,
							Type:     pluginsdk.TypeString,
						},
						"id": {
							Computed: true,
							Type:     pluginsdk.TypeString,
						},
						"name": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"paths": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Required: true,
							Type:     pluginsdk.TypeList,
						},
						"redirect_configuration_id": {
							Computed: true,
							Type:     pluginsdk.TypeString,
						},
						"redirect_configuration_name": {
							Optional: true,
							Type:     pluginsdk.TypeString,
						},
						"rewrite_rule_set_id": {
							Computed: true,
							Type:     pluginsdk.TypeString,
						},
						"rewrite_rule_set_name": {
							Optional: true,
							Type:     pluginsdk.TypeString,
						},
					}},
					Required: true,
					Type:     pluginsdk.TypeList,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"waf_configuration": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"disabled_rule_group": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"rule_group_name": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"rules": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeInt},
							Optional: true,
							Type:     pluginsdk.TypeList,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"enabled": {
					Required: true,
					Type:     pluginsdk.TypeBool,
				},
				"exclusion": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"match_variable": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"selector": {
							Optional: true,
							Type:     pluginsdk.TypeString,
						},
						"selector_match_operator": {
							Optional: true,
							Type:     pluginsdk.TypeString,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"file_upload_limit_mb": {
					Optional: true,
					Type:     pluginsdk.TypeInt,
				},
				"firewall_mode": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"max_request_body_size_kb": {
					Optional: true,
					Type:     pluginsdk.TypeInt,
				},
				"request_body_check": {
					Optional: true,
					Type:     pluginsdk.TypeBool,
				},
				"rule_set_type": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"rule_set_version": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"zones": {
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			ForceNew: true,
			Optional: true,
			Type:     pluginsdk.TypeSet,
		},
	}

}

func (ApplicationGatewayV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Printf("[Debug] start upgrade application gateway id")
		oldID := rawState["id"].(string)
		if newID, err := normalizeApplicationGatewayID(oldID); err != nil {
			return nil, err
		} else if newID != nil {
			rawState["id"] = *newID
		}
		return rawState, nil
	}
}

func normalizeApplicationGatewayID(id string) (*string, error) {
	if id == "" {
		return nil, nil
	}
	parseID, err := applicationgateways.ParseApplicationGatewayIDInsensitively(id)
	if err != nil {
		return nil, fmt.Errorf("prase id: %v", err)
	}
	normalizedID := parseID.ID()
	return &normalizedID, nil
}

func applicationGatewayHttpListenerHash(v interface{}) int {
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
