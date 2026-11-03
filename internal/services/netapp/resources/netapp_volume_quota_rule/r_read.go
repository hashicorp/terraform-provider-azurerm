// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_quota_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumequotarules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetApp.VolumeQuotaRules

			id, err := volumequotarules.ParseVolumeQuotaRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state models.NetAppVolumeQuotaRuleModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, pointer.From(id))
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			volumeID := volumequotarules.NewVolumeID(id.SubscriptionId, id.ResourceGroupName, id.NetAppAccountName, id.CapacityPoolName, id.VolumeName)

			model := models.NetAppVolumeQuotaRuleModel{
				Name:           id.VolumeQuotaRuleName,
				VolumeID:       volumeID.ID(),
				Location:       location.Normalize(existing.Model.Location),
				QuotaTarget:    pointer.From(existing.Model.Properties.QuotaTarget),
				QuotaSizeInKiB: pointer.From(existing.Model.Properties.QuotaSizeInKiBs),
				QuotaType:      pointer.FromEnum(existing.Model.Properties.QuotaType),
			}

			metadata.SetID(id)

			return metadata.Encode(&model)
		},
	}
}

func waitForQuotaRuleCreateOrUpdate(ctx context.Context, client *volumequotarules.VolumeQuotaRulesClient, id volumequotarules.VolumeQuotaRuleId) error {
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
		Refresh:                   netAppVolumeQuotaRuleStateRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to finish creating/updating: %+v", id, err)
	}

	return nil
}

func netAppVolumeQuotaRuleStateRefreshFunc(ctx context.Context, client *volumequotarules.VolumeQuotaRulesClient, id volumequotarules.VolumeQuotaRuleId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(res.HttpResponse) {
				return nil, "", fmt.Errorf("retrieving %s: %s", id, err)
			}
		}

		statusCode := "dropped connection"
		if res.HttpResponse != nil {
			statusCode = fmt.Sprintf("%d", res.HttpResponse.StatusCode)
		}
		return res, statusCode, nil
	}
}

func waitForQuotaRuleDelete(ctx context.Context, client *volumequotarules.VolumeQuotaRulesClient, id volumequotarules.VolumeQuotaRuleId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 5,
		Delay:                     10 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"200", "202"},
		Target:                    []string{"204", "404"},
		Refresh:                   netAppVolumeQuotaRuleStateRefreshFunc(ctx, client, id),
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to finish deleting: %+v", id, err)
	}

	return nil
}
