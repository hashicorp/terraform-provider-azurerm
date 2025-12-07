// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SpringCloudNewRelicApplicationPerformanceMonitoringModel struct {
	Name                         string            `tfschema:"name"`
	SpringCloudServiceId         string            `tfschema:"spring_cloud_service_id"`
	GloballyEnabled              bool              `tfschema:"globally_enabled"`
	AppName                      string            `tfschema:"app_name"`
	AgentEnabled                 bool              `tfschema:"agent_enabled"`
	AppServerPort                int64             `tfschema:"app_server_port"`
	AuditModeEnabled             bool              `tfschema:"audit_mode_enabled"`
	AutoAppNamingEnabled         bool              `tfschema:"auto_app_naming_enabled"`
	AutoTransactionNamingEnabled bool              `tfschema:"auto_transaction_naming_enabled"`
	CustomTracingEnabled         bool              `tfschema:"custom_tracing_enabled"`
	Labels                       map[string]string `tfschema:"labels"`
	LicenseKey                   string            `tfschema:"license_key"`
}

type SpringCloudNewRelicApplicationPerformanceMonitoringResource struct{}

func (s SpringCloudNewRelicApplicationPerformanceMonitoringResource) DeprecationMessage() string {
	return features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_new_relic_application_performance_monitoring` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.")
}

var (
	_ sdk.ResourceWithUpdate                      = SpringCloudNewRelicApplicationPerformanceMonitoringResource{}
	_ sdk.ResourceWithDeprecationAndNoReplacement = SpringCloudNewRelicApplicationPerformanceMonitoringResource{}
)

func (s SpringCloudNewRelicApplicationPerformanceMonitoringResource) ResourceType() string {
	return "azurerm_spring_cloud_new_relic_application_performance_monitoring"
}

func (s SpringCloudNewRelicApplicationPerformanceMonitoringResource) ModelObject() interface{} {
	return &SpringCloudNewRelicApplicationPerformanceMonitoringModel{}
}

func (s SpringCloudNewRelicApplicationPerformanceMonitoringResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return appplatform.ValidateApmID
}

func (s SpringCloudNewRelicApplicationPerformanceMonitoringResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"spring_cloud_service_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.SpringCloudServiceId{}),

		"app_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"license_key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"agent_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"app_server_port": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 65535),
		},

		"audit_mode_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"auto_app_naming_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"auto_transaction_naming_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"custom_tracing_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"labels": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},

		"globally_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
	}
}

func (s SpringCloudNewRelicApplicationPerformanceMonitoringResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (s SpringCloudNewRelicApplicationPerformanceMonitoringResource) Create() sdk.ResourceFunc {
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

			existing, err := client.ApmsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(s.ResourceType(), id)
			}

			resource := appplatform.ApmResource{
				Properties: &appplatform.ApmProperties{
					Type: "NewRelic",
					Properties: pointer.To(map[string]string{
						"app_name":                       model.AppName,
						"agent_enabled":                  fmt.Sprintf("%t", model.AgentEnabled),
						"appserver_port":                 fmt.Sprintf("%d", model.AppServerPort),
						"audit_mode":                     fmt.Sprintf("%t", model.AuditModeEnabled),
						"enable_auto_app_naming":         fmt.Sprintf("%t", model.AutoAppNamingEnabled),
						"enable_auto_transaction_naming": fmt.Sprintf("%t", model.AutoTransactionNamingEnabled),
						"enable_custom_tracing":          fmt.Sprintf("%t", model.CustomTracingEnabled),
						"labels":                         expandNewRelicLabels(model.Labels),
					}),
					Secrets: pointer.To(map[string]string{
						"license_key": model.LicenseKey,
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

func (s SpringCloudNewRelicApplicationPerformanceMonitoringResource) Update() sdk.ResourceFunc {
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

func (s SpringCloudNewRelicApplicationPerformanceMonitoringResource) Read() sdk.ResourceFunc {
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

			var model SpringCloudNewRelicApplicationPerformanceMonitoringModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := SpringCloudNewRelicApplicationPerformanceMonitoringModel{
				Name:                 id.ApmName,
				SpringCloudServiceId: springId.ID(),
				GloballyEnabled:      globallyEnabled,
				LicenseKey:           model.LicenseKey,
			}

			if props := resp.Model.Properties; props != nil {
				if props.Type != "NewRelic" {
					return fmt.Errorf("retrieving %s: type was %s, expected NewRelic", *id, props.Type)
				}
				if props.Properties != nil {
					if value, ok := (*props.Properties)["app_name"]; ok {
						state.AppName = value
					}
					if value, ok := (*props.Properties)["agent_enabled"]; ok {
						state.AgentEnabled = value == "true"
					}
					if value, ok := (*props.Properties)["appserver_port"]; ok {
						if v, err := strconv.ParseInt(value, 10, 0); err == nil {
							state.AppServerPort = v
						}
					}
					if value, ok := (*props.Properties)["audit_mode"]; ok {
						state.AuditModeEnabled = value == "true"
					}
					if value, ok := (*props.Properties)["enable_auto_app_naming"]; ok {
						state.AutoAppNamingEnabled = value == "true"
					}
					if value, ok := (*props.Properties)["enable_auto_transaction_naming"]; ok {
						state.AutoTransactionNamingEnabled = value == "true"
					}
					if value, ok := (*props.Properties)["enable_custom_tracing"]; ok {
						state.CustomTracingEnabled = value == "true"
					}
					if value, ok := (*props.Properties)["labels"]; ok {
						state.Labels = flattenNewRelicLabels(value)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (s SpringCloudNewRelicApplicationPerformanceMonitoringResource) Delete() sdk.ResourceFunc {
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

func flattenNewRelicLabels(input string) map[string]string {
	if input == "" {
		return nil
	}
	labels := make(map[string]string)
	for _, label := range strings.Split(input, ";") {
		parts := strings.Split(label, ":")
		if len(parts) == 2 {
			labels[parts[0]] = parts[1]
		}
	}
	return labels
}
