// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_new_relic_application_performance_monitoring

import (
	"context"
	"fmt"
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

			var model SpringCloudNewRelicApplicationPerformanceMonitoringModel
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

			if metadata.ResourceData.HasChange("app_name") {
				(*properties.Properties)["app_name"] = model.AppName
			}

			if metadata.ResourceData.HasChange("agent_enabled") {
				(*properties.Properties)["agent_enabled"] = fmt.Sprintf("%t", model.AgentEnabled)
			}

			if metadata.ResourceData.HasChange("app_server_port") {
				(*properties.Properties)["appserver_port"] = fmt.Sprintf("%d", model.AppServerPort)
			}

			if metadata.ResourceData.HasChange("audit_mode_enabled") {
				(*properties.Properties)["audit_mode"] = fmt.Sprintf("%t", model.AuditModeEnabled)
			}

			if metadata.ResourceData.HasChange("auto_app_naming_enabled") {
				(*properties.Properties)["enable_auto_app_naming"] = fmt.Sprintf("%t", model.AutoAppNamingEnabled)
			}

			if metadata.ResourceData.HasChange("auto_transaction_naming_enabled") {
				(*properties.Properties)["enable_auto_transaction_naming"] = fmt.Sprintf("%t", model.AutoTransactionNamingEnabled)
			}

			if metadata.ResourceData.HasChange("custom_tracing_enabled") {
				(*properties.Properties)["enable_custom_tracing"] = fmt.Sprintf("%t", model.CustomTracingEnabled)
			}

			if metadata.ResourceData.HasChange("labels") {
				(*properties.Properties)["labels"] = expandNewRelicLabels(model.Labels)
			}

			if metadata.ResourceData.HasChange("license_key") {
				(*properties.Secrets)["license_key"] = model.LicenseKey
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
