// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_share

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/account"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/synchronizationsetting"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataShareCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataShare.SharesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	syncClient := meta.(*clients.Client).DataShare.SynchronizationClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	accountId, err := account.ParseAccountID(d.Get("account_id").(string))
	if err != nil {
		return err
	}

	id := share.NewShareID(subscriptionId, accountId.ResourceGroupName, accountId.AccountName, name)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_data_share", id.ID())
		}
	}

	share := share.Share{
		Properties: &share.ShareProperties{
			ShareKind:   pointer.ToEnum[share.ShareKind](d.Get("kind").(string)),
			Description: pointer.To(d.Get("description").(string)),
			Terms:       pointer.To(d.Get("terms").(string)),
		},
	}

	if _, err := client.Create(ctx, id, share); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	if d.HasChange("snapshot_schedule") {
		// only one dependent sync setting is allowed in one data share
		o, _ := d.GetChange("snapshot_schedule")
		if origins := o.([]any); len(origins) > 0 {
			origin := origins[0].(map[string]any)
			if originName, ok := origin["name"].(string); ok && originName != "" {
				syncId := synchronizationsetting.NewSynchronizationSettingID(id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ShareName, originName)
				if err := syncClient.DeleteThenPoll(ctx, syncId); err != nil {
					return fmt.Errorf("deleting datashare snapshot schedule %s: %+v", syncId, err)
				}
			}
		}
	}

	if snapshotSchedule := expandAzureRmDataShareSnapshotSchedule(d.Get("snapshot_schedule").([]any)); snapshotSchedule != nil {
		syncId := synchronizationsetting.NewSynchronizationSettingID(id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ShareName, d.Get("snapshot_schedule.0.name").(string))
		if _, err := syncClient.Create(ctx, syncId, snapshotSchedule); err != nil {
			return fmt.Errorf("creating datashare snapshot schedule %s: %+v", syncId, err)
		}
	}

	return resourceDataShareRead(d, meta)
}

func expandAzureRmDataShareSnapshotSchedule(input []any) *synchronizationsetting.ScheduledSynchronizationSetting {
	if len(input) == 0 {
		return nil
	}

	snapshotSchedule := input[0].(map[string]any)

	startTime, _ := time.Parse(time.RFC3339, snapshotSchedule["start_time"].(string))

	syncTime := date.Time{Time: startTime}.String()

	return &synchronizationsetting.ScheduledSynchronizationSetting{
		Properties: synchronizationsetting.ScheduledSynchronizationSettingProperties{
			RecurrenceInterval:  synchronizationsetting.RecurrenceInterval(snapshotSchedule["recurrence"].(string)),
			SynchronizationTime: syncTime,
		},
	}
}
