package keyvault

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/keyvault/mgmt/2020-04-01-preview/keyvault"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	commonValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	networkParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceKeyVaultNetworkAclRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKeyVaultNetworkAclRuleCreate,
		Read:   resourceKeyVaultNetworkAclRuleRead,
		Update: nil,
		Delete: resourceKeyVaultNetworkAclDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.NetworkAclRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: func() map[string]*pluginsdk.Schema {
			rSchema := map[string]*pluginsdk.Schema{
				"key_vault_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validate.VaultID,
				},
				"source": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.Any(
						commonValidate.IPv4Address,
						commonValidate.CIDR,
						networkValidate.SubnetID,
					),
				},
			}

			return rSchema
		}(),
	}
}

func resourceKeyVaultNetworkAclRuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Print("[INFO] Preparing arguments for Key Vault Network ACL Rule creation.")

	key_vault_id := d.Get("key_vault_id").(string)
	source := d.Get("source").(string)

	id, err := parse.NewNetworkAclRuleId(key_vault_id, source)
	if err != nil {
		return err
	}

	// Locking to prevent parallel changes causing issues
	locks.ByName(id.VaultId.Name, keyVaultResourceName)
	defer locks.UnlockByName(id.VaultId.Name, keyVaultResourceName)

	vault, err := client.Get(ctx, id.VaultId.ResourceGroup, id.VaultId.Name)
	if err != nil {
		// If the key vault does not exist but this is not a new resource, the rule
		// which previously existed was deleted with the key vault, so reflect that in
		// state. If this is a new resource and key vault does not exist, it's likely
		// a bad ID was given.
		if utils.ResponseWasNotFound(vault.Response) && !d.IsNewResource() {
			log.Printf("[DEBUG] Parent %q was not found - removing from state!", id.VaultId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id.VaultId, err)
	}

	vaultProperties := vault.Properties
	if vaultProperties == nil {
		return fmt.Errorf("parsing Key Vault: `properties` was nil")
	}
	if vaultProperties.NetworkAcls == nil {
		return fmt.Errorf("parsing Key Vault: `properties.networkAcls` was nil")
	}

	networkRuleSet := keyvault.NetworkRuleSet{
		Bypass:        vault.Properties.NetworkAcls.Bypass,
		DefaultAction: vault.Properties.NetworkAcls.DefaultAction,
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"notfound"},
		Target:                    []string{"found"},
		Delay:                     30 * time.Second,
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 10,
		Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
	}

	sourceRule := id.Rule()

	switch id.Type {
	case parse.CIDR, parse.Ip:
		if FindKeyVaultNetworkAclIpRule(*vaultProperties.NetworkAcls.IPRules, sourceRule) != nil {
			return tf.ImportAsExistsError("azurerm_key_vault_network_acl_rule", id.ID())
		}

		stateConf.Refresh = networkAclIpRuleRefreshFunc(ctx, client, id.VaultId.ResourceGroup, id.VaultId.Name, sourceRule)

		ipRules := append(*vault.Properties.NetworkAcls.IPRules, keyvault.IPRule{Value: &source})
		networkRuleSet.IPRules = &ipRules
	case parse.VirtualNetwork:
		subnetId, err := networkParse.SubnetIDInsensitively(id.Source)
		if err != nil {
			return err
		}

		if FindKeyVaultNetworkAclVirtualNetworkRule(*vaultProperties.NetworkAcls.VirtualNetworkRules, sourceRule) != nil {
			return tf.ImportAsExistsError("azurerm_key_vault_network_acl_rule", id.ID())
		}

		stateConf.Refresh = networkAclVirtualNetworkRuleRefreshFunc(ctx, client, id.VaultId.ResourceGroup, id.VaultId.Name, sourceRule)

		virtualNetworkRules := append(*vault.Properties.NetworkAcls.VirtualNetworkRules, keyvault.VirtualNetworkRule{ID: &source})
		networkRuleSet.VirtualNetworkRules = &virtualNetworkRules

		// also lock on the Virtual Network ID's since modifications in the networking stack are exclusive
		virtualNetworkNames := []string{subnetId.VirtualNetworkName}
		locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
		defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	default:
		return fmt.Errorf("unexpected network ACL type %s", id.Type)
	}

	parameters := keyvault.VaultPatchParameters{
		Tags: vault.Tags,
		Properties: &keyvault.VaultPatchProperties{
			NetworkAcls: &networkRuleSet,
		},
	}
	if _, err = client.Update(ctx, id.VaultId.ResourceGroup, id.VaultId.Name, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Waiting for %q to become available", id)
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for scaling of %q: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourceKeyVaultNetworkAclRuleRead(d, meta)
}

func resourceKeyVaultNetworkAclRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NetworkAclRuleID(d.Id())
	if err != nil {
		return err
	}

	vault, err := client.Get(ctx, id.VaultId.ResourceGroup, id.VaultId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(vault.Response) {
			log.Printf("[DEBUG] Parent %q was not found - removing from state!", id.VaultId)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id.VaultId, err)
	}

	vaultProperties := vault.Properties
	if vaultProperties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}
	if vaultProperties.NetworkAcls == nil {
		return fmt.Errorf("retrieving %s: `properties.networkAcls` was nil", *id)
	}

	meta.(*clients.Client).KeyVault.AddToCache(id.VaultId, *vault.Properties.VaultURI)

	sourceRule := id.Rule()

	switch id.Type {
	case parse.CIDR, parse.Ip:
		if FindKeyVaultNetworkAclIpRule(*vaultProperties.NetworkAcls.IPRules, sourceRule) == nil {
			existingIpRules := make([]string, len(*vaultProperties.NetworkAcls.IPRules))
			for idx, ipRule := range *vaultProperties.NetworkAcls.IPRules {
				existingIpRules[idx] = *ipRule.Value
			}
			return fmt.Errorf("couldn't find ip rule %s in %v", id.Source, existingIpRules)
		}
	case parse.VirtualNetwork:
		if FindKeyVaultNetworkAclVirtualNetworkRule(*vaultProperties.NetworkAcls.VirtualNetworkRules, sourceRule) == nil {
			existingVirtualNetworkRules := make([]string, len(*vaultProperties.NetworkAcls.VirtualNetworkRules))
			for idx, virtualNetworkRule := range *vaultProperties.NetworkAcls.VirtualNetworkRules {
				existingVirtualNetworkRules[idx] = *virtualNetworkRule.ID
			}
			return fmt.Errorf("couldn't find virtual network subnet rule %s in %v", id.Source, existingVirtualNetworkRules)
		}
	default:
		return fmt.Errorf("unexpected network ACL type %s", id.Type)
	}

	d.Set("key_vault_id", id.VaultId.ID())
	d.Set("source", id.Source)

	return nil
}

func resourceKeyVaultNetworkAclDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Print("[INFO] Preparing arguments for Key Vault Network ACL Rule creation.")

	id, err := parse.NetworkAclRuleID(d.Id())
	if err != nil {
		return err
	}

	// Locking to prevent parallel changes causing issues
	locks.ByName(id.VaultId.Name, keyVaultResourceName)
	defer locks.UnlockByName(id.VaultId.Name, keyVaultResourceName)

	vault, err := client.Get(ctx, id.VaultId.ResourceGroup, id.VaultId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(vault.Response) && !d.IsNewResource() {
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", id.VaultId, err)
	}

	vaultProperties := vault.Properties
	if vaultProperties == nil {
		return fmt.Errorf("parsing Key Vault: `properties` was nil")
	}
	if vaultProperties.NetworkAcls == nil {
		return fmt.Errorf("parsing Key Vault: `properties.NetworkAcls` was nil")
	}

	networkRuleSet := keyvault.NetworkRuleSet{
		Bypass:        vaultProperties.NetworkAcls.Bypass,
		DefaultAction: vaultProperties.NetworkAcls.DefaultAction,
	}
	stateConf := pluginsdk.StateChangeConf{
		Pending:                   []string{"found"},
		Target:                    []string{"notfound"},
		Delay:                     30 * time.Second,
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 10,
		Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
	}

	sourceRule := id.Rule()

	switch id.Type {
	case parse.CIDR, parse.Ip:
		idx := FindKeyVaultNetworkAclIpRule(*vaultProperties.NetworkAcls.IPRules, sourceRule)
		if idx == nil {
			return nil
		}

		stateConf.Refresh = networkAclIpRuleRefreshFunc(ctx, client, id.VaultId.ResourceGroup, id.VaultId.Name, sourceRule)

		ipRules := append((*vault.Properties.NetworkAcls.IPRules)[:*idx], (*vault.Properties.NetworkAcls.IPRules)[*idx+1:]...)
		networkRuleSet.IPRules = &ipRules
	case parse.VirtualNetwork:
		subnetId, err := networkParse.SubnetIDInsensitively(id.Source)
		if err != nil {
			return err
		}

		idx := FindKeyVaultNetworkAclVirtualNetworkRule(*vaultProperties.NetworkAcls.VirtualNetworkRules, sourceRule)
		if idx == nil {
			return nil
		}

		stateConf.Refresh = networkAclVirtualNetworkRuleRefreshFunc(ctx, client, id.VaultId.ResourceGroup, id.VaultId.Name, sourceRule)

		virtualNetworkRules := append((*vault.Properties.NetworkAcls.VirtualNetworkRules)[:*idx], (*vault.Properties.NetworkAcls.VirtualNetworkRules)[*idx+1:]...)
		networkRuleSet.VirtualNetworkRules = &virtualNetworkRules

		// also lock on the Virtual Network ID's since modifications in the networking stack are exclusive
		virtualNetworkNames := []string{subnetId.VirtualNetworkName}
		locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
		defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	default:
		return fmt.Errorf("unexpected network ACL type %s", id.Type)
	}

	parameters := keyvault.VaultPatchParameters{
		Tags: vault.Tags,
		Properties: &keyvault.VaultPatchProperties{
			NetworkAcls: &networkRuleSet,
		},
	}
	if _, err = client.Update(ctx, id.VaultId.ResourceGroup, id.VaultId.Name, parameters); err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	log.Printf("[DEBUG] Waiting for %q to become available", id)
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for scaling of %q: %+v", id, err)
	}
	d.SetId(id.ID())

	return nil
}

func FindKeyVaultNetworkAclIpRule(rules []keyvault.IPRule, policyValue string) *int {
	for idx, ipRule := range rules {
		if ipRule.Value == nil {
			continue
		}
		if policyValue == *ipRule.Value {
			return &idx
		}
	}
	return nil
}

func FindKeyVaultNetworkAclVirtualNetworkRule(rules []keyvault.VirtualNetworkRule, policyId string) *int {
	for idx, virtualNetworkRule := range rules {
		if virtualNetworkRule.ID == nil {
			continue
		}
		if strings.EqualFold(policyId, *virtualNetworkRule.ID) {
			return &idx
		}
	}
	return nil
}

func networkAclIpRuleRefreshFunc(ctx context.Context, client *keyvault.VaultsClient, resourceGroup string, vaultName string, source string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking for completion of Network ACL IP Rule create/update")

		read, err := client.Get(ctx, resourceGroup, vaultName)
		if err != nil && utils.ResponseWasNotFound(read.Response) {
			return "vaultnotfound", "vaultnotfound", fmt.Errorf("failed to find vault %q (resource group %q)", vaultName, resourceGroup)
		}

		if read.Properties != nil && read.Properties.NetworkAcls != nil && read.Properties.NetworkAcls.IPRules != nil && FindKeyVaultNetworkAclIpRule(*read.Properties.NetworkAcls.IPRules, source) != nil {
			return "found", "found", nil
		}

		return "notfound", "notfound", nil
	}
}

func networkAclVirtualNetworkRuleRefreshFunc(ctx context.Context, client *keyvault.VaultsClient, resourceGroup string, vaultName string, source string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking for completion of Network ACL IP Rule create/update")

		read, err := client.Get(ctx, resourceGroup, vaultName)
		if err != nil && utils.ResponseWasNotFound(read.Response) {
			return "vaultnotfound", "vaultnotfound", fmt.Errorf("failed to find vault %q (resource group %q)", vaultName, resourceGroup)
		}

		if read.Properties != nil && read.Properties.NetworkAcls != nil && read.Properties.NetworkAcls.IPRules != nil && FindKeyVaultNetworkAclVirtualNetworkRule(*read.Properties.NetworkAcls.VirtualNetworkRules, source) != nil {
			return "found", "found", nil
		}

		return "notfound", "notfound", nil
	}
}
