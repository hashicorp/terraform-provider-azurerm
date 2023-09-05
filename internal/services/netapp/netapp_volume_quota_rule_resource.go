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
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-05-01/volumequotarules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2022-05-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	netAppModels "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetAppVolumeQuotaRuleResource struct{}

var _ sdk.Resource = NetAppVolumeQuotaRuleResource{}

func (r NetAppVolumeQuotaRuleResource) ModelObject() interface{} {
	return &netAppModels.NetAppVolumeQuotaRuleModel{}
}

func (r NetAppVolumeQuotaRuleResource) ResourceType() string {
	return "azurerm_netapp_volume_quota_rule"
}

func (r NetAppVolumeQuotaRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return volumequotarules.ValidateVolumeQuotaRuleID
}

func (r NetAppVolumeQuotaRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.VolumeQuotaRuleName,
		},

		"location": commonschema.Location(),

		"volume_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.VolumeID,
		},

		"quota_type": {
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validation.StringInSlice(volumequotarules.PossibleValuesForType(), false),
		},

		"quota_target": {
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"quota_size_in_kib": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(4),
		},
	}
}

func (r NetAppVolumeQuotaRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NetAppVolumeQuotaRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.VolumeQuotaRules
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model netAppModels.NetAppVolumeQuotaRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			volumeID, err := volumes.ParseVolumeID(model.VolumeID)
			if err != nil {
				return fmt.Errorf("error parsing volume id %s: %+v", model.VolumeID, err)
			}

			id := volumequotarules.NewVolumeQuotaRuleID(subscriptionId, volumeID.ResourceGroupName, volumeID.NetAppAccountName, volumeID.CapacityPoolName, volumeID.VolumeName, model.Name)

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, id)
			if err != nil && existing.HttpResponse.StatusCode != http.StatusNotFound {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			// Performing some validations that are not possible in the schema
			if errorList := netAppValidate.ValidateNetAppVolumeQuotaRule(ctx, pointer.From(volumeID), metadata.Client, pointer.To(model)); len(errorList) > 0 {
				return fmt.Errorf("one or more issues found while performing deeper validations for %s:\n%+v", id, errorList)
			}

			parameters := volumequotarules.VolumeQuotaRule{
				Location: location.Normalize(model.Location),
				Properties: &volumequotarules.VolumeQuotaRulesProperties{
					QuotaSizeInKiBs: pointer.To(model.QuotaSizeInKiB),
					QuotaType:       pointer.To(volumequotarules.Type(model.QuotaType)),
					QuotaTarget:     utils.String(model.QuotaTarget),
				},
			}

			err = client.CreateThenPoll(ctx, id, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r NetAppVolumeQuotaRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.VolumeQuotaRules

			id, err := volumequotarules.ParseVolumeQuotaRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state netAppModels.NetAppVolumeQuotaRuleModel
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

func (r NetAppVolumeQuotaRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {

			client := metadata.Client.NetApp.VolumeQuotaRules

			id, err := volumequotarules.ParseVolumeQuotaRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state netAppModels.NetAppVolumeQuotaRuleModel
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

			model := netAppModels.NetAppVolumeQuotaRuleModel{
				Name:           id.VolumeQuotaRuleName,
				VolumeID:       volumes.NewVolumeID(id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName).String(),
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

func (r NetAppVolumeQuotaRuleResource) Delete() sdk.ResourceFunc {
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
