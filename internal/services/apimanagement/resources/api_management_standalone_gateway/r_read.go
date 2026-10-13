// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_standalone_gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apigateway"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiGatewayClient

			id, err := apigateway.ParseGatewayID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ApiManagementStandaloneGatewayModel{
				Name:              id.GatewayName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				if backend := model.Properties.Backend; backend != nil {
					if subnet := backend.Subnet; subnet != nil {
						state.BackendSubnetId = pointer.From(subnet.Id)
					}
				}

				virtualNetworkType := ""
				if v := model.Properties.VirtualNetworkType; v != nil && *v != apigateway.VirtualNetworkTypeNone {
					virtualNetworkType = string(*v)
				}
				state.VirtualNetworkType = virtualNetworkType

				state.Sku = flattenGatewaySkuModel(&model.Sku)
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenGatewaySkuModel(input *apigateway.ApiManagementGatewaySkuProperties) []GatewaySkuModel {
	outputList := make([]GatewaySkuModel, 0)
	if input == nil {
		return outputList
	}
	output := GatewaySkuModel{
		Name:     input.Name,
		Capacity: pointer.From(input.Capacity),
	}

	return append(outputList, output)
}
