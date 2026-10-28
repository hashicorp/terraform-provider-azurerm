// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_query_pack_query

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2019-09-01/querypackqueries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.QueryPackQueriesClient

			id, err := querypackqueries.ParseQueryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model LogAnalyticsQueryPackQueryModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := &querypackqueries.LogAnalyticsQueryPackQuery{
				Properties: &querypackqueries.LogAnalyticsQueryPackQueryProperties{
					Body:        model.Body,
					DisplayName: model.DisplayName,
					Related:     &querypackqueries.LogAnalyticsQueryPackQueryPropertiesRelated{},
				},
			}

			if metadata.ResourceData.HasChange("description") {
				parameters.Properties.Description = &model.Description
			}

			if metadata.ResourceData.HasChange("additional_settings_json") {
				var additionalSettingsJson any
				if err := json.Unmarshal([]byte(model.AdditionalSettingsJson), &additionalSettingsJson); err != nil {
					return fmt.Errorf("parsing JSON: %+v", err)
				}
				parameters.Properties.Properties = &additionalSettingsJson
			}

			if metadata.ResourceData.HasChange("categories") {
				parameters.Properties.Related.Categories = &model.Categories
			}

			if metadata.ResourceData.HasChange("resource_types") {
				parameters.Properties.Related.ResourceTypes = &model.ResourceTypes
			}

			if metadata.ResourceData.HasChange("solutions") {
				parameters.Properties.Related.Solutions = &model.Solutions
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Properties.Tags = expandLogAnalyticsQueryPackQueryTags(model.Tags)
			}

			if _, err := client.QueriesUpdate(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
