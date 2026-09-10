// Copyright IBM Corp. 2014, 2026
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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/elasticsan/2023-01-01/volumegroups"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name data_protection_backup_instance_elastic_san_volume_group -service-package-name dataprotection -properties "name" -compare-values "subscription_id:vault_id,resource_group_name:vault_id,backup_vault_name:vault_id"

type BackupInstanceElasticSanVolumeGroupModel struct {
	Name                      string `tfschema:"name"`
	Location                  string `tfschema:"location"`
	VaultId                   string `tfschema:"vault_id"`
	BackupPolicyId            string `tfschema:"backup_policy_id"`
	ElasticSanVolumeGroupId   string `tfschema:"elastic_san_volume_group_id"`
	SnapshotResourceGroupName string `tfschema:"snapshot_resource_group_name"`
	VolumeName                string `tfschema:"volume_name"`
	ProtectionState           string `tfschema:"protection_state"`
}

type DataProtectionBackupInstanceElasticSanVolumeGroupResource struct{}

var (
	_ sdk.Resource             = DataProtectionBackupInstanceElasticSanVolumeGroupResource{}
	_ sdk.ResourceWithIdentity = DataProtectionBackupInstanceElasticSanVolumeGroupResource{}
)

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) Identity() resourceids.ResourceId {
	return &backupinstanceresources.BackupInstanceId{}
}

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) ResourceType() string {
	return "azurerm_data_protection_backup_instance_elastic_san_volume_group"
}

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) ModelObject() interface{} {
	return &BackupInstanceElasticSanVolumeGroupModel{}
}

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backupinstanceresources.ValidateBackupInstanceID
}

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) Arguments() map[string]*pluginsdk.Schema {
	snapshotResourceGroupName := commonschema.ResourceGroupName()
	snapshotResourceGroupName.ForceNew = true

	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&backupvaultresources.BackupVaultId{}),

		"backup_policy_id": commonschema.ResourceIDReferenceRequired(&basebackuppolicyresources.BackupPolicyId{}),

		"elastic_san_volume_group_id": commonschema.ResourceIDReferenceRequiredForceNew(&volumegroups.VolumeGroupId{}),

		"snapshot_resource_group_name": snapshotResourceGroupName,

		"volume_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"protection_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupInstanceElasticSanVolumeGroupModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.DataProtection.BackupInstanceClient20260601
			vaultId, err := backupvaultresources.ParseBackupVaultID(model.VaultId)
			if err != nil {
				return err
			}
			id := backupinstanceresources.NewBackupInstanceID(vaultId.SubscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.BackupInstancesGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			volumeGroupId, err := volumegroups.ParseVolumeGroupID(model.ElasticSanVolumeGroupId)
			if err != nil {
				return err
			}
			policyId, err := basebackuppolicyresources.ParseBackupPolicyID(model.BackupPolicyId)
			if err != nil {
				return err
			}

			elasticSanId := volumegroups.NewElasticSanID(volumeGroupId.SubscriptionId, volumeGroupId.ResourceGroupName, volumeGroupId.ElasticSanName)
			snapshotResourceGroupId := commonids.NewResourceGroupID(metadata.Client.Account.SubscriptionId, model.SnapshotResourceGroupName)
			resourceLocation := pointer.To(location.Normalize(model.Location))
			parameters := backupinstanceresources.BackupInstanceResource{
				Properties: &backupinstanceresources.BackupInstance{
					DataSourceInfo: backupinstanceresources.Datasource{
						DatasourceType:   pointer.To("Microsoft.ElasticSan/elasticSans/volumeGroups"),
						ObjectType:       pointer.To("Datasource"),
						ResourceID:       volumeGroupId.ID(),
						ResourceLocation: resourceLocation,
						ResourceName:     pointer.To(volumeGroupId.VolumeGroupName),
						ResourceType:     pointer.To("Microsoft.ElasticSan/elasticSans/volumeGroups"),
						ResourceUri:      pointer.To(volumeGroupId.ID()),
					},
					DataSourceSetInfo: &backupinstanceresources.DatasourceSet{
						DatasourceType:   pointer.To("Microsoft.ElasticSan/elasticSans/volumeGroups"),
						ObjectType:       pointer.To("DatasourceSet"),
						ResourceID:       elasticSanId.ID(),
						ResourceLocation: resourceLocation,
						ResourceName:     pointer.To(elasticSanId.ElasticSanName),
						ResourceType:     pointer.To("Microsoft.ElasticSan/elasticSans"),
						ResourceUri:      pointer.To(elasticSanId.ID()),
					},
					FriendlyName: pointer.To(id.BackupInstanceName),
					ObjectType:   "BackupInstance",
					PolicyInfo: backupinstanceresources.PolicyInfo{
						PolicyId: policyId.ID(),
						PolicyParameters: &backupinstanceresources.PolicyParameters{
							BackupDatasourceParametersList: expandElasticSanVolumeGroupBackupDatasourceParameters(model.VolumeName),
							DataStoreParametersList: &[]backupinstanceresources.DataStoreParameters{
								backupinstanceresources.AzureOperationalStoreParameters{
									ResourceGroupId: pointer.To(snapshotResourceGroupId.ID()),
									DataStoreType:   backupinstanceresources.DataStoreTypesOperationalStore,
									ObjectType:      "AzureOperationalStoreParameters",
								},
							},
						},
					},
				},
			}

			if err := client.BackupInstancesCreateOrUpdateCallbackThenPoll(ctx, id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions(), metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return pollElasticSanVolumeGroupBackupInstance(ctx, client, id, backupinstanceresources.CurrentProtectionStateConfiguringProtection)
		},
	}
}

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient20260601
			id, err := backupinstanceresources.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.BackupInstancesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) flatten(metadata sdk.ResourceMetaData, id *backupinstanceresources.BackupInstanceId, model *backupinstanceresources.BackupInstanceResource) error {
	state := BackupInstanceElasticSanVolumeGroupModel{
		Name:    id.BackupInstanceName,
		VaultId: backupvaultresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName).ID(),
	}
	if model != nil && model.Properties != nil {
		properties := model.Properties
		state.Location = location.NormalizeNilable(properties.DataSourceInfo.ResourceLocation)
		volumeGroupId, err := volumegroups.ParseVolumeGroupIDInsensitively(properties.DataSourceInfo.ResourceID)
		if err != nil {
			return err
		}
		state.ElasticSanVolumeGroupId = volumeGroupId.ID()

		policyId, err := basebackuppolicyresources.ParseBackupPolicyIDInsensitively(properties.PolicyInfo.PolicyId)
		if err != nil {
			return err
		}
		state.BackupPolicyId = policyId.ID()
		state.ProtectionState = pointer.FromEnum(properties.CurrentProtectionState)

		if policyParameters := properties.PolicyInfo.PolicyParameters; policyParameters != nil {
			if dataStoreParameters := policyParameters.DataStoreParametersList; dataStoreParameters != nil {
				for _, item := range *dataStoreParameters {
					if parameter, ok := item.(backupinstanceresources.AzureOperationalStoreParameters); ok && parameter.ResourceGroupId != nil {
						resourceGroupId, err := commonids.ParseResourceGroupIDInsensitively(*parameter.ResourceGroupId)
						if err != nil {
							return err
						}
						state.SnapshotResourceGroupName = resourceGroupId.ResourceGroupName
						break
					}
				}
			}
			state.VolumeName, err = flattenElasticSanVolumeGroupBackupDatasourceParameters(policyParameters.BackupDatasourceParametersList)
			if err != nil {
				return err
			}
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient20260601
			id, err := backupinstanceresources.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BackupInstanceElasticSanVolumeGroupModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			existing, err := client.BackupInstancesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			parameters := *existing.Model
			if metadata.ResourceData.HasChange("backup_policy_id") {
				policyId, err := basebackuppolicyresources.ParseBackupPolicyID(model.BackupPolicyId)
				if err != nil {
					return err
				}
				parameters.Properties.PolicyInfo.PolicyId = policyId.ID()
			}
			if metadata.ResourceData.HasChange("volume_name") {
				if parameters.Properties.PolicyInfo.PolicyParameters == nil {
					parameters.Properties.PolicyInfo.PolicyParameters = &backupinstanceresources.PolicyParameters{}
				}
				parameters.Properties.PolicyInfo.PolicyParameters.BackupDatasourceParametersList = expandElasticSanVolumeGroupBackupDatasourceParameters(model.VolumeName)
			}

			if err := client.BackupInstancesCreateOrUpdateThenPoll(ctx, *id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return pollElasticSanVolumeGroupBackupInstance(ctx, client, *id, backupinstanceresources.CurrentProtectionStateUpdatingProtection)
		},
	}
}

func (r DataProtectionBackupInstanceElasticSanVolumeGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient20260601
			id, err := backupinstanceresources.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			if err := client.BackupInstancesDeleteThenPoll(ctx, *id, backupinstanceresources.DefaultBackupInstancesDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func expandElasticSanVolumeGroupBackupDatasourceParameters(volumeName string) *[]backupinstanceresources.BackupDatasourceParameters {
	return &[]backupinstanceresources.BackupDatasourceParameters{
		backupinstanceresources.GenericBackupDatasourceParameters{
			ObjectType:        "GenericBackupDatasourceParameters",
			ResourceSelectors: []string{volumeName},
		},
	}
}

func flattenElasticSanVolumeGroupBackupDatasourceParameters(input *[]backupinstanceresources.BackupDatasourceParameters) (string, error) {
	if input == nil {
		return "", nil
	}
	for _, item := range *input {
		parameter, ok := item.(backupinstanceresources.GenericBackupDatasourceParameters)
		if !ok {
			continue
		}
		if len(parameter.ResourceSelectors) != 1 {
			return "", fmt.Errorf("flattening Elastic SAN volume selector: expected one selector but got %d", len(parameter.ResourceSelectors))
		}
		return parameter.ResourceSelectors[0], nil
	}
	return "", nil
}

func pollElasticSanVolumeGroupBackupInstance(ctx context.Context, client *backupinstanceresources.BackupInstanceResourcesClient, id backupinstanceresources.BackupInstanceId, pendingState backupinstanceresources.CurrentProtectionState) error {
	pollerType := custompollers.NewDataProtectionBackupInstance20260601Poller(client, id, backupinstanceresources.CurrentProtectionStateProtectionConfigured, []backupinstanceresources.CurrentProtectionState{pendingState})
	poller := pollers.NewPoller(pollerType, time.Minute, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %+v", id, err)
	}
	return nil
}
