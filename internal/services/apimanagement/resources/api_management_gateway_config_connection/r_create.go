// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_gateway_config_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apigateway"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apigatewayconfigconnection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiGatewayConfigConnectionClient

			var model ApiManagementGatewayConfigConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			gatewayId, err := apigateway.ParseGatewayID(model.ApiManagementGatewayId)
			if err != nil {
				return err
			}

			id := apigatewayconfigconnection.NewConfigConnectionID(gatewayId.SubscriptionId, gatewayId.ResourceGroupName, gatewayId.GatewayName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := apigatewayconfigconnection.ApiManagementGatewayConfigConnectionResource{
				Properties: apigatewayconfigconnection.GatewayConfigConnectionBaseProperties{
					SourceId: pointer.To(model.ApiManagementWorkspaceId),
				},
			}

			if len(model.Hostnames) > 0 {
				payload.Properties.Hostnames = pointer.To(model.Hostnames)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
