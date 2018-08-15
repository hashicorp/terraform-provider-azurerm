package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2016-06-01/recoveryservices"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRecoveryServicesVault() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRecoveryServicesVaultCreateUpdate,
		Read:   resourceArmRecoveryServicesVaultRead,
		Update: resourceArmRecoveryServicesVaultCreateUpdate,
		Delete: resourceArmRecoveryServicesVaultDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{1,49}$"),
					"Recovery Service Vault name must be 2 - 50 characters long, start with a letter, contain only letters, numbers and hyphens.",
				),
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"tags": tagsSchema(),

			"sku": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(recoveryservices.RS0),
					string(recoveryservices.Standard),
				}, true),
			},
		},
	}
}

func resourceArmRecoveryServicesVaultCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).recoveryServicesVaultsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	log.Printf("[DEBUG] Creating/updating Recovery Service Vault %q (resource group %q)", name, resourceGroup)

	//build vault struct
	vault := recoveryservices.Vault{
		Location: utils.String(location),
		Tags:     expandTags(tags),
		Sku: &recoveryservices.Sku{
			Name: recoveryservices.SkuName(d.Get("sku").(string)),
		},
		Properties: &recoveryservices.VaultProperties{},
	}

	//create recovery services vault
	vault, err := client.CreateOrUpdate(ctx, resourceGroup, name, vault)
	if err != nil {
		return fmt.Errorf("Error creating/updating Recovery Service Vault %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*vault.ID)

	return resourceArmRecoveryServicesVaultRead(d, meta)
}

func resourceArmRecoveryServicesVaultRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).recoveryServicesVaultsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["vaults"]
	resourceGroup := id.ResourceGroup

	log.Printf("[DEBUG] Reading Recovery Service Vault %q (resource group %q)", name, resourceGroup)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Recovery Service Vault %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmRecoveryServicesVaultDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).recoveryServicesVaultsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["vaults"]
	resourceGroup := id.ResourceGroup

	log.Printf("[DEBUG] Deleting Recovery Service Vault %q (resource group %q)", name, resourceGroup)

	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error issuing delete request for Recovery Service Vault %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}
