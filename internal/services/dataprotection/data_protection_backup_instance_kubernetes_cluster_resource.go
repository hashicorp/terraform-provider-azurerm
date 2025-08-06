// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backuppolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	resourceParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type BackupInstanceKubernatesClusterModel struct {
	Name                       string                       `tfschema:"name"`
	Location                   string                       `tfschema:"location"`
	VaultId                    string                       `tfschema:"vault_id"`
	BackupPolicyId             string                       `tfschema:"backup_policy_id"`
	KubernetesClusterId        string                       `tfschema:"kubernetes_cluster_id"`
	SnapshotResourceGroupName  string                       `tfschema:"snapshot_resource_group_name"`
	BackupDatasourceParameters []BackupDatasourceParameters `tfschema:"backup_datasource_parameters"`
}

type BackupDatasourceParameters struct {
	IncludedNamespaces          []string `tfschema:"included_namespaces"`
	IncludedResourceTypes       []string `tfschema:"included_resource_types"`
	ExcludedNamespaces          []string `tfschema:"excluded_namespaces"`
	ExcludedResourceTypes       []string `tfschema:"excluded_resource_types"`
	LabelSelectors              []string `tfschema:"label_selectors"`
	VolumeSnapshotEnabled       bool     `tfschema:"volume_snapshot_enabled"`
	ClusterScopeResourceEnabled bool     `tfschema:"cluster_scoped_resources_enabled"`
}

type DataProtectionBackupInstanceKubernatesClusterResource struct{}

var _ sdk.Resource = DataProtectionBackupInstanceKubernatesClusterResource{}

func (r DataProtectionBackupInstanceKubernatesClusterResource) ResourceType() string {
	return "azurerm_data_protection_backup_instance_kubernetes_cluster"
}

func (r DataProtectionBackupInstanceKubernatesClusterResource) ModelObject() interface{} {
	return &BackupInstanceKubernatesClusterModel{}
}

func (r DataProtectionBackupInstanceKubernatesClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backupinstances.ValidateBackupInstanceID
}

func (r DataProtectionBackupInstanceKubernatesClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": commonschema.Location(),

		"vault_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: backupinstances.ValidateBackupVaultID,
		},

		"backup_policy_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: backuppolicies.ValidateBackupPolicyID,
		},

		"kubernetes_cluster_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKubernetesClusterID,
		},

		"snapshot_resource_group_name": commonschema.ResourceGroupName(),

		"backup_datasource_parameters": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"excluded_namespaces": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"excluded_resource_types": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"cluster_scoped_resources_enabled": {
						Type:     pluginsdk.TypeBool,
						ForceNew: true,
						Optional: true,
						Default:  false,
					},
					"included_namespaces": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"included_resource_types": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"label_selectors": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"volume_snapshot_enabled": {
						Type:     pluginsdk.TypeBool,
						ForceNew: true,
						Optional: true,
						Default:  false,
					},
				},
			},
		},
	}
}

func (r DataProtectionBackupInstanceKubernatesClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataProtectionBackupInstanceKubernatesClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupInstanceKubernatesClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.DataProtection.BackupInstanceClient

			vaultId, err := backupinstances.ParseBackupVaultID(model.VaultId)
			if err != nil {
				return err
			}

			id := backupinstances.NewBackupInstanceID(vaultId.SubscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			policyId, err := backuppolicies.ParseBackupPolicyID(model.BackupPolicyId)
			if err != nil {
				return err
			}

			aksId, err := commonids.ParseKubernetesClusterID(model.KubernetesClusterId)
			if err != nil {
				return err
			}

			snapshotResourceGroupId := resourceParse.NewResourceGroupID(metadata.Client.Account.SubscriptionId, model.SnapshotResourceGroupName)
			parameters := backupinstances.BackupInstanceResource{
				Properties: &backupinstances.BackupInstance{
					DataSourceInfo: backupinstances.Datasource{
						DatasourceType:   pointer.To("Microsoft.ContainerService/managedClusters"),
						ObjectType:       pointer.To("Datasource"),
						ResourceID:       aksId.ID(),
						ResourceLocation: pointer.To(location.Normalize(model.Location)),
						ResourceName:     pointer.To(aksId.ManagedClusterName),
						ResourceType:     pointer.To("Microsoft.ContainerService/managedClusters"),
						ResourceUri:      pointer.To(aksId.ID()),
					},
					DataSourceSetInfo: &backupinstances.DatasourceSet{
						DatasourceType:   pointer.To("Microsoft.ContainerService/managedClusters"),
						ObjectType:       pointer.To("DatasourceSet"),
						ResourceID:       aksId.ID(),
						ResourceLocation: pointer.To(location.Normalize(model.Location)),
						ResourceName:     pointer.To(aksId.ManagedClusterName),
						ResourceType:     pointer.To("Microsoft.ContainerService/managedClusters"),
						ResourceUri:      pointer.To(aksId.ID()),
					},
					FriendlyName: pointer.To(id.BackupInstanceName),
					ObjectType:   "BackupInstance",
					PolicyInfo: backupinstances.PolicyInfo{
						PolicyId: policyId.ID(),
						PolicyParameters: &backupinstances.PolicyParameters{
							DataStoreParametersList: &[]backupinstances.DataStoreParameters{
								backupinstances.AzureOperationalStoreParameters{
									ResourceGroupId: pointer.To(snapshotResourceGroupId.ID()),
									DataStoreType:   backupinstances.DataStoreTypesOperationalStore,
								},
							},
							BackupDatasourceParametersList: expandBackupDatasourceParameters(model.BackupDatasourceParameters),
						},
					},
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameters, backupinstances.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DataProtectionBackupInstanceKubernatesClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient

			id, err := backupinstances.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			vaultId := backupinstances.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)

			state := BackupInstanceKubernatesClusterModel{
				Name:    id.BackupInstanceName,
				VaultId: vaultId.ID(),
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					state.Location = location.NormalizeNilable(properties.DataSourceInfo.ResourceLocation)
					state.BackupPolicyId = properties.PolicyInfo.PolicyId
					state.KubernetesClusterId = properties.DataSourceInfo.ResourceID

					if policyParameters := properties.PolicyInfo.PolicyParameters; policyParameters != nil {
						if dataStorePara := policyParameters.DataStoreParametersList; dataStorePara != nil {
							if dsp := pointer.From(dataStorePara); len(dsp) > 0 {
								if parameter, ok := dsp[0].(backupinstances.AzureOperationalStoreParameters); ok && parameter.ResourceGroupId != nil {
									resourceGroupId, err := resourceParse.ResourceGroupID(*parameter.ResourceGroupId)
									if err != nil {
										return err
									}
									state.SnapshotResourceGroupName = resourceGroupId.ResourceGroup
								}
							}
						}
						if backupDsp := policyParameters.BackupDatasourceParametersList; backupDsp != nil {
							if v := flattenBackupDatasourceParameters(*backupDsp); v != nil {
								state.BackupDatasourceParameters = pointer.From(v)
							}
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DataProtectionBackupInstanceKubernatesClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient

			id, err := backupinstances.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.DeleteThenPoll(ctx, *id, backupinstances.DefaultDeleteOperationOptions())
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandBackupDatasourceParameters(input []BackupDatasourceParameters) *[]backupinstances.BackupDatasourceParameters {
	if len(input) == 0 {
		return nil
	}
	results := make([]backupinstances.BackupDatasourceParameters, 0)
	results = append(results, backupinstances.KubernetesClusterBackupDatasourceParameters{
		ExcludedNamespaces:           pointer.To(input[0].ExcludedNamespaces),
		ExcludedResourceTypes:        pointer.To(input[0].ExcludedResourceTypes),
		IncludeClusterScopeResources: input[0].ClusterScopeResourceEnabled,
		IncludedNamespaces:           pointer.To(input[0].IncludedNamespaces),
		IncludedResourceTypes:        pointer.To(input[0].IncludedResourceTypes),
		LabelSelectors:               pointer.To(input[0].LabelSelectors),
		SnapshotVolumes:              input[0].VolumeSnapshotEnabled,
	})
	return &results
}

func flattenBackupDatasourceParameters(input []backupinstances.BackupDatasourceParameters) *[]BackupDatasourceParameters {
	results := make([]BackupDatasourceParameters, 0)
	if len(input) == 0 {
		return &results
	}

	if item, ok := input[0].(backupinstances.KubernetesClusterBackupDatasourceParameters); ok {
		results = append(results, BackupDatasourceParameters{
			ExcludedNamespaces:          pointer.From(item.ExcludedNamespaces),
			ExcludedResourceTypes:       pointer.From(item.ExcludedResourceTypes),
			ClusterScopeResourceEnabled: item.IncludeClusterScopeResources,
			IncludedNamespaces:          pointer.From(item.IncludedNamespaces),
			IncludedResourceTypes:       pointer.From(item.IncludedResourceTypes),
			LabelSelectors:              pointer.From(item.LabelSelectors),
			VolumeSnapshotEnabled:       item.SnapshotVolumes,
		})
	}
	return &results
}
