// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_network_rules

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageAccountNetworkRulesDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	// We can't delete a network rule set so we'll just update it back to the default instead
	payload := storageaccounts.StorageAccountUpdateParameters{
		Properties: &storageaccounts.StorageAccountPropertiesUpdateParameters{
			NetworkAcls: &storageaccounts.NetworkRuleSet{
				Bypass:              pointer.To(storageaccounts.BypassAzureServices),
				DefaultAction:       storageaccounts.DefaultActionAllow,
				IPRules:             pointer.To(make([]storageaccounts.IPRule, 0)),
				ResourceAccessRules: pointer.To(make([]storageaccounts.ResourceAccessRule, 0)),
				VirtualNetworkRules: pointer.To(make([]storageaccounts.VirtualNetworkRule, 0)),
			},
		},
	}

	if _, err := client.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("removing Network Rules for %s: %+v", *id, err)
	}

	return nil
}
