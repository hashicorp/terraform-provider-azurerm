// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_elastic_application_performance_monitoring

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
)

func (s Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := springcloudClient.From(metadata.Client).AppPlatformClient

			id, err := appplatform.ParseApmID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SpringCloudElasticApplicationPerformanceMonitoringModel
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

			if metadata.ResourceData.HasChange("service_name") {
				(*properties.Properties)["service_name"] = model.ServiceName
			}

			if metadata.ResourceData.HasChange("application_packages") {
				(*properties.Properties)["application_packages"] = strings.Join(model.ApplicationPackages, ",")
			}

			if metadata.ResourceData.HasChange("server_url") {
				(*properties.Properties)["server_url"] = model.ServerUrl
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
