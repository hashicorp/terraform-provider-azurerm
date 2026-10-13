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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiGatewayClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ApiManagementStandaloneGatewayModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := apigateway.NewGatewayID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			virtualNetworkType := string(apigateway.VirtualNetworkTypeNone)
			if v := model.VirtualNetworkType; v != "" {
				virtualNetworkType = v
			}

			properties := &apigateway.ApiManagementGatewayResource{
				Location: location.Normalize(model.Location),
				Properties: apigateway.ApiManagementGatewayBaseProperties{
					VirtualNetworkType: pointer.ToEnum[apigateway.VirtualNetworkType](virtualNetworkType),
				},
				Sku:  pointer.From(expandGatewaySkuModel(model.Sku)),
				Tags: pointer.To(model.Tags),
			}

			if model.BackendSubnetId != "" {
				properties.Properties.Backend = &apigateway.BackendConfiguration{
					Subnet: &apigateway.BackendSubnetConfiguration{
						Id: pointer.To(model.BackendSubnetId),
					},
				}
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, *properties, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandGatewaySkuModel(inputList []GatewaySkuModel) *apigateway.ApiManagementGatewaySkuProperties {
	if len(inputList) == 0 {
		return nil
	}

	input := &inputList[0]
	return &apigateway.ApiManagementGatewaySkuProperties{
		Capacity: pointer.To(input.Capacity),
		Name:     input.Name,
	}
}
