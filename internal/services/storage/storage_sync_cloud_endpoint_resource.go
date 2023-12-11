// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/cloudendpointresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/syncgroupresource"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceStorageSyncCloudEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageSyncCloudEndpointCreate,
		Read:   resourceStorageSyncCloudEndpointRead,
		Delete: resourceStorageSyncCloudEndpointDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := cloudendpointresource.ParseCloudEndpointID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(45 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(45 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageSyncName,
			},

			"storage_sync_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: syncgroupresource.ValidateSyncGroupID,
			},

			"file_share_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageShareName,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"storage_account_tenant_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceStorageSyncCloudEndpointCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncCloudEndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	groupId, err := syncgroupresource.ParseSyncGroupID(d.Get("storage_sync_group_id").(string))
	if err != nil {
		return err
	}

	id := cloudendpointresource.NewCloudEndpointID(groupId.SubscriptionId, groupId.ResourceGroupName, groupId.StorageSyncServiceName, groupId.SyncGroupName, d.Get("name").(string))
	existing, err := client.CloudEndpointsGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_storage_sync_cloud_endpoint", id.ID())
	}

	parameters := cloudendpointresource.CloudEndpointCreateParameters{
		Properties: &cloudendpointresource.CloudEndpointCreateParametersProperties{
			StorageAccountResourceId: utils.String(d.Get("storage_account_id").(string)),
			AzureFileShareName:       utils.String(d.Get("file_share_name").(string)),
		},
	}

	tenantId := meta.(*clients.Client).Account.TenantId
	if v, ok := d.GetOk("storage_account_tenant_id"); ok {
		tenantId = v.(string)
	}
	parameters.Properties.StorageAccountTenantId = &tenantId

	if err := client.CloudEndpointsCreateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceStorageSyncCloudEndpointRead(d, meta)
}

func resourceStorageSyncCloudEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncCloudEndpointsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cloudendpointresource.ParseCloudEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.CloudEndpointsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.CloudEndpointName)

	groupId := syncgroupresource.NewSyncGroupID(id.SubscriptionId, id.ResourceGroupName, id.StorageSyncServiceName, id.SyncGroupName)
	d.Set("storage_sync_group_id", groupId.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("file_share_name", props.AzureFileShareName)
			d.Set("storage_account_id", props.StorageAccountResourceId)
			d.Set("storage_account_tenant_id", props.StorageAccountTenantId)
		}
	}
	return nil
}

func resourceStorageSyncCloudEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.SyncCloudEndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cloudendpointresource.ParseCloudEndpointID(d.Id())
	if err != nil {
		return err
	}

	if err = client.CloudEndpointsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
