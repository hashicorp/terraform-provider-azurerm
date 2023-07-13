// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/storagesyncservicesresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/syncgroupresource"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageSyncGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageSyncGroupCreate,
		Read:   resourceStorageSyncGroupRead,
		Delete: resourceStorageSyncGroupDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := syncgroupresource.ParseSyncGroupID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageSyncName,
			},

			"storage_sync_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: storagesyncservicesresource.ValidateStorageSyncServiceID,
			},
		},
	}
}

func resourceStorageSyncGroupCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncGroupsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serviceId, err := syncgroupresource.ParseStorageSyncServiceID(d.Get("storage_sync_id").(string))
	if err != nil {
		return err
	}

	id := syncgroupresource.NewSyncGroupID(serviceId.SubscriptionId, serviceId.ResourceGroupName, serviceId.StorageSyncServiceName, d.Get("name").(string))
	existing, err := client.SyncGroupsGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_storage_sync_group", id.ID())
	}

	if _, err := client.SyncGroupsCreate(ctx, id, syncgroupresource.SyncGroupCreateParameters{}); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceStorageSyncGroupRead(d, meta)
}

func resourceStorageSyncGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := syncgroupresource.ParseSyncGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SyncGroupsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.SyncGroupName)

	serviceId := syncgroupresource.NewStorageSyncServiceID(id.SubscriptionId, id.ResourceGroupName, id.StorageSyncServiceName)
	d.Set("storage_sync_id", serviceId.ID())

	return nil
}

func resourceStorageSyncGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := syncgroupresource.ParseSyncGroupID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.SyncGroupsDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
