// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_table

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageTableDelete(d *pluginsdk.ResourceData, meta any) error {
	storageClient := meta.(*clients.Client).Storage
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	rmId, err := parse.StorageTableResourceManagerID(d.Id())
	if err != nil {
		return err
	}

	var account *client.AccountDetails
	if meta.(*clients.Client).Storage.StorageUseAzureAD {
		account = &client.AccountDetails{
			StorageAccountId: commonids.NewStorageAccountID(rmId.SubscriptionId, rmId.ResourceGroup, rmId.StorageAccountName),
		}
	} else {
		account, err = storageClient.GetAccount(ctx, commonids.NewStorageAccountID(rmId.SubscriptionId, rmId.ResourceGroup, rmId.StorageAccountName))
		if err != nil {
			return fmt.Errorf("retrieving Storage Account %q for Table %q: %v", rmId.StorageAccountName, rmId.TableName, err)
		}
	}
	if account == nil {
		return fmt.Errorf("locating Storage Account %q", rmId.StorageAccountName)
	}

	client, err := storageClient.TablesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
	if err != nil {
		return fmt.Errorf("building Tables Client: %v", err)
	}

	if err = client.Delete(ctx, rmId.TableName); err != nil {
		if strings.Contains(err.Error(), "unexpected status 40") {
			return nil
		}
		return fmt.Errorf("deleting table %q: %v", rmId.TableName, err)
	}

	return nil
}
