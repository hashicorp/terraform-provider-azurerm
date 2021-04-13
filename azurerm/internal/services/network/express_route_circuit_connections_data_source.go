package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceExpressRouteCircuitConnection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceExpressRouteCircuitConnectionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"circuit_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func dataSourceExpressRouteCircuitConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitConnectionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	circuitName := d.Get("circuit_name").(string)

	resp, err := client.Get(ctx, resourceGroup, circuitName, "AzurePrivatePeering", name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf(" ExpressRouteCircuitConnection %q (Resource Group %q / circuitName %q) does not exist", name, resourceGroup, circuitName)
		}
		return fmt.Errorf("retrieving ExpressRouteCircuitConnection %q (Resource Group %q / circuitName %q): %+v", name, resourceGroup, circuitName, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for ExpressRouteCircuitConnection %q (Resource Group %q / circuitName %q) ID", name, resourceGroup, circuitName)
	}

	d.SetId(*resp.ID)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("circuit_name", circuitName)
	return nil
}
