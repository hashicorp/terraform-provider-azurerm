package azurerm

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	aznet "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPrivateLinkEndpointConnection() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: `The 'azurerm_private_link_endpoint_connection' resource is being deprecated in favour of the renamed version 'azurerm_private_endpoint_connection'.

Information on migrating to the renamed resource can be found here: https://terraform.io/docs/providers/azurerm/guides/migrating-between-renamed-resources.html

As such the existing 'azurerm_private_link_endpoint_connection' resource is deprecated and will be removed in the next major version of the AzureRM Provider (2.0).
`,

		Read: dataSourceArmPrivateLinkEndpointConnectionRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: aznet.ValidatePrivateLinkName,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"private_service_connection": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"request_response": {
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
		},
	}
}

func dataSourceArmPrivateLinkEndpointConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.PrivateEndpointClient
	nicsClient := meta.(*ArmClient).Network.InterfacesClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

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
		privateIpAddress := ""

		if nics := props.NetworkInterfaces; nics != nil && len(*nics) > 0 {
			nic := (*nics)[0]
			if nic.ID != nil && *nic.ID != "" {
				privateIpAddress = getPrivateIpAddress(ctx, nicsClient, *nic.ID)
			}
		}

		if err := d.Set("private_service_connection", dataSourceFlattenArmPrivateLinkEndpointServiceConnection(props.PrivateLinkServiceConnections, props.ManualPrivateLinkServiceConnections, privateIpAddress)); err != nil {
			return fmt.Errorf("Error setting `private_service_connection`: %+v", err)
		}
	}

	return nil
}

func dataSourceFlattenArmPrivateLinkEndpointServiceConnection(serviceConnections *[]network.PrivateLinkServiceConnection, manualServiceConnections *[]network.PrivateLinkServiceConnection, privateIpAddress string) []interface{} {
	results := make([]interface{}, 0)
	if serviceConnections == nil && manualServiceConnections == nil {
		return results
	}

	if serviceConnections != nil {
		for _, item := range *serviceConnections {
			result := make(map[string]interface{})
			result["private_ip_address"] = privateIpAddress

			if v := item.Name; v != nil {
				result["name"] = *v
			}
			if props := item.PrivateLinkServiceConnectionProperties; props != nil {
				if v := props.PrivateLinkServiceConnectionState; v != nil {
					if s := v.Status; s != nil {
						result["status"] = *s
					}
					if d := v.Description; d != nil {
						result["request_response"] = *d
					}
				}
			}

			results = append(results, result)
		}
	}

	if manualServiceConnections != nil {
		for _, item := range *manualServiceConnections {
			result := make(map[string]interface{})
			result["private_ip_address"] = privateIpAddress

			if v := item.Name; v != nil {
				result["name"] = *v
			}
			if props := item.PrivateLinkServiceConnectionProperties; props != nil {
				if v := props.PrivateLinkServiceConnectionState; v != nil {
					if s := v.Status; s != nil {
						result["status"] = *s
					}
					if d := v.Description; d != nil {
						result["request_response"] = *d
					}
				}
			}

			results = append(results, result)
		}
	}

	return results
}
