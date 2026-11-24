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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.LoggerClient_v2024_05_01

			var model ApiManagementWorkspaceLoggerApplicationInsightsModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := logger.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			id := logger.NewWorkspaceLoggerID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.ServiceName, workspaceId.WorkspaceId, model.Name)
			existing, err := client.WorkspaceLoggerGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := logger.LoggerContract{
				Properties: &logger.LoggerContractProperties{
					LoggerType:  logger.LoggerTypeApplicationInsights,
					IsBuffered:  pointer.To(model.BufferingEnabled),
					Credentials: expandApiManagementWorkspaceLoggerApplicationInsights(model.ApplicationInsights),
				},
			}

			if model.Description != "" {
				parameters.Properties.Description = pointer.To(model.Description)
			}

			if model.ResourceId != "" {
				parameters.Properties.ResourceId = pointer.To(model.ResourceId)
			}

			if _, err := client.WorkspaceLoggerCreateOrUpdate(ctx, id, parameters, logger.WorkspaceLoggerCreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandApiManagementWorkspaceLoggerApplicationInsights(inputs []ApplicationInsightsModel) *map[string]string {
	if len(inputs) == 0 {
		return nil
	}

	input := &inputs[0]
	result := make(map[string]string)
	if input.InstrumentationKey != "" {
		result["instrumentationKey"] = input.InstrumentationKey
	}
	if input.ConnectionString != "" {
		result["connectionString"] = input.ConnectionString
	}
	return &result
}
