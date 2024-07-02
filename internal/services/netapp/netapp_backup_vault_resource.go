// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-11-01/backupvaults"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetAppBackupVaultResource struct{}

var _ sdk.Resource = NetAppBackupVaultResource{}

func (r NetAppBackupVaultResource) ModelObject() interface{} {
	return &netAppModels.NetAppBackupVaultModel{}
}

func (r NetAppBackupVaultResource) ResourceType() string {
	return "azurerm_netapp_backup_vault"
}

func (r NetAppBackupVaultResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backupvaults.ValidateBackupVaultID
}

func (r NetAppBackupVaultResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.VolumeQuotaRuleName,
		},

		"location": commonschema.Location(),

		"account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.AccountName,
		},

		"tags": commonschema.Tags(),
	}
}

func (r NetAppBackupVaultResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetAppBackupVaultResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BackupVaultsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model netAppModels.NetAppBackupVaultModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			backupVaultID, err := backupvaults.ParseBackupVaultID(model.Id)
			if err != nil {
				return fmt.Errorf("error parsing backup vault id %s: %+v", model.Id, err)
			}

			id := backupvaults.NewBackupVaultID(subscriptionId, backupVaultID.ResourceGroupName, backupVaultID.NetAppAccountName, model.Name)

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(r.ResourceType(), id.ID())
			}

			parameters := backupvaults.BackupVault{
				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
			}

			err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r NetAppBackupVaultResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.BackupVaultsClient

			id, err := backupvaults.ParseBackupVaultID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state netAppModels.NetAppBackupVaultModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			metadata.Logger.Infof("Updating %s", id)

			if metadata.ResourceData.HasChange("Tags") {

				update := backupvaults.BackupVaultPatch{
					Tags: pointer.To(state.Tags),
				}

				if err := client.UpdateThenPoll(ctx, pointer.From(id), update); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}

				metadata.SetID(id)
			}

			return nil
		},
	}
}

func (r NetAppBackupVaultResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.NetApp.BackupVaultsClient

			id, err := backupvaults.ParseBackupVaultID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state netAppModels.NetAppBackupVaultModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, pointer.From(id))
			if err != nil {
				if existing.HttpResponse.StatusCode == http.StatusNotFound {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := netAppModels.NetAppBackupVaultModel{
				Name:     id.BackupVaultName,
				Location: location.NormalizeNilable(pointer.To(existing.Model.Location)),
				Tags:     pointer.From(existing.Model.Tags),
			}

			metadata.SetID(id)

			return metadata.Encode(&model)
		},
	}
}

func (r NetAppBackupVaultResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.NetApp.BackupVaultsClient

			id, err := backupvaults.ParseBackupVaultID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, pointer.From(id))
			if err != nil {
				if existing.HttpResponse.StatusCode == http.StatusNotFound {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			if err = client.DeleteThenPoll(ctx, pointer.From(id)); err != nil {
				return fmt.Errorf("deleting %s: %+v", pointer.From(id), err)
			}

			return nil
		},
	}
}
