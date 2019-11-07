package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	aznet "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPrivateLinkEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPrivateLinkEndpointRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew:     true,
				ValidateFunc: aznet.ValidatePrivateLinkEndpointName,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_service_connection": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_manual_connection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"private_connection_resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subresource_names": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"request_message": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"network_interface_ids": {
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

func dataSourceArmPrivateLinkEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.PrivateEndpointClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Private Link Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on Private Link Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)
	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.PrivateEndpointProperties; props != nil {
		if subnet := props.Subnet; subnet != nil {
			d.Set("subnet_id", subnet.ID)
		}
		if err := d.Set("private_service_connection", flattenArmPrivateLinkEndpointServiceConnection(props.PrivateLinkServiceConnections, props.ManualPrivateLinkServiceConnections)); err != nil {
			return fmt.Errorf("Error setting `private_service_connection`: %+v", err)
		}
		if props.NetworkInterfaces != nil {
			if err := d.Set("network_interface_ids", flattenArmPrivateLinkEndpointInterface(props.NetworkInterfaces)); err != nil {
				return fmt.Errorf("Error setting `network_interface_ids`: %+v", err)
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
