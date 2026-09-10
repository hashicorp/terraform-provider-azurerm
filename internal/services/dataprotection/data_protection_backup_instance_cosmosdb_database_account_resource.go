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

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name data_protection_backup_instance_cosmosdb_database_account -service-package-name dataprotection -properties "name" -compare-values "subscription_id:vault_id,resource_group_name:vault_id,backup_vault_name:vault_id"

type BackupInstanceCosmosDBDatabaseAccountModel struct {
	Name              string `tfschema:"name"`
	Location          string `tfschema:"location"`
	VaultId           string `tfschema:"vault_id"`
	BackupPolicyId    string `tfschema:"backup_policy_id"`
	CosmosDBAccountId string `tfschema:"cosmosdb_account_id"`
	ProtectionState   string `tfschema:"protection_state"`
}

type DataProtectionBackupInstanceCosmosDBDatabaseAccountResource struct{}

var (
	_ sdk.Resource             = DataProtectionBackupInstanceCosmosDBDatabaseAccountResource{}
	_ sdk.ResourceWithIdentity = DataProtectionBackupInstanceCosmosDBDatabaseAccountResource{}
)

func (r DataProtectionBackupInstanceCosmosDBDatabaseAccountResource) Identity() resourceids.ResourceId {
	return &backupinstanceresources.BackupInstanceId{}
}

func (r DataProtectionBackupInstanceCosmosDBDatabaseAccountResource) ResourceType() string {
	return "azurerm_data_protection_backup_instance_cosmosdb_database_account"
}

func (r DataProtectionBackupInstanceCosmosDBDatabaseAccountResource) ModelObject() interface{} {
	return &BackupInstanceCosmosDBDatabaseAccountModel{}
}

func (r DataProtectionBackupInstanceCosmosDBDatabaseAccountResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backupinstanceresources.ValidateBackupInstanceID
}

func (r DataProtectionBackupInstanceCosmosDBDatabaseAccountResource) Arguments() map[string]*pluginsdk.Schema {
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

		"cosmosdb_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&cosmosdb.DatabaseAccountId{}),
	}
}

func (r DataProtectionBackupInstanceCosmosDBDatabaseAccountResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"protection_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r DataProtectionBackupInstanceCosmosDBDatabaseAccountResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupInstanceCosmosDBDatabaseAccountModel
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
					FriendlyName: pointer.To(accountId.DatabaseAccountName),
					ObjectType:   "BackupInstance",
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

			return pollCosmosDBDatabaseAccountBackupInstance(ctx, client, id, backupinstanceresources.CurrentProtectionStateConfiguringProtection)
		},
	}
}

func (r DataProtectionBackupInstanceCosmosDBDatabaseAccountResource) Read() sdk.ResourceFunc {
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

func (r DataProtectionBackupInstanceCosmosDBDatabaseAccountResource) flatten(metadata sdk.ResourceMetaData, id *backupinstanceresources.BackupInstanceId, model *backupinstanceresources.BackupInstanceResource) error {
	state := BackupInstanceCosmosDBDatabaseAccountModel{
		Name:    id.BackupInstanceName,
		VaultId: backupvaultresources.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName).ID(),
	}
	if model != nil && model.Properties != nil {
		properties := model.Properties
		state.Location = location.NormalizeNilable(properties.DataSourceInfo.ResourceLocation)
		accountId, err := cosmosdb.ParseDatabaseAccountIDInsensitively(properties.DataSourceInfo.ResourceID)
		if err != nil {
			return err
		}
		state.CosmosDBAccountId = accountId.ID()

		policyId, err := basebackuppolicyresources.ParseBackupPolicyIDInsensitively(properties.PolicyInfo.PolicyId)
		if err != nil {
			return err
		}
		state.BackupPolicyId = policyId.ID()
		state.ProtectionState = pointer.FromEnum(properties.CurrentProtectionState)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}

func (r DataProtectionBackupInstanceCosmosDBDatabaseAccountResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient20260601
			id, err := backupinstanceresources.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			var model BackupInstanceCosmosDBDatabaseAccountModel
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

			if err := client.BackupInstancesCreateOrUpdateThenPoll(ctx, *id, parameters, backupinstanceresources.DefaultBackupInstancesCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return pollCosmosDBDatabaseAccountBackupInstance(ctx, client, *id, backupinstanceresources.CurrentProtectionStateUpdatingProtection)
		},
	}
}

func (r DataProtectionBackupInstanceCosmosDBDatabaseAccountResource) Delete() sdk.ResourceFunc {
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

func pollCosmosDBDatabaseAccountBackupInstance(ctx context.Context, client *backupinstanceresources.BackupInstanceResourcesClient, id backupinstanceresources.BackupInstanceId, pendingState backupinstanceresources.CurrentProtectionState) error {
	pollerType := custompollers.NewDataProtectionBackupInstance20260601Poller(client, id, backupinstanceresources.CurrentProtectionStateProtectionConfigured, []backupinstanceresources.CurrentProtectionState{pendingState})
	poller := pollers.NewPoller(pollerType, time.Minute, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %+v", id, err)
	}
	return nil
}
