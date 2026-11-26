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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.QueryPackQueriesClient

			id, err := querypackqueries.ParseQueryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.QueriesGet(ctx, *id)
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

			props := model.Properties
			if props == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			state := LogAnalyticsQueryPackQueryModel{
				Name:        id.QueryName,
				QueryPackId: querypackqueries.NewQueryPackID(id.SubscriptionId, id.ResourceGroupName, id.QueryPackName).ID(),
				Body:        props.Body,
				DisplayName: props.DisplayName,
			}

			if props.Properties != nil {
				propsJson, jsonErr := json.Marshal(props.Properties)
				if jsonErr != nil {
					return fmt.Errorf("parsing JSON for Log Analytics Query Pack Query Properties: %+v", jsonErr)
				}
				state.AdditionalSettingsJson = string(propsJson)
			}

			if props.Description != nil {
				state.Description = *props.Description
			}

			if additionalSettings := props.Related; additionalSettings != nil {
				if additionalSettings.Categories != nil {
					state.Categories = *additionalSettings.Categories
				}

				if additionalSettings.ResourceTypes != nil {
					state.ResourceTypes = *additionalSettings.ResourceTypes
				}

				if additionalSettings.Solutions != nil {
					state.Solutions = *additionalSettings.Solutions
				}
			}

			if tags := props.Tags; tags != nil {
				state.Tags = flattenLogAnalyticsQueryPackQueryTags(*tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenLogAnalyticsQueryPackQueryTags(input map[string][]string) map[string]string {
	if input == nil {
		return map[string]string{}
	}

	results := make(map[string]string)
	for k, v := range input {
		results[k] = strings.Join(v, ",")
	}

	return results
}
