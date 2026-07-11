package testdata

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VirtualHubRoutingIntentModel struct {
	Name         string `tfschema:"name"`
	VirtualHubId string `tfschema:"virtual_hub_id"`
}

type VirtualHubRoutingIntentResource struct{}

var _ sdk.ResourceWithIdentity = VirtualHubRoutingIntentResource{}

func (r VirtualHubRoutingIntentResource) ResourceType() string {
	return "azurerm_virtual_hub_routing_intent"
}

func (r VirtualHubRoutingIntentResource) ModelObject() interface{} {
	return &VirtualHubRoutingIntentModel{}
}

func (r VirtualHubRoutingIntentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return virtualwans.ValidateRoutingIntentID
}

func (r VirtualHubRoutingIntentResource) Identity() resourceids.ResourceId {
	return &virtualwans.RoutingIntentId{}
}

func (r VirtualHubRoutingIntentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
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

			properties := &virtualwans.RoutingIntent{}
			if err := client.RoutingIntentCreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
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

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r VirtualHubRoutingIntentResource) flatten(metadata sdk.ResourceMetaData, id *virtualwans.RoutingIntentId, model *virtualwans.RoutingIntent) error {
	state := VirtualHubRoutingIntentModel{
		Name:         id.RoutingIntentName,
		VirtualHubId: virtualwans.NewVirtualHubID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName).ID(),
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
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
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
