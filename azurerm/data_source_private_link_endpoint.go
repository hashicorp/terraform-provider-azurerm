package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPrivateLinkEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPrivateLinkEndpointRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"manual_private_link_service_connection": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_link_service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"state_action_required": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"request_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"private_link_service_connection": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_link_service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"state_action_required": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"request_message": {
							Type:     schema.TypeString,
							Computed: true,
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
			
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmPrivateLinkEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PrivateEndpointClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Private Endpoint %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading Private Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if privateEndpointProperties := resp.PrivateEndpointProperties; privateEndpointProperties != nil {
		if err := d.Set("manual_private_link_service_connection", flattenArmPrivateEndpointPrivateLinkServiceConnection(privateEndpointProperties.ManualPrivateLinkServiceConnections)); err != nil {
			return fmt.Errorf("Error setting `manual_private_link_service_connection`: %+v", err)
		}
		if err := d.Set("network_interfaces", flattenArmPrivateEndpointInterface(privateEndpointProperties.NetworkInterfaces)); err != nil {
			return fmt.Errorf("Error setting `network_interfaces`: %+v", err)
		}
		if err := d.Set("private_link_service_connection", flattenArmPrivateEndpointPrivateLinkServiceConnection(privateEndpointProperties.PrivateLinkServiceConnections)); err != nil {
			return fmt.Errorf("Error setting `private_link_service_connection`: %+v", err)
		}
		if subnet := privateEndpointProperties.Subnet; subnet != nil {
			d.Set("subnet_id", subnet.ID)
		}
	}
	tags.FlattenAndSet(d, resp.Tags)

	return nil
}
