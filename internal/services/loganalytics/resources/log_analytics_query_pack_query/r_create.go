// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_query_pack_query

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2019-09-01/querypackqueries"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LogAnalyticsQueryPackQueryModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.LogAnalytics.QueryPackQueriesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			queryPackId, _ := querypackqueries.ParseQueryPackID(model.QueryPackId)

			if model.Name == "" {
				uuid, err := uuid.GenerateUUID()
				if err != nil {
					return fmt.Errorf("generating UUID for Log Analytics Query Pack Query: %+v", err)
				}

				model.Name = uuid
			}

			id := querypackqueries.NewQueryID(subscriptionId, queryPackId.ResourceGroupName, queryPackId.QueryPackName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.QueriesGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			parameters := &querypackqueries.LogAnalyticsQueryPackQuery{
				Properties: &querypackqueries.LogAnalyticsQueryPackQueryProperties{
					Body:        model.Body,
					DisplayName: model.DisplayName,
					Related:     &querypackqueries.LogAnalyticsQueryPackQueryPropertiesRelated{},
				},
			}

			if model.Description != "" {
				parameters.Properties.Description = &model.Description
			}

			if model.AdditionalSettingsJson != "" {
				var additionalSettingsJson any
				if err := json.Unmarshal([]byte(model.AdditionalSettingsJson), &additionalSettingsJson); err != nil {
					return fmt.Errorf("parsing JSON: %+v", err)
				}
				parameters.Properties.Properties = &additionalSettingsJson
			}

			if model.Categories != nil {
				parameters.Properties.Related.Categories = &model.Categories
			}

			if model.ResourceTypes != nil {
				parameters.Properties.Related.ResourceTypes = &model.ResourceTypes
			}

			if model.Solutions != nil {
				parameters.Properties.Related.Solutions = &model.Solutions
			}

			if model.Tags != nil {
				parameters.Properties.Tags = expandLogAnalyticsQueryPackQueryTags(model.Tags)
			}

			if _, err := client.QueriesPut(ctx, id, *parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandLogAnalyticsQueryPackQueryTags(input map[string]string) *map[string][]string {
	if input == nil {
		return nil
	}

	output := make(map[string][]string)
	for k, v := range input {
		output[k] = strings.Split(v, ",")
	}

	return &output
}
