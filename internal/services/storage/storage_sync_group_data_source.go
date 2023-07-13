// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/storagesyncservicesresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/syncgroupresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceStorageSyncGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceStorageSyncGroupRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.StorageSyncName,
			},

			"storage_sync_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: storagesyncservicesresource.ValidateStorageSyncServiceID,
			},
		},
	}
}

func dataSourceStorageSyncGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serviceId, err := syncgroupresource.ParseStorageSyncServiceID(d.Get("storage_sync_id").(string))
	if err != nil {
		return err
	}

	id := syncgroupresource.NewSyncGroupID(serviceId.SubscriptionId, serviceId.ResourceGroupName, serviceId.StorageSyncServiceName, d.Get("name").(string))
	resp, err := client.SyncGroupsGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.SyncGroupName)
	d.Set("storage_sync_id", serviceId.ID())

	return nil
}
