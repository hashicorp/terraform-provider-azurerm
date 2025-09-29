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

type ApiManagementWorkspaceLoggerModel struct {
	Name                     string                     `tfschema:"name"`
	ApiManagementWorkspaceId string                     `tfschema:"api_management_workspace_id"`
	ApplicationInsights      []ApplicationInsightsModel `tfschema:"application_insights"`
	EventHub                 []EventHubModel            `tfschema:"eventhub"`
	BufferingEnabled         bool                       `tfschema:"buffering_enabled"`
	Description              string                     `tfschema:"description"`
	ResourceId               string                     `tfschema:"resource_id"`
}

type ApplicationInsightsModel struct {
	InstrumentationKey string `tfschema:"instrumentation_key"`
	ConnectionString   string `tfschema:"connection_string"`
}

type EventHubModel struct {
	Name                         string `tfschema:"name"`
	ConnectionString             string `tfschema:"connection_string"`
	EndpointUri                  string `tfschema:"endpoint_uri"`
	UserAssignedIdentityClientId string `tfschema:"user_assigned_identity_client_id"`
}

type ApiManagementWorkspaceLoggerResource struct{}

var _ sdk.ResourceWithUpdate = ApiManagementWorkspaceLoggerResource{}

func (r ApiManagementWorkspaceLoggerResource) ResourceType() string {
	return "azurerm_api_management_workspace_logger"
}

func (r ApiManagementWorkspaceLoggerResource) ModelObject() interface{} {
	return &ApiManagementWorkspaceLoggerModel{}
}

func (r ApiManagementWorkspaceLoggerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return logger.ValidateWorkspaceLoggerID
}

func (r ApiManagementWorkspaceLoggerResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementChildName(),

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&logger.WorkspaceId{}),

		"application_insights": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ForceNew:     true,
			ExactlyOneOf: []string{"application_insights", "eventhub"},
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

		"eventhub": {
			Type:         pluginsdk.TypeList,
			MaxItems:     1,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"application_insights", "eventhub"},
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

		"resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: azure.ValidateResourceID,
		},
	}
}

func (r ApiManagementWorkspaceLoggerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiManagementWorkspaceLoggerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.LoggerClient_v2024_05_01

			var model ApiManagementWorkspaceLoggerModel
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
					IsBuffered: pointer.To(model.BufferingEnabled),
				},
			}

			if len(model.ApplicationInsights) > 0 {
				parameters.Properties.LoggerType = logger.LoggerTypeApplicationInsights
				parameters.Properties.Credentials = expandApiManagementWorkspaceLoggerApplicationInsights(model.ApplicationInsights)
			}

			if len(model.EventHub) > 0 {
				parameters.Properties.LoggerType = logger.LoggerTypeAzureEventHub
				parameters.Properties.Credentials = expandApiManagementWorkspaceLoggerEventHub(model.EventHub)
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

func (r ApiManagementWorkspaceLoggerResource) Read() sdk.ResourceFunc {
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

			model := ApiManagementWorkspaceLoggerModel{
				Name:                     id.LoggerId,
				ApiManagementWorkspaceId: logger.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if respModel := resp.Model; respModel != nil {
				if props := respModel.Properties; props != nil {
					model.BufferingEnabled = pointer.From(props.IsBuffered)
					model.Description = pointer.From(props.Description)
					model.ResourceId = pointer.From(props.ResourceId)

					if credentials := props.Credentials; credentials != nil {
						var config ApiManagementWorkspaceLoggerModel
						if err := metadata.Decode(&config); err != nil {
							return fmt.Errorf("decoding: %+v", err)
						}

						switch props.LoggerType {
						case logger.LoggerTypeApplicationInsights:
							// The `application_insights.0.instrumentation_key` and `application_insights.0.connection_string` returned by the Azure API is intentionally masked
							// (e.g. "{{Logger-Credentials--<hash>}}") and does not match the original value provided during creation/update.
							// This is by design to prevent exposing sensitive credentials in API responses.
							// Therefore, the `application_insights` is sourced from the state.
							model.ApplicationInsights = config.ApplicationInsights
						case logger.LoggerTypeAzureEventHub:
							model.EventHub = flattenApiManagementWorkspaceLoggerEventHub(config, props)
						}
					}
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r ApiManagementWorkspaceLoggerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.LoggerClient_v2024_05_01

			var model ApiManagementWorkspaceLoggerModel
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
			if len(model.ApplicationInsights) > 0 {
				// `application_insights.0,instrumentation_key` and `application_insights.0.connection_string` read from config since the API masks `instrumentationKey ` and `connectionString`.
				payload.Properties.Credentials = expandApiManagementWorkspaceLoggerApplicationInsights(model.ApplicationInsights)
			}

			if len(model.EventHub) > 0 {
				if metadata.ResourceData.HasChange("eventhub") {
					payload.Properties.Credentials = expandApiManagementWorkspaceLoggerEventHub(model.EventHub)
				}
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

func (r ApiManagementWorkspaceLoggerResource) Delete() sdk.ResourceFunc {
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

func expandApiManagementWorkspaceLoggerEventHub(inputs []EventHubModel) *map[string]string {
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

func flattenApiManagementWorkspaceLoggerEventHub(model ApiManagementWorkspaceLoggerModel, input *logger.LoggerContractProperties) []EventHubModel {
	outputList := make([]EventHubModel, 0)
	if input == nil || input.Credentials == nil {
		return outputList
	}

	output := EventHubModel{}

	if name, ok := (*input.Credentials)["name"]; ok {
		output.Name = name
	}

	if endpoint, ok := (*input.Credentials)["endpointAddress"]; ok {
		output.EndpointUri = endpoint
	}

	if eventhub := model.EventHub; len(eventhub) > 0 {
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
