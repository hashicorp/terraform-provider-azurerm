package servicenetworking

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/trafficcontrollerinterface"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ALBResource struct{}

type ContainerApplicationGatewayModel struct {
	Name                   string            `tfschema:"name"`
	ResourceGroupName      string            `tfschema:"resource_group_name"`
	Location               string            `tfschema:"location"`
	ConfigurationEndpoints []string          `tfschema:"configuration_endpoint"`
	Tags                   map[string]string `tfschema:"tags"`
}

var _ sdk.ResourceWithUpdate = ALBResource{}

func (t ALBResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"tags": tags.Schema(),
	}
}

func (t ALBResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"configuration_endpoint": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (t ALBResource) ModelObject() interface{} {
	return &ContainerApplicationGatewayModel{}
}

func (t ALBResource) ResourceType() string {
	return "azurerm_alb"
}

func (t ALBResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return trafficcontrollerinterface.ValidateTrafficControllerID
}

func (t ALBResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan ContainerApplicationGatewayModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			client := metadata.Client.ServiceNetworking.ServiceNetworkingClient
			SubscriptionId := metadata.Client.Account.SubscriptionId

			id := trafficcontrollerinterface.NewTrafficControllerID(SubscriptionId, plan.ResourceGroupName, plan.Name)

			existing, err := client.TrafficControllerInterface.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(t.ResourceType(), id)
			}

			controller := trafficcontrollerinterface.TrafficController{
				Location: location.Normalize(plan.Location),
				Tags:     pointer.To(plan.Tags),
			}

			if err = client.TrafficControllerInterface.CreateOrUpdateThenPoll(ctx, id, controller); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (t ALBResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.ServiceNetworkingClient

			id, err := trafficcontrollerinterface.ParseTrafficControllerID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			resp, err := client.TrafficControllerInterface.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", metadata.ResourceData.Id(), err)
			}

			state := ContainerApplicationGatewayModel{
				Name:              id.TrafficControllerName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = model.Location
				state.Tags = pointer.From(model.Tags)

				if prop := model.Properties; prop != nil {
					state.ConfigurationEndpoints = pointer.From(prop.ConfigurationEndpoints)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (t ALBResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan ContainerApplicationGatewayModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %v", err)
			}
			client := metadata.Client.ServiceNetworking.ServiceNetworkingClient

			id, err := trafficcontrollerinterface.ParseTrafficControllerID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			existing, err := client.TrafficControllerInterface.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retreiving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("existing Traffic Controller %s has no model", id)
			}

			controller := existing.Model

			if metadata.ResourceData.HasChange("tags") {
				if len(plan.Tags) > 0 {
					controller.Tags = pointer.To(plan.Tags)
				} else {
					controller.Tags = nil
				}
			}

			if err = client.TrafficControllerInterface.CreateOrUpdateThenPoll(ctx, *id, *controller); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (t ALBResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := trafficcontrollerinterface.ParseTrafficControllerID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			client := metadata.Client.ServiceNetworking.ServiceNetworkingClient.TrafficControllerInterface

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", metadata.ResourceData.Id(), err)
			}

			return nil
		},
	}
}
