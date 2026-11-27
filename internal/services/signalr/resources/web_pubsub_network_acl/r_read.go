// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_pubsub_network_acl

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	signalrClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceWebPubsubNetworkACLRead(d *pluginsdk.ResourceData, meta any) error {
	client := signalrClient.From(meta).WebPubSubClient.WebPubSub
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := webpubsub.ParseWebPubSubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	d.Set("web_pubsub_id", id.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.NetworkACLs != nil {
				defaultAction := ""
				if props.NetworkACLs.DefaultAction != nil && *props.NetworkACLs.DefaultAction != "" {
					defaultAction = string(*props.NetworkACLs.DefaultAction)
				}
				d.Set("default_action", defaultAction)

				if err := d.Set("public_network", flattenWebpubsubPublicNetwork(props.NetworkACLs.PublicNetwork)); err != nil {
					return fmt.Errorf("setting `public_network`: %+v", err)
				}

				if err := d.Set("private_endpoint", flattenWebpubsubPrivateEndpoint(props.NetworkACLs.PrivateEndpoints, props.PrivateEndpointConnections)); err != nil {
					return fmt.Errorf("setting `private_endpoint`: %+v", err)
				}
			}
		}
	}

	return nil
}

func flattenWebpubsubPublicNetwork(input *webpubsub.NetworkACL) []any {
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

func flattenWebpubsubPrivateEndpoint(input *[]webpubsub.PrivateEndpointACL, privateEndpointConnections *[]webpubsub.PrivateEndpointConnection) []any {
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
				if !strings.EqualFold(*privateEndpointConnection.Name, item.Name) {
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

var defaultRequestTypes = []webpubsub.WebPubSubRequestType{
	webpubsub.WebPubSubRequestTypeRESTAPI,
	webpubsub.WebPubSubRequestTypeClientConnection,
	webpubsub.WebPubSubRequestTypeServerConnection,
	webpubsub.WebPubSubRequestTypeTrace,
}

func isNewNetworkACL(existing webpubsub.WebPubSubResource) bool {
	if existing.Properties == nil || existing.Properties.NetworkACLs == nil {
		return true
	}

	if existing.Properties.NetworkACLs.DefaultAction == nil || *existing.Properties.NetworkACLs.DefaultAction != webpubsub.ACLActionDeny {
		return false
	}

	if existing.Properties.NetworkACLs.PublicNetwork == nil ||
		existing.Properties.NetworkACLs.PublicNetwork.Deny != nil ||
		(existing.Properties.NetworkACLs.PublicNetwork.Allow == nil || len(*existing.Properties.NetworkACLs.PublicNetwork.Allow) != len(defaultRequestTypes)) {
		return false
	}

	defaultRequestTypeMap := map[webpubsub.WebPubSubRequestType]bool{
		webpubsub.WebPubSubRequestTypeClientConnection: true,
		webpubsub.WebPubSubRequestTypeRESTAPI:          true,
		webpubsub.WebPubSubRequestTypeServerConnection: true,
		webpubsub.WebPubSubRequestTypeTrace:            true,
	}

	for _, allowType := range *existing.Properties.NetworkACLs.PublicNetwork.Allow {
		if !defaultRequestTypeMap[allowType] {
			return false
		}
	}

	if existing.Properties.NetworkACLs.PrivateEndpoints != nil {
		for _, peItem := range *existing.Properties.NetworkACLs.PrivateEndpoints {
			if peItem.Allow == nil || len(*peItem.Allow) != len(defaultRequestTypes) || peItem.Deny != nil {
				return false
			}
			for _, allowType := range *peItem.Allow {
				if !defaultRequestTypeMap[allowType] {
					return false
				}
			}
		}
	}
	return true
}
