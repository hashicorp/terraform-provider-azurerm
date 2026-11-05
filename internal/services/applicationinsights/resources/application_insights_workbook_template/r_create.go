// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_workbook_template

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-11-20/workbooktemplatesapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ApplicationInsightsWorkbookTemplateModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppInsights.WorkbookTemplateClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := workbooktemplatesapis.NewWorkbookTemplateID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.WorkbookTemplatesGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			var templateDataValue any
			if err := json.Unmarshal([]byte(model.TemplateData), &templateDataValue); err != nil {
				return err
			}

			properties := &workbooktemplatesapis.WorkbookTemplate{
				Location: location.Normalize(model.Location),
				Properties: &workbooktemplatesapis.WorkbookTemplateProperties{
					Priority:     &model.Priority,
					TemplateData: templateDataValue,
				},

				Tags: &model.Tags,
			}

			if model.Author != "" {
				properties.Properties.Author = &model.Author
			}

			if model.Localized != "" {
				var localizedValue map[string][]workbooktemplatesapis.WorkbookTemplateLocalizedGallery
				if err := json.Unmarshal([]byte(model.Localized), &localizedValue); err != nil {
					return err
				}

				properties.Properties.Localized = &localizedValue
			}

			galleriesValue := expandWorkbookTemplateGalleryModel(model.Galleries)
			if galleriesValue != nil {
				properties.Properties.Galleries = *galleriesValue
			}

			if _, err := client.WorkbookTemplatesCreateOrUpdate(ctx, id, *properties); err != nil {
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

func expandWorkbookTemplateGalleryModel(inputList []WorkbookTemplateGalleryModel) *[]workbooktemplatesapis.WorkbookTemplateGallery {
	outputList := make([]workbooktemplatesapis.WorkbookTemplateGallery, 0, len(inputList))
	for _, input := range inputList {
		output := workbooktemplatesapis.WorkbookTemplateGallery{
			Category:     pointer.To(input.Category),
			Name:         pointer.To(input.Name),
			Order:        pointer.To(input.Order),
			ResourceType: pointer.To(input.ResourceType),
			Type:         pointer.To(input.Type),
		}

		outputList = append(outputList, output)
	}

	return &outputList
}
