package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmAvailabilitySet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmAvailabilitySetRead,
		Schema: map[string]*schema.Schema{
			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"platform_update_domain_count": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"platform_fault_domain_count": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"managed": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceArmAvailabilitySetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).availSetClient
	ctx := meta.(*ArmClient).StopContext

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Availability Set %q (Resource Group %q) was not found", name, resGroup)
		}

		return fmt.Errorf("Error making Read request on Availability Set %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.SetId(*resp.ID)
	d.Set("location", resp.Location)
	d.Set("platform_update_domain_count", resp.PlatformUpdateDomainCount)
	d.Set("platform_fault_domain_count", resp.PlatformFaultDomainCount)
	flattenAndSetTags(d, resp.Tags)

	return nil
}
