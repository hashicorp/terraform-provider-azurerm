// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageAccountDelete(d *pluginsdk.ResourceData, meta any) error {
	storageUtils := meta.(*clients.Client).Storage
	client := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// the networking api's only allow a single change to be made to a network layout at once, so let's lock to handle that
	virtualNetworkNames := make([]string, 0)
	if model := existing.Model; model != nil && model.Properties != nil {
		if acls := model.Properties.NetworkAcls; acls != nil {
			if vnr := acls.VirtualNetworkRules; vnr != nil {
				for _, v := range *vnr {
					subnetId, err := commonids.ParseSubnetIDInsensitively(v.Id)
					if err != nil {
						return err
					}

					networkName := subnetId.VirtualNetworkName
					for _, virtualNetworkName := range virtualNetworkNames {
						if networkName == virtualNetworkName {
							continue
						}
					}
					virtualNetworkNames = append(virtualNetworkNames, networkName)
				}
			}
		}
	}

	locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// remove this from the cache
	storageUtils.RemoveAccountFromCache(*id)

	return nil
}
