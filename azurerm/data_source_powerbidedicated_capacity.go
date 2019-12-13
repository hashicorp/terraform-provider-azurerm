package azurerm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	azpowerbidedicated "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/powerbidedicated"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPowerBIDedicatedCapacity() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPowerBIDedicatedCapacityRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azpowerbidedicated.ValidateCapacityName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"administrators": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"sku": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmPowerBIDedicatedCapacityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).PowerBIDedicated.CapacityClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.GetDetails(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Capacity %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading Capacity %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Error retrieving Capacity %q (Resource Group %q): ID was nil or empty", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if props := resp.DedicatedCapacityProperties; props != nil {
		if err := d.Set("administrators", utils.FlattenStringSlice(props.Administration.Members)); err != nil {
			return fmt.Errorf("Error setting `administrators`: %+v", err)
		}
	}
	if err := d.Set("sku", resp.Sku.Name); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}

	return nil
}
