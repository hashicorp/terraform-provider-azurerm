// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_snapshot_policy

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/snapshotpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetAppSnapshotPolicyCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetApp.SnapshotPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := snapshotpolicies.NewSnapshotPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id.ID(), err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_netapp_snapshot_policy", id.ID())
		}
	}

	parameters := snapshotpolicies.SnapshotPolicy{
		Location: location.Normalize(d.Get("location").(string)),
		Name:     pointer.To(id.SnapshotPolicyName),
		Properties: snapshotpolicies.SnapshotPolicyProperties{
			HourlySchedule:  expandNetAppSnapshotPolicyHourlySchedule(d.Get("hourly_schedule").([]any)),
			DailySchedule:   expandNetAppSnapshotPolicyDailySchedule(d.Get("daily_schedule").([]any)),
			WeeklySchedule:  expandNetAppSnapshotPolicyWeeklySchedule(d.Get("weekly_schedule").([]any)),
			MonthlySchedule: expandNetAppSnapshotPolicyMonthlySchedule(d.Get("monthly_schedule").([]any)),
			Enabled:         pointer.To(d.Get("enabled").(bool)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if _, err := client.Create(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// Waiting for snapshot policy be completely provisioned
	log.Printf("[DEBUG] Waiting for %s to complete", id)
	if err := waitForSnapshotPolicyCreation(ctx, client, id, d.Timeout(pluginsdk.TimeoutCreate)); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceNetAppSnapshotPolicyRead(d, meta)
}
