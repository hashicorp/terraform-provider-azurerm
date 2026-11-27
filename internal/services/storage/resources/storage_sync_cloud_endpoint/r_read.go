// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_sync_cloud_endpoint

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/cloudendpointresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/syncgroupresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	storageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageSyncCloudEndpointRead(d *pluginsdk.ResourceData, meta any) error {
	client := storageClient.From(meta).SyncCloudEndpointsClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
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
	return pluginsdk.SetResourceIdentityData(d, id)
}
