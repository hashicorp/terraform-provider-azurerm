// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementWorkspaceModel struct {
	Name            string `tfschema:"name"`
	ApiManagementId string `tfschema:"api_management_id"`
	DisplayName     string `tfschema:"display_name"`
	Description     string `tfschema:"description"`
}

type ApiManagementWorkspaceResource struct{}

var _ sdk.Resource = ApiManagementWorkspaceResource{}

func (r ApiManagementWorkspaceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_management_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apimanagementservice.ValidateServiceID,
			Description:  "The ID of the API Management Service in which this Workspace should be created.",
		},
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9-]{1,80}$`),
				"Workspace name must be 1 - 80 characters long, contain only letters, numbers and hyphens.",
			),
			Description: "The name of the API Management Workspace.",
		},
		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The display name of the API Management Workspace.",
		},
		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The description of the API Management Workspace.",
		},
	}
}

func (r ApiManagementWorkspaceResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiManagementWorkspaceResource) ModelObject() interface{} {
	return &ApiManagementWorkspaceModel{}
}

func (r ApiManagementWorkspaceResource) ResourceType() string {
	return "azurerm_api_management_workspace"
}

func (r ApiManagementWorkspaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return workspace.ValidateWorkspaceID
}

func (r ApiManagementWorkspaceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 45 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.WorkspaceClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ApiManagementWorkspaceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			serviceId, err := apimanagementservice.ParseServiceID(model.ApiManagementId)
			if err != nil {
				return err
			}

			id := workspace.NewWorkspaceID(subscriptionId, serviceId.ResourceGroupName, serviceId.ServiceName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := workspace.WorkspaceContract{
				Properties: &workspace.WorkspaceContractProperties{
					DisplayName: model.DisplayName,
					Description: &model.Description,
				},
			}

			if _, err := client.CreateOrUpdate(ctx, id, properties, workspace.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ApiManagementWorkspaceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.WorkspaceClient

			id, err := workspace.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			apiManagementId := apimanagementservice.NewServiceID(subscriptionId, id.ResourceGroupName, id.ServiceName)

			state := ApiManagementWorkspaceModel{
				Name:            id.WorkspaceId,
				ApiManagementId: apiManagementId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.DisplayName = props.DisplayName
					state.Description = pointer.ToString(props.Description)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ApiManagementWorkspaceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 45 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.WorkspaceClient

			id, err := workspace.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.Delete(ctx, *id, workspace.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 45 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.WorkspaceClient

			id, err := workspace.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ApiManagementWorkspaceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			properties := workspace.WorkspaceContract{
				Properties: &workspace.WorkspaceContractProperties{
					DisplayName: model.DisplayName,
					Description: &model.Description,
				},
			}

			if _, err := client.CreateOrUpdate(ctx, *id, properties, workspace.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
