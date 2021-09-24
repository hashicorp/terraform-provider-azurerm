package storage

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-01-01/storage"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
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
			_, err := parse.StorageAccountID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"storage_account_id": {
				Type: pluginsdk.TypeString,
				// TODO: Make required in 3.0
				Optional: true,
				// TODO: Remove in 3.0
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageAccountID,
				// TODO: Remove in 3.0
				ConflictsWith: []string{
					"resource_group_name",
					"storage_account_name",
				},
			},

			// TODO: remove in 3.0
			"resource_group_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceGroupName,
				Deprecated:   "Deprecated in favour of `storage_account_id`",
				RequiredWith: []string{
					"storage_account_name",
				},
				ConflictsWith: []string{
					"storage_account_id",
				},
			},

			// TODO: remove in 3.0
			"storage_account_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageAccountName,
				Deprecated:   "Deprecated in favour of `storage_account_id`",
				RequiredWith: []string{
					"resource_group_name",
				},
				ConflictsWith: []string{
					"storage_account_id",
				},
			},

			"bypass": {
				Type:       pluginsdk.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(storage.AzureServices),
						string(storage.Logging),
						string(storage.Metrics),
						string(storage.None),
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
					Type: pluginsdk.TypeString,
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
		},
	}
}

func resourceStorageAccountNetworkRulesCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	tenantId := meta.(*clients.Client).Account.TenantId
	client := meta.(*clients.Client).Storage.AccountsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageAccountName := d.Get("storage_account_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	raw, ok := d.GetOk("storage_account_id")
	if ok {
		parsedStorageAccountId, err := parse.StorageAccountID(raw.(string))
		if err != nil {
			return err
		}

		storageAccountName = parsedStorageAccountId.Name
		resourceGroup = parsedStorageAccountId.ResourceGroup
	}

	locks.ByName(storageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(storageAccountName, storageAccountResourceName)

	storageAccount, err := client.GetProperties(ctx, resourceGroup, storageAccountName, "")
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
			return tf.ImportAsExistsError("azurerm_storage_account_network_rule", *storageAccount.ID)
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

	d.SetId(*storageAccount.ID)

	return resourceStorageAccountNetworkRulesRead(d, meta)
}

func resourceStorageAccountNetworkRulesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	storageAccountName := id.Path["storageAccounts"]

	storageAccount, err := client.GetProperties(ctx, resourceGroup, storageAccountName, "")
	if err != nil {
		if utils.ResponseWasNotFound(storageAccount.Response) {
			log.Printf("[INFO] Storage Account Network Rules %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Storage Account Network Rules %q (Resource Group %q): %+v", storageAccountName, resourceGroup, err)
	}

	d.Set("storage_account_id", d.Id())
	d.Set("storage_account_name", storageAccountName)
	d.Set("resource_group_name", resourceGroup)

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

	parsedStorageAccountNetworkRuleId, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := parsedStorageAccountNetworkRuleId.ResourceGroup
	storageAccountName := parsedStorageAccountNetworkRuleId.Path["storageAccounts"]

	locks.ByName(storageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(storageAccountName, storageAccountResourceName)

	storageAccount, err := client.GetProperties(ctx, resourceGroup, storageAccountName, "")
	if err != nil {
		if utils.ResponseWasNotFound(storageAccount.Response) {
			return fmt.Errorf("Storage Account %q (Resource Group %q) was not found", storageAccountName, resourceGroup)
		}

		return fmt.Errorf("loading Storage Account %q (Resource Group %q): %+v", storageAccountName, resourceGroup, err)
	}

	if storageAccount.NetworkRuleSet == nil {
		return nil
	}

	// We can't delete a network rule set so we'll just update it back to the default instead
	opts := storage.AccountUpdateParameters{
		AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
			NetworkRuleSet: &storage.NetworkRuleSet{
				Bypass:        storage.AzureServices,
				DefaultAction: storage.DefaultActionAllow,
			},
		},
	}

	if _, err := client.Update(ctx, resourceGroup, storageAccountName, opts); err != nil {
		return fmt.Errorf("deleting Azure Storage Account Network Rule %q (Resource Group %q): %+v", storageAccountName, resourceGroup, err)
	}

	return nil
}

// To make sure that someone isn't overriding their existing network rules, we'll check for a non default network rule
func checkForNonDefaultStorageAccountNetworkRule(rule *storage.NetworkRuleSet) bool {
	if rule == nil {
		return false
	}

	if (rule.IPRules != nil && len(*rule.IPRules) != 0) ||
		(rule.VirtualNetworkRules != nil && len(*rule.VirtualNetworkRules) != 0) ||
		rule.Bypass != "AzureServices" || rule.DefaultAction != "Allow" {
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
			Action:           storage.Allow,
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
			Action:                   storage.Allow,
		}
		virtualNetworks[i] = virtualNetwork
	}

	return &virtualNetworks
}
