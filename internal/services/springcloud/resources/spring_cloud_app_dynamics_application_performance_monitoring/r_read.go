// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_app_dynamics_application_performance_monitoring

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
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
)

func (s Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := springcloudClient.From(metadata.Client).AppPlatformClient

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

			var model SpringCloudAppDynamicsApplicationPerformanceMonitoringModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := SpringCloudAppDynamicsApplicationPerformanceMonitoringModel{
				Name:                  id.ApmName,
				SpringCloudServiceId:  springId.ID(),
				GloballyEnabled:       globallyEnabled,
				AgentAccountName:      model.AgentAccountName,
				AgentAccountAccessKey: model.AgentAccountAccessKey,
			}

			if props := resp.Model.Properties; props != nil {
				if props.Type != "AppDynamics" {
					return fmt.Errorf("retrieving %s: expected type AppDynamics, got %s", *id, props.Type)
				}
				if props.Properties != nil {
					if value, ok := (*props.Properties)["agent_application_name"]; ok {
						state.AgentApplicationName = value
					}
					if value, ok := (*props.Properties)["agent_tier_name"]; ok {
						state.AgentTierName = value
					}
					if value, ok := (*props.Properties)["agent_node_name"]; ok {
						state.AgentNodeName = value
					}
					if value, ok := (*props.Properties)["agent_unique_host_id"]; ok {
						state.AgentUniqueHostId = value
					}
					if value, ok := (*props.Properties)["controller_host_name"]; ok {
						state.ControllerHostName = value
					}
					if value, ok := (*props.Properties)["controller_ssl_enabled"]; ok {
						state.ControllerSslEnabled = value == "true"
					}
					if value, ok := (*props.Properties)["controller_port"]; ok {
						if v, err := strconv.ParseInt(value, 10, 32); err == nil {
							state.ControllerPort = v
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}
