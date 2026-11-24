// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_logger_eventhub

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/logger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.LoggerClient_v2024_05_01

			id, err := logger.ParseWorkspaceLoggerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkspaceLoggerGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := ApiManagementWorkspaceLoggerEventhubModel{
				Name:                     id.LoggerId,
				ApiManagementWorkspaceId: logger.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if respModel := resp.Model; respModel != nil {
				if props := respModel.Properties; props != nil {
					if props.LoggerType != logger.LoggerTypeAzureEventHub {
						return fmt.Errorf("expected Logger Type to be %q but got %q", string(logger.LoggerTypeAzureEventHub), string(props.LoggerType))
					}

					model.BufferingEnabled = pointer.From(props.IsBuffered)
					model.Description = pointer.From(props.Description)
					model.ResourceId = pointer.From(props.ResourceId)

					if credentials := props.Credentials; credentials != nil {
						var config ApiManagementWorkspaceLoggerEventhubModel
						if err := metadata.Decode(&config); err != nil {
							return fmt.Errorf("decoding: %+v", err)
						}

						model.Eventhub = flattenApiManagementWorkspaceLoggerEventhub(config, props)
					}
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func flattenApiManagementWorkspaceLoggerEventhub(model ApiManagementWorkspaceLoggerEventhubModel, input *logger.LoggerContractProperties) []EventhubModel {
	outputList := make([]EventhubModel, 0)
	if input == nil || input.Credentials == nil {
		return outputList
	}

	output := EventhubModel{}

	if name, ok := (*input.Credentials)["name"]; ok {
		output.Name = name
	}

	if endpoint, ok := (*input.Credentials)["endpointAddress"]; ok {
		output.EndpointUri = endpoint
	}

	if eventhub := model.Eventhub; len(eventhub) > 0 {
		// The `eventhub.0.connection_string` returned by the Azure API is intentionally masked
		// (e.g. "{{Logger-Credentials--<hash>}}") and does not match the original value provided during creation/update.
		// This is by design to prevent exposing sensitive credentials in API responses.
		// Therefore, the `connection_string` is sourced from the state.
		output.ConnectionString = eventhub[0].ConnectionString

		// The API return `user_assigned_identity_client_id` as an internal identifier (e.g., hex string),
		// not the original AAD `client_id` (UUID). Therefore, we read `user_assigned_identity_client_id` from config.
		if clientId := eventhub[0].UserAssignedIdentityClientId; clientId != "SystemAssigned" {
			output.UserAssignedIdentityClientId = clientId
		}
	}
	outputList = append(outputList, output)

	return outputList
}
