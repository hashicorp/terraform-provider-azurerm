// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/trafficcontrollerinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (t Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceNetworking.TrafficControllerInterface

			id, err := trafficcontrollerinterface.ParseTrafficControllerID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", metadata.ResourceData.Id(), err)
			}

			state := ApplicationLoadBalancerModel{
				Name:              id.TrafficControllerName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = model.Location
				state.Tags = pointer.From(model.Tags)

				if prop := model.Properties; prop != nil {
					if endpoint := prop.ConfigurationEndpoints; endpoint != nil && len(*endpoint) > 0 {
						state.PrimaryConfigurationEndpoints = (*endpoint)[0]
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}
