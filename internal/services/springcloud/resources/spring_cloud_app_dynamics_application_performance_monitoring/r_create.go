// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_app_dynamics_application_performance_monitoring

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (s Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudAppDynamicsApplicationPerformanceMonitoringModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppPlatform.AppPlatformClient
			springId, err := commonids.ParseSpringCloudServiceID(model.SpringCloudServiceId)
			if err != nil {
				return err
			}
			id := appplatform.NewApmID(springId.SubscriptionId, springId.ResourceGroupName, springId.ServiceName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ApmsGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(s.ResourceType(), id)
				}
			}

			resource := appplatform.ApmResource{
				Properties: &appplatform.ApmProperties{
					Type: "AppDynamics",
					Properties: pointer.To(map[string]string{
						"agent_application_name": model.AgentApplicationName,
						"agent_tier_name":        model.AgentTierName,
						"agent_node_name":        model.AgentNodeName,
						"agent_unique_host_id":   model.AgentUniqueHostId,
						"controller_host_name":   model.ControllerHostName,
						"controller_ssl_enabled": fmt.Sprintf("%t", model.ControllerSslEnabled),
						"controller_port":        fmt.Sprintf("%d", model.ControllerPort),
					}),
					Secrets: pointer.To(map[string]string{
						"agent_account_name":       model.AgentAccountName,
						"agent_account_access_key": model.AgentAccountAccessKey,
					}),
				},
			}

			if err := client.ApmsCreateOrUpdateCallbackThenPoll(ctx, id, resource, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			if model.GloballyEnabled {
				apmReference := appplatform.ApmReference{
					ResourceId: id.ID(),
				}
				if err = client.ServicesEnableApmGloballyThenPoll(ctx, *springId, apmReference); err != nil {
					return fmt.Errorf("enabling %s globally: %+v", id, err)
				}
			}

			return nil
		},
	}
}
