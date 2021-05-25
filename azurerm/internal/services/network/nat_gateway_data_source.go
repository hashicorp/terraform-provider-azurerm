package network

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceNatGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceNatGatewayRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.NatGatewayName,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"idle_timeout_in_minutes": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"public_ip_address_ids": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"public_ip_prefix_ids": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"resource_guid": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"zones": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceNatGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Nat Gateway %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading Nat Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read NAT Gateway %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", resp.Sku.Name)
	}

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.NatGatewayPropertiesFormat; props != nil {
		d.Set("idle_timeout_in_minutes", props.IdleTimeoutInMinutes)
		d.Set("resource_guid", props.ResourceGUID)

		if err := d.Set("public_ip_address_ids", flattenNetworkSubResourceID(props.PublicIPAddresses)); err != nil {
			return fmt.Errorf("Error setting `public_ip_address_ids`: %+v", err)
		}

		if err := d.Set("public_ip_prefix_ids", flattenNetworkSubResourceID(props.PublicIPPrefixes)); err != nil {
			return fmt.Errorf("Error setting `public_ip_prefix_ids`: %+v", err)
		}
	}

	if err := d.Set("zones", utils.FlattenStringSlice(resp.Zones)); err != nil {
		return fmt.Errorf("Error setting `zones`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
