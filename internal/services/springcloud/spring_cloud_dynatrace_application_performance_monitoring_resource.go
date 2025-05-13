// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SpringCloudDynatraceApplicationPerformanceMonitoringModel struct {
	Name                 string `tfschema:"name"`
	SpringCloudServiceId string `tfschema:"spring_cloud_service_id"`
	GloballyEnabled      bool   `tfschema:"globally_enabled"`
	ApiUrl               string `tfschema:"api_url"`
	ApiToken             string `tfschema:"api_token"`
	ConnectionPoint      string `tfschema:"connection_point"`
	EnvironmentId        string `tfschema:"environment_id"`
	Tenant               string `tfschema:"tenant"`
	TenantToken          string `tfschema:"tenant_token"`
}

type SpringCloudDynatraceApplicationPerformanceMonitoringResource struct{}

func (s SpringCloudDynatraceApplicationPerformanceMonitoringResource) DeprecationMessage() string {
	return features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_dynatrace_application_performance_monitoring` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.")
}

var (
	_ sdk.ResourceWithUpdate                      = SpringCloudDynatraceApplicationPerformanceMonitoringResource{}
	_ sdk.ResourceWithDeprecationAndNoReplacement = SpringCloudDynatraceApplicationPerformanceMonitoringResource{}
)

func (s SpringCloudDynatraceApplicationPerformanceMonitoringResource) ResourceType() string {
	return "azurerm_spring_cloud_dynatrace_application_performance_monitoring"
}

func (s SpringCloudDynatraceApplicationPerformanceMonitoringResource) ModelObject() interface{} {
	return &SpringCloudDynatraceApplicationPerformanceMonitoringModel{}
}

func (s SpringCloudDynatraceApplicationPerformanceMonitoringResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return appplatform.ValidateApmID
}

func (s SpringCloudDynatraceApplicationPerformanceMonitoringResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"spring_cloud_service_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.SpringCloudServiceId{}),

		"connection_point": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tenant": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tenant_token": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"api_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},

		"api_token": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"environment_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"globally_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func (s SpringCloudDynatraceApplicationPerformanceMonitoringResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (s SpringCloudDynatraceApplicationPerformanceMonitoringResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudDynatraceApplicationPerformanceMonitoringModel
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
					Type: "Dynatrace",
					Properties: pointer.To(map[string]string{
						"api-url":          model.ApiUrl,
						"connection_point": model.ConnectionPoint,
						"environment-id":   model.EnvironmentId,
					}),
					Secrets: pointer.To(map[string]string{
						"api-token":   model.ApiToken,
						"tenanttoken": model.TenantToken,
						"tenant":      model.Tenant,
					}),
				},
			}
			err = client.ApmsCreateOrUpdateThenPoll(ctx, id, resource)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model.GloballyEnabled {
				apmReference := appplatform.ApmReference{
					ResourceId: id.ID(),
				}
				err = client.ServicesEnableApmGloballyThenPoll(ctx, *springId, apmReference)
				if err != nil {
					return fmt.Errorf("enabling %s globally: %+v", id, err)
				}
			}

			return nil
		},
	}
}

func (s SpringCloudDynatraceApplicationPerformanceMonitoringResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseApmID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SpringCloudDynatraceApplicationPerformanceMonitoringModel
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

			if metadata.ResourceData.HasChange("api_url") {
				(*properties.Properties)["api-url"] = model.ApiUrl
			}

			if metadata.ResourceData.HasChange("api_token") {
				(*properties.Secrets)["api-token"] = model.ApiToken
			}

			if metadata.ResourceData.HasChange("connection_point") {
				(*properties.Properties)["connection_point"] = model.ConnectionPoint
			}

			if metadata.ResourceData.HasChange("environment_id") {
				(*properties.Properties)["environment-id"] = model.EnvironmentId
			}

			if metadata.ResourceData.HasChange("tenant") {
				(*properties.Secrets)["tenant"] = model.Tenant
			}

			if metadata.ResourceData.HasChange("tenant_token") {
				(*properties.Secrets)["tenanttoken"] = model.TenantToken
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

func (s SpringCloudDynatraceApplicationPerformanceMonitoringResource) Read() sdk.ResourceFunc {
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

			var model SpringCloudDynatraceApplicationPerformanceMonitoringModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := SpringCloudDynatraceApplicationPerformanceMonitoringModel{
				Name:                 id.ApmName,
				SpringCloudServiceId: springId.ID(),
				GloballyEnabled:      globallyEnabled,
				ApiToken:             model.ApiToken,
				TenantToken:          model.TenantToken,
				Tenant:               model.Tenant,
			}

			if props := resp.Model.Properties; props != nil {
				if props.Type != "Dynatrace" {
					return fmt.Errorf("retrieving %s: type was not Dynatrace", *id)
				}
				if props.Properties != nil {
					if value, ok := (*props.Properties)["api-url"]; ok {
						state.ApiUrl = value
					}
					if value, ok := (*props.Properties)["connection_point"]; ok {
						state.ConnectionPoint = value
					}
					if value, ok := (*props.Properties)["environment-id"]; ok {
						state.EnvironmentId = value
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (s SpringCloudDynatraceApplicationPerformanceMonitoringResource) Delete() sdk.ResourceFunc {
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
