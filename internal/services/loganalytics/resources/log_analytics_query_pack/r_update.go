// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_query_pack

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2019-09-01/querypacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.QueryPacksClient

			id, err := querypacks.ParseQueryPackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model LogAnalyticsQueryPackModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if resp, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				// update check logic once the issue https://github.com/Azure/azure-rest-api-specs/issues/19603 is fixed
				if !response.WasStatusCode(resp.HttpResponse, http.StatusCreated) {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
