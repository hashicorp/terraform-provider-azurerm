package apimanagement

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceApiManagementGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceApiManagementGatewayRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"gateway_id": schemaz.SchemaApiManagementChildDataSourceName(),

			"api_management_name": schemaz.SchemaApiManagementDataSourceName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceApiManagementGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	gatewayId := d.Get("gateway_id").(string)

	resp, err := client.Get(ctx, resourceGroup, serviceName, gatewayId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Gateway %q was not found in API Management Service %q / Resource Group %q", gatewayId, serviceName, resourceGroup)
		}

		return fmt.Errorf("making Read request on Gateway %q (API Management Service %q / Resource Group %q): %+v", gatewayId, serviceName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	if props := resp.GatewayContractProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("location", props.LocationData.Name)
	}

	return nil
}
