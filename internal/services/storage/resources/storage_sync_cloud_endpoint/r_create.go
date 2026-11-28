// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_sync_cloud_endpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/cloudendpointresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/syncgroupresource"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageSyncCloudEndpointCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Storage.SyncCloudEndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	groupId, err := syncgroupresource.ParseSyncGroupID(d.Get("storage_sync_group_id").(string))
	if err != nil {
		return err
	}

	id := cloudendpointresource.NewCloudEndpointID(groupId.SubscriptionId, groupId.ResourceGroupName, groupId.StorageSyncServiceName, groupId.SyncGroupName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.CloudEndpointsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_storage_sync_cloud_endpoint", id.ID())
		}
	}

	parameters := cloudendpointresource.CloudEndpointCreateParameters{
		Properties: &cloudendpointresource.CloudEndpointCreateParametersProperties{
			StorageAccountResourceId: pointer.To(d.Get("storage_account_id").(string)),
			AzureFileShareName:       pointer.To(d.Get("file_share_name").(string)),
		},
	}

	tenantId := meta.(*clients.Client).Account.TenantId
	if v, ok := d.GetOk("storage_account_tenant_id"); ok {
		tenantId = v.(string)
	}
	parameters.Properties.StorageAccountTenantId = &tenantId

	if err := client.CloudEndpointsCreateCallbackThenPoll(ctx, id, parameters, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceStorageSyncCloudEndpointRead(d, meta)
}
