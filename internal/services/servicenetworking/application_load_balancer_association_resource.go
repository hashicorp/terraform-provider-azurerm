package servicenetworking

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/associationsinterface"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-05-01-preview/trafficcontrollerinterface"
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

func (t AssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,62}[a-zA-Z0-9]$`), "the name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens. The value must be 1-64 characters long."),
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

func (t AssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
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

			albId, err := associationsinterface.ParseTrafficControllerID(config.ApplicationLoadBalancerId)
			if err != nil {
				return err
			}

			id := associationsinterface.NewAssociationID(albId.SubscriptionId, albId.ResourceGroupName, albId.TrafficControllerName, config.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of exisiting %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(t.ResourceType(), id)
			}

			controllerId := trafficcontrollerinterface.NewTrafficControllerID(albId.SubscriptionId, albId.ResourceGroupName, albId.TrafficControllerName)
			controller, err := trafficControllerClient.Get(ctx, controllerId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", controllerId, err)
			}

			association := associationsinterface.Association{
				Properties: &associationsinterface.AssociationProperties{
					Subnet: &associationsinterface.AssociationSubnet{
						Id: config.SubnetId,
					},
					AssociationType: associationsinterface.AssociationTypeSubnets,
				},
			}

			if controller.Model != nil {
				association.Location = location.Normalize(controller.Model.Location)
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
						parsedSubnetId, err := commonids.ParseSubnetID(prop.Subnet.Id)
						if err != nil {
							return err
						}
						state.SubnetId = parsedSubnetId.ID()
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

			// Thought `AssociationSubnetUpdate` defined in the SDK contains the `subnetId`, while per testing it can not be updated
			associationUpdate := associationsinterface.AssociationUpdate{}

			if metadata.ResourceData.HasChange("tags") {
				associationUpdate.Tags = &plan.Tags
			}

			if _, err = client.Update(ctx, *id, associationUpdate); err != nil {
				return fmt.Errorf("updating %s: %v", id.ID(), err)
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
