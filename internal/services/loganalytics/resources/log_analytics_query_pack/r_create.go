// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_query_pack

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2019-09-01/querypacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LogAnalyticsQueryPackModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.LogAnalytics.QueryPacksClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := querypacks.NewQueryPackID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := &querypacks.LogAnalyticsQueryPack{
				Location:   location.Normalize(model.Location),
				Properties: querypacks.LogAnalyticsQueryPackProperties{},
				Tags:       &model.Tags,
			}

			if resp, err := client.CreateOrUpdate(ctx, id, *properties); err != nil {
				// update check logic once the issue https://github.com/Azure/azure-rest-api-specs/issues/19603 is fixed
				if !response.WasStatusCode(resp.HttpResponse, http.StatusCreated) {
					return fmt.Errorf("creating %s: %+v", id, err)
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}
