package apimanagement

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceApiManagementGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceApiManagementGroupRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementChildDataSourceName(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"api_management_name": schemaz.SchemaApiManagementDataSourceName(),

			"display_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"external_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceApiManagementGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
