// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func machineLearningWorkspaceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.WorkspaceName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"application_insights_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: componentsapis.ValidateComponentID,
			// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"key_vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKeyVaultID,
			// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
			// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

		"kind": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Default",
				"FeatureStore",
			}, false),
			Default: "Default",
		},

		"feature_store": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"computer_spark_runtime_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"offline_connection_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"online_connection_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"primary_user_assigned_identity": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"container_registry_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: registries.ValidateRegistryID,
			// TODO -- remove when issue https://github.com/Azure/azure-rest-api-specs/issues/8323 is addressed
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"image_build_compute_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"encryption": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),
					"key_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
					"user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						// TODO: remove this
						DiffSuppressFunc: suppress.CaseDifference,
					},
				},
			},
		},

		"managed_network": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"isolation_mode": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.StringInSlice(workspaces.PossibleValuesForIsolationMode(), false),
					},
					"provision_on_creation_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
						Default:  false,
					},
				},
			},
		},

		"friendly_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"high_business_impact": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(Basic),
			ValidateFunc: validation.StringInSlice([]string{string(Basic)}, false),
		},

		"service_side_encryption_enabled": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			ForceNew:     true,
			Default:      false,
			RequiredWith: []string{"encryption"},
		},

		"v1_legacy_mode_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"storage_account_access_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  workspaces.SystemDatastoresAuthModeAccessKey,
			ValidateFunc: validation.StringInSlice([]string{
				string(workspaces.SystemDatastoresAuthModeAccessKey),
				string(workspaces.SystemDatastoresAuthModeIdentity),
			}, false),
		},

		"serverless_compute": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
					"public_ip_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"discovery_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}
