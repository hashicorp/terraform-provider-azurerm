// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_workspace_logger_application_insights

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

			model := ApiManagementWorkspaceLoggerApplicationInsightsModel{
				Name:                     id.LoggerId,
				ApiManagementWorkspaceId: logger.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if respModel := resp.Model; respModel != nil {
				if props := respModel.Properties; props != nil {
					if props.LoggerType != logger.LoggerTypeApplicationInsights {
						return fmt.Errorf("expected Logger Type to be %q but got %q", string(logger.LoggerTypeApplicationInsights), string(props.LoggerType))
					}

					model.BufferingEnabled = pointer.From(props.IsBuffered)
					model.Description = pointer.From(props.Description)
					model.ResourceId = pointer.From(props.ResourceId)

					var config ApiManagementWorkspaceLoggerApplicationInsightsModel
					if err := metadata.Decode(&config); err != nil {
						return fmt.Errorf("decoding: %+v", err)
					}

					// The `application_insights.0.instrumentation_key` and `application_insights.0.connection_string` returned by the Azure API is intentionally masked
					// (e.g. "{{Logger-Credentials--<hash>}}") and does not match the original value provided during creation/update.
					// This is by design to prevent exposing sensitive credentials in API responses.
					// Therefore, the `application_insights` is sourced from the state.
					model.ApplicationInsights = config.ApplicationInsights
				}
			}

			return metadata.Encode(&model)
		},
	}
}
