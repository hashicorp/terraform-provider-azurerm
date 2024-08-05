// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicenetworking

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-11-01/frontendsinterface"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-11-01/trafficcontrollerinterface"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type FrontendsResource struct{}

type FrontendsModel struct {
	Name                      string                 `tfschema:"name"`
	ApplicationLoadBalancerId string                 `tfschema:"application_load_balancer_id"`
	Fqdn                      string                 `tfschema:"fully_qualified_domain_name"`
	Tags                      map[string]interface{} `tfschema:"tags"`
}

var _ sdk.Resource = FrontendsResource{}

func (f FrontendsResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"application_load_balancer_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: frontendsinterface.ValidateTrafficControllerID,
		},

		"tags": commonschema.Tags(),
	}
}

func (f FrontendsResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"fully_qualified_domain_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (f FrontendsResource) ModelObject() interface{} {
	return &FrontendsModel{}
}

func (f FrontendsResource) ResourceType() string {
	return "azurerm_application_load_balancer_frontend"
}

func (f FrontendsResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return frontendsinterface.ValidateFrontendID
}

func (f FrontendsResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			trafficControllerClient := metadata.Client.ServiceNetworking.TrafficControllerInterface
			client := metadata.Client.ServiceNetworking.FrontendsInterface

			var config FrontendsModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			trafficControllerId, err := frontendsinterface.ParseTrafficControllerID(config.ApplicationLoadBalancerId)
			if err != nil {
				return err
			}

			controllerId := trafficcontrollerinterface.NewTrafficControllerID(trafficControllerId.SubscriptionId, trafficControllerId.ResourceGroupName, trafficControllerId.TrafficControllerName)
			controller, err := trafficControllerClient.Get(ctx, controllerId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", controllerId, err)
			}

			if controller.Model == nil {
				return fmt.Errorf("retrieving %s: Model was nil", controllerId)
			}

			loc := controller.Model.Location

			id := frontendsinterface.NewFrontendID(trafficControllerId.SubscriptionId, trafficControllerId.ResourceGroupName, trafficControllerId.TrafficControllerName, config.Name)
			resp, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(resp.HttpResponse) {
				return tf.ImportAsExistsError(f.ResourceType(), id.ID())
			}

			frontend := frontendsinterface.Frontend{
				Location:   loc,
				Properties: &frontendsinterface.FrontendProperties{},
				Tags:       tags.Expand(config.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, frontend); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (f FrontendsResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.FrontendsInterface

			id, err := frontendsinterface.ParseFrontendID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", metadata.ResourceData.Id(), err)
			}

			trafficControllerId := frontendsinterface.NewTrafficControllerID(id.SubscriptionId, id.ResourceGroupName, id.TrafficControllerName)
			state := FrontendsModel{
				Name:                      id.FrontendName,
				ApplicationLoadBalancerId: trafficControllerId.ID(),
			}

			if model := resp.Model; model != nil {
				state.Tags = tags.Flatten(model.Tags)

				if prop := model.Properties; prop != nil {
					state.Fqdn = pointer.From(prop.Fqdn)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (f FrontendsResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.FrontendsInterface

			id, err := frontendsinterface.ParseFrontendID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config FrontendsModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			update := frontendsinterface.FrontendUpdate{}

			if metadata.ResourceData.HasChange("tags") {
				update.Tags = tags.Expand(config.Tags)
			}
			if _, err := client.Update(ctx, *id, update); err != nil {
				return fmt.Errorf("updating `azurerm_application_load_balancer_frontend` %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (f FrontendsResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.FrontendsInterface

			id, err := frontendsinterface.ParseFrontendID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %q: %+v", id.ID(), err)
			}

			return nil
		},
	}
}
