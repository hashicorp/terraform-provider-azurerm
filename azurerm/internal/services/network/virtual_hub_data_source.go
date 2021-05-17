package network

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceVirtualHub() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVirtualHubRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.VirtualHubName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"address_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"virtual_wan_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_route_table": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"labels": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"route": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"destinations": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},

									"destinations_type": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"next_hop": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"next_hop_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceVirtualHubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualHubClient
	routeTableClient := meta.(*clients.Client).Network.HubRouteTableClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Virtual Hub %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.VirtualHubProperties; props != nil {
		d.Set("address_prefix", props.AddressPrefix)

		var virtualWanId *string
		if props.VirtualWan != nil {
			virtualWanId = props.VirtualWan.ID
		}
		d.Set("virtual_wan_id", virtualWanId)
	}

	defaultRouteResp, err := routeTableClient.Get(ctx, resourceGroup, *resp.Name, defaultRouteTable)
	if err != nil {
		if !utils.ResponseWasForbidden(defaultRouteResp.Response) && !utils.ResponseWasNotFound(defaultRouteResp.Response) {
			return fmt.Errorf("retrieving `defaultRouteTable` for VirtualHub: %+v", err)
		}
	}
	flattenedDefaultRouteTable := flattenDefaultRouteTable(defaultRouteResp.HubRouteTableProperties)
	if err := d.Set("default_route_table", flattenedDefaultRouteTable); err != nil {
		return fmt.Errorf("setting `default_route_table` for VirtualHub: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
