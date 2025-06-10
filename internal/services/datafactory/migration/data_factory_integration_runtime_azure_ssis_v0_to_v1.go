// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DataFactoryIntegrationRuntimeAzureSsisV0ToV1 struct{}

var _ pluginsdk.StateUpgrade = DataFactoryIntegrationRuntimeAzureSsisV0ToV1{}

func (DataFactoryIntegrationRuntimeAzureSsisV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"data_factory_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"node_size": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"number_of_nodes": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"max_parallel_executions_per_node": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"credential_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"edition": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"copy_compute_scale": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"data_integration_unit": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"time_to_live": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"express_vnet_integration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"license_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"vnet_integration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"vnet_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"subnet_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"public_ips": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"custom_setup_script": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"blob_container_uri": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"sas_token": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"catalog_info": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"server_endpoint": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"administrator_login": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"administrator_password": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"pricing_tier": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"elastic_pool_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"dual_standby_pair_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"express_custom_setup": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"environment": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"powershell_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"command_key": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"target_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"user_name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"password": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"key_vault_password": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"linked_service_name": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},

											"secret_name": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},

											"parameters": {
												Type:     pluginsdk.TypeMap,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},

											"secret_version": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},
										},
									},
								},
							},
						},
					},

					"component": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"license": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"key_vault_license": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"linked_service_name": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},

											"secret_name": {
												Type:     pluginsdk.TypeString,
												Required: true,
											},

											"parameters": {
												Type:     pluginsdk.TypeMap,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},

											"secret_version": {
												Type:     pluginsdk.TypeString,
												Optional: true,
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

		"package_store": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"linked_service_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"pipeline_external_compute_scale": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"number_of_external_nodes": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"number_of_pipeline_nodes": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"time_to_live": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"proxy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"self_hosted_integration_runtime_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"staging_storage_linked_service_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"path": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
	}
}

func (DataFactoryIntegrationRuntimeAzureSsisV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// Migration to update ID segment from lowercase to camelCase (integrationruntimes to integrationRuntimes)

		oldId := rawState["id"].(string)
		parsedId, err := integrationruntimes.ParseIntegrationRuntimeIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		newId := parsedId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
