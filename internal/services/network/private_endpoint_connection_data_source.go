// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/privateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
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

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

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
	client := meta.(*clients.Client).Network.PrivateEndpoints
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	nicsClient := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := privateendpoints.NewPrivateEndpointID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id, privateendpoints.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.PrivateEndpointName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			networkInterfaceId := ""
			privateIpAddress := ""

			if nics := props.NetworkInterfaces; nics != nil && len(*nics) > 0 {
				nic := (*nics)[0]
				if nic.Id != nil && *nic.Id != "" {
					networkInterfaceId = *nic.Id
					privateIpAddress = getPrivateIpAddress(ctx, nicsClient, networkInterfaceId)
				}
			}

			if err := d.Set("network_interface", flattenNetworkInterface(networkInterfaceId)); err != nil {
				return fmt.Errorf("setting `network_interface`: %+v", err)
			}

			if err := d.Set("private_service_connection", dataSourceFlattenPrivateEndpointServiceConnection(props.PrivateLinkServiceConnections, props.ManualPrivateLinkServiceConnections, privateIpAddress)); err != nil {
				return fmt.Errorf("setting `private_service_connection`: %+v", err)
			}
		}
	}

	return nil
}

func flattenNetworkInterface(networkInterfaceId string) interface{} {
	id, err := parse.NetworkInterfaceID(networkInterfaceId)
	if err != nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"id":   id.ID(),
			"name": id.Name,
		},
	}
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

func dataSourceFlattenPrivateEndpointServiceConnection(serviceConnections *[]privateendpoints.PrivateLinkServiceConnection, manualServiceConnections *[]privateendpoints.PrivateLinkServiceConnection, privateIpAddress string) []interface{} {
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
			if props := item.Properties; props != nil {
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
			if props := item.Properties; props != nil {
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
