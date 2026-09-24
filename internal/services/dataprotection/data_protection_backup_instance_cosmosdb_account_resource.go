// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupvaultresources"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type BackupInstanceCosmosDBAccountModel struct {
	Name                        string `tfschema:"name"`
	DataProtectionBackupVaultId string `tfschema:"data_protection_backup_vault_id"`
	Location                    string `tfschema:"location"`
	BackupPolicyId              string `tfschema:"backup_policy_cosmosdb_account_id"`
	CosmosDBAccountId           string `tfschema:"cosmosdb_account_id"`
	ProtectionState             string `tfschema:"protection_state"`
}

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name data_protection_backup_instance_cosmosdb_account -service-package-name dataprotection -properties "name" -compare-values "subscription_id:data_protection_backup_vault_id,resource_group_name:data_protection_backup_vault_id,backup_vault_name:data_protection_backup_vault_id"

type DataProtectionBackupInstanceCosmosDBAccountResource struct{}

var (
	_ sdk.Resource             = DataProtectionBackupInstanceCosmosDBAccountResource{}
	_ sdk.ResourceWithIdentity = DataProtectionBackupInstanceCosmosDBAccountResource{}
)

func (r DataProtectionBackupInstanceCosmosDBAccountResource) Identity() resourceids.ResourceId {
	return &backupinstanceresources.BackupInstanceId{}
}

func (r DataProtectionBackupInstanceCosmosDBAccountResource) ResourceType() string {
	return "azurerm_data_protection_backup_instance_cosmosdb_account"
}

func (r DataProtectionBackupInstanceCosmosDBAccountResource) ModelObject() interface{} {
	return &BackupInstanceCosmosDBAccountModel{}
}

func (r DataProtectionBackupInstanceCosmosDBAccountResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backupinstanceresources.ValidateBackupInstanceID
}

func (r DataProtectionBackupInstanceCosmosDBAccountResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"data_protection_backup_vault_id": commonschema.ResourceIDReferenceRequiredForceNew(&backupvaultresources.BackupVaultId{}),

		"location": commonschema.Location(),

		"backup_policy_cosmosdb_account_id": commonschema.ResourceIDReferenceRequired(&basebackuppolicyresources.BackupPolicyId{}),

		"cosmosdb_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&cosmosdb.DatabaseAccountId{}),
	}
}

func (r DataProtectionBackupInstanceCosmosDBAccountResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"protection_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r DataProtectionBackupInstanceCosmosDBAccountResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupInstanceCosmosDBAccountModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.DataProtection.BackupInstanceClient20260601

			vaultId, err := backupvaultresources.ParseBackupVaultID(model.DataProtectionBackupVaultId)
			if err != nil {
				return err
			}

			id := backupinstanceresources.NewBackupInstanceID(vaultId.SubscriptionId, vaultId.ResourceGroupName, vaultId.BackupVaultName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.BackupInstancesGet(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			accountId, err := cosmosdb.ParseDatabaseAccountID(model.CosmosDBAccountId)
			if err != nil {
				return err
			}

			policyId, err := basebackuppolicyresources.ParseBackupPolicyID(model.BackupPolicyId)
			if err != nil {
				return err
			}

			parameters := backupinstanceresources.BackupInstanceResource{
				Properties: &backupinstanceresources.BackupInstance{
					DataSourceInfo: backupinstanceresources.Datasource{
						DatasourceType:   pointer.To("Microsoft.DocumentDB/databaseAccounts"),
						ObjectType:       pointer.To("Datasource"),
						ResourceID:       accountId.ID(),
						ResourceLocation: pointer.To(location.Normalize(model.Location)),
						ResourceName:     pointer.To(accountId.DatabaseAccountName),
						ResourceType:     pointer.To("Microsoft.DocumentDB/databaseAccounts"),
						ResourceUri:      pointer.To(accountId.ID()),
					},
					FriendlyName: pointer.To(id.BackupInstanceName),
					PolicyInfo: backupinstanceresources.PolicyInfo{
						PolicyId: policyId.ID(),
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

			// the built-in poller is for the LRO to be finished, but the service requires additional time to finish the configure for backup
			// Tracked on https://github.com/Azure/azure-rest-api-specs/issues/41986
			pollerType := custompollers.NewDataProtectionBackupInstance20260601Poller(client, id, backupinstanceresources.CurrentProtectionStateProtectionConfigured, []backupinstanceresources.CurrentProtectionState{
				backupinstanceresources.CurrentProtectionStateConfiguringProtection,
			})
			poller := pollers.NewPoller(pollerType, 1*time.Minute, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become available: %+v", id, err)
			}

			return nil
		},
	}
}

func (r DataProtectionBackupInstanceCosmosDBAccountResource) Read() sdk.ResourceFunc {
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

// flatten is shared by Read and the list resource, which encodes state directly from
// list API responses without a per-item BackupInstancesGet round trip.
func (r DataProtectionBackupInstanceCosmosDBAccountResource) flatten(metadata sdk.ResourceMetaData, id *backupinstanceresources.BackupInstanceId, model *backupinstanceresources.BackupInstanceResource) error {
	state := BackupInstanceCosmosDBAccountModel{
		Name:                        id.BackupInstanceName,
		DataProtectionBackupVaultId: backupvaultresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName).ID(),
	}

	if model != nil {
		if props := model.Properties; props != nil {
			state.Location = location.NormalizeNilable(props.DataSourceInfo.ResourceLocation)

			accountId, err := cosmosdb.ParseDatabaseAccountIDInsensitively(props.DataSourceInfo.ResourceID)
			if err != nil {
				return err
			}
			state.CosmosDBAccountId = accountId.ID()

			backupPolicyId, err := basebackuppolicyresources.ParseBackupPolicyIDInsensitively(props.PolicyInfo.PolicyId)
			if err != nil {
				return err
			}
			state.BackupPolicyId = backupPolicyId.ID()

			state.ProtectionState = pointer.FromEnum(props.CurrentProtectionState)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}

func (r DataProtectionBackupInstanceCosmosDBAccountResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient20260601

			id, err := backupinstanceresources.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BackupInstanceCosmosDBAccountModel
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
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			if metadata.ResourceData.HasChange("backup_policy_cosmosdb_account_id") {
				policyId, err := basebackuppolicyresources.ParseBackupPolicyID(model.BackupPolicyId)
				if err != nil {
					return err
				}
				parameters.Properties.PolicyInfo.PolicyId = policyId.ID()
			}

			if err := client.BackupInstancesCreateOrUpdateThenPoll(ctx, *id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			// the built-in poller is for the LRO to be finished, but the service requires additional time to finish the configure for backup
			// Tracked on https://github.com/Azure/azure-rest-api-specs/issues/41986
			pollerType := custompollers.NewDataProtectionBackupInstance20260601Poller(client, *id, backupinstanceresources.CurrentProtectionStateProtectionConfigured, []backupinstanceresources.CurrentProtectionState{
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

func (r DataProtectionBackupInstanceCosmosDBAccountResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient20260601

			id, err := backupinstanceresources.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.BackupInstancesDeleteThenPoll(ctx, *id, backupinstanceresources.DefaultBackupInstancesDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
