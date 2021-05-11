package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-01-01/storage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceStorageAccountNetworkRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageAccountNetworkRulesCreateUpdate,
		Read:   resourceStorageAccountNetworkRulesRead,
		Update: resourceStorageAccountNetworkRulesCreateUpdate,
		Delete: resourceStorageAccountNetworkRulesDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"storage_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageAccountName,
			},

			"bypass": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(storage.AzureServices),
						string(storage.Logging),
						string(storage.Metrics),
						string(storage.None),
					}, false),
				},
				Set: schema.HashString,
			},

			"ip_rules": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"virtual_network_subnet_ids": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
				Set: schema.HashString,
			},

			"default_action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.DefaultActionAllow),
					string(storage.DefaultActionDeny),
				}, false),
			},

			"private_link_access": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_resource_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: networkValidate.PrivateEndpointID,
						},

						"endpoint_tenant_id": {
							Type:         schema.TypeString,
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

func resourceStorageAccountNetworkRulesCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	tenantId := meta.(*clients.Client).Account.TenantId
	client := meta.(*clients.Client).Storage.AccountsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageAccountName := d.Get("storage_account_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(storageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(storageAccountName, storageAccountResourceName)

	storageAccount, err := client.GetProperties(ctx, resourceGroup, storageAccountName, "")
	if err != nil {
		if utils.ResponseWasNotFound(storageAccount.Response) {
			return fmt.Errorf("Storage Account %q (Resource Group %q) was not found", storageAccountName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Storage Account %q (Resource Group %q): %+v", storageAccountName, resourceGroup, err)
	}

	if d.IsNewResource() {
		if storageAccount.AccountProperties == nil {
			return fmt.Errorf("Error retrieving Storage Account %q (Resource Group %q): `properties` was nil", storageAccountName, resourceGroup)
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
	rules.Bypass = expandStorageAccountNetworkRuleBypass(d.Get("bypass").(*schema.Set).List())
	rules.IPRules = expandStorageAccountNetworkRuleIpRules(d.Get("ip_rules").(*schema.Set).List())
	rules.VirtualNetworkRules = expandStorageAccountNetworkRuleVirtualRules(d.Get("virtual_network_subnet_ids").(*schema.Set).List())
	rules.ResourceAccessRules = expandStorageAccountPrivateLinkAccess(d.Get("private_link_access").([]interface{}), tenantId)

	opts := storage.AccountUpdateParameters{
		AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
			NetworkRuleSet: rules,
		},
	}

	if _, err := client.Update(ctx, resourceGroup, storageAccountName, opts); err != nil {
		return fmt.Errorf("Error updating Azure Storage Account Network Rules %q (Resource Group %q): %+v", storageAccountName, resourceGroup, err)
	}

	d.SetId(*storageAccount.ID)

	return resourceStorageAccountNetworkRulesRead(d, meta)
}

func resourceStorageAccountNetworkRulesRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error reading Storage Account Network Rules %q (Resource Group %q): %+v", storageAccountName, resourceGroup, err)
	}

	d.Set("storage_account_name", storageAccountName)
	d.Set("resource_group_name", resourceGroup)

	if rules := storageAccount.NetworkRuleSet; rules != nil {
		if err := d.Set("ip_rules", schema.NewSet(schema.HashString, flattenStorageAccountIPRules(rules.IPRules))); err != nil {
			return fmt.Errorf("Error setting `ip_rules`: %+v", err)
		}
		if err := d.Set("virtual_network_subnet_ids", schema.NewSet(schema.HashString, flattenStorageAccountVirtualNetworks(rules.VirtualNetworkRules))); err != nil {
			return fmt.Errorf("Error setting `virtual_network_subnet_ids`: %+v", err)
		}
		if err := d.Set("bypass", schema.NewSet(schema.HashString, flattenStorageAccountBypass(rules.Bypass))); err != nil {
			return fmt.Errorf("Error setting `bypass`: %+v", err)
		}
		d.Set("default_action", string(rules.DefaultAction))
		if err := d.Set("private_link_access", flattenStorageAccountPrivateLinkAccess(rules.ResourceAccessRules)); err != nil {
			return fmt.Errorf("setting `private_link_access`: %+v", err)
		}
	}

	return nil
}

func resourceStorageAccountNetworkRulesDelete(d *schema.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("Error loading Storage Account %q (Resource Group %q): %+v", storageAccountName, resourceGroup, err)
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
		return fmt.Errorf("Error deleting Azure Storage Account Network Rule %q (Resource Group %q): %+v", storageAccountName, resourceGroup, err)
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
