// Copyright (c) HashiCorp, Inc.
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backuppolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2022-12-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2023-06-01-preview/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type BackupInstancePostgreSQLFlexibleServerModel struct {
	Name           string `tfschema:"name"`
	Location       string `tfschema:"location"`
	VaultId        string `tfschema:"vault_id"`
	BackupPolicyId string `tfschema:"backup_policy_id"`
	DatabaseId     string `tfschema:"database_id"`
}

type DataProtectionBackupInstancePostgreSQLFlexibleServerResource struct{}

var _ sdk.Resource = DataProtectionBackupInstancePostgreSQLFlexibleServerResource{}

func (r DataProtectionBackupInstancePostgreSQLFlexibleServerResource) ResourceType() string {
	return "azurerm_data_protection_backup_instance_postgresql_flexible_server"
}

func (r DataProtectionBackupInstancePostgreSQLFlexibleServerResource) ModelObject() interface{} {
	return pointer.To(BackupInstancePostgreSQLFlexibleServerModel{})
}

func (r DataProtectionBackupInstancePostgreSQLFlexibleServerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backupinstances.ValidateBackupInstanceID
}

func (r DataProtectionBackupInstancePostgreSQLFlexibleServerResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": commonschema.Location(),

		"vault_id": commonschema.ResourceIDReferenceRequiredForceNew(pointer.To(backuppolicies.BackupVaultId{})),

		"backup_policy_id": commonschema.ResourceIDReferenceRequired(pointer.To(backuppolicies.BackupPolicyId{})),

		"database_id": commonschema.ResourceIDReferenceRequiredForceNew(pointer.To(databases.DatabaseId{})),
	}
}

func (r DataProtectionBackupInstancePostgreSQLFlexibleServerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataProtectionBackupInstancePostgreSQLFlexibleServerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BackupInstancePostgreSQLFlexibleServerModel
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

			databaseId, err := databases.ParseDatabaseID(model.DatabaseId)
			if err != nil {
				return err
			}

			serverId := servers.NewFlexibleServerID(databaseId.SubscriptionId, databaseId.ResourceGroupName, databaseId.FlexibleServerName)

			policyId, err := backuppolicies.ParseBackupPolicyID(model.BackupPolicyId)
			if err != nil {
				return err
			}

			parameters := backupinstances.BackupInstanceResource{
				Properties: &backupinstances.BackupInstance{
					DataSourceInfo: backupinstances.Datasource{
						DatasourceType:   pointer.To("Microsoft.DBforPostgreSQL/flexibleServers/databases"),
						ObjectType:       pointer.To("Datasource"),
						ResourceID:       databaseId.ID(),
						ResourceLocation: pointer.To(location.Normalize(model.Location)),
						ResourceName:     pointer.To(databaseId.DatabaseName),
						ResourceType:     pointer.To("Microsoft.DBforPostgreSQL/flexibleServers/databases"),
						ResourceUri:      pointer.To(""),
					},
					DataSourceSetInfo: &backupinstances.DatasourceSet{
						DatasourceType:   pointer.To("Microsoft.DBForPostgreSQL/flexibleServers"),
						ObjectType:       pointer.To("DatasourceSet"),
						ResourceID:       serverId.ID(),
						ResourceLocation: pointer.To(location.Normalize(model.Location)),
						ResourceName:     pointer.To(serverId.FlexibleServerName),
						ResourceType:     pointer.To("Microsoft.DBForPostgreSQL/flexibleServers"),
						ResourceUri:      pointer.To(""),
					},
					FriendlyName: pointer.To(id.BackupInstanceName),
					PolicyInfo: backupinstances.PolicyInfo{
						PolicyId: policyId.ID(),
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

func (r DataProtectionBackupInstancePostgreSQLFlexibleServerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient

			id, err := backupinstances.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, pointer.From(id))
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(pointer.From(id))
				}

				return fmt.Errorf("retrieving %s: %+v", pointer.From(id), err)
			}

			vaultId := backupinstances.NewBackupVaultID(id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName)

			state := BackupInstancePostgreSQLFlexibleServerModel{
				Name:    id.BackupInstanceName,
				VaultId: vaultId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Location = location.Normalize(pointer.From(props.DataSourceInfo.ResourceLocation))

					databaseId, err := databases.ParseDatabaseID(props.DataSourceInfo.ResourceID)
					if err != nil {
						return err
					}
					state.DatabaseId = databaseId.ID()

					backupPolicyId, err := backuppolicies.ParseBackupPolicyID(props.PolicyInfo.PolicyId)
					if err != nil {
						return err
					}
					state.BackupPolicyId = backupPolicyId.ID()
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DataProtectionBackupInstancePostgreSQLFlexibleServerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient

			id, err := backupinstances.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model BackupInstancePostgreSQLFlexibleServerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, pointer.From(id))
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("reading %s: %+v", pointer.From(id), err)
			}

			parameters := pointer.From(existing.Model)

			if metadata.ResourceData.HasChange("backup_policy_id") {
				policyId, err := backuppolicies.ParseBackupPolicyID(model.BackupPolicyId)
				if err != nil {
					return err
				}
				parameters.Properties.PolicyInfo.PolicyId = policyId.ID()
			}

			if err := client.CreateOrUpdateThenPoll(ctx, pointer.From(id), parameters, backupinstances.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r DataProtectionBackupInstancePostgreSQLFlexibleServerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DataProtection.BackupInstanceClient

			id, err := backupinstances.ParseBackupInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.DeleteThenPoll(ctx, pointer.From(id), backupinstances.DefaultDeleteOperationOptions())
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", pointer.From(id), err)
			}

			return nil
		},
	}
}
