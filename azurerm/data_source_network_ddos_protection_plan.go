package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceNetworkDDoSProtectionPlan() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkDDoSProtectionPlanRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"virtual_network_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func dataSourceNetworkDDoSProtectionPlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.DDOSProtectionPlansClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	plan, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(plan.Response) {
			return fmt.Errorf("Error DDoS Protection Plan %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on DDoS Protection Plan %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", plan.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := plan.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := plan.DdosProtectionPlanPropertiesFormat; props != nil {
		vNetIDs := flattenArmNetworkDDoSProtectionPlanVirtualNetworkIDs(props.VirtualNetworks)
		if err := d.Set("virtual_network_ids", vNetIDs); err != nil {
			return fmt.Errorf("Error setting `virtual_network_ids`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, plan.Tags)
}
