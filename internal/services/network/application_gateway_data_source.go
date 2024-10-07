// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/applicationgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceApplicationGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceApplicationGatewayRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"backend_address_pool": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"fqdns": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"ip_addresses": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
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
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"path": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"port": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"protocol": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"cookie_based_affinity": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"affinity_cookie_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"host_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"pick_host_name_from_backend_address": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"request_timeout": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"authentication_certificate": {
							Type:     pluginsdk.TypeList,
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

						"trusted_root_certificate_names": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"connection_draining": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"enabled": {
										Type:     pluginsdk.TypeBool,
										Computed: true,
									},

									"drain_timeout_sec": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},
								},
							},
						},

						"probe_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"subnet_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"private_ip_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"public_ip_address_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"private_ip_address_allocation": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"private_link_configuration_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"port": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
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
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"subnet_id": {
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

			"global": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"request_buffering_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"response_buffering_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"http_listener": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"frontend_ip_configuration_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"frontend_port_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"protocol": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"host_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"host_names": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"ssl_certificate_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"require_sni": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
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
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"status_code": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"custom_error_page_url": {
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

						"firewall_policy_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"ssl_profile_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"fips_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"private_endpoint_connection": {
				Type:     pluginsdk.TypeList,
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
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"ip_configuration": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"subnet_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"private_ip_address": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"private_ip_address_allocation": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"primary": {
										Type:     pluginsdk.TypeBool,
										Computed: true,
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
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"rule_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"http_listener_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"backend_address_pool_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"backend_http_settings_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"url_path_map_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"redirect_configuration_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"rewrite_rule_set_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"priority": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
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
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"redirect_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"target_listener_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"target_url": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"include_path": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"include_query_string": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
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
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"min_capacity": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"max_capacity": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"sku": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tier": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"capacity": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"authentication_certificate": {
				Type:     pluginsdk.TypeList,
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

			"trusted_root_certificate": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"key_vault_secret_id": {
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

			"ssl_policy": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"disabled_protocols": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"policy_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"policy_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"cipher_suites": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"min_protocol_version": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"http2_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"force_firewall_policy_association": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"probe": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"protocol": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"path": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"host": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"interval": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"timeout": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"unhealthy_threshold": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"port": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"pick_host_name_from_backend_http_settings": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"minimum_servers": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"match": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"body": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"status_code": {
										Type:     pluginsdk.TypeList,
										Computed: true,
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
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"rewrite_rule": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"rule_sequence": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},

									"condition": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"variable": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"pattern": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"ignore_case": {
													Type:     pluginsdk.TypeBool,
													Computed: true,
												},
												"negate": {
													Type:     pluginsdk.TypeBool,
													Computed: true,
												},
											},
										},
									},

									"request_header_configuration": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"header_name": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"header_value": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
											},
										},
									},

									"response_header_configuration": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"header_name": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"header_value": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
											},
										},
									},

									"url": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"path": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},
												"query_string": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},

												"components": {
													Type:     pluginsdk.TypeString,
													Computed: true,
												},

												"reroute": {
													Type:     pluginsdk.TypeBool,
													Computed: true,
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
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"key_vault_secret_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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

			"trusted_client_certificate": {
				Type:     pluginsdk.TypeList,
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

						"data": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"ssl_profile": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"trusted_client_certificate_names": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"verify_client_certificate_issuer_dn": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"verify_client_certificate_revocation": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"ssl_policy": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"disabled_protocols": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},

									"policy_type": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"policy_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"cipher_suites": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},

									"min_protocol_version": {
										Type:     pluginsdk.TypeString,
										Computed: true,
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

			"url_path_map": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"default_backend_address_pool_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"default_backend_http_settings_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"default_redirect_configuration_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"default_rewrite_rule_set_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"path_rule": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"paths": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},

									"backend_address_pool_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"backend_http_settings_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"redirect_configuration_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"rewrite_rule_set_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
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
										Type:     pluginsdk.TypeString,
										Computed: true,
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
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"firewall_mode": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"rule_set_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"rule_set_version": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"file_upload_limit_mb": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
						"request_body_check": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"max_request_body_size_kb": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
						"disabled_rule_group": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"rule_group_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"rules": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeInt,
										},
									},
								},
							},
						},
						"exclusion": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"match_variable": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"selector_match_operator": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"selector": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"firewall_policy_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"custom_error_configuration": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"status_code": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"custom_error_page_url": {
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

			"location": commonschema.LocationComputed(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"zones": commonschema.ZonesMultipleComputed(),

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceApplicationGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ApplicationGateways
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := applicationgateways.NewApplicationGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.ApplicationGatewayName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("zones", zones.FlattenUntyped(model.Zones))
		d.Set("location", location.NormalizeNilable(model.Location))

		identityFlattened, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err = d.Set("identity", identityFlattened); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
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

			d.Set("http2_enabled", props.EnableHTTP2)
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

			sslProfiles, err := flattenApplicationGatewayDataSourceSslProfiles(props.SslProfiles)
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

// TODO: 4.0 remove this after the r/app_gateway schema `verify_client_cert_issuer_dn` is changed to `verify_client_certificate_issuer_dn`, and reuse the flatten function in r/app_gateway file.
func flattenApplicationGatewayDataSourceSslProfiles(input *[]applicationgateways.ApplicationGatewaySslProfile) ([]interface{}, error) {
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

		verifyClientCertIssuerDn := false
		verifyClientCertificateRevocation := ""
		if props := v.Properties; props != nil {
			if clientAuthConfig := props.ClientAuthConfiguration; clientAuthConfig != nil {
				verifyClientCertIssuerDn = pointer.From(clientAuthConfig.VerifyClientCertIssuerDN)
				if clientAuthConfig.VerifyClientRevocation != nil && *clientAuthConfig.VerifyClientRevocation != applicationgateways.ApplicationGatewayClientRevocationOptionsNone {
					verifyClientCertificateRevocation = string(pointer.From(clientAuthConfig.VerifyClientRevocation))
				}
			}
			output["verify_client_certificate_issuer_dn"] = verifyClientCertIssuerDn
			output["verify_client_certificate_revocation"] = verifyClientCertificateRevocation

			output["ssl_policy"] = flattenApplicationGatewaySslPolicy(props.SslPolicy)
		}

		if props := v.Properties; props != nil {
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
		}

		results = append(results, output)
	}

	return results, nil
}
