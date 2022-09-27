package automanage

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceAutomanageConfigurationProfileAssignment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAutomanageConfigurationProfileAssignmentRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"vm_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceAutomanageConfigurationProfileAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfileAssignmentClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	vmName := d.Get("vm_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name, vmName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Automanage ConfigurationProfileAssignment %q (Resource Group %q / vmName %q) does not exist", name, resourceGroup, vmName)
		}
		return fmt.Errorf("retrieving Automanage ConfigurationProfileAssignment %q (Resource Group %q / vmName %q): %+v", name, resourceGroup, vmName, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Automanage ConfigurationProfileAssignment %q (Resource Group %q / vmName %q) ID", name, resourceGroup, vmName)
	}

	d.SetId(*resp.ID)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("vm_name", vmName)
	return nil
}
