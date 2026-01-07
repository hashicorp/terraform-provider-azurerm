// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiversionset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apiversionsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name api_management_workspace_api_version_set -service-package-name apimanagement -properties "name" -compare-values "resource_group_name:api_management_workspace_id,service_name:api_management_workspace_id,workspace_id:api_management_workspace_id" -known-values "subscription_id:data.Subscriptions.Primary" -test-name basic

type ApiManagementWorkspaceApiVersionSetModel struct {
	Name                     string `tfschema:"name"`
	ApiManagementWorkspaceId string `tfschema:"api_management_workspace_id"`
	DisplayName              string `tfschema:"display_name"`
	VersioningScheme         string `tfschema:"versioning_scheme"`
	Description              string `tfschema:"description"`
	VersionHeaderName        string `tfschema:"version_header_name"`
	VersionQueryName         string `tfschema:"version_query_name"`
}

type ApiManagementWorkspaceApiVersionSetResource struct{}

var (
	_ sdk.ResourceWithUpdate        = ApiManagementWorkspaceApiVersionSetResource{}
	_ sdk.ResourceWithCustomizeDiff = ApiManagementWorkspaceApiVersionSetResource{}
	_ sdk.ResourceWithIdentity      = ApiManagementWorkspaceApiVersionSetResource{}
)

func (r ApiManagementWorkspaceApiVersionSetResource) ResourceType() string {
	return "azurerm_api_management_workspace_api_version_set"
}

func (r ApiManagementWorkspaceApiVersionSetResource) Identity() resourceids.ResourceId {
	return &apiversionset.WorkspaceApiVersionSetId{}
}

func (r ApiManagementWorkspaceApiVersionSetResource) ModelObject() interface{} {
	return &ApiManagementWorkspaceApiVersionSetModel{}
}

func (r ApiManagementWorkspaceApiVersionSetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return apiversionset.ValidateWorkspaceApiVersionSetID
}

func (r ApiManagementWorkspaceApiVersionSetResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9]([a-zA-Z0-9-_]{0,78}[a-zA-Z0-9])?$"),
				"The `name` can only contain alphanumeric characters, underscores and dashes up to 80 characters in length.",
			),
		},

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&workspace.WorkspaceId{}),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"versioning_scheme": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(apiversionset.PossibleValuesForVersioningScheme(), false),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"version_header_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"version_query_name"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},

		"version_query_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"version_header_name"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},
	}
}

func (r ApiManagementWorkspaceApiVersionSetResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiManagementWorkspaceApiVersionSetResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiVersionSetClient_v2024_05_01

			var model ApiManagementWorkspaceApiVersionSetModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := workspace.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			id := apiversionset.NewWorkspaceApiVersionSetID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.ServiceName, workspaceId.WorkspaceId, model.Name)
			existing, err := client.WorkspaceApiVersionSetGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := apiversionset.ApiVersionSetContract{
				Properties: &apiversionset.ApiVersionSetContractProperties{
					DisplayName:      model.DisplayName,
					VersioningScheme: apiversionset.VersioningScheme(model.VersioningScheme),
				},
			}

			if model.Description != "" {
				parameters.Properties.Description = pointer.To(model.Description)
			}

			if model.VersionHeaderName != "" {
				parameters.Properties.VersionHeaderName = pointer.To(model.VersionHeaderName)
			}

			if model.VersionQueryName != "" {
				parameters.Properties.VersionQueryName = pointer.To(model.VersionQueryName)
			}

			_, err = client.WorkspaceApiVersionSetCreateOrUpdate(ctx, id, parameters, apiversionset.DefaultWorkspaceApiVersionSetCreateOrUpdateOperationOptions())
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceApiVersionSetResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiVersionSetClient_v2024_05_01

			var model ApiManagementWorkspaceApiVersionSetModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := apiversionset.ParseWorkspaceApiVersionSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkspaceApiVersionSetGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			properties := resp.Model.Properties
			if metadata.ResourceData.HasChange("display_name") {
				properties.DisplayName = model.DisplayName
			}

			if metadata.ResourceData.HasChange("versioning_scheme") {
				properties.VersioningScheme = apiversionset.VersioningScheme(model.VersioningScheme)
			}

			if metadata.ResourceData.HasChange("description") {
				properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("version_header_name") {
				properties.VersionHeaderName = nil
				if model.VersionHeaderName != "" {
					properties.VersionHeaderName = pointer.To(model.VersionHeaderName)
				}
			}

			if metadata.ResourceData.HasChange("version_query_name") {
				properties.VersionQueryName = nil
				if model.VersionQueryName != "" {
					properties.VersionQueryName = pointer.To(model.VersionQueryName)
				}
			}

			resp.Model.Properties = properties
			_, err = client.WorkspaceApiVersionSetCreateOrUpdate(ctx, *id, *resp.Model, apiversionset.DefaultWorkspaceApiVersionSetCreateOrUpdateOperationOptions())
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceApiVersionSetResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiVersionSetClient_v2024_05_01

			id, err := apiversionset.ParseWorkspaceApiVersionSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkspaceApiVersionSetGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := ApiManagementWorkspaceApiVersionSetModel{
				Name:                     id.VersionSetId,
				ApiManagementWorkspaceId: workspace.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if resp.Model != nil {
				if properties := resp.Model.Properties; properties != nil {
					model.DisplayName = properties.DisplayName
					model.VersioningScheme = string(properties.VersioningScheme)
					model.Description = pointer.From(properties.Description)
					model.VersionHeaderName = pointer.From(properties.VersionHeaderName)
					model.VersionQueryName = pointer.From(properties.VersionQueryName)
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&model)
		},
	}
}

func (r ApiManagementWorkspaceApiVersionSetResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.ApiVersionSetsClient_v2024_05_01

			id, err := apiversionsets.ParseWorkspaceApiVersionSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			_, err = client.WorkspaceApiVersionSetDelete(ctx, *id, apiversionsets.DefaultWorkspaceApiVersionSetDeleteOperationOptions())
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceApiVersionSetResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			var headerSet, querySet bool
			if v, ok := rd.GetOk("version_header_name"); ok {
				headerSet = v.(string) != ""
			}
			if v, ok := rd.GetOk("version_query_name"); ok {
				querySet = v.(string) != ""
			}

			versioningScheme := apiversionset.VersioningScheme(rd.Get("versioning_scheme").(string))
			switch schema := versioningScheme; schema {
			case apiversionset.VersioningSchemeHeader:
				if !headerSet {
					return errors.New("`version_header_name` must be set if `versioning_scheme` is `Header`")
				}

			case apiversionset.VersioningSchemeQuery:
				if headerSet {
					return errors.New("`version_header_name` can not be set if `versioning_scheme` is `Query`")
				}
				if !querySet {
					return errors.New("`version_query_name` must be set if `versioning_scheme` is `Query`")
				}

			case apiversionset.VersioningSchemeSegment:
				if headerSet {
					return errors.New("`version_header_name` can not be set if `versioning_scheme` is `Segment`")
				}
				if querySet {
					return errors.New("`version_query_name` can not be set if `versioning_scheme` is `Segment`")
				}
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
