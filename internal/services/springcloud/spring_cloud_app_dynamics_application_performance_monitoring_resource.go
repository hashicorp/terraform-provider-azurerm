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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SpringCloudAppDynamicsApplicationPerformanceMonitoringModel struct {
	Name                  string `tfschema:"name"`
	SpringCloudServiceId  string `tfschema:"spring_cloud_service_id"`
	GloballyEnabled       bool   `tfschema:"globally_enabled"`
	AgentApplicationName  string `tfschema:"agent_application_name"`
	AgentTierName         string `tfschema:"agent_tier_name"`
	AgentNodeName         string `tfschema:"agent_node_name"`
	AgentUniqueHostId     string `tfschema:"agent_unique_host_id"`
	ControllerHostName    string `tfschema:"controller_host_name"`
	ControllerSslEnabled  bool   `tfschema:"controller_ssl_enabled"`
	ControllerPort        int64  `tfschema:"controller_port"`
	AgentAccountName      string `tfschema:"agent_account_name"`
	AgentAccountAccessKey string `tfschema:"agent_account_access_key"`
}

type SpringCloudAppDynamicsApplicationPerformanceMonitoringResource struct{}

func (s SpringCloudAppDynamicsApplicationPerformanceMonitoringResource) DeprecationMessage() string {
	return features.DeprecatedInFivePointOh("Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_app_dynamics_application_performance_monitoring` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.")
}

var (
	_ sdk.ResourceWithUpdate                      = SpringCloudAppDynamicsApplicationPerformanceMonitoringResource{}
	_ sdk.ResourceWithDeprecationAndNoReplacement = SpringCloudAppDynamicsApplicationPerformanceMonitoringResource{}
)

func (s SpringCloudAppDynamicsApplicationPerformanceMonitoringResource) ResourceType() string {
	return "azurerm_spring_cloud_app_dynamics_application_performance_monitoring"
}

func (s SpringCloudAppDynamicsApplicationPerformanceMonitoringResource) ModelObject() interface{} {
	return &SpringCloudAppDynamicsApplicationPerformanceMonitoringModel{}
}

func (s SpringCloudAppDynamicsApplicationPerformanceMonitoringResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return appplatform.ValidateApmID
}

func (s SpringCloudAppDynamicsApplicationPerformanceMonitoringResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"spring_cloud_service_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.SpringCloudServiceId{}),

		"agent_account_access_key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"agent_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"controller_host_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"agent_application_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"agent_tier_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"agent_node_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"agent_unique_host_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"controller_ssl_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"controller_port": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
		},

		"globally_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
	}
}

func (s SpringCloudAppDynamicsApplicationPerformanceMonitoringResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (s SpringCloudAppDynamicsApplicationPerformanceMonitoringResource) Create() sdk.ResourceFunc {
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

			existing, err := client.ApmsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(s.ResourceType(), id)
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

func (s SpringCloudAppDynamicsApplicationPerformanceMonitoringResource) Update() sdk.ResourceFunc {
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
				(*properties.Properties)["controller_ssl_enabled"] = fmt.Sprintf("%t", model.ControllerSslEnabled)
			}

			if metadata.ResourceData.HasChange("controller_port") {
				(*properties.Properties)["controller_port"] = fmt.Sprintf("%d", model.ControllerPort)
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

func (s SpringCloudAppDynamicsApplicationPerformanceMonitoringResource) Read() sdk.ResourceFunc {
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

			var model SpringCloudAppDynamicsApplicationPerformanceMonitoringModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := SpringCloudAppDynamicsApplicationPerformanceMonitoringModel{
				Name:                  id.ApmName,
				SpringCloudServiceId:  springId.ID(),
				GloballyEnabled:       globallyEnabled,
				AgentAccountName:      model.AgentAccountName,
				AgentAccountAccessKey: model.AgentAccountAccessKey,
			}

			if props := resp.Model.Properties; props != nil {
				if props.Type != "AppDynamics" {
					return fmt.Errorf("retrieving %s: expected type AppDynamics, got %s", *id, props.Type)
				}
				if props.Properties != nil {
					if value, ok := (*props.Properties)["agent_application_name"]; ok {
						state.AgentApplicationName = value
					}
					if value, ok := (*props.Properties)["agent_tier_name"]; ok {
						state.AgentTierName = value
					}
					if value, ok := (*props.Properties)["agent_node_name"]; ok {
						state.AgentNodeName = value
					}
					if value, ok := (*props.Properties)["agent_unique_host_id"]; ok {
						state.AgentUniqueHostId = value
					}
					if value, ok := (*props.Properties)["controller_host_name"]; ok {
						state.ControllerHostName = value
					}
					if value, ok := (*props.Properties)["controller_ssl_enabled"]; ok {
						state.ControllerSslEnabled = value == "true"
					}
					if value, ok := (*props.Properties)["controller_port"]; ok {
						if v, err := strconv.ParseInt(value, 10, 32); err == nil {
							state.ControllerPort = v
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (s SpringCloudAppDynamicsApplicationPerformanceMonitoringResource) Delete() sdk.ResourceFunc {
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
