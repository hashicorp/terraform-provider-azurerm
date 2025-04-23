// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"context"
	"fmt"
	"strings"
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

type SpringCloudElasticApplicationPerformanceMonitoringModel struct {
	Name                 string   `tfschema:"name"`
	SpringCloudServiceId string   `tfschema:"spring_cloud_service_id"`
	GloballyEnabled      bool     `tfschema:"globally_enabled"`
	ServiceName          string   `tfschema:"service_name"`
	ApplicationPackages  []string `tfschema:"application_packages"`
	ServerUrl            string   `tfschema:"server_url"`
}

type SpringCloudElasticApplicationPerformanceMonitoringResource struct{}

func (s SpringCloudElasticApplicationPerformanceMonitoringResource) DeprecationMessage() string {
	return features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_elastic_application_performance_monitoring` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.")
}

var (
	_ sdk.ResourceWithUpdate                      = SpringCloudElasticApplicationPerformanceMonitoringResource{}
	_ sdk.ResourceWithDeprecationAndNoReplacement = SpringCloudElasticApplicationPerformanceMonitoringResource{}
)

func (s SpringCloudElasticApplicationPerformanceMonitoringResource) ResourceType() string {
	return "azurerm_spring_cloud_elastic_application_performance_monitoring"
}

func (s SpringCloudElasticApplicationPerformanceMonitoringResource) ModelObject() interface{} {
	return &SpringCloudElasticApplicationPerformanceMonitoringModel{}
}

func (s SpringCloudElasticApplicationPerformanceMonitoringResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return appplatform.ValidateApmID
}

func (s SpringCloudElasticApplicationPerformanceMonitoringResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"spring_cloud_service_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.SpringCloudServiceId{}),

		"application_packages": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"service_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"server_url": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"globally_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func (s SpringCloudElasticApplicationPerformanceMonitoringResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (s SpringCloudElasticApplicationPerformanceMonitoringResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudElasticApplicationPerformanceMonitoringModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppPlatform.AppPlatformClient
			springId, err := commonids.ParseSpringCloudServiceID(model.SpringCloudServiceId)
			if err != nil {
				return err
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
					Type: "ElasticAPM",
					Properties: pointer.To(map[string]string{
						"service_name":         model.ServiceName,
						"application_packages": strings.Join(model.ApplicationPackages, ","),
						"server_url":           model.ServerUrl,
					}),
					Secrets: pointer.To(map[string]string{}),
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

func (s SpringCloudElasticApplicationPerformanceMonitoringResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

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

func (s SpringCloudElasticApplicationPerformanceMonitoringResource) Read() sdk.ResourceFunc {
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

			var model SpringCloudElasticApplicationPerformanceMonitoringModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := SpringCloudElasticApplicationPerformanceMonitoringModel{
				Name:                 id.ApmName,
				SpringCloudServiceId: springId.ID(),
				GloballyEnabled:      globallyEnabled,
			}

			if props := resp.Model.Properties; props != nil {
				if props.Type != "ElasticAPM" {
					return fmt.Errorf("retrieving %s: type was not ElasticAPM", *id)
				}
				if props.Properties != nil {
					if value, ok := (*props.Properties)["service_name"]; ok {
						state.ServiceName = value
					}
					if value, ok := (*props.Properties)["application_packages"]; ok {
						state.ApplicationPackages = strings.Split(value, ",")
					}
					if value, ok := (*props.Properties)["server_url"]; ok {
						state.ServerUrl = value
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (s SpringCloudElasticApplicationPerformanceMonitoringResource) Delete() sdk.ResourceFunc {
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
