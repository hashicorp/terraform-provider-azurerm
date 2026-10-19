// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_workbook_template

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-11-20/workbooktemplatesapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppInsights.WorkbookTemplateClient

			id, err := workbooktemplatesapis.ParseWorkbookTemplateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ApplicationInsightsWorkbookTemplateModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.WorkbookTemplatesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("author") {
				properties.Properties.Author = &model.Author
			}

			if metadata.ResourceData.HasChange("galleries") {
				galleriesValue := expandWorkbookTemplateGalleryModel(model.Galleries)
				if galleriesValue != nil {
					properties.Properties.Galleries = *galleriesValue
				}
			}

			if metadata.ResourceData.HasChange("priority") {
				properties.Properties.Priority = &model.Priority
			}

			if metadata.ResourceData.HasChange("template_data") {
				var templateDataValue any
				if err := json.Unmarshal([]byte(model.TemplateData), &templateDataValue); err != nil {
					return err
				}

				properties.Properties.TemplateData = templateDataValue
			}

			if metadata.ResourceData.HasChange("localized") {
				var localizedValue map[string][]workbooktemplatesapis.WorkbookTemplateLocalizedGallery
				if err := json.Unmarshal([]byte(model.Localized), &localizedValue); err != nil {
					return err
				}

				properties.Properties.Localized = &localizedValue
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if _, err := client.WorkbookTemplatesCreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
