// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2024-04-01/backupvaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2023-05-01/volumequotarules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
			// client := metadata.Client.NetApp.AccountClient
			// subscriptionId := metadata.Client.Account.SubscriptionId

			// var model netAppModels.NetAppBackupVaultModel
			// if err := metadata.Decode(&model); err != nil {
			// 	return fmt.Errorf("decoding: %+v", err)
			// }

			// backupVaultID, err := backupvaults.ParseBackupVaultID(model.Id)
			// if err != nil {
			// 	return fmt.Errorf("error parsing backup vault id %s: %+v", model.Id, err)
			// }

			// id := backupvaults.NewBackupVaultID(subscriptionId, backupVaultID.ResourceGroupName, model.Name)

			// metadata.Logger.Infof("Import check for %s", id)
			// existing, err := client.(ctx, id)
			// if err != nil {
			// 	if !response.WasNotFound(existing.HttpResponse) {
			// 		return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			// 	}
			// }

			// if !response.WasNotFound(existing.HttpResponse) {
			// 	return tf.ImportAsExistsError(r.ResourceType(), id.ID())
			// }

			// // Performing some validations that are not possible in the schema
			// if errorList := netAppValidate.ValidateNetAppVolumeQuotaRule(ctx, pointer.From(volumeID), metadata.Client, pointer.To(model)); len(errorList) > 0 {
			// 	return fmt.Errorf("one or more issues found while performing deeper validations for %s:\n%+v", id, errorList)
			// }

			// parameters := volumequotarules.VolumeQuotaRule{
			// 	Location: location.Normalize(model.Location),
			// 	Properties: &volumequotarules.VolumeQuotaRulesProperties{
			// 		QuotaSizeInKiBs: pointer.To(model.QuotaSizeInKiB),
			// 		QuotaType:       pointer.To(volumequotarules.Type(model.QuotaType)),
			// 		QuotaTarget:     utils.String(model.QuotaTarget),
			// 	},
			// }

			// err = client.CreateThenPoll(ctx, id, parameters)
			// if err != nil {
			// 	return fmt.Errorf("creating %s: %+v", id, err)
			// }

			// metadata.SetID(id)

			return nil
		},
	}
}

func (r NetAppBackupVaultResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.VolumeQuotaRules

			id, err := volumequotarules.ParseVolumeQuotaRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state netAppModels.NetAppBackupVaultModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			metadata.Logger.Infof("Updating %s", id)

			update := volumequotarules.VolumeQuotaRulePatch{
				Properties: &volumequotarules.VolumeQuotaRulesProperties{},
			}

			update.Properties.QuotaSizeInKiBs = utils.Int64(state.QuotaSizeInKiB)
			if err := client.UpdateThenPoll(ctx, pointer.From(id), update); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r NetAppBackupVaultResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.NetApp.VolumeQuotaRules

			id, err := volumequotarules.ParseVolumeQuotaRuleID(metadata.ResourceData.Id())
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

			volumeID := volumequotarules.NewVolumeID(id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName)

			model := netAppModels.NetAppBackupVaultModel{
				Name:           id.VolumeQuotaRuleName,
				VolumeID:       volumeID.ID(),
				Location:       location.NormalizeNilable(pointer.To(existing.Model.Location)),
				QuotaTarget:    pointer.From(existing.Model.Properties.QuotaTarget),
				QuotaSizeInKiB: pointer.From(existing.Model.Properties.QuotaSizeInKiBs),
				QuotaType:      string(pointer.From(existing.Model.Properties.QuotaType)),
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

			client := metadata.Client.NetApp.VolumeQuotaRules

			id, err := volumequotarules.ParseVolumeQuotaRuleID(metadata.ResourceData.Id())
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
