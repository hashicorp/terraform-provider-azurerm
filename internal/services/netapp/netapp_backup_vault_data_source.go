// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/backupvaults"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetAppBackupVaultDataSource struct{}

var _ sdk.DataSource = NetAppBackupVaultDataSource{}

func (r NetAppBackupVaultDataSource) ResourceType() string {
	return "azurerm_netapp_backup_vault"
}

func (r NetAppBackupVaultDataSource) ModelObject() interface{} {
	return &netAppModels.NetAppBackupVaultModel{}
}

func (r NetAppBackupVaultDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backupvaults.ValidateBackupVaultID
}

func (r NetAppBackupVaultDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"account_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (r NetAppBackupVaultDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"tags": commonschema.TagsDataSource(),
	}
}

func (r NetAppBackupVaultDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BackupVaultsClient

			var state netAppModels.NetAppBackupVaultModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			backupVaultID := backupvaults.NewBackupVaultID(metadata.Client.Account.SubscriptionId, state.ResourceGroupName, state.AccountName, state.Name)

			existing, err := client.Get(ctx, backupVaultID)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", backupVaultID)
				}
				return fmt.Errorf("retrieving %s: %v", backupVaultID, err)
			}

			if model := existing.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)
			}

			metadata.SetID(backupVaultID)

			return metadata.Encode(&state)
		},
	}
}
