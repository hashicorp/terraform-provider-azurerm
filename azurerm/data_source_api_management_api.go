package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceApiManagementApi() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApiManagementApiRead,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"location": locationForDataSourceSchema(),

			// "tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceApiManagementApiRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementApiClient
	ctx := meta.(*ArmClient).StopContext

	resGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("service_name").(string)
	apiId := d.Get("name").(string)

	resp, err := client.Get(ctx, resGroup, serviceName, apiId)

	log.Printf("Response:")
	log.Printf("%+v\n", resp)

	if err != nil {
		return fmt.Errorf("Error making Read request on API Management Service %q (Resource Group %q): %+v", serviceName, resGroup, err)
	}

	if utils.ResponseWasNotFound(resp.Response) {
		return fmt.Errorf("Error: API Management Service %q (Resource Group %q) was not found", serviceName, resGroup)
	}

	d.SetId(*resp.ID)

	d.Set("service_name", serviceName)
	d.Set("resource_group_name", resGroup)
	d.Set("name", apiId)

	// if location := resp.Location(); location != nil {
	// 	d.Set("location", azureRMNormalizeLocation(*location))
	// }

	// flattenAndSetTags(d, resp.tag)

	return nil
}
