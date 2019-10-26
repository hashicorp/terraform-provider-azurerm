package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPrivateLinkEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPrivateLinkEndpointRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
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
								Type:     schema.TypeString,
								Computed: true,
							},
						},
						"request_message": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"provisioning_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_address": {
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
					Type:     schema.TypeString,
					Computed: true,
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmPrivateLinkEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.PrivateEndpointClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["privateEndpoints"]

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Private Link Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.PrivateEndpointProperties; props != nil {
		if subnet := props.Subnet; subnet != nil {
			d.Set("subnet_id", subnet.ID)
		}

		privateIpAddress := ""

		if props.NetworkInterfaces != nil {
			if err := d.Set("network_interface_ids", flattenArmPrivateLinkEndpointInterface(props.NetworkInterfaces)); err != nil {
				return fmt.Errorf("Error setting `network_interface_ids`: %+v", err)
			}

			// now we need to get the nic to get the private ip address for the private link endpoint
			client := meta.(*ArmClient).Network.InterfacesClient
			ctx := meta.(*ArmClient).StopContext

			nic := d.Get("network_interface_ids").([]interface{})

			nicId, err := azure.ParseAzureResourceID(nic[0].(string))
			if err != nil {
				return err
			}
			nicName := nicId.Path["networkInterfaces"]

			nicResp, err := client.Get(ctx, resourceGroup, nicName, "")
			if err != nil {
				if utils.ResponseWasNotFound(nicResp.Response) {
					return fmt.Errorf("Azure Network Interface %q (Resource Group %q): %+v", nicName, resourceGroup, err)
				}
				return fmt.Errorf("Error making Read request on Azure Network Interface %q (Resource Group %q): %+v", nicName, resourceGroup, err)
			}

			if nicProps := nicResp.InterfacePropertiesFormat; nicProps != nil {
				if configs := nicProps.IPConfigurations; configs != nil {
					for i, config := range *nicProps.IPConfigurations {
						if ipProps := config.InterfaceIPConfigurationPropertiesFormat; ipProps != nil {
							if v := ipProps.PrivateIPAddress; v != nil {
								if i == 0 {
									privateIpAddress = *v
								}
							}
						}
					}
				}
			}
		}

		if err := d.Set("private_service_connection", flattenArmPrivateLinkEndpointServiceConnection(props.PrivateLinkServiceConnections, props.ManualPrivateLinkServiceConnections, privateIpAddress)); err != nil {
			return fmt.Errorf("Error setting `private_service_connection`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
