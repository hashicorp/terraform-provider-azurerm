// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/basebackuppolicyresources"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/custompollers"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type BackupInstanceDataLakeStorageModel struct {
	Name                        string   `tfschema:"name"`
	DataProtectionBackupVaultId string   `tfschema:"data_protection_backup_vault_id"`
	Location                    string   `tfschema:"location"`
	BackupPolicyId              string   `tfschema:"backup_policy_id"`
	StorageAccountId            string   `tfschema:"storage_account_id"`
	StorageContainerNames       []string `tfschema:"storage_container_names"`
	ProtectionState             string   `tfschema:"protection_state"`
}

type DataProtectionBackupInstanceDataLakeStorageResource struct{}

var _ sdk.Resource = DataProtectionBackupInstanceDataLakeStorageResource{}

func (r DataProtectionBackupInstanceDataLakeStorageResource) ResourceType() string {
	return "azurerm_data_protection_backup_instance_data_lake_storage"
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) ModelObject() interface{} {
	return &BackupInstanceDataLakeStorageModel{}
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backupinstanceresources.ValidateBackupInstanceID
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"data_protection_backup_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&backupvaultresources.BackupVaultId{}),

		"location": commonschema.Location(),

		"backup_policy_id": commonschema.ResourceIDReferenceRequired(&basebackuppolicyresources.BackupPolicyId{}),

		"storage_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.StorageAccountId{}),

		"storage_container_names": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: storageValidate.StorageContainerName,
			},
		},
	}
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"protection_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupInstanceDataLakeStorageModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.DataProtection.BackupInstanceClient

			vaultId, err := backupvaultresources.ParseBackupVaultID(model.DataProtectionBackupVaultId)
			if err != nil {
				return err
			}

			id := backupinstanceresources.NewBackupInstanceID(vaultId.SubscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, model.Name)

			existing, err := client.BackupInstancesGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			storageAccountId, err := commonids.ParseStorageAccountID(model.StorageAccountId)
			if err != nil {
				return err
			}

			policyId, err := basebackuppolicyresources.ParseBackupPolicyID(model.BackupPolicyId)
			if err != nil {
				return err
			}

			dataSourceLocation := pointer.To(location.Normalize(model.Location))

			parameters := backupinstanceresources.BackupInstanceResource{
				Properties: &backupinstanceresources.BackupInstance{
					DataSourceInfo: backupinstanceresources.Datasource{
						DatasourceType:   pointer.To("Microsoft.Storage/storageAccounts/adlsBlobServices"),
						ObjectType:       pointer.To("Datasource"),
						ResourceID:       storageAccountId.ID(),
						ResourceLocation: dataSourceLocation,
						ResourceName:     pointer.To(storageAccountId.StorageAccountName),
						ResourceType:     pointer.To("Microsoft.Storage/storageAccounts"),
						ResourceUri:      pointer.To(storageAccountId.ID()),
					},
					DataSourceSetInfo: &backupinstanceresources.DatasourceSet{
						DatasourceType:   pointer.To("Microsoft.Storage/storageAccounts"),
						ObjectType:       pointer.To("DatasourceSet"),
						ResourceID:       storageAccountId.ID(),
						ResourceLocation: dataSourceLocation,
						ResourceName:     pointer.To(storageAccountId.StorageAccountName),
						ResourceType:     pointer.To("Microsoft.Storage/storageAccounts"),
						ResourceUri:      pointer.To(storageAccountId.ID()),
					},
					FriendlyName: pointer.To(id.BackupInstanceName),
					PolicyInfo: backupinstanceresources.PolicyInfo{
						PolicyId: policyId.ID(),
						PolicyParameters: &backupinstanceresources.PolicyParameters{
							BackupDatasourceParametersList: &[]backupinstanceresources.BackupDatasourceParameters{
								backupinstanceresources.AdlsBlobBackupDatasourceParameters{
									ContainersList: model.StorageContainerNames,
								},
							},
						},
					},
				},
			}

			if err := client.BackupInstancesCreateOrUpdateThenPoll(ctx, id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			pollerType := custompollers.NewDataProtectionBackupInstancePoller(client, id, backupinstanceresources.CurrentProtectionStateProtectionConfigured, []backupinstanceresources.CurrentProtectionState{
				backupinstanceresources.CurrentProtectionStateConfiguringProtection,
			})
			poller := pollers.NewPoller(pollerType, 1*time.Minute, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become available: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient

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

			vaultId := backupvaultresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)

			state := BackupInstanceDataLakeStorageModel{
				Name:                        id.BackupInstanceName,
				DataProtectionBackupVaultId: vaultId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Location = location.NormalizeNilable(props.DataSourceInfo.ResourceLocation)

					storageAccountId, err := commonids.ParseStorageAccountIDInsensitively(props.DataSourceInfo.ResourceID)
					if err != nil {
						return err
					}
					state.StorageAccountId = storageAccountId.ID()

					backupPolicyId, err := basebackuppolicyresources.ParseBackupPolicyIDInsensitively(props.PolicyInfo.PolicyId)
					if err != nil {
						return err
					}
					state.BackupPolicyId = backupPolicyId.ID()

					state.ProtectionState = pointer.FromEnum(props.CurrentProtectionState)

					if policyParas := props.PolicyInfo.PolicyParameters; policyParas != nil {
						if dataStoreParas := policyParas.BackupDatasourceParametersList; dataStoreParas != nil {
							if dsp := pointer.From(dataStoreParas); len(dsp) > 0 {
								if parameter, ok := dsp[0].(backupinstanceresources.AdlsBlobBackupDatasourceParameters); ok {
									state.StorageContainerNames = parameter.ContainersList
								}
							}
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient

			id, err := backupinstanceresources.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BackupInstanceDataLakeStorageModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.BackupInstancesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			parameters := *existing.Model
			if parameters.Properties == nil {
				parameters.Properties = &backupinstanceresources.BackupInstance{
					PolicyInfo: backupinstanceresources.PolicyInfo{},
				}
			}

			if metadata.ResourceData.HasChange("backup_policy_id") {
				policyId, err := basebackuppolicyresources.ParseBackupPolicyID(model.BackupPolicyId)
				if err != nil {
					return err
				}
				parameters.Properties.PolicyInfo.PolicyId = policyId.ID()
			}

			if metadata.ResourceData.HasChange("storage_container_names") {
				parameters.Properties.PolicyInfo.PolicyParameters = &backupinstanceresources.PolicyParameters{
					BackupDatasourceParametersList: &[]backupinstanceresources.BackupDatasourceParameters{
						backupinstanceresources.AdlsBlobBackupDatasourceParameters{
							ContainersList: model.StorageContainerNames,
						},
					},
				}
			}

			if err := client.BackupInstancesCreateOrUpdateThenPoll(ctx, *id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			pollerType := custompollers.NewDataProtectionBackupInstancePoller(client, *id, backupinstanceresources.CurrentProtectionStateProtectionConfigured, []backupinstanceresources.CurrentProtectionState{
				backupinstanceresources.CurrentProtectionStateUpdatingProtection,
			})
			poller := pollers.NewPoller(pollerType, 1*time.Minute, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become available: %+v", id, err)
			}

			return nil
		},
	}
}

func (r DataProtectionBackupInstanceDataLakeStorageResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient

			id, err := backupinstanceresources.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.BackupInstancesDeleteThenPoll(ctx, *id, backupinstanceresources.DefaultBackupInstancesDeleteOperationOptions())
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
