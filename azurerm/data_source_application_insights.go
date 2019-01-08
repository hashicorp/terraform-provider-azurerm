package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmApplicationInsights() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmApplicationInsightsRead,
		Schema: map[string]*schema.Schema{
			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"instrumentation_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmApplicationInsightsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).appInsightsClient
	ctx := meta.(*ArmClient).StopContext

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Application Insights bucket %q (Resource Group %q) was not found", name, resGroup)
		}

		return fmt.Errorf("Error making Read request on Application Insights bucket %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.SetId(*resp.ID)
	d.Set("instrumentation_key", *resp.InstrumentationKey)

	return nil
}
