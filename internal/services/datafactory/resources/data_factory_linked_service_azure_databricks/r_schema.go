// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_azure_databricks

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2026-01-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryLinkedServiceAzureDatabricksSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LinkedServiceDatasetName,
		},

		"data_factory_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},

		// Authentication types
		"msi_workspace_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
			ExactlyOneOf: []string{"access_token", "msi_workspace_id", "key_vault_password"},
		},

		"access_token": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"access_token", "msi_workspace_id", "key_vault_password"},
		},

		"key_vault_password": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"linked_service_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.LinkedServiceDatasetName,
					},

					"secret_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.LinkedServiceDatasetName,
					},
				},
			},
			ExactlyOneOf: []string{"access_token", "msi_workspace_id", "key_vault_password"},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"adb_domain": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		// Cluster types [existing cluster, new cluster, interactive pools]
		"existing_cluster_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"existing_cluster_id", "new_cluster_config", "instance_pool"},
		},

		"new_cluster_config": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"node_type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"custom_tags": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"min_number_of_workers": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntBetween(1, 25000),
					},
					"max_number_of_workers": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 25000),
					},
					"cluster_version": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"spark_config": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"spark_environment_variables": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"log_destination": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"init_scripts": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"driver_node_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			ExactlyOneOf: []string{"existing_cluster_id", "new_cluster_config", "instance_pool"},
		},

		"instance_pool": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"min_number_of_workers": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntBetween(1, 25000),
					},
					"max_number_of_workers": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 25000),
					},
					"instance_pool_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"cluster_version": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			ExactlyOneOf: []string{"existing_cluster_id", "new_cluster_config", "instance_pool"},
		},

		"integration_runtime_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"parameters": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"annotations": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"additional_properties": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
