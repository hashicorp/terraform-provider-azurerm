// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type VirtualHubRoutingIntentModel struct {
	Name            string          `tfschema:"name"`
	VirtualHubId    string          `tfschema:"virtual_hub_id"`
	RoutingPolicies []RoutingPolicy `tfschema:"routing_policy"`
}

type RoutingPolicy struct {
	Destinations []string `tfschema:"destinations"`
	Name         string   `tfschema:"name"`
	NextHop      string   `tfschema:"next_hop"`
}

type VirtualHubRoutingIntentResource struct{}

var _ sdk.ResourceWithUpdate = VirtualHubRoutingIntentResource{}

func (r VirtualHubRoutingIntentResource) ResourceType() string {
	return "azurerm_virtual_hub_routing_intent"
}

func (r VirtualHubRoutingIntentResource) ModelObject() interface{} {
	return &VirtualHubRoutingIntentModel{}
}

func (r VirtualHubRoutingIntentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return virtualwans.ValidateRoutingIntentID
}

func (r VirtualHubRoutingIntentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"virtual_hub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: virtualwans.ValidateVirtualHubID,
		},

		"routing_policy": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"destinations": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"Internet",
								"PrivateTraffic",
							}, false),
						},
					},

					"next_hop": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},
				},
			},
		},
	}
}

func (r VirtualHubRoutingIntentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r VirtualHubRoutingIntentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model VirtualHubRoutingIntentModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.VirtualWANs
			virtualHubId, err := virtualwans.ParseVirtualHubID(model.VirtualHubId)
			if err != nil {
				return err
			}

			id := virtualwans.NewRoutingIntentID(virtualHubId.SubscriptionId, virtualHubId.ResourceGroupName, virtualHubId.VirtualHubName, model.Name)
			existing, err := client.RoutingIntentGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &virtualwans.RoutingIntent{
				Properties: &virtualwans.RoutingIntentProperties{
					RoutingPolicies: expandRoutingPolicy(model.RoutingPolicies),
				},
			}

			if err := client.RoutingIntentCreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r VirtualHubRoutingIntentResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VirtualWANs

			id, err := virtualwans.ParseRoutingIntentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model VirtualHubRoutingIntentModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.RoutingIntentGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("routing_policy") {
				properties.Properties.RoutingPolicies = expandRoutingPolicy(model.RoutingPolicies)
			}

			if err := client.RoutingIntentCreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r VirtualHubRoutingIntentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VirtualWANs

			id, err := virtualwans.ParseRoutingIntentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.RoutingIntentGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := VirtualHubRoutingIntentModel{
				Name:         id.RoutingIntentName,
				VirtualHubId: virtualwans.NewVirtualHubID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.RoutingPolicies = flattenRoutingPolicy(props.RoutingPolicies)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r VirtualHubRoutingIntentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VirtualWANs

			id, err := virtualwans.ParseRoutingIntentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.RoutingIntentDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandRoutingPolicy(input []RoutingPolicy) *[]virtualwans.RoutingPolicy {
	result := make([]virtualwans.RoutingPolicy, 0)
	if len(input) == 0 {
		return &result
	}

	for _, v := range input {
		routingPolicy := virtualwans.RoutingPolicy{
			Destinations: v.Destinations,
			Name:         v.Name,
			NextHop:      v.NextHop,
		}

		result = append(result, routingPolicy)
	}

	return &result
}

func flattenRoutingPolicy(input *[]virtualwans.RoutingPolicy) []RoutingPolicy {
	var result []RoutingPolicy
	if input == nil {
		return result
	}

	for _, v := range *input {
		routingPolicy := RoutingPolicy{
			Destinations: v.Destinations,
			Name:         v.Name,
			NextHop:      v.NextHop,
		}

		result = append(result, routingPolicy)
	}

	return result
}
