// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datashare

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/account"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/synchronizationsetting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceDataShare() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDataShareRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ShareName(),
			},

			"account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: account.ValidateAccountID,
			},

			"kind": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"snapshot_schedule": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"recurrence": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"start_time": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"terms": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDataShareRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.SharesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	syncClient := meta.(*clients.Client).DataShare.SynchronizationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	accountId, err := account.ParseAccountID(d.Get("account_id").(string))
	if err != nil {
		return err
	}

	id := share.NewShareID(subscriptionId, accountId.ResourceGroupName, accountId.AccountName, name)

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", name)
	d.Set("account_id", accountId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("kind", string(pointer.From(props.ShareKind)))
			d.Set("description", props.Description)
			d.Set("terms", props.Terms)
		}
	}

	settings := make([]synchronizationsetting.ScheduledSynchronizationSetting, 0)
	snapshotSchedules, err := syncClient.ListByShareComplete(ctx, synchronizationsetting.NewShareID(id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ShareName))
	if err != nil {
		return fmt.Errorf("listing snapshot schedules for %s: %+v", id, err)
	}
	for _, item := range snapshotSchedules.Items {
		if s, ok := item.(synchronizationsetting.ScheduledSynchronizationSetting); ok {
			settings = append(settings, s)
		}
	}

	if err := d.Set("snapshot_schedule", flattenDataShareDataSourceSnapshotSchedule(settings)); err != nil {
		return fmt.Errorf("setting `snapshot_schedule`: %+v", err)
	}

	return nil
}

func flattenDataShareDataSourceSnapshotSchedule(input []synchronizationsetting.ScheduledSynchronizationSetting) []interface{} {
	output := make([]interface{}, 0)

	for _, setting := range input {
		props := setting.Properties
		name := ""
		if props.UserName != nil {
			name = *props.UserName
		}

		output = append(output, map[string]interface{}{
			"name":       name,
			"recurrence": string(props.RecurrenceInterval),
			"start_time": props.SynchronizationTime,
		})
	}

	return output
}
