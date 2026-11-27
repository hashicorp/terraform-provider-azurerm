// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_new_relic_application_performance_monitoring

import (
	"context"
	"fmt"
	"strconv"
	"strings"
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
			var model SpringCloudNewRelicApplicationPerformanceMonitoringModel
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
					Type: "NewRelic",
					Properties: pointer.To(map[string]string{
						"app_name":                       model.AppName,
						"agent_enabled":                  strconv.FormatBool(model.AgentEnabled),
						"appserver_port":                 strconv.FormatInt(model.AppServerPort, 10),
						"audit_mode":                     strconv.FormatBool(model.AuditModeEnabled),
						"enable_auto_app_naming":         strconv.FormatBool(model.AutoAppNamingEnabled),
						"enable_auto_transaction_naming": strconv.FormatBool(model.AutoTransactionNamingEnabled),
						"enable_custom_tracing":          strconv.FormatBool(model.CustomTracingEnabled),
						"labels":                         expandNewRelicLabels(model.Labels),
					}),
					Secrets: pointer.To(map[string]string{
						"license_key": model.LicenseKey,
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

func expandNewRelicLabels(input map[string]string) string {
	if len(input) == 0 {
		return ""
	}
	labels := make([]string, 0)
	for k, v := range input {
		labels = append(labels, fmt.Sprintf("%s:%s", k, v))
	}
	return strings.Join(labels, ";")
}
