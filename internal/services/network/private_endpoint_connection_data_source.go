package network

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourcePrivateEndpointConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourcePrivateEndpointConnectionRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.PrivateLinkName,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"network_interface": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"private_service_connection": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"request_response": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"status": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"private_ip_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePrivateEndpointConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateEndpointClient
	nicsClient := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Private Endpoint %q was not found in Resource Group %q", name, resourceGroup)
		}
		return fmt.Errorf("reading Private Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on Private Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)
	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.PrivateEndpointProperties; props != nil {
		networkInterfaceId := ""
		privateIpAddress := ""

		if nics := props.NetworkInterfaces; nics != nil && len(*nics) > 0 {
			nic := (*nics)[0]
			if nic.ID != nil && *nic.ID != "" {
				networkInterfaceId = *nic.ID
				privateIpAddress = getPrivateIpAddress(ctx, nicsClient, networkInterfaceId)
			}
		}

		if err := d.Set("network_interface", getNetworkInterface(networkInterfaceId)); err != nil {
			return fmt.Errorf("setting `network_interface`: %+v", err)
		}

		if err := d.Set("private_service_connection", dataSourceFlattenPrivateEndpointServiceConnection(props.PrivateLinkServiceConnections, props.ManualPrivateLinkServiceConnections, privateIpAddress)); err != nil {
			return fmt.Errorf("setting `private_service_connection`: %+v", err)
		}
	}

	return nil
}

func getNetworkInterface(networkInterfaceId string) interface{} {
	results := make([]interface{}, 0)

	id, err := parse.NetworkInterfaceID(networkInterfaceId)
	if err != nil {
		return results
	}

	elem := map[string]string{}

	elem["id"] = id.ID()
	elem["name"] = id.Name

	return append(results, elem)
}

func getPrivateIpAddress(ctx context.Context, client *network.InterfacesClient, networkInterfaceId string) string {
	privateIpAddress := ""
	id, err := parse.NetworkInterfaceID(networkInterfaceId)
	if err != nil {
		return privateIpAddress
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return privateIpAddress
	}

	if props := resp.InterfacePropertiesFormat; props != nil {
		if configs := props.IPConfigurations; configs != nil {
			for i, config := range *configs {
				if propFmt := config.InterfaceIPConfigurationPropertiesFormat; propFmt != nil {
					if propFmt.PrivateIPAddress != nil && *propFmt.PrivateIPAddress != "" && i == 0 {
						privateIpAddress = *propFmt.PrivateIPAddress
					}
					break
				}
			}
		}
	}

	return privateIpAddress
}

func dataSourceFlattenPrivateEndpointServiceConnection(serviceConnections *[]network.PrivateLinkServiceConnection, manualServiceConnections *[]network.PrivateLinkServiceConnection, privateIpAddress string) []interface{} {
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
