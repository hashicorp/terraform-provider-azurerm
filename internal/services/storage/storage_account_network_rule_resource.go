package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-04-01/storage"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceStorageAccountNetworkRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageAccountNetworkRuleCreate,
		Read:   resourceStorageAccountNetworkRuleRead,
		Delete: resourceStorageAccountNetworkRuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.StorageAccountNetworkRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(10 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(10 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageAccountID,
			},

			"ip_rule": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"virtual_network_rule", "resource_access_rule"},
				ValidateFunc:  validate.StorageAccountIpRule,
			},

			"virtual_network_rule": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"ip_rule", "resource_access_rule"},
				ValidateFunc:  networkValidate.SubnetID,
			},

			"resource_access_rule": {
				Type:          pluginsdk.TypeSet,
				Optional:      true,
				ForceNew:      true,
				MaxItems:      1,
				ConflictsWith: []string{"ip_rule", "virtual_network_rule"},
				Elem:          resourceAccessRuleResource(),
			},
		},
	}
}

func resourceAccessRuleResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceStorageAccountNetworkRuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	tenantID := meta.(*clients.Client).Account.TenantId
	client := meta.(*clients.Client).Storage.AccountsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageAccountIDRaw := d.Get("storage_account_id").(string)
	storageAccountID, err := parse.StorageAccountID(storageAccountIDRaw)
	if err != nil {
		return err
	}

	locks.ByName(storageAccountID.Name, storageAccountResourceName)
	defer locks.UnlockByName(storageAccountID.Name, storageAccountResourceName)

	existing, err := client.GetProperties(ctx, storageAccountID.ResourceGroup, storageAccountID.Name, "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for the presence of %s: %s", storageAccountID, err)
		}
	}

	rules := existing.NetworkRuleSet
	if rules == nil {
		rules = &storage.NetworkRuleSet{}
	}

	id := parse.StorageAccountNetworkRuleId{
		StorageAccountId: storageAccountID,
	}

	if ipRange, ok := d.GetOk("ip_rule"); ok {
		id.IPRule = &storage.IPRule{
			IPAddressOrRange: utils.String(ipRange.(string)),
			Action:           storage.ActionAllow,
		}
		*rules.IPRules = append(*rules.IPRules, *id.IPRule)
	}

	if subnetID, ok := d.GetOk("virtual_network_rule"); ok {
		id.VirtualNetworkRule = &storage.VirtualNetworkRule{
			VirtualNetworkResourceID: utils.String(subnetID.(string)),
			Action:                   storage.ActionAllow,
		}
		*rules.VirtualNetworkRules = append(*rules.VirtualNetworkRules, *id.VirtualNetworkRule)
	}

	if resourceAccessRule, ok := d.GetOk("resource_access_rule"); ok {
		id.ResourceAccessRule = expandStorageAccountNetworkRuleResourceAccessRule(resourceAccessRule.(*pluginsdk.Set).List(), tenantID)
		if rules.ResourceAccessRules == nil {
			r := make([]storage.ResourceAccessRule, 0)
			rules.ResourceAccessRules = &r
		}
		*rules.ResourceAccessRules = append(*rules.ResourceAccessRules, *id.ResourceAccessRule)
	}

	opts := storage.AccountUpdateParameters{
		AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
			NetworkRuleSet: rules,
		},
	}

	if _, err := client.Update(ctx, storageAccountID.ResourceGroup, storageAccountID.Name, opts); err != nil {
		return fmt.Errorf("creating %s: %+v", id.String(), err)
	}

	d.SetId(id.ID())

	return resourceStorageAccountNetworkRuleRead(d, meta)
}

func resourceStorageAccountNetworkRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	tenantID := meta.(*clients.Client).Account.TenantId
	client := meta.(*clients.Client).Storage.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageAccountNetworkRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetProperties(ctx, id.StorageAccountId.ResourceGroup, id.StorageAccountId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading the state of Storage Account %q: %+v", id.StorageAccountId.Name, err)
	}

	d.Set("storage_account_id", id.StorageAccountId.ID())
	if rules := resp.NetworkRuleSet; rules != nil {
		if ipRule, _, exists := FindStorageAccountNetworkIPRule(rules.IPRules, id.IPRule); exists {
			if err := d.Set("ip_rule", ipRule.IPAddressOrRange); err != nil {
				return fmt.Errorf("setting `ip_rule`: %+v", err)
			}
		}

		if vnetRule, _, exists := FindStorageAccountVirtualNetworkRule(rules.VirtualNetworkRules, id.VirtualNetworkRule); exists {
			if err := d.Set("virtual_network_rule", vnetRule.VirtualNetworkResourceID); err != nil {
				return fmt.Errorf("setting `virtual_network_rule`: %+v", err)
			}
		}

		if resourceAccessRule, _, exists := FindStorageAccountNetworkResourceAccessRule(rules.ResourceAccessRules, id.ResourceAccessRule); exists {
			flattenedResourceAccessRule := flattenStorageAccountNetworkRuleResourceAccessRule(resourceAccessRule, tenantID)
			if err := d.Set("resource_access_rule", pluginsdk.NewSet(pluginsdk.HashResource(resourceAccessRuleResource()), flattenedResourceAccessRule)); err != nil {
				return fmt.Errorf("setting `resource_access_rule`: %+v", err)
			}
		}
	}

	return nil
}

func resourceStorageAccountNetworkRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StorageAccountNetworkRuleID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountId.Name, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountId.Name, storageAccountResourceName)

	existing, err := client.GetProperties(ctx, id.StorageAccountId.ResourceGroup, id.StorageAccountId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return nil
		}
		return fmt.Errorf("retrieving %+v: %+v", *id, err)
	}

	if rules := existing.NetworkRuleSet; rules != nil {
		if _, index, exists := FindStorageAccountNetworkIPRule(rules.IPRules, id.IPRule); exists {
			r := *rules.IPRules
			r = append(r[:index], r[index+1:]...)
			rules.IPRules = &r
		}

		if _, index, exists := FindStorageAccountVirtualNetworkRule(rules.VirtualNetworkRules, id.VirtualNetworkRule); exists {
			r := *rules.VirtualNetworkRules
			r = append(r[:index], r[index+1:]...)
			rules.VirtualNetworkRules = &r
		}

		if _, index, exists := FindStorageAccountNetworkResourceAccessRule(rules.ResourceAccessRules, id.ResourceAccessRule); exists {
			r := *rules.ResourceAccessRules
			r = append(r[:index], r[index+1:]...)
			rules.ResourceAccessRules = &r
		}

		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				NetworkRuleSet: rules,
			},
		}

		if _, err := client.Update(ctx, id.StorageAccountId.ResourceGroup, id.StorageAccountId.Name, opts); err != nil {
			return fmt.Errorf("deleting %s: %+v", id.String(), err)
		}
	}

	return nil
}

func FindStorageAccountNetworkIPRule(items *[]storage.IPRule, item *storage.IPRule) (*storage.IPRule, int, bool) {
	if items != nil && item != nil {
		for idx, rule := range *items {
			if strings.EqualFold(*rule.IPAddressOrRange, *item.IPAddressOrRange) {
				return &rule, idx, true
			}
		}
	}
	return nil, -1, false
}

func FindStorageAccountVirtualNetworkRule(items *[]storage.VirtualNetworkRule, item *storage.VirtualNetworkRule) (*storage.VirtualNetworkRule, int, bool) {
	if items != nil && item != nil {
		for idx, rule := range *items {
			if strings.EqualFold(*rule.VirtualNetworkResourceID, *item.VirtualNetworkResourceID) {
				return &rule, idx, true
			}
		}
	}
	return nil, -1, false
}

func FindStorageAccountNetworkResourceAccessRule(items *[]storage.ResourceAccessRule, item *storage.ResourceAccessRule) (*storage.ResourceAccessRule, int, bool) {
	if items != nil && item != nil {
		for idx, rule := range *items {
			if strings.EqualFold(*rule.ResourceID, *item.ResourceID) && strings.EqualFold(*rule.TenantID, *item.TenantID) {
				return &rule, idx, true
			}
		}
	}
	return nil, -1, false
}

func flattenStorageAccountNetworkRuleResourceAccessRule(input *storage.ResourceAccessRule, tenantId string) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})
	if input.TenantID != nil {
		tenantId = *input.TenantID
	}
	result["resource_id"] = *input.ResourceID
	result["tenant_id"] = tenantId
	return []interface{}{result}
}

func expandStorageAccountNetworkRuleResourceAccessRule(input []interface{}, tenantId string) *storage.ResourceAccessRule {
	if len(input) == 0 {
		return nil
	}

	item := input[0].(map[string]interface{})
	if v := item["tenant_id"].(string); v != "" {
		tenantId = v
	}
	return &storage.ResourceAccessRule{
		ResourceID: utils.String(item["resource_id"].(string)),
		TenantID:   utils.String(tenantId),
	}
}
