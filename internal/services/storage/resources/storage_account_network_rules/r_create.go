// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_network_rules

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	storageClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageAccountNetworkRulesCreate(d *pluginsdk.ResourceData, meta any) error {
	tenantId := account.From(meta).TenantId
	client := storageClient.From(meta).ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
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
	usesNonDefaultStorageAccountRules := false
	if acls := resp.Model.Properties.NetworkAcls; acls != nil {
		// The default action "Allow" is set in the creation of the storage account resource as default value.
		hasIPRules := acls.IPRules != nil && len(*acls.IPRules) > 0
		defaultActionConfigured := acls.DefaultAction != storageaccounts.DefaultActionAllow
		hasVirtualNetworkRules := acls.VirtualNetworkRules != nil && len(*acls.VirtualNetworkRules) > 0
		if hasIPRules || defaultActionConfigured || hasVirtualNetworkRules {
			usesNonDefaultStorageAccountRules = true
		}
	}
	if usesNonDefaultStorageAccountRules {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			return tf.ImportAsExistsError("azurerm_storage_account_network_rule", id.ID())
		}
	}

	acls := resp.Model.Properties.NetworkAcls
	if acls == nil {
		acls = &storageaccounts.NetworkRuleSet{}
	}

	acls.DefaultAction = storageaccounts.DefaultAction(d.Get("default_action").(string))
	acls.Bypass = helpers.ExpandAccountNetworkRuleBypass(d.Get("bypass").(*pluginsdk.Set).List())
	acls.IPRules = helpers.ExpandAccountNetworkRuleIPRules(d.Get("ip_rules").(*pluginsdk.Set).List())
	acls.VirtualNetworkRules = helpers.ExpandAccountNetworkRuleVirtualNetworkRules(d.Get("virtual_network_subnet_ids").(*pluginsdk.Set).List())
	acls.ResourceAccessRules = helpers.ExpandAccountNetworkRulePrivateLinkAccess(d.Get("private_link_access").([]any), tenantId)

	payload := storageaccounts.StorageAccountUpdateParameters{
		Properties: &storageaccounts.StorageAccountPropertiesUpdateParameters{
			NetworkAcls: acls,
		},
	}
	if _, err = client.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("creating Network Rules for %s: %+v", *id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, id, pluginsdk.ResourceTypeForIdentityVirtual); err != nil {
		return err
	}

	return resourceStorageAccountNetworkRulesRead(d, meta)
}
