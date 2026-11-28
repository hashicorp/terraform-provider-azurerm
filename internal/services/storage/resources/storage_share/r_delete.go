// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_share

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/fileshares"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageShareDelete(d *pluginsdk.ResourceData, meta any) error {
	fileSharesClient := meta.(*clients.Client).Storage.ResourceManager.FileShares
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := fileshares.ParseShareID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := fileSharesClient.Delete(ctx, *id, fileshares.DefaultDeleteOperationOptions()); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %v", id, err)
		}
	}

	return nil
}
