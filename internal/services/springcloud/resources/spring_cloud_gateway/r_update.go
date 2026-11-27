// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (s Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudGatewayModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := appplatform.ParseGatewayID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.AppPlatform.AppPlatformClient
			resp, err := client.GatewaysGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			properties := resp.Model.Properties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			sku := resp.Model.Sku
			if sku == nil {
				return fmt.Errorf("retrieving %s: sku was nil", id)
			}

			if metadata.ResourceData.HasChange("client_authorization") {
				properties.ClientAuth = expandGatewayClientAuth(model.ClientAuthorization)
			}

			if metadata.ResourceData.HasChange("api_metadata") {
				properties.ApiMetadataProperties = expandGatewayGatewayAPIMetadataProperties(model.ApiMetadata)
			}

			if metadata.ResourceData.HasChange("application_performance_monitoring_ids") {
				properties.Apms = expandGatewayApms(model.ApplicationPerformanceMonitoringIds)
			}

			if metadata.ResourceData.HasChange("application_performance_monitoring_types") {
				properties.ApmTypes = expandGatewayGatewayApmTypes(model.ApplicationPerformanceMonitoringTypes)
			}

			if metadata.ResourceData.HasChange("cors") {
				properties.CorsProperties = expandGatewayGatewayCorsProperties(model.Cors)
			}

			if metadata.ResourceData.HasChanges("environment_variables", "sensitive_environment_variables") {
				properties.EnvironmentVariables = expandGatewayGatewayEnvironmentVariables(model.EnvironmentVariables, model.SensitiveEnvironmentVariables)
			}

			if metadata.ResourceData.HasChange("https_only") {
				properties.HTTPSOnly = pointer.To(model.HttpsOnly)
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				properties.Public = pointer.To(model.PublicNetworkAccessEnabled)
			}

			if metadata.ResourceData.HasChange("quota") {
				properties.ResourceRequests = expandGatewayGatewayResourceRequests(model.Quota)
			}

			if metadata.ResourceData.HasChange("sso") {
				properties.SsoProperties = expandGatewaySsoProperties(model.Sso)
			}

			if metadata.ResourceData.HasChanges("local_response_cache_per_instance", "local_response_cache_per_route") {
				properties.ResponseCacheProperties = expandGatewayResponseCacheProperties(model)
			}

			if metadata.ResourceData.HasChange("instance_count") {
				sku.Capacity = pointer.To(model.InstanceCount)
			}
			resource := appplatform.GatewayResource{
				Properties: properties,
				Sku:        sku,
			}

			if err = client.GatewaysCreateOrUpdateThenPoll(ctx, *id, resource); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
