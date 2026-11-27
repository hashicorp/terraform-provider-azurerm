// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer_subnet_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/associationsinterface"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/trafficcontrollerinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	servicenetworkingClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicenetworking/client"
)

func (t Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			trafficControllerClient := servicenetworkingClient.From(metadata.Client).TrafficControllerInterface
			client := servicenetworkingClient.From(metadata.Client).AssociationsInterface

			var config AssociationResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %v", err)
			}

			albId, err := trafficcontrollerinterface.ParseTrafficControllerID(config.ApplicationLoadBalancerId)
			if err != nil {
				return err
			}

			id := associationsinterface.NewAssociationID(albId.SubscriptionId, albId.ResourceGroupName, albId.TrafficControllerName, config.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(t.ResourceType(), id)
				}
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

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, association, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
