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

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiGatewayConfigConnectionClient

			id, err := apigatewayconfigconnection.ParseConfigConnectionID(metadata.ResourceData.Id())
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

			state := ApiManagementGatewayConfigConnectionModel{
				Name:                   id.ConfigConnectionName,
				ApiManagementGatewayId: apigateway.NewGatewayID(id.SubscriptionId, id.ResourceGroupName, id.GatewayName).ID(),
			}

			if model := resp.Model; model != nil {
				state.ApiManagementWorkspaceId = pointer.From(model.Properties.SourceId)
				state.Hostnames = pointer.From(model.Properties.Hostnames)
			}

			return metadata.Encode(&state)
		},
	}
}
