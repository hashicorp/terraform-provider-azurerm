// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_application_insights_application_performance_monitoring

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (s Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseApmID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ApmsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			springId := commonids.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroupName, id.SpringName)
			result, err := client.ServicesListGloballyEnabledApms(ctx, springId)
			if err != nil {
				return fmt.Errorf("listing globally enabled apms: %+v", err)
			}
			globallyEnabled := false
			if result.Model != nil && result.Model.Value != nil {
				for _, value := range *result.Model.Value {
					apmId, err := appplatform.ParseApmIDInsensitively(value)
					if err == nil && resourceids.Match(apmId, id) {
						globallyEnabled = true
						break
					}
				}
			}

			var model SpringCloudApplicationInsightsApplicationPerformanceMonitoringModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := SpringCloudApplicationInsightsApplicationPerformanceMonitoringModel{
				Name:                 id.ApmName,
				SpringCloudServiceId: springId.ID(),
				ConnectionString:     model.ConnectionString,
				GloballyEnabled:      globallyEnabled,
			}

			if props := resp.Model.Properties; props != nil {
				if props.Type != "ApplicationInsights" {
					return fmt.Errorf("retrieving %s: type was not ApplicationInsights", *id)
				}
				if props.Properties != nil {
					if value, ok := (*props.Properties)["role_name"]; ok {
						state.RoleName = value
					}
					if value, ok := (*props.Properties)["role_instance"]; ok {
						state.RoleInstance = value
					}
					if value, ok := (*props.Properties)["sampling_requests_per_second"]; ok {
						if samplingRequestsPerSecond, err := strconv.ParseInt(value, 10, 64); err == nil {
							state.SamplingRequestsPerSecond = samplingRequestsPerSecond
						}
					}
					if value, ok := (*props.Properties)["sampling_percentage"]; ok {
						if samplingPercentage, err := strconv.ParseInt(value, 10, 64); err == nil {
							state.SamplingPercentage = samplingPercentage
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}
