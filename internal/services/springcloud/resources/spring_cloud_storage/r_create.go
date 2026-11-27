// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_storage

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudStorageCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := account.From(meta).SubscriptionId
	client := springcloudClient.From(meta).StoragesClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	springCloudId, err := commonids.ParseSpringCloudServiceIDInsensitively(d.Get("spring_cloud_service_id").(string))
	if err != nil {
		return err
	}

	id := appplatform_rm.NewStorageID(subscriptionId, springCloudId.ResourceGroupName, springCloudId.ServiceName, d.Get("name").(string))

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.StorageName)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for existing %q: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_spring_cloud_storage", id.ID())
			}
		}
	}

	storageResource := appplatform.StorageResource{
		Properties: &appplatform.StorageAccount{
			AccountName: pointer.To(d.Get("storage_account_name").(string)),
			AccountKey:  pointer.To(d.Get("storage_account_key").(string)),
			StorageType: appplatform.StorageTypeStorageAccount,
		},
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.SpringName, id.StorageName, storageResource)
	if err != nil {
		return fmt.Errorf("creating/updating %q: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
	}

	return resourceSpringCloudStorageRead(d, meta)
}
