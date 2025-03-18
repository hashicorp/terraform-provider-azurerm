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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementWorkspaceModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	ServiceName       string `tfschema:"service_name"`
	WorkspaceName     string `tfschema:"workspace_name"`
}

type ApiManagementWorkspaceResource struct{}

var _ sdk.Resource = ApiManagementWorkspaceResource{}

func (r ApiManagementWorkspaceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),
		"service_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The name of the API Management Service in which this Workspace should be created.",
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
		"workspace_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The display name of the API Management Workspace.",
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

			var model ApiManagementWorkspaceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			newId := workspace.NewWorkspaceID(subscriptionId, model.ResourceGroupName, model.ServiceName, model.Name)

			existing, err := client.Get(ctx, newId)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", newId, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), newId)
			}

			properties := workspace.WorkspaceContract{
				Properties: &workspace.WorkspaceContractProperties{
					DisplayName: model.WorkspaceName,
				},
			}

			if _, err := client.CreateOrUpdate(ctx, newId, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", newId, err)
			}

			metadata.SetID(newId)
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

			state := ApiManagementWorkspaceModel{
				Name:              id.WorkspaceId,
				ServiceName:       id.ServiceName,
				ResourceGroupName: id.ResourceGroup,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.WorkspaceName = props.DisplayName
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

			if _, err = client.Delete(ctx, *id); err != nil {
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
					DisplayName: model.WorkspaceName,
				},
			}

			if _, err := client.CreateOrUpdate(ctx, *id, properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
