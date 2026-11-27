// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_new_relic_application_performance_monitoring

import (
	"context"
	"fmt"
	"strconv"
	"strings"
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

			var model SpringCloudNewRelicApplicationPerformanceMonitoringModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := SpringCloudNewRelicApplicationPerformanceMonitoringModel{
				Name:                 id.ApmName,
				SpringCloudServiceId: springId.ID(),
				GloballyEnabled:      globallyEnabled,
				LicenseKey:           model.LicenseKey,
			}

			if props := resp.Model.Properties; props != nil {
				if props.Type != "NewRelic" {
					return fmt.Errorf("retrieving %s: type was %s, expected NewRelic", *id, props.Type)
				}
				if props.Properties != nil {
					if value, ok := (*props.Properties)["app_name"]; ok {
						state.AppName = value
					}
					if value, ok := (*props.Properties)["agent_enabled"]; ok {
						state.AgentEnabled = value == "true"
					}
					if value, ok := (*props.Properties)["appserver_port"]; ok {
						if v, err := strconv.ParseInt(value, 10, 0); err == nil {
							state.AppServerPort = v
						}
					}
					if value, ok := (*props.Properties)["audit_mode"]; ok {
						state.AuditModeEnabled = value == "true"
					}
					if value, ok := (*props.Properties)["enable_auto_app_naming"]; ok {
						state.AutoAppNamingEnabled = value == "true"
					}
					if value, ok := (*props.Properties)["enable_auto_transaction_naming"]; ok {
						state.AutoTransactionNamingEnabled = value == "true"
					}
					if value, ok := (*props.Properties)["enable_custom_tracing"]; ok {
						state.CustomTracingEnabled = value == "true"
					}
					if value, ok := (*props.Properties)["labels"]; ok {
						state.Labels = flattenNewRelicLabels(value)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenNewRelicLabels(input string) map[string]string {
	labels := make(map[string]string)
	if input == "" {
		return labels
	}
	for label := range strings.SplitSeq(input, ";") {
		parts := strings.Split(label, ":")
		if len(parts) == 2 {
			labels[parts[0]] = parts[1]
		}
	}
	return labels
}
