// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_workbook

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-04-01/workbooksapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppInsights.WorkbookClient

			id, err := workbooksapis.ParseWorkbookID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ApplicationInsightsWorkbookModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.WorkbooksGet(ctx, *id, workbooksapis.WorkbooksGetOperationOptions{CanFetchContent: pointer.To(true)})
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

			if _, err := client.WorkbooksCreateOrUpdate(ctx, *id, *properties, workbooksapis.WorkbooksCreateOrUpdateOperationOptions{SourceId: &model.SourceId}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
