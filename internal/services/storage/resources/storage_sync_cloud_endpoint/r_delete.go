// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_sync_cloud_endpoint

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/cloudendpointresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	storageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageSyncCloudEndpointDelete(d *pluginsdk.ResourceData, meta any) error {
	client := storageClient.From(meta).SyncCloudEndpointsClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
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
