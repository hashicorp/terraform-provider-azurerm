// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datashare

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/account"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/synchronizationsetting"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataShare() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataShareCreateUpdate,
		Read:   resourceDataShareRead,
		Update: resourceDataShareCreateUpdate,
		Delete: resourceDataShareDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := share.ParseShareID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ShareName(),
			},

			"account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: account.ValidateAccountID,
			},

			"kind": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(share.ShareKindCopyBased),
					string(share.ShareKindInPlace),
				}, false),
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"snapshot_schedule": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.SnapshotScheduleName(),
						},

						"recurrence": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(synchronizationsetting.RecurrenceIntervalDay),
								string(synchronizationsetting.RecurrenceIntervalHour),
							}, false),
						},

						"start_time": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							ValidateFunc:     validation.IsRFC3339Time,
							DiffSuppressFunc: suppress.RFC3339Time,
						},
					},
				},
			},

			"terms": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDataShareCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
			ShareKind:   pointer.To(share.ShareKind(d.Get("kind").(string))),
			Description: utils.String(d.Get("description").(string)),
			Terms:       utils.String(d.Get("terms").(string)),
		},
	}

	if _, err := client.Create(ctx, id, share); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if d.HasChange("snapshot_schedule") {
		// only one dependent sync setting is allowed in one data share
		o, _ := d.GetChange("snapshot_schedule")
		if origins := o.([]interface{}); len(origins) > 0 {
			origin := origins[0].(map[string]interface{})
			if originName, ok := origin["name"].(string); ok && originName != "" {
				syncId := synchronizationsetting.NewSynchronizationSettingID(id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ShareName, originName)
				if err := syncClient.DeleteThenPoll(ctx, syncId); err != nil {
					return fmt.Errorf("deleting datashare snapshot schedule %s: %+v", syncId, err)
				}
			}
		}
	}

	if snapshotSchedule := expandAzureRmDataShareSnapshotSchedule(d.Get("snapshot_schedule").([]interface{})); snapshotSchedule != nil {
		syncId := synchronizationsetting.NewSynchronizationSettingID(id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ShareName, d.Get("snapshot_schedule.0.name").(string))
		if _, err := syncClient.Create(ctx, syncId, snapshotSchedule); err != nil {
			return fmt.Errorf("creating datashare snapshot schedule %s: %+v", syncId, err)
		}
	}

	return resourceDataShareRead(d, meta)
}

func resourceDataShareRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			d.Set("kind", string(pointer.From(props.ShareKind)))
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

	return nil
}

func resourceDataShareDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.SharesClient
	syncClient := meta.(*clients.Client).DataShare.SynchronizationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := share.ParseShareID(d.Id())
	if err != nil {
		return err
	}

	// sync setting will not automatically be deleted after the data share is deleted
	if _, ok := d.GetOk("snapshot_schedule"); ok {
		syncId := synchronizationsetting.NewSynchronizationSettingID(id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ShareName, d.Get("snapshot_schedule.0.name").(string))
		if err := syncClient.DeleteThenPoll(ctx, syncId); err != nil {
			return fmt.Errorf("deleting datashare snapshot schedule %s: %+v", syncId, err)
		}
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandAzureRmDataShareSnapshotSchedule(input []interface{}) *synchronizationsetting.ScheduledSynchronizationSetting {
	if len(input) == 0 {
		return nil
	}

	snapshotSchedule := input[0].(map[string]interface{})

	startTime, _ := time.Parse(time.RFC3339, snapshotSchedule["start_time"].(string))

	syncTime := date.Time{Time: startTime}.String()

	return &synchronizationsetting.ScheduledSynchronizationSetting{
		Properties: synchronizationsetting.ScheduledSynchronizationSettingProperties{
			RecurrenceInterval:  synchronizationsetting.RecurrenceInterval(snapshotSchedule["recurrence"].(string)),
			SynchronizationTime: syncTime,
		},
	}
}

func flattenAzureRmDataShareSnapshotSchedule(input []synchronizationsetting.ScheduledSynchronizationSetting) []interface{} {
	output := make([]interface{}, 0)

	for _, setting := range input {
		props := setting.Properties
		name := ""
		if setting.Name != nil {
			name = *setting.Name
		}

		output = append(output, map[string]interface{}{
			"name":       name,
			"recurrence": string(props.RecurrenceInterval),
			"start_time": props.SynchronizationTime,
		})
	}

	return output
}
