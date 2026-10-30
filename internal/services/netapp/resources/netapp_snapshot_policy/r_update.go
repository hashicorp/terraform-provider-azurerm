// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_snapshot_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/snapshotpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetAppSnapshotPolicyUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetApp.SnapshotPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := snapshotpolicies.ParseSnapshotPolicyID(d.Id())
	if err != nil {
		return err
	}

	parameters := snapshotpolicies.SnapshotPolicyPatch{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Name:     pointer.To(id.SnapshotPolicyName),
		Properties: &snapshotpolicies.SnapshotPolicyProperties{
			HourlySchedule:  expandNetAppSnapshotPolicyHourlySchedule(d.Get("hourly_schedule").([]any)),
			DailySchedule:   expandNetAppSnapshotPolicyDailySchedule(d.Get("daily_schedule").([]any)),
			WeeklySchedule:  expandNetAppSnapshotPolicyWeeklySchedule(d.Get("weekly_schedule").([]any)),
			MonthlySchedule: expandNetAppSnapshotPolicyMonthlySchedule(d.Get("monthly_schedule").([]any)),
			Enabled:         pointer.To(d.Get("enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err = client.UpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceNetAppSnapshotPolicyRead(d, meta)
}
