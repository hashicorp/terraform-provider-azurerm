package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmVirtualMachineRead,
		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func dataSourceArmVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).compute.VMClient
	ctx := meta.(*ArmClient).StopContext

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Virtual Machine %q (Resource Group %q) was not found", name, resGroup)
		}

		return fmt.Errorf("Error making Read request on Virtual Machine %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.SetId(*resp.ID)

	return nil
}
