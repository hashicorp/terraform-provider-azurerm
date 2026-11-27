// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_network_rules

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceStorageAccountNetworkRulesFlatten(d *pluginsdk.ResourceData, id *commonids.StorageAccountId, model *storageaccounts.StorageAccount) error {
	d.Set("storage_account_id", d.Id())

	if model != nil {
		if props := model.Properties; props != nil {
			if rules := props.NetworkAcls; rules != nil {
				if err := d.Set("ip_rules", pluginsdk.NewSet(pluginsdk.HashString, helpers.FlattenAccountNetworkRuleIPRules(rules.IPRules))); err != nil {
					return fmt.Errorf("setting `ip_rules`: %+v", err)
				}
				if err := d.Set("virtual_network_subnet_ids", pluginsdk.NewSet(pluginsdk.HashString, helpers.FlattenAccountNetworkRuleVirtualNetworkRules(rules.VirtualNetworkRules))); err != nil {
					return fmt.Errorf("setting `virtual_network_subnet_ids`: %+v", err)
				}
				if err := d.Set("bypass", pluginsdk.NewSet(pluginsdk.HashString, helpers.FlattenAccountNetworkRuleBypass(rules.Bypass))); err != nil {
					return fmt.Errorf("setting `bypass`: %+v", err)
				}
				d.Set("default_action", string(rules.DefaultAction))
				if err := d.Set("private_link_access", helpers.FlattenAccountNetworkRulePrivateLinkAccess(rules.ResourceAccessRules)); err != nil {
					return fmt.Errorf("setting `private_link_access`: %+v", err)
				}
			}
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id, pluginsdk.ResourceTypeForIdentityVirtual)
}
