package azurerm

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmDedicatedHostGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDedicatedHostGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"platform_fault_domain_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmDedicatedHostGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Compute.DedicatedHostGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Dedicated Host Group %q (Resource Group %q) was not found", name, resourceGroupName)
		}
		return fmt.Errorf("Error reading Dedicated Host Group %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroupName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if dedicatedHostGroupProperties := resp.DedicatedHostGroupProperties; dedicatedHostGroupProperties != nil {
		if err := d.Set("hosts", flattenArmDedicatedHostGroupSubResourceReadOnly(dedicatedHostGroupProperties.Hosts)); err != nil {
			return fmt.Errorf("Error setting `hosts`: %+v", err)
		}
		d.Set("platform_fault_domain_count", int(*dedicatedHostGroupProperties.PlatformFaultDomainCount))
	}
	d.Set("type", resp.Type)
	if resp.Zones != nil {
		d.Set("zones", utils.FlattenStringSlice(resp.Zones))
	}

	return nil
}
