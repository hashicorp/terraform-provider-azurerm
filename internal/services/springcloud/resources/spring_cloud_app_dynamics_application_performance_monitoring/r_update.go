// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_app_dynamics_application_performance_monitoring

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (s Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseApmID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SpringCloudAppDynamicsApplicationPerformanceMonitoringModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.ApmsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model.Properties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}
			if properties.Properties == nil {
				properties.Properties = pointer.To(map[string]string{})
			}
			if properties.Secrets == nil {
				properties.Secrets = pointer.To(map[string]string{})
			}

			if metadata.ResourceData.HasChange("agent_application_name") {
				(*properties.Properties)["agent_application_name"] = model.AgentApplicationName
			}

			if metadata.ResourceData.HasChange("agent_tier_name") {
				(*properties.Properties)["agent_tier_name"] = model.AgentTierName
			}

			if metadata.ResourceData.HasChange("agent_node_name") {
				(*properties.Properties)["agent_node_name"] = model.AgentNodeName
			}

			if metadata.ResourceData.HasChange("agent_unique_host_id") {
				(*properties.Properties)["agent_unique_host_id"] = model.AgentUniqueHostId
			}

			if metadata.ResourceData.HasChange("controller_host_name") {
				(*properties.Properties)["controller_host_name"] = model.ControllerHostName
			}

			if metadata.ResourceData.HasChange("controller_ssl_enabled") {
				(*properties.Properties)["controller_ssl_enabled"] = strconv.FormatBool(model.ControllerSslEnabled)
			}

			if metadata.ResourceData.HasChange("controller_port") {
				(*properties.Properties)["controller_port"] = strconv.FormatInt(model.ControllerPort, 10)
			}

			if metadata.ResourceData.HasChange("agent_account_name") {
				(*properties.Secrets)["agent_account_name"] = model.AgentAccountName
			}

			if metadata.ResourceData.HasChange("agent_account_access_key") {
				(*properties.Secrets)["agent_account_access_key"] = model.AgentAccountAccessKey
			}

			resource := appplatform.ApmResource{
				Properties: properties,
			}

			if err = client.ApmsCreateOrUpdateThenPoll(ctx, *id, resource); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("globally_enabled") {
				apmReference := appplatform.ApmReference{
					ResourceId: id.ID(),
				}
				springId := commonids.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroupName, id.SpringName)
				if model.GloballyEnabled {
					if err := client.ServicesEnableApmGloballyThenPoll(ctx, springId, apmReference); err != nil {
						return fmt.Errorf("enabling %s globally: %+v", id, err)
					}
				} else {
					if err := client.ServicesDisableApmGloballyThenPoll(ctx, springId, apmReference); err != nil {
						return fmt.Errorf("disabling %s globally: %+v", id, err)
					}
				}
			}

			return nil
		},
	}
}
