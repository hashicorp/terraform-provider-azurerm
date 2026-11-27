// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer_frontend

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/frontendsinterface"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/trafficcontrollerinterface"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	servicenetworkingClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicenetworking/client"
)

func (f Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			trafficControllerClient := servicenetworkingClient.From(metadata.Client).TrafficControllerInterface
			client := servicenetworkingClient.From(metadata.Client).FrontendsInterface

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

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				resp, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(resp.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(resp.HttpResponse) {
					return tf.ImportAsExistsError(f.ResourceType(), id.ID())
				}
			}

			frontend := frontendsinterface.Frontend{
				Location:   loc,
				Properties: &frontendsinterface.FrontendProperties{},
				Tags:       tags.Expand(config.Tags),
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, frontend, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}
