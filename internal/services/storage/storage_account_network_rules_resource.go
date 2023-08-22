// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceStorageAccountNetworkRules() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageAccountNetworkRulesCreateUpdate,
		Read:   resourceStorageAccountNetworkRulesRead,
		Update: resourceStorageAccountNetworkRulesCreateUpdate,
		Delete: resourceStorageAccountNetworkRulesDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseStorageAccountID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: resourceStorageAccountNetworkRulesSchema(),
	}
}

func resourceStorageAccountNetworkRulesSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// lintignore: S013
		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
		},

		"bypass": {
			Type:       pluginsdk.TypeSet,
			Optional:   true,
			Computed:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.BypassAzureServices),
					string(storage.BypassLogging),
					string(storage.BypassMetrics),
					string(storage.BypassNone),
				}, false),
			},
			Set: pluginsdk.HashString,
		},

		"ip_rules": {
			Type:       pluginsdk.TypeSet,
			Optional:   true,
			Computed:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.StorageAccountIpRule,
			},
			Set: pluginsdk.HashString,
		},

		"virtual_network_subnet_ids": {
			Type:       pluginsdk.TypeSet,
			Optional:   true,
			Computed:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
			},
			Set: pluginsdk.HashString,
		},

		"default_action": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(storage.DefaultActionAllow),
				string(storage.DefaultActionDeny),
			}, false),
		},

		"private_link_access": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"endpoint_resource_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"endpoint_tenant_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},
	}
}

func resourceStorageAccountNetworkRulesCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	tenantId := meta.(*clients.Client).Account.TenantId
	client := meta.(*clients.Client).Storage.AccountsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	var storageAccountName string
	var resourceGroup string
	raw, ok := d.GetOk("storage_account_id")
	if ok {
		parsedStorageAccountId, err := commonids.ParseStorageAccountID(raw.(string))
		if err != nil {
			return err
		}

		storageAccountName = parsedStorageAccountId.StorageAccountName
		resourceGroup = parsedStorageAccountId.ResourceGroupName
	}

	locks.ByName(storageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(storageAccountName, storageAccountResourceName)

	id := commonids.NewStorageAccountID(subscriptionId, resourceGroup, storageAccountName)
	storageAccount, err := client.GetProperties(ctx, id.ResourceGroupName, id.StorageAccountName, "")
	if err != nil {
		if utils.ResponseWasNotFound(storageAccount.Response) {
			return fmt.Errorf("Storage Account %q (Resource Group %q) was not found", storageAccountName, resourceGroup)
		}

		return fmt.Errorf("retrieving Storage Account %q (Resource Group %q): %+v", storageAccountName, resourceGroup, err)
	}

	if d.IsNewResource() {
		if storageAccount.AccountProperties == nil {
			return fmt.Errorf("retrieving Storage Account %q (Resource Group %q): `properties` was nil", storageAccountName, resourceGroup)
		}

		if checkForNonDefaultStorageAccountNetworkRule(storageAccount.AccountProperties.NetworkRuleSet) {
			return tf.ImportAsExistsError("azurerm_storage_account_network_rule", id.ID())
		}
	}

	rules := storageAccount.NetworkRuleSet
	if rules == nil {
		rules = &storage.NetworkRuleSet{}
	}

	rules.DefaultAction = storage.DefaultAction(d.Get("default_action").(string))
	rules.Bypass = expandStorageAccountNetworkRuleBypass(d.Get("bypass").(*pluginsdk.Set).List())
	rules.IPRules = expandStorageAccountNetworkRuleIpRules(d.Get("ip_rules").(*pluginsdk.Set).List())
	rules.VirtualNetworkRules = expandStorageAccountNetworkRuleVirtualRules(d.Get("virtual_network_subnet_ids").(*pluginsdk.Set).List())
	rules.ResourceAccessRules = expandStorageAccountPrivateLinkAccess(d.Get("private_link_access").([]interface{}), tenantId)

	opts := storage.AccountUpdateParameters{
		AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
			NetworkRuleSet: rules,
		},
	}

	if _, err := client.Update(ctx, resourceGroup, storageAccountName, opts); err != nil {
		return fmt.Errorf("updating Azure Storage Account Network Rules %q (Resource Group %q): %+v", storageAccountName, resourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceStorageAccountNetworkRulesRead(d, meta)
}

func resourceStorageAccountNetworkRulesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	storageAccount, err := client.GetProperties(ctx, id.ResourceGroupName, id.StorageAccountName, "")
	if err != nil {
		if utils.ResponseWasNotFound(storageAccount.Response) {
			log.Printf("[INFO] Storage Account Network Rules %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Storage Account Network Rules %s : %+v", *id, err)
	}

	d.Set("storage_account_id", d.Id())

	if rules := storageAccount.NetworkRuleSet; rules != nil {
		if err := d.Set("ip_rules", pluginsdk.NewSet(pluginsdk.HashString, flattenStorageAccountIPRules(rules.IPRules))); err != nil {
			return fmt.Errorf("setting `ip_rules`: %+v", err)
		}
		if err := d.Set("virtual_network_subnet_ids", pluginsdk.NewSet(pluginsdk.HashString, flattenStorageAccountVirtualNetworks(rules.VirtualNetworkRules))); err != nil {
			return fmt.Errorf("setting `virtual_network_subnet_ids`: %+v", err)
		}
		if err := d.Set("bypass", pluginsdk.NewSet(pluginsdk.HashString, flattenStorageAccountBypass(rules.Bypass))); err != nil {
			return fmt.Errorf("setting `bypass`: %+v", err)
		}
		d.Set("default_action", string(rules.DefaultAction))
		if err := d.Set("private_link_access", flattenStorageAccountPrivateLinkAccess(rules.ResourceAccessRules)); err != nil {
			return fmt.Errorf("setting `private_link_access`: %+v", err)
		}
	}

	return nil
}

func resourceStorageAccountNetworkRulesDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsedStorageAccountNetworkRuleId, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(parsedStorageAccountNetworkRuleId.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(parsedStorageAccountNetworkRuleId.StorageAccountName, storageAccountResourceName)

	storageAccount, err := client.GetProperties(ctx, parsedStorageAccountNetworkRuleId.ResourceGroupName, parsedStorageAccountNetworkRuleId.StorageAccountName, "")
	if err != nil {
		if utils.ResponseWasNotFound(storageAccount.Response) {
			return fmt.Errorf("%s was not found", *parsedStorageAccountNetworkRuleId)
		}

		return fmt.Errorf("loading %s: %+v", *parsedStorageAccountNetworkRuleId, err)
	}

	if storageAccount.NetworkRuleSet == nil {
		return nil
	}

	// We can't delete a network rule set so we'll just update it back to the default instead
	virtualNetworkRules := make([]storage.VirtualNetworkRule, 0)
	ipRules := make([]storage.IPRule, 0)
	opts := storage.AccountUpdateParameters{
		AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
			NetworkRuleSet: &storage.NetworkRuleSet{
				Bypass:              storage.BypassAzureServices,
				VirtualNetworkRules: &virtualNetworkRules,
				IPRules:             &ipRules,
				DefaultAction:       storage.DefaultActionAllow,
			},
		},
	}

	if _, err := client.Update(ctx, parsedStorageAccountNetworkRuleId.ResourceGroupName, parsedStorageAccountNetworkRuleId.StorageAccountName, opts); err != nil {
		return fmt.Errorf("deleting Azure %s: %+v", *parsedStorageAccountNetworkRuleId, err)
	}

	return nil
}

// To make sure that someone isn't overriding their existing network rules, we'll check for a non default network rule
func checkForNonDefaultStorageAccountNetworkRule(rule *storage.NetworkRuleSet) bool {
	if rule == nil {
		return false
	}

	// THe default action "Allow" is set in the creation of the storage account resource as default value.
	if (rule.IPRules != nil && len(*rule.IPRules) != 0) ||
		(rule.VirtualNetworkRules != nil && len(*rule.VirtualNetworkRules) != 0) ||
		rule.DefaultAction != storage.DefaultActionAllow {
		return true
	}

	return false
}

func expandStorageAccountNetworkRuleBypass(bypass []interface{}) storage.Bypass {
	var bypassValues []string
	for _, bypassConfig := range bypass {
		bypassValues = append(bypassValues, bypassConfig.(string))
	}

	return storage.Bypass(strings.Join(bypassValues, ", "))
}

func expandStorageAccountNetworkRuleIpRules(ipRulesInfo []interface{}) *[]storage.IPRule {
	ipRules := make([]storage.IPRule, len(ipRulesInfo))

	for i, ipRuleConfig := range ipRulesInfo {
		attrs := ipRuleConfig.(string)
		ipRule := storage.IPRule{
			IPAddressOrRange: utils.String(attrs),
			Action:           storage.ActionAllow,
		}
		ipRules[i] = ipRule
	}

	return &ipRules
}

func expandStorageAccountNetworkRuleVirtualRules(virtualNetworkInfo []interface{}) *[]storage.VirtualNetworkRule {
	virtualNetworks := make([]storage.VirtualNetworkRule, len(virtualNetworkInfo))

	for i, virtualNetworkConfig := range virtualNetworkInfo {
		attrs := virtualNetworkConfig.(string)
		virtualNetwork := storage.VirtualNetworkRule{
			VirtualNetworkResourceID: utils.String(attrs),
			Action:                   storage.ActionAllow,
		}
		virtualNetworks[i] = virtualNetwork
	}

	return &virtualNetworks
}
