// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStorageAccountNetworkRules() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageAccountNetworkRulesCreate,
		Read:   resourceStorageAccountNetworkRulesRead,
		Update: resourceStorageAccountNetworkRulesUpdate,
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

		Schema: map[string]*pluginsdk.Schema{
			// lintignore: S013
			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"bypass": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: true, // Defaults to storageaccounts.BypassAzureServices in the API, but schema does not support defaults for lists/sets.
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(storageaccounts.BypassAzureServices),
						string(storageaccounts.BypassLogging),
						string(storageaccounts.BypassMetrics),
						string(storageaccounts.BypassNone),
					}, false),
				},
				Set: pluginsdk.HashString,
			},

			"ip_rules": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.StorageAccountIpRule,
				},
				Set: pluginsdk.HashString,
			},

			"virtual_network_subnet_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
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
					string(storageaccounts.DefaultActionAllow),
					string(storageaccounts.DefaultActionDeny),
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
		},
	}
}

func resourceStorageAccountNetworkRulesCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	tenantId := meta.(*clients.Client).Account.TenantId
	client := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
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
		return tf.ImportAsExistsError("azurerm_storage_account_network_rule", id.ID())
	}

	acls := resp.Model.Properties.NetworkAcls
	if acls == nil {
		acls = &storageaccounts.NetworkRuleSet{}
	}

	acls.DefaultAction = storageaccounts.DefaultAction(d.Get("default_action").(string))
	acls.Bypass = expandAccountNetworkRuleBypass(d.Get("bypass").(*pluginsdk.Set).List())
	acls.IPRules = expandAccountNetworkRuleIPRules(d.Get("ip_rules").(*pluginsdk.Set).List())
	acls.VirtualNetworkRules = expandAccountNetworkRuleVirtualNetworkRules(d.Get("virtual_network_subnet_ids").(*pluginsdk.Set).List())
	acls.ResourceAccessRules = expandAccountNetworkRulePrivateLinkAccess(d.Get("private_link_access").([]interface{}), tenantId)

	payload := storageaccounts.StorageAccountUpdateParameters{
		Properties: &storageaccounts.StorageAccountPropertiesUpdateParameters{
			NetworkAcls: acls,
		},
	}
	if _, err = client.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("creating Network Rules for %s: %+v", *id, err)
	}

	d.SetId(id.ID())

	return resourceStorageAccountNetworkRulesRead(d, meta)
}

func resourceStorageAccountNetworkRulesUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
		acls.Bypass = expandAccountNetworkRuleBypass(d.Get("bypass").(*pluginsdk.Set).List())
	}
	if d.HasChange("ip_rules") {
		acls.IPRules = expandAccountNetworkRuleIPRules(d.Get("ip_rules").(*pluginsdk.Set).List())
	}
	if d.HasChange("virtual_network_subnet_ids") {
		acls.VirtualNetworkRules = expandAccountNetworkRuleVirtualNetworkRules(d.Get("virtual_network_subnet_ids").(*pluginsdk.Set).List())
	}

	if d.HasChange("private_link_access") {
		acls.ResourceAccessRules = expandAccountNetworkRulePrivateLinkAccess(d.Get("private_link_access").([]interface{}), tenantId)
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

func resourceStorageAccountNetworkRulesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("storage_account_id", d.Id())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if rules := props.NetworkAcls; rules != nil {
				if err := d.Set("ip_rules", pluginsdk.NewSet(pluginsdk.HashString, flattenAccountNetworkRuleIPRules(rules.IPRules))); err != nil {
					return fmt.Errorf("setting `ip_rules`: %+v", err)
				}
				if err := d.Set("virtual_network_subnet_ids", pluginsdk.NewSet(pluginsdk.HashString, flattenAccountNetworkRuleVirtualNetworkRules(rules.VirtualNetworkRules))); err != nil {
					return fmt.Errorf("setting `virtual_network_subnet_ids`: %+v", err)
				}
				if err := d.Set("bypass", pluginsdk.NewSet(pluginsdk.HashString, flattenAccountNetworkRuleBypass(rules.Bypass))); err != nil {
					return fmt.Errorf("setting `bypass`: %+v", err)
				}
				d.Set("default_action", string(rules.DefaultAction))
				if err := d.Set("private_link_access", flattenAccountNetworkRulePrivateLinkAccess(rules.ResourceAccessRules)); err != nil {
					return fmt.Errorf("setting `private_link_access`: %+v", err)
				}
			}
		}
	}

	return nil
}

func resourceStorageAccountNetworkRulesDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
