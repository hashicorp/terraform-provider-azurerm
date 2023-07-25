// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SpringCloudV0ToV1 struct{}

func (s SpringCloudV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "S0",
			ForceNew: true,
		},

		"build_agent_pool_size": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"log_stream_public_endpoint_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"network": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"app_subnet_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"service_runtime_subnet_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"cidr_ranges": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MinItems: 3,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"app_network_resource_group": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
					},

					"read_timeout_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"service_runtime_network_resource_group": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
					},
				},
			},
		},

		"config_server_git_setting": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"uri": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"label": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"search_paths": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"http_basic_auth": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"username": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"password": {
									Type:      pluginsdk.TypeString,
									Required:  true,
									Sensitive: true,
								},
							},
						},
					},

					"ssh_auth": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"private_key": {
									Type:      pluginsdk.TypeString,
									Required:  true,
									Sensitive: true,
								},

								"host_key": {
									Type:      pluginsdk.TypeString,
									Optional:  true,
									Sensitive: true,
								},

								"host_key_algorithm": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"strict_host_key_checking_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  true,
								},
							},
						},
					},

					"repository": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"uri": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"label": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"pattern": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"search_paths": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"http_basic_auth": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"username": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},

											"password": {
												Type:      pluginsdk.TypeString,
												Required:  true,
												Sensitive: true,
											},
										},
									},
								},

								"ssh_auth": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"private_key": {
												Type:      pluginsdk.TypeString,
												Required:  true,
												Sensitive: true,
											},

											"host_key": {
												Type:      pluginsdk.TypeString,
												Optional:  true,
												Sensitive: true,
											},

											"host_key_algorithm": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},

											"strict_host_key_checking_enabled": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												Default:  true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},

		"trace": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"connection_string": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"sample_rate": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
						Default:  10,
					},
				},
			},
		},

		"service_registry_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"outbound_public_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"required_network_traffic_rules": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"protocol": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"port": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"ip_addresses": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"fqdns": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"direction": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"zone_redundant": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"service_registry_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (s SpringCloudV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.SpringCloudServiceIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
