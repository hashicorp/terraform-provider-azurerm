// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_quota_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/volumequotarules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/models"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 120 * time.Minute,
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

			if metadata.ResourceData.HasChange("quota_size_in_kib") {
				update := volumequotarules.VolumeQuotaRulePatch{
					Properties: &volumequotarules.VolumeQuotaRulesProperties{},
				}

				update.Properties.QuotaSizeInKiBs = pointer.To(state.QuotaSizeInKiB)

				if err := client.UpdateThenPoll(ctx, pointer.From(id), update); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}

				// Waiting for quota rule to be updated
				if err := waitForQuotaRuleCreateOrUpdate(ctx, client, pointer.From(id)); err != nil {
					return fmt.Errorf("waiting update %s: %+v", id, err)
				}
			}

			return nil
		},
	}
}
