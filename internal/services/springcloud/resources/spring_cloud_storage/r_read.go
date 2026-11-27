// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_storage

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSpringCloudStorageRead(d *pluginsdk.ResourceData, meta any) error {
	client := springcloudClient.From(meta).StoragesClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	id, err := appplatform.ParseStorageIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.StorageName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] %q does not exist - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}
	d.Set("name", id.StorageName)
	d.Set("spring_cloud_service_id", commonids.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroupName, id.SpringName).ID())
	if resp.Properties != nil {
		if props, ok := resp.Properties.AsStorageAccount(); ok {
			d.Set("storage_account_name", props.AccountName)
		}
	}
	return nil
}
