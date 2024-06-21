// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package servicenetworking

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-11-01/associationsinterface"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-11-01/trafficcontrollerinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicenetworking/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApplicationLoadBalancerSubnetAssociationResource struct{}

type AssociationResourceModel struct {
	Name                      string                 `tfschema:"name"`
	ApplicationLoadBalancerId string                 `tfschema:"application_load_balancer_id"`
	SubnetId                  string                 `tfschema:"subnet_id"`
	Tags                      map[string]interface{} `tfschema:"tags"`
}

var _ sdk.ResourceWithUpdate = ApplicationLoadBalancerSubnetAssociationResource{}

func (t ApplicationLoadBalancerSubnetAssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ApplicationLoadBalancerSubnetAssociationName(),
		},

		"application_load_balancer_id": commonschema.ResourceIDReferenceRequiredForceNew(&associationsinterface.TrafficControllerId{}),

		"subnet_id": commonschema.ResourceIDReferenceRequired(&commonids.SubnetId{}),

		"tags": commonschema.Tags(),
	}
}

func (t ApplicationLoadBalancerSubnetAssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (t ApplicationLoadBalancerSubnetAssociationResource) ModelObject() interface{} {
	return &AssociationResourceModel{}
}

func (t ApplicationLoadBalancerSubnetAssociationResource) ResourceType() string {
	return "azurerm_application_load_balancer_subnet_association"
}

func (t ApplicationLoadBalancerSubnetAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return associationsinterface.ValidateAssociationID
}

func (t ApplicationLoadBalancerSubnetAssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			trafficControllerClient := metadata.Client.ServiceNetworking.TrafficControllerInterface
			client := metadata.Client.ServiceNetworking.AssociationsInterface

			var config AssociationResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			albId, err := trafficcontrollerinterface.ParseTrafficControllerID(config.ApplicationLoadBalancerId)
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

			controller, err := trafficControllerClient.Get(ctx, *albId)
			if err != nil {
				return fmt.Errorf("retrieving parent %s: %+v", *albId, err)
			}

			if controller.Model == nil {
				return fmt.Errorf("retrieving parent %s: model was nil", *albId)
			}

			association := associationsinterface.Association{
				Location: location.Normalize(controller.Model.Location),
				Properties: &associationsinterface.AssociationProperties{
					Subnet: &associationsinterface.AssociationSubnet{
						Id: config.SubnetId,
					},
					AssociationType: associationsinterface.AssociationTypeSubnets,
				},
				Tags: tags.Expand(config.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, association); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (t ApplicationLoadBalancerSubnetAssociationResource) Read() sdk.ResourceFunc {
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
				return fmt.Errorf("retreiving %s: %v", *id, err)
			}

			trafficControllerId := associationsinterface.NewTrafficControllerID(id.SubscriptionId, id.ResourceGroupName, id.TrafficControllerName)
			state := AssociationResourceModel{
				Name:                      id.AssociationName,
				ApplicationLoadBalancerId: trafficControllerId.ID(),
			}

			if model := resp.Model; model != nil {
				state.Tags = tags.Flatten(model.Tags)

				if prop := model.Properties; prop != nil {
					if prop.Subnet != nil {
						parsedSubnetId, err := commonids.ParseSubnetIDInsensitively(prop.Subnet.Id)
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

func (t ApplicationLoadBalancerSubnetAssociationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.AssociationsInterface

			var config AssociationResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			id, err := associationsinterface.ParseAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Although `AssociationSubnetUpdate` defined in the SDK contains the `subnetId`, while per testing it can not be updated
			// Tracked on https://github.com/Azure/azure-rest-api-specs/issues/26657
			associationUpdate := associationsinterface.AssociationUpdate{}

			if metadata.ResourceData.HasChange("tags") {
				associationUpdate.Tags = tags.Expand(config.Tags)
			}

			if _, err = client.Update(ctx, *id, associationUpdate); err != nil {
				return fmt.Errorf("updating %s: %v", *id, err)
			}

			return nil
		},
	}
}

func (t ApplicationLoadBalancerSubnetAssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.AssociationsInterface

			id, err := associationsinterface.ParseAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", *id, err)
			}

			return nil
		},
	}
}
