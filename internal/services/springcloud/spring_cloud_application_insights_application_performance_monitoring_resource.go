// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SpringCloudApplicationInsightsApplicationPerformanceMonitoringModel struct {
	Name                      string `tfschema:"name"`
	SpringCloudServiceId      string `tfschema:"spring_cloud_service_id"`
	GloballyEnabled           bool   `tfschema:"globally_enabled"`
	ConnectionString          string `tfschema:"connection_string"`
	RoleName                  string `tfschema:"role_name"`
	RoleInstance              string `tfschema:"role_instance"`
	SamplingRequestsPerSecond int64  `tfschema:"sampling_requests_per_second"`
	SamplingPercentage        int64  `tfschema:"sampling_percentage"`
}

type SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource struct{}

func (s SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource) DeprecationMessage() string {
	return features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_application_insights_application_performance_monitoring` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.")
}

var (
	_ sdk.ResourceWithUpdate                      = SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource{}
	_ sdk.ResourceWithDeprecationAndNoReplacement = SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource{}
)

func (s SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource) ResourceType() string {
	return "azurerm_spring_cloud_application_insights_application_performance_monitoring"
}

func (s SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource) ModelObject() interface{} {
	return &SpringCloudApplicationInsightsApplicationPerformanceMonitoringModel{}
}

func (s SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return appplatform.ValidateApmID
}

func (s SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"spring_cloud_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SpringCloudServiceID,
		},

		"globally_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"connection_string": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"role_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"role_instance": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"sampling_requests_per_second": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"sampling_percentage": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 100),
		},
	}
}

func (s SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (s SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudApplicationInsightsApplicationPerformanceMonitoringModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppPlatform.AppPlatformClient
			springId, err := commonids.ParseSpringCloudServiceID(model.SpringCloudServiceId)
			if err != nil {
				return fmt.Errorf("parsing spring service ID: %+v", err)
			}
			id := appplatform.NewApmID(springId.SubscriptionId, springId.ResourceGroupName, springId.ServiceName, model.Name)

			existing, err := client.ApmsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(s.ResourceType(), id)
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
			err = client.ApmsCreateOrUpdateThenPoll(ctx, id, resource)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if model.GloballyEnabled {
				apmReference := appplatform.ApmReference{
					ResourceId: id.ID(),
				}
				err = client.ServicesEnableApmGloballyThenPoll(ctx, *springId, apmReference)
				if err != nil {
					return fmt.Errorf("enabling %s globally: %+v", id, err)
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (s SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseApmID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SpringCloudApplicationInsightsApplicationPerformanceMonitoringModel
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

			if metadata.ResourceData.HasChange("role_name") {
				(*properties.Properties)["role_name"] = model.RoleName
			}

			if metadata.ResourceData.HasChange("role_instance") {
				(*properties.Properties)["role_instance"] = model.RoleInstance
			}

			if metadata.ResourceData.HasChange("sampling_requests_per_second") {
				(*properties.Properties)["sampling_requests_per_second"] = fmt.Sprintf("%d", model.SamplingRequestsPerSecond)
			}

			if metadata.ResourceData.HasChange("sampling_percentage") {
				(*properties.Properties)["sampling_percentage"] = fmt.Sprintf("%d", model.SamplingPercentage)
			}

			if metadata.ResourceData.HasChange("connection_string") {
				(*properties.Secrets)["connection_string"] = model.ConnectionString
			}

			resource := appplatform.ApmResource{
				Properties: properties,
			}

			err = client.ApmsCreateOrUpdateThenPoll(ctx, *id, resource)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("globally_enabled") {
				apmReference := appplatform.ApmReference{
					ResourceId: id.ID(),
				}
				springId := commonids.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroupName, id.SpringName)
				if model.GloballyEnabled {
					err := client.ServicesEnableApmGloballyThenPoll(ctx, springId, apmReference)
					if err != nil {
						return fmt.Errorf("enabling %s globally: %+v", id, err)
					}
				} else {
					err := client.ServicesDisableApmGloballyThenPoll(ctx, springId, apmReference)
					if err != nil {
						return fmt.Errorf("disabling %s globally: %+v", id, err)
					}
				}
			}

			return nil
		},
	}
}

func (s SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource) Read() sdk.ResourceFunc {
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

func (s SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseApmID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.ApmsDeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
