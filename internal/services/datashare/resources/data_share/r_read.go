// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_share

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/account"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/synchronizationsetting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataShareRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataShare.SharesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	syncClient := meta.(*clients.Client).DataShare.SynchronizationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := share.ParseShareID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	accountId := account.NewAccountID(subscriptionId, id.ResourceGroupName, id.AccountName)

	d.Set("name", id.ShareName)
	d.Set("account_id", accountId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("kind", pointer.FromEnum(props.ShareKind))
			d.Set("description", props.Description)
			d.Set("terms", props.Terms)
		}
	}

	settings := make([]synchronizationsetting.ScheduledSynchronizationSetting, 0)
	snapshotSchedules, err := syncClient.ListByShareComplete(ctx, synchronizationsetting.NewShareID(id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ShareName))
	if err != nil {
		return fmt.Errorf("listing snapshot schedules for %s: %+v", *id, err)
	}
	for _, item := range snapshotSchedules.Items {
		if s, ok := item.(synchronizationsetting.ScheduledSynchronizationSetting); ok {
			settings = append(settings, s)
		}
	}
	if err := d.Set("snapshot_schedule", flattenAzureRmDataShareSnapshotSchedule(settings)); err != nil {
		return fmt.Errorf("setting `snapshot_schedule`: %+v", err)
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenAzureRmDataShareSnapshotSchedule(input []synchronizationsetting.ScheduledSynchronizationSetting) []any {
	output := make([]any, 0)

	for _, setting := range input {
		props := setting.Properties
		output = append(output, map[string]any{
			"name":       pointer.From(setting.Name),
			"recurrence": string(props.RecurrenceInterval),
			"start_time": props.SynchronizationTime,
		})
	}

	return output
}
