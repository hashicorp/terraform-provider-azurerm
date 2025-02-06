// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	workbooks "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-04-01/workbooksapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ApplicationInsightsWorkbookModel struct {
	Name               string            `tfschema:"name"`
	ResourceGroupName  string            `tfschema:"resource_group_name"`
	Category           string            `tfschema:"category"`
	Description        string            `tfschema:"description"`
	DisplayName        string            `tfschema:"display_name"`
	Location           string            `tfschema:"location"`
	DataJson           string            `tfschema:"data_json"`
	SourceId           string            `tfschema:"source_id"`
	StorageContainerId string            `tfschema:"storage_container_id"`
	Tags               map[string]string `tfschema:"tags"`
}

type ApplicationInsightsWorkbookResource struct{}

var _ sdk.ResourceWithUpdate = ApplicationInsightsWorkbookResource{}

func (r ApplicationInsightsWorkbookResource) ResourceType() string {
	return "azurerm_application_insights_workbook"
}

func (r ApplicationInsightsWorkbookResource) ModelObject() interface{} {
	return &ApplicationInsightsWorkbookModel{}
}

func (r ApplicationInsightsWorkbookResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return workbooks.ValidateWorkbookID
}

func (r ApplicationInsightsWorkbookResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.IsUUID,
				validate.StringDoesNotContainUpperCaseLetter,
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"data_json": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"source_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "azure monitor",
			ValidateFunc: validate.StringDoesNotContainUpperCaseLetter,
		},

		"category": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "workbook",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptionalForceNew(),

		"storage_container_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{
				"identity",
			},
		},

		"tags": {
			Type:         pluginsdk.TypeMap,
			Optional:     true,
			ValidateFunc: validate.WorkbookTags,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r ApplicationInsightsWorkbookResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApplicationInsightsWorkbookResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ApplicationInsightsWorkbookModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppInsights.WorkbookClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := workbooks.NewWorkbookID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.WorkbooksGet(ctx, id, workbooks.WorkbooksGetOperationOptions{CanFetchContent: utils.Bool(true)})
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			kindValue := workbooks.WorkbookSharedTypeKindShared
			properties := &workbooks.Workbook{
				Identity: identityValue,
				Kind:     &kindValue,
				Location: location.Normalize(model.Location),
				Properties: &workbooks.WorkbookProperties{
					Category:       model.Category,
					DisplayName:    model.DisplayName,
					SerializedData: model.DataJson,
					SourceId:       &model.SourceId,
				},

				Tags: &model.Tags,
			}

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			if model.StorageContainerId != "" {
				properties.Properties.StorageUri = &model.StorageContainerId
			}

			if _, err := client.WorkbooksCreateOrUpdate(ctx, id, *properties, workbooks.WorkbooksCreateOrUpdateOperationOptions{SourceId: &model.SourceId}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ApplicationInsightsWorkbookResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppInsights.WorkbookClient

			id, err := workbooks.ParseWorkbookID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ApplicationInsightsWorkbookModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.WorkbooksGet(ctx, *id, workbooks.WorkbooksGetOperationOptions{CanFetchContent: utils.Bool(true)})
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil || properties.Properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("category") {
				properties.Properties.Category = model.Category
			}

			if metadata.ResourceData.HasChange("description") {
				properties.Properties.Description = &model.Description
			}

			if metadata.ResourceData.HasChange("display_name") {
				properties.Properties.DisplayName = model.DisplayName
				if properties.Tags != nil {
					delete(*properties.Tags, "hidden-title")
				}
			}

			if metadata.ResourceData.HasChange("data_json") {
				properties.Properties.SerializedData = model.DataJson
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if _, err := client.WorkbooksCreateOrUpdate(ctx, *id, *properties, workbooks.WorkbooksCreateOrUpdateOperationOptions{SourceId: &model.SourceId}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ApplicationInsightsWorkbookResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppInsights.WorkbookClient

			id, err := workbooks.ParseWorkbookID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkbooksGet(ctx, *id, workbooks.WorkbooksGetOperationOptions{CanFetchContent: utils.Bool(true)})
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := ApplicationInsightsWorkbookModel{
				Name:              id.WorkbookName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			identityValue, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}

			if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if properties := model.Properties; properties != nil {
				state.Category = properties.Category

				if properties.Description != nil {
					state.Description = *properties.Description
				}

				state.DisplayName = properties.DisplayName

				state.DataJson = properties.SerializedData

				if properties.SourceId != nil {
					state.SourceId = *properties.SourceId
				}

				if properties.StorageUri != nil {
					state.StorageContainerId = *properties.StorageUri
				}
			}

			if model.Tags != nil {
				// The backend returns a tags with key `hidden-title` by default. Since it has the same value with `display_name` and will cause inconsistency with user's configuration, remove it as a workaround.
				delete(*model.Tags, "hidden-title")
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ApplicationInsightsWorkbookResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppInsights.WorkbookClient

			id, err := workbooks.ParseWorkbookID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.WorkbooksDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
