package apimanagement

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceApiManagementGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApiManagementGroupRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": azure.SchemaApiManagementChildDataSourceName(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"api_management_name": azure.SchemaApiManagementDataSourceName(),

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceApiManagementGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GroupClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroup, serviceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Group %q (Resource Group %q / API Management Service %q) was not found", name, resourceGroup, serviceName)
		}

		return fmt.Errorf("making Read request for Group %q (Resource Group %q / API Management Service %q): %+v", name, resourceGroup, serviceName, err)
	}

	id, err := parse.GroupID(*resp.ID)
	if err != nil {
		return fmt.Errorf("Error parsing Group ID %q", *resp.ID)
	}

	d.SetId(id.ID())

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	if properties := resp.GroupContractProperties; properties != nil {
		d.Set("display_name", properties.DisplayName)
		d.Set("description", properties.Description)
		d.Set("external_id", properties.ExternalID)
		d.Set("type", string(properties.Type))
	}

	return nil
}
