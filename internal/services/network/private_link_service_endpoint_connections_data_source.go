// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/privateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func dataSourcePrivateLinkServiceEndpointConnections() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourcePrivateLinkServiceEndpointConnectionsRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"service_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"service_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"location": commonschema.LocationComputed(),

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"private_endpoint_connections": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"connection_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"connection_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"private_endpoint_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"private_endpoint_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"action_required": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"description": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"status": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePrivateLinkServiceEndpointConnectionsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateLinkServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serviceId := d.Get("service_id").(string)

	id, err := parse.PrivateLinkServiceID(serviceId)
	if err != nil {
		return fmt.Errorf("parsing %q: %s", serviceId, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: %s was not found", *id)
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on %s: %+v", *id, err)
	}

	d.Set("service_id", serviceId)
	d.Set("service_name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.PrivateLinkServiceProperties; props != nil {
		if err := d.Set("private_endpoint_connections", dataSourceflattenPrivateLinkServicePrivateEndpointConnections(props.PrivateEndpointConnections)); err != nil {
			return fmt.Errorf("setting `private_endpoint_connections`: %+v", err)
		}
	}

	d.SetId(fmt.Sprintf("%s/privateLinkServiceEndpointConnections/%s", *resp.ID, id.Name))

	return nil
}

func dataSourceflattenPrivateLinkServicePrivateEndpointConnections(input *[]network.PrivateEndpointConnection) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})
		if id := item.ID; id != nil {
			v["connection_id"] = *id
		}
		if name := item.Name; name != nil {
			v["connection_name"] = *name
		}

		if props := item.PrivateEndpointConnectionProperties; props != nil {
			if p := props.PrivateEndpoint; p != nil {
				if id := p.ID; id != nil {
					v["private_endpoint_id"] = *id

					id, _ := privateendpoints.ParsePrivateEndpointID(*id)
					if id.PrivateEndpointName != "" {
						v["private_endpoint_name"] = id.PrivateEndpointName
					}
				}
			}

			if s := props.PrivateLinkServiceConnectionState; s != nil {
				if a := s.ActionsRequired; a != nil {
					v["action_required"] = *a
				} else {
					v["action_required"] = "none"
				}
				if d := s.Description; d != nil {
					v["description"] = *d
				}
				if t := s.Status; t != nil {
					v["status"] = *t
				}
			}
		}

		results = append(results, v)
	}

	return results
}
