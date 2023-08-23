// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
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

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.AccountName,
		},

		"pool_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.PoolName,
		},

		"volume_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: netAppValidate.VolumeName,
		},

		"quota_target": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"quota_size_in_kib": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(4),
		},

		"quota_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(volumequotarules.PossibleValuesForType(), false),
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

			id := volumequotarules.NewVolumeQuotaRuleID(subscriptionId, model.ResourceGroupName, model.AccountName, model.CapacityPoolName, model.VolumeName, model.Name)
			volumeID := volumes.NewVolumeID(subscriptionId, model.ResourceGroupName, model.AccountName, model.CapacityPoolName, model.VolumeName)

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, id)
			if err != nil && existing.HttpResponse.StatusCode != http.StatusNotFound {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			// Performing some validations that are not possible in the schema
			if errorList := netAppValidate.ValidateNetAppVolumeQuotaRule(ctx, volumeID, metadata.Client, pointer.To(model)); len(errorList) > 0 {
				return fmt.Errorf("one or more issues found while performing deeper validations for %s:\n%+v", id, errorList)
			}

			parameters := volumequotarules.VolumeQuotaRule{
				Location: location.Normalize(model.Location),
				Properties: &volumequotarules.VolumeQuotaRulesProperties{
					QuotaSizeInKiBs: pointer.To(int64(model.QuotaSizeInKiB)),
					QuotaType:       pointer.To(volumequotarules.Type(model.QuotaType)),
					QuotaTarget:     utils.String(model.QuotaTarget),
				},
			}

			err = client.CreateThenPoll(ctx, id, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// Waiting for volume quota rule to be completely provisioned
			if err := waitForVolumeQuotaRuleCreateOrUpdate(ctx, client, id); err != nil {
				return err
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

			shouldUpdate := false
			update := volumequotarules.VolumeQuotaRulePatch{
				Properties: &volumequotarules.VolumeQuotaRulesProperties{},
			}

			if metadata.ResourceData.HasChange("quota_size_in_kib") {
				shouldUpdate = true
				update.Properties.QuotaSizeInKiBs = utils.Int64(int64(state.QuotaSizeInKiB))
			}

			if shouldUpdate {
				if err := client.UpdateThenPoll(ctx, pointer.From(id), update); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}

				// Waiting for volume quota rule to be completely provisioned
				if err := waitForVolumeQuotaRuleCreateOrUpdate(ctx, client, pointer.From(id)); err != nil {
					return err
				}
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
				Name:              id.VolumeQuotaRuleName,
				AccountName:       id.NetAppAccountName,
				CapacityPoolName:  id.CapacityPoolName,
				VolumeName:        id.VolumeName,
				Location:          location.NormalizeNilable(pointer.To(existing.Model.Location)),
				ResourceGroupName: id.ResourceGroupName,
				QuotaTarget:       pointer.From(existing.Model.Properties.QuotaTarget),
				QuotaSizeInKiB:    int64(pointer.From(existing.Model.Properties.QuotaSizeInKiBs)),
				QuotaType:         string(pointer.From(existing.Model.Properties.QuotaType)),
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

func waitForVolumeQuotaRuleCreateOrUpdate(ctx context.Context, client *volumequotarules.VolumeQuotaRulesClient, id volumequotarules.VolumeQuotaRuleId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"204", "404"},
		Target:                    []string{"200", "202"},
		Refresh:                   netappVolumeQuotaRuleStateRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to finish creating: %+v", id, err)
	}

	return nil
}

func netappVolumeQuotaRuleStateRefreshFunc(ctx context.Context, client *volumequotarules.VolumeQuotaRulesClient, id volumequotarules.VolumeQuotaRuleId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(res.HttpResponse) {
				return nil, "", fmt.Errorf("retrieving %s: %s", id, err)
			}
		}

		statusCode := "dropped connection"
		if res.HttpResponse != nil {
			statusCode = strconv.Itoa(res.HttpResponse.StatusCode)
		}
		return res, statusCode, nil
	}
}
