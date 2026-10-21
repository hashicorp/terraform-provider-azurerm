// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_share

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/synchronizationsetting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataShareDelete(d *pluginsdk.ResourceData, meta any) error {
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
