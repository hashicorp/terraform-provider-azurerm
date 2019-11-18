package azurerm

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmNetAppAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmNetAppAccountRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`(^[\da-zA-Z])([-\da-zA-Z]{1,62})([\da-zA-Z]$)`),
					`The name must be between 3 and 64 characters in length and begin with a letter or number, end with a letter or number and may contain only letters, numbers or hyphens.`,
				),
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			// Handles tags being interface{} until https://github.com/Azure/azure-rest-api-specs/issues/7447 is fixed
		},
	}
}

func dataSourceArmNetAppAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Netapp.AccountClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: NetApp Account %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading NetApp Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	return nil
}
