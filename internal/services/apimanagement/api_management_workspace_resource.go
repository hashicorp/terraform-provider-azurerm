// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementWorkspaceModel struct {
	Name            string `tfschema:"name"`
	ApiManagementId string `tfschema:"api_management_id"`
	Description     string `tfschema:"description"`
	DisplayName     string `tfschema:"display_name"`
}

type ApiManagementWorkspaceResource struct{}

var _ sdk.ResourceWithUpdate = ApiManagementWorkspaceResource{}

func (r ApiManagementWorkspaceResource) ResourceType() string {
	return "azurerm_api_management_workspace"
}

func (r ApiManagementWorkspaceResource) ModelObject() interface{} {
	return &ApiManagementWorkspaceModel{}
}

func (r ApiManagementWorkspaceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return workspace.ValidateWorkspaceID
}

func (r ApiManagementWorkspaceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,78}[a-zA-Z0-9])?$"),
				"The `name` must be between 1 and 80 characters long and can only include letters, numbers, and hyphens sign when preceded and followed by number or a letter.",
			),
		},

		"api_management_id": commonschema.ResourceIDReferenceRequiredForceNew(&apimanagementservice.ServiceId{}),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ApiManagementWorkspaceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiManagementWorkspaceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.WorkspaceClient

			var model ApiManagementWorkspaceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			serviceId, err := apimanagementservice.ParseServiceID(model.ApiManagementId)
			if err != nil {
				return err
			}

			id := workspace.NewWorkspaceID(serviceId.SubscriptionId, serviceId.ResourceGroupName, serviceId.ServiceName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := workspace.WorkspaceContract{
				Properties: &workspace.WorkspaceContractProperties{
					DisplayName: model.DisplayName,
				},
			}

			if model.Description != "" {
				properties.Properties.Description = pointer.To(model.Description)
			}

			if _, err := client.CreateOrUpdate(ctx, id, properties, workspace.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ApiManagementWorkspaceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.WorkspaceClient

			id, err := workspace.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ApiManagementWorkspaceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			properties := resp.Model
			if metadata.ResourceData.HasChange("description") {
				if model.Description != "" {
					properties.Properties.Description = pointer.To(model.Description)
				} else {
					properties.Properties.Description = nil
				}
			}

			if metadata.ResourceData.HasChange("display_name") {
				properties.Properties.DisplayName = model.DisplayName
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties, workspace.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

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
					return metadata.MarkAsGone(*id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ApiManagementWorkspaceModel{
				Name:            id.WorkspaceId,
				ApiManagementId: apimanagementservice.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName).ID(),
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					state.Description = pointer.From(properties.Description)
					state.DisplayName = properties.DisplayName
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ApiManagementWorkspaceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.WorkspaceClient

			id, err := workspace.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id, workspace.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
