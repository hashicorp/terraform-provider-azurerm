// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_network_rules

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageAccountNetworkRulesUpdate(d *pluginsdk.ResourceData, meta any) error {
	tenantId := meta.(*clients.Client).Account.TenantId
	client := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	resp, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if resp.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", *id)
	}

	acls := resp.Model.Properties.NetworkAcls
	if acls == nil {
		acls = &storageaccounts.NetworkRuleSet{}
	}

	if d.HasChange("default_action") {
		acls.DefaultAction = storageaccounts.DefaultAction(d.Get("default_action").(string))
	}
	if d.HasChange("bypass") {
		acls.Bypass = helpers.ExpandAccountNetworkRuleBypass(d.Get("bypass").(*pluginsdk.Set).List())
	}
	if d.HasChange("ip_rules") {
		acls.IPRules = helpers.ExpandAccountNetworkRuleIPRules(d.Get("ip_rules").(*pluginsdk.Set).List())
	}
	if d.HasChange("virtual_network_subnet_ids") {
		acls.VirtualNetworkRules = helpers.ExpandAccountNetworkRuleVirtualNetworkRules(d.Get("virtual_network_subnet_ids").(*pluginsdk.Set).List())
	}

	if d.HasChange("private_link_access") {
		acls.ResourceAccessRules = helpers.ExpandAccountNetworkRulePrivateLinkAccess(d.Get("private_link_access").([]any), tenantId)
	}

	payload := storageaccounts.StorageAccountUpdateParameters{
		Properties: &storageaccounts.StorageAccountPropertiesUpdateParameters{
			NetworkAcls: acls,
		},
	}

	if _, err := client.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating Network Rules for %s: %+v", *id, err)
	}

	return resourceStorageAccountNetworkRulesRead(d, meta)
}
