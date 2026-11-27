// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer_subnet_association

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/associationsinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (t Resource) Read() sdk.ResourceFunc {
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
				return fmt.Errorf("retrieving %s: %v", *id, err)
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
