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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementWorkspaceLoggerEventhubModel struct {
	Name                     string          `tfschema:"name"`
	ApiManagementWorkspaceId string          `tfschema:"api_management_workspace_id"`
	Eventhub                 []EventhubModel `tfschema:"eventhub"`
	BufferingEnabled         bool            `tfschema:"buffering_enabled"`
	Description              string          `tfschema:"description"`
	ResourceId               string          `tfschema:"resource_id"`
}

type EventhubModel struct {
	Name                         string `tfschema:"name"`
	ConnectionString             string `tfschema:"connection_string"`
	EndpointUri                  string `tfschema:"endpoint_uri"`
	UserAssignedIdentityClientId string `tfschema:"user_assigned_identity_client_id"`
}

type ApiManagementWorkspaceLoggerEventhubResource struct{}

var _ sdk.ResourceWithUpdate = ApiManagementWorkspaceLoggerEventhubResource{}

func (r ApiManagementWorkspaceLoggerEventhubResource) ResourceType() string {
	return "azurerm_api_management_workspace_logger_eventhub"
}

func (r ApiManagementWorkspaceLoggerEventhubResource) ModelObject() interface{} {
	return &ApiManagementWorkspaceLoggerEventhubModel{}
}

func (r ApiManagementWorkspaceLoggerEventhubResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return logger.ValidateWorkspaceLoggerID
}

func (r ApiManagementWorkspaceLoggerEventhubResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementChildName(),

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&logger.WorkspaceId{}),

		"eventhub": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ValidateEventHubName(),
					},

					"connection_string": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
						ExactlyOneOf: []string{
							"eventhub.0.connection_string",
							"eventhub.0.endpoint_uri",
						},
						ConflictsWith: []string{
							"eventhub.0.user_assigned_identity_client_id",
						},
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"endpoint_uri": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{
							"eventhub.0.connection_string",
							"eventhub.0.endpoint_uri",
						},
					},

					"user_assigned_identity_client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
						ConflictsWith: []string{
							"eventhub.0.connection_string",
						},
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

func (r ApiManagementWorkspaceLoggerEventhubResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiManagementWorkspaceLoggerEventhubResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.LoggerClient_v2024_05_01

			var model ApiManagementWorkspaceLoggerEventhubModel
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
					LoggerType:  logger.LoggerTypeAzureEventHub,
					IsBuffered:  pointer.To(model.BufferingEnabled),
					Credentials: expandApiManagementWorkspaceLoggerEventhub(model.Eventhub),
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

func (r ApiManagementWorkspaceLoggerEventhubResource) Read() sdk.ResourceFunc {
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

func (r ApiManagementWorkspaceLoggerEventhubResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.LoggerClient_v2024_05_01

			var model ApiManagementWorkspaceLoggerEventhubModel
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

			if metadata.ResourceData.HasChange("eventhub") {
				payload.Properties.Credentials = expandApiManagementWorkspaceLoggerEventhub(model.Eventhub)
			}

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

func (r ApiManagementWorkspaceLoggerEventhubResource) Delete() sdk.ResourceFunc {
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

func expandApiManagementWorkspaceLoggerEventhub(inputs []EventhubModel) *map[string]string {
	if len(inputs) == 0 {
		return nil
	}

	input := inputs[0]
	result := make(map[string]string)
	result["name"] = input.Name
	if len(input.ConnectionString) > 0 {
		result["connectionString"] = input.ConnectionString
	} else if len(input.EndpointUri) > 0 {
		result["endpointAddress"] = input.EndpointUri

		// This field is required by the API and only accepts either a valid UUID or `SystemAssigned` as a value, so we default this to `SystemAssigned` in the creation/update if the field is omitted/removed
		result["identityClientId"] = "SystemAssigned"
		if input.UserAssignedIdentityClientId != "" {
			result["identityClientId"] = input.UserAssignedIdentityClientId
		}
	}

	return &result
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
