// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/logger"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementWorkspaceLoggerApplicationInsightsModel struct {
	Name                     string                     `tfschema:"name"`
	ApiManagementWorkspaceId string                     `tfschema:"api_management_workspace_id"`
	ApplicationInsights      []ApplicationInsightsModel `tfschema:"application_insights"`
	BufferingEnabled         bool                       `tfschema:"buffering_enabled"`
	Description              string                     `tfschema:"description"`
	ResourceId               string                     `tfschema:"resource_id"`
}

type ApplicationInsightsModel struct {
	InstrumentationKey string `tfschema:"instrumentation_key"`
	ConnectionString   string `tfschema:"connection_string"`
}

type ApiManagementWorkspaceLoggerApplicationInsightsResource struct{}

var _ sdk.ResourceWithUpdate = ApiManagementWorkspaceLoggerApplicationInsightsResource{}

func (r ApiManagementWorkspaceLoggerApplicationInsightsResource) ResourceType() string {
	return "azurerm_api_management_workspace_logger_application_insights"
}

func (r ApiManagementWorkspaceLoggerApplicationInsightsResource) ModelObject() interface{} {
	return &ApiManagementWorkspaceLoggerApplicationInsightsModel{}
}

func (r ApiManagementWorkspaceLoggerApplicationInsightsResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return logger.ValidateWorkspaceLoggerID
}

func (r ApiManagementWorkspaceLoggerApplicationInsightsResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementChildName(),

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&logger.WorkspaceId{}),

		"application_insights": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"connection_string": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ExactlyOneOf: []string{"application_insights.0.connection_string", "application_insights.0.instrumentation_key"},
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"instrumentation_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ExactlyOneOf: []string{"application_insights.0.connection_string", "application_insights.0.instrumentation_key"},
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"buffering_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: azure.ValidateResourceID,
		},
	}
}

func (r ApiManagementWorkspaceLoggerApplicationInsightsResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiManagementWorkspaceLoggerApplicationInsightsResource) Create() sdk.ResourceFunc {
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

func (r ApiManagementWorkspaceLoggerApplicationInsightsResource) Read() sdk.ResourceFunc {
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

func (r ApiManagementWorkspaceLoggerApplicationInsightsResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.LoggerClient_v2024_05_01

			var model ApiManagementWorkspaceLoggerApplicationInsightsModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := logger.ParseWorkspaceLoggerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkspaceLoggerGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			payload := resp.Model
			payload.Properties.Credentials = expandApiManagementWorkspaceLoggerApplicationInsights(model.ApplicationInsights)

			if metadata.ResourceData.HasChange("buffering_enabled") {
				payload.Properties.IsBuffered = pointer.To(model.BufferingEnabled)
			}

			if metadata.ResourceData.HasChange("description") {
				payload.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("resource_id") {
				payload.Properties.ResourceId = pointer.To(model.ResourceId)
			}

			if _, err := client.WorkspaceLoggerCreateOrUpdate(ctx, *id, *payload, logger.WorkspaceLoggerCreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceLoggerApplicationInsightsResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.LoggerClient_v2024_05_01

			id, err := logger.ParseWorkspaceLoggerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.WorkspaceLoggerDelete(ctx, *id, logger.WorkspaceLoggerDeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

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
