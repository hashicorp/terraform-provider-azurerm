// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package signalr_service_network_acl

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSignalRServiceNetworkACLRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSignalRID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("signalr_service_id", id.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil && props.NetworkACLs != nil {
			defaultAction := ""
			if props.NetworkACLs.DefaultAction != nil {
				defaultAction = string(*props.NetworkACLs.DefaultAction)
			}
			d.Set("default_action", defaultAction)

			if err := d.Set("public_network", flattenSignalRServicePublicNetwork(props.NetworkACLs.PublicNetwork)); err != nil {
				return fmt.Errorf("setting `public_network`: %+v", err)
			}

			if err := d.Set("private_endpoint", flattenSignalRServicePrivateEndpoint(props.NetworkACLs.PrivateEndpoints, props.PrivateEndpointConnections)); err != nil {
				return fmt.Errorf("setting `private_endpoint`: %+v", err)
			}
		}
	}

	return nil
}

func flattenSignalRServicePublicNetwork(input *signalr.NetworkACL) []any {
	if input == nil {
		return make([]any, 0)
	}

	allowRequestTypes := make([]string, 0)
	if input.Allow != nil {
		for _, item := range *input.Allow {
			allowRequestTypes = append(allowRequestTypes, string(item))
		}
	}
	allow := pluginsdk.FlattenSlice(&allowRequestTypes)

	deniedRequestTypes := make([]string, 0)
	if input.Deny != nil {
		for _, item := range *input.Deny {
			deniedRequestTypes = append(deniedRequestTypes, string(item))
		}
	}
	deny := pluginsdk.FlattenSlice(&deniedRequestTypes)

	return []any{
		map[string]any{
			"allowed_request_types": allow,
			"denied_request_types":  deny,
		},
	}
}

func flattenSignalRServicePrivateEndpoint(input *[]signalr.PrivateEndpointACL, privateEndpointConnections *[]signalr.PrivateEndpointConnection) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if privateEndpointConnections != nil {
			for _, privateEndpointConnection := range *privateEndpointConnections {
				if privateEndpointConnection.Name == nil || privateEndpointConnection.Properties == nil {
					continue
				}
				if *privateEndpointConnection.Name != item.Name {
					continue
				}
				props := privateEndpointConnection.Properties
				if props.PrivateEndpoint == nil || props.PrivateEndpoint.Id == nil {
					continue
				}

				allowedRequestTypes := make([]string, 0)
				if item.Allow != nil {
					for _, item := range *item.Allow {
						allowedRequestTypes = append(allowedRequestTypes, string(item))
					}
				}
				allow := pluginsdk.FlattenSlice(&allowedRequestTypes)

				deniedRequestTypes := make([]string, 0)
				if item.Deny != nil {
					for _, item := range *item.Deny {
						deniedRequestTypes = append(deniedRequestTypes, string(item))
					}
				}
				deny := pluginsdk.FlattenSlice(&deniedRequestTypes)

				results = append(results, map[string]any{
					"id":                    *props.PrivateEndpoint.Id,
					"allowed_request_types": allow,
					"denied_request_types":  deny,
				})

				break
			}
		}
	}

	return results
}
