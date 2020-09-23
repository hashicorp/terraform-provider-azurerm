package keyvault

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/mgmt/2019-09-01/keyvault"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/set"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmKeyVaultNetworkAcls() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmKeyVaultNetworkAclsCreateUpdate,
		Read:   resourceArmKeyVaultNetworkAclsRead,
		Update: resourceArmKeyVaultNetworkAclsCreateUpdate,
		Delete: resourceArmKeyVaultNetworkAclsDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		//Setting custom timeouts
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"key_vault_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.KeyVaultName,
			},

			"network_acls": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_action": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(keyvault.Allow),
								string(keyvault.Deny),
							}, false),
						},

						"bypass": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(keyvault.None),
								string(keyvault.AzureServices),
							}, false),
						},

						"ip_rules": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},

						"virtual_network_subnet_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: azure.ValidateResourceID,
							},
							Set: set.HashStringIgnoreCase,
						},
					},
				},
			},
		},
	}
}

func resourceArmKeyVaultNetworkAclsCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	keyVaultName := d.Get("key_vault_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	//Lock the resource to prevent external changes
	locks.ByName(keyVaultName, keyVaultResourceName)
	defer locks.UnlockByName(keyVaultName, keyVaultResourceName)

	d.Partial(true)

	keyVault, err := client.Get(ctx, resourceGroup, keyVaultName)
	if err != nil {
		if utils.ResponseWasNotFound(keyVault.Response) {
			return fmt.Errorf("Key Vault %q (Resource Group %q) was not found", keyVaultName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Key Vault %q (Resource Group %q): %+v", keyVaultName, resourceGroup, err)
	}

	if keyVault.Properties == nil {
		return fmt.Errorf("Error retrieving Key Vault %q (Resource Group %q): `properties` was nil", keyVaultName, resourceGroup)
	}

	//Generate a Resource ID for the ACLs as Azure Resource manager doesn't maintain a Resource ID for this resource
	resourceId := fmt.Sprintf("%s/aclsId/rule", *keyVault.ID)
	update := keyvault.VaultPatchParameters{}

	if update.Properties == nil {
		update.Properties = &keyvault.VaultPatchProperties{}
	}

	networkAclsRaw := d.Get("network_acls").([]interface{})
	networkAcls, subnetIds := expandKeyVaultNetworkAcls(networkAclsRaw)

	//Lock on the Virtual Network ID's since modifications in the networking stack are exclusive
	virtualNetworkNames := make([]string, 0)
	for _, v := range subnetIds {
		id, err2 := azure.ParseAzureResourceID(v)
		if err2 != nil {
			return err2
		}

		virtualNetworkName := id.Path["virtualNetworks"]
		if !azure.SliceContainsValue(virtualNetworkNames, virtualNetworkName) {
			virtualNetworkNames = append(virtualNetworkNames, virtualNetworkName)
		}
	}

	locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

	update.Properties.NetworkAcls = networkAcls

	if _, err := client.Update(ctx, resourceGroup, keyVaultName, update); err != nil {
		return fmt.Errorf("Error Updating Azure Key Vault Network ACLs %q (Resource Group %q): %+v", keyVaultName, resourceGroup, err)
	}

	if d.IsNewResource() {
		d.SetId(resourceId)
	}

	d.Partial(false)

	return resourceArmKeyVaultNetworkAclsRead(d, meta)
}

func resourceArmKeyVaultNetworkAclsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	keyVaultName := id.Path["vaults"]

	keyVault, err := client.Get(ctx, resourceGroup, keyVaultName)
	if err != nil {
		if utils.ResponseWasNotFound(keyVault.Response) {
			return fmt.Errorf("Key Vault %q (Resource Group %q) was not found", keyVaultName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Key Vault %q (Resource Group %q): %+v", keyVaultName, resourceGroup, err)
	}

	d.Set("key_vault_name", keyVaultName)
	d.Set("resource_group_name", resourceGroup)
	if location := keyVault.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if rules := keyVault.Properties.NetworkAcls; rules != nil {
		if err := d.Set("network_acls", flattenKeyVaultNetworkAcls(rules)); err != nil {
			return fmt.Errorf("Error setting `network_acls` for KeyVault %q: %+v", keyVaultName, err)
		}
	}

	return nil
}

func resourceArmKeyVaultNetworkAclsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.VaultsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	keyVaultName := id.Path["vaults"]

	locks.ByName(keyVaultName, keyVaultResourceName)
	defer locks.UnlockByName(keyVaultName, keyVaultResourceName)

	keyVault, err := client.Get(ctx, resourceGroup, keyVaultName)
	if err != nil {
		if utils.ResponseWasNotFound(keyVault.Response) {
			return fmt.Errorf("Key Vault %q (Resource Group %q) was not found", keyVaultName, resourceGroup)
		}

		return fmt.Errorf("Error loading Key Vault %q (Resource Group %q): %+v", keyVaultName, resourceGroup, err)
	}

	if keyVault.Properties.NetworkAcls == nil {
		return nil
	}

	//Delete operation doesn't delete anything, instead it resets it to default
	update := keyvault.VaultPatchParameters{
		Properties: &keyvault.VaultPatchProperties{
			NetworkAcls: &keyvault.NetworkRuleSet{
				Bypass:        keyvault.AzureServices,
				DefaultAction: keyvault.Allow,
			},
		},
	}

	if _, err := client.Update(ctx, resourceGroup, keyVaultName, update); err != nil {
		return fmt.Errorf("Error deleting Azure Key Vault Network ACLs %q (Resource Group %q): %+v", keyVaultName, resourceGroup, err)
	}

	return nil
}
