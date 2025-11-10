// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/group"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementWorkspaceGroupModel struct {
	Name                     string `tfschema:"name"`
	ApiManagementWorkspaceId string `tfschema:"api_management_workspace_id"`
	DisplayName              string `tfschema:"display_name"`
	Description              string `tfschema:"description"`
	ExternalId               string `tfschema:"external_id"`
	Type                     string `tfschema:"type"`
}

type ApiManagementWorkspaceGroupResource struct{}

var _ sdk.ResourceWithUpdate = ApiManagementWorkspaceGroupResource{}

var _ sdk.ResourceWithCustomizeDiff = ApiManagementWorkspaceGroupResource{}

func (r ApiManagementWorkspaceGroupResource) ResourceType() string {
	return "azurerm_api_management_workspace_group"
}

func (r ApiManagementWorkspaceGroupResource) ModelObject() interface{} {
	return &ApiManagementWorkspaceGroupModel{}
}

func (r ApiManagementWorkspaceGroupResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return group.ValidateWorkspaceGroupID
}

func (r ApiManagementWorkspaceGroupResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": schemaz.SchemaApiManagementChildName(),

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&group.WorkspaceId{}),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"external_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(group.GroupTypeCustom),
			// `group.GroupTypeSystem` is excluded here because system groups are predefined by API Management
			// and cannot be created through the API.
			ValidateFunc: validation.StringInSlice([]string{
				string(group.GroupTypeCustom),
				string(group.GroupTypeExternal),
			}, false),
		},
	}
}

func (r ApiManagementWorkspaceGroupResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiManagementWorkspaceGroupResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.GroupClient_v2024_05_01

			var model ApiManagementWorkspaceGroupModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := group.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			id := group.NewWorkspaceGroupID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.ServiceName, workspaceId.WorkspaceId, model.Name)

			existing, err := client.WorkspaceGroupGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := group.GroupCreateParameters{
				Properties: &group.GroupCreateParametersProperties{
					DisplayName: model.DisplayName,
					Type:        pointer.ToEnum[group.GroupType](model.Type),
				},
			}

			if model.Description != "" {
				properties.Properties.Description = pointer.To(model.Description)
			}

			if model.ExternalId != "" {
				properties.Properties.ExternalId = pointer.To(model.ExternalId)
			}

			if _, err := client.WorkspaceGroupCreateOrUpdate(ctx, id, properties, group.WorkspaceGroupCreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ApiManagementWorkspaceGroupResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.GroupClient_v2024_05_01

			id, err := group.ParseWorkspaceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkspaceGroupGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ApiManagementWorkspaceGroupModel{
				Name:                     id.GroupId,
				ApiManagementWorkspaceId: group.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.DisplayName = props.DisplayName
					state.Description = pointer.From(props.Description)
					state.ExternalId = pointer.From(props.ExternalId)
					state.Type = pointer.FromEnum(props.Type)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ApiManagementWorkspaceGroupResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.GroupClient_v2024_05_01

			var model ApiManagementWorkspaceGroupModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := group.ParseWorkspaceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkspaceGroupGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			parameters := group.GroupCreateParameters{
				Properties: &group.GroupCreateParametersProperties{
					Description: resp.Model.Properties.Description,
					DisplayName: resp.Model.Properties.DisplayName,
					ExternalId:  resp.Model.Properties.ExternalId,
					Type:        resp.Model.Properties.Type,
				},
			}

			if metadata.ResourceData.HasChange("display_name") {
				parameters.Properties.DisplayName = model.DisplayName
			}

			if metadata.ResourceData.HasChange("description") {
				parameters.Properties.Description = pointer.To(model.Description)
			}

			if _, err := client.WorkspaceGroupCreateOrUpdate(ctx, *id, parameters, group.WorkspaceGroupCreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceGroupResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.GroupClient_v2024_05_01

			id, err := group.ParseWorkspaceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.WorkspaceGroupDelete(ctx, *id, group.DefaultWorkspaceGroupDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceGroupResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff
			groupType := rd.Get("type")
			isExternalIdUnset := rd.GetRawConfig().AsValueMap()["external_id"].IsNull()

			if groupType == string(group.GroupTypeExternal) && isExternalIdUnset {
				return errors.New("`external_id` must be specified when `type` is set to `external`")
			}
			if groupType == string(group.GroupTypeCustom) && !isExternalIdUnset {
				return errors.New("`type` must be set to `external` when `external_id` is specified`")
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
