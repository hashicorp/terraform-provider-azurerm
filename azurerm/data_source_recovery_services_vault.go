package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmRecoveryServicesVault() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmRecoveryServicesVaultRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationForDataSourceSchema(),

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"tags": tagsForDataSourceSchema(),

			"sku": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmRecoveryServicesVaultRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).recoveryServicesVaultsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	log.Printf("[DEBUG] Reading Recovery Service Vault %q (resource group %q)", name, resourceGroup)

	vault, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(vault.Response) {
			return fmt.Errorf("Error: Recovery Services Vault %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on Recovery Service Vault %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*vault.ID)
	d.Set("name", vault.Name)
	d.Set("location", azureRMNormalizeLocation(*vault.Location))
	d.Set("resource_group_name", resourceGroup)

	if sku := vault.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	flattenAndSetTags(d, vault.Tags)
	return nil
}
