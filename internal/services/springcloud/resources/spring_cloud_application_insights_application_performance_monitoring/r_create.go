// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_application_insights_application_performance_monitoring

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
)

func (s Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudApplicationInsightsApplicationPerformanceMonitoringModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := springcloudClient.From(metadata.Client).AppPlatformClient
			springId, err := commonids.ParseSpringCloudServiceID(model.SpringCloudServiceId)
			if err != nil {
				return fmt.Errorf("parsing spring service ID: %+v", err)
			}
			id := appplatform.NewApmID(springId.SubscriptionId, springId.ResourceGroupName, springId.ServiceName, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
					Type: "ApplicationInsights",
					Properties: pointer.To(map[string]string{
						"role_name":                    model.RoleName,
						"role_instance":                model.RoleInstance,
						"sampling_requests_per_second": fmt.Sprintf("%d", model.SamplingRequestsPerSecond),
						"sampling_percentage":          fmt.Sprintf("%d", model.SamplingPercentage),
					}),
					Secrets: pointer.To(map[string]string{
						"connection_string": model.ConnectionString,
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
