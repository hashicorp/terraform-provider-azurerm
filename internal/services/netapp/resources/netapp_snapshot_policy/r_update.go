// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_snapshot_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/snapshotpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model NetAppSnapshotPolicyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.NetApp.SnapshotPoliciesClient

			id, err := snapshotpolicies.ParseSnapshotPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			parameters := snapshotpolicies.SnapshotPolicyPatch{
				Location: pointer.To(location.Normalize(model.Location)),
				Name:     pointer.To(id.SnapshotPolicyName),
				Properties: &snapshotpolicies.SnapshotPolicyProperties{
					HourlySchedule:  expandNetAppSnapshotPolicyHourlySchedule(model.HourlySchedule),
					DailySchedule:   expandNetAppSnapshotPolicyDailySchedule(model.DailySchedule),
					WeeklySchedule:  expandNetAppSnapshotPolicyWeeklySchedule(model.WeeklySchedule),
					MonthlySchedule: expandNetAppSnapshotPolicyMonthlySchedule(model.MonthlySchedule),
					Enabled:         pointer.To(model.Enabled),
				},
				Tags: pointer.To(model.Tags),
			}

			if err = client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
