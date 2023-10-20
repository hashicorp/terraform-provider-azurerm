package servicenetworking

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/associationsinterface"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/trafficcontrollerinterface"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AssociationResource struct{}

type AssociationModel struct {
	Name                      string            `tfschema:"name"`
	ApplicationLoadBalancerId string            `tfschema:"application_load_balancer_id"`
	SubnetId                  string            `tfschema:"subnet_id"`
	Tags                      map[string]string `tfschema:"tags"`
}

var _ sdk.ResourceWithUpdate = AssociationResource{}

func (t AssociationResource) Arguments() map[string]*schema.Schema {
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
			ValidateFunc: associationsinterface.ValidateTrafficControllerID,
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"tags": commonschema.Tags(),
	}
}

func (t AssociationResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (t AssociationResource) ModelObject() interface{} {
	return &AssociationModel{}
}

func (t AssociationResource) ResourceType() string {
	return "azurerm_application_load_balancer_association"
}

func (t AssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return associationsinterface.ValidateAssociationID
}
func (t AssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			trafficControllerClient := metadata.Client.ServiceNetworking.TrafficControllerInterface
			client := metadata.Client.ServiceNetworking.AssociationsInterface

			var config AssociationModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			parsedTrafficControllerId, err := associationsinterface.ParseTrafficControllerID(config.ApplicationLoadBalancerId)
			if err != nil {
				return err
			}

			controllerId := trafficcontrollerinterface.NewTrafficControllerID(parsedTrafficControllerId.SubscriptionId, parsedTrafficControllerId.ResourceGroupName, parsedTrafficControllerId.TrafficControllerName)
			controller, err := trafficControllerClient.Get(ctx, controllerId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", controllerId, err)
			}

			if controller.Model == nil {
				return fmt.Errorf("retrieving %s: Model was nil", controllerId)
			}

			loc := controller.Model.Location

			id := associationsinterface.NewAssociationID(parsedTrafficControllerId.SubscriptionId, parsedTrafficControllerId.ResourceGroupName, parsedTrafficControllerId.TrafficControllerName, config.Name)
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of exisiting %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(t.ResourceType(), id.ID())
			}

			association := associationsinterface.Association{
				Location: location.Normalize(loc),
				Properties: &associationsinterface.AssociationProperties{
					Subnet: &associationsinterface.AssociationSubnet{
						Id: config.SubnetId,
					},
					AssociationType: associationsinterface.AssociationTypeSubnets,
				},
			}

			if len(config.Tags) > 0 {
				association.Tags = &config.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, association); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (t AssociationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.AssociationsInterface

			id, err := associationsinterface.ParseAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retreiving %s: %v", id.ID(), err)
			}

			trafficControllerId := associationsinterface.NewTrafficControllerID(id.SubscriptionId, id.ResourceGroupName, id.TrafficControllerName)
			state := AssociationModel{
				Name:                      id.AssociationName,
				ApplicationLoadBalancerId: trafficControllerId.ID(),
			}

			if model := resp.Model; model != nil {
				state.Tags = pointer.From(model.Tags)

				if prop := model.Properties; prop != nil {
					if prop.Subnet != nil {
						state.SubnetId = prop.Subnet.Id
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (t AssociationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.AssociationsInterface

			var plan AssociationModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			id, err := associationsinterface.ParseAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing id %v", err)
			}

			// thought `AssociationSubnetUpdate` is defined in the SDK, per testing the subnet id can not be updated.
			associationUpdate := associationsinterface.AssociationUpdate{}

			if metadata.ResourceData.HasChange("tags") {
				associationUpdate.Tags = &plan.Tags
			}

			if _, err = client.Update(ctx, *id, associationUpdate); err != nil {
				return fmt.Errorf("updating `azurerm_application_load_balancer_association` %s: %v", id.ID(), err)
			}

			return nil
		},
	}
}

func (t AssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.AssociationsInterface

			id, err := associationsinterface.ParseAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", id.ID(), err)
			}

			return nil
		},
	}
}
