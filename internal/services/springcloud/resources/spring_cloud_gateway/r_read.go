// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_gateway

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

func (s Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.AppPlatformClient

			id, err := appplatform.ParseGatewayID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GatewaysGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			springId := commonids.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroupName, id.SpringName)

			var model SpringCloudGatewayModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := SpringCloudGatewayModel{
				Name:                          id.GatewayName,
				SpringCloudServiceId:          springId.ID(),
				SensitiveEnvironmentVariables: model.SensitiveEnvironmentVariables,
			}

			if resp.Model != nil {
				if props := resp.Model.Properties; props != nil {
					state.ApiMetadata = flattenGatewayGatewayAPIMetadataProperties(props.ApiMetadataProperties)
					apms, err := flattenGatewayApms(props.Apms)
					if err != nil {
						return err
					}
					state.ApplicationPerformanceMonitoringIds = apms
					state.ApplicationPerformanceMonitoringTypes = flattenGatewayGatewayApmTypes(props.ApmTypes)
					state.ClientAuthorization = flattenGatewayClientAuth(props.ClientAuth)
					state.Cors = flattenGatewayGatewayCorsProperties(props.CorsProperties)
					if props.EnvironmentVariables != nil {
						state.EnvironmentVariables = pointer.From(props.EnvironmentVariables.Properties)
					}
					state.HttpsOnly = pointer.From(props.HTTPSOnly)
					state.PublicNetworkAccessEnabled = pointer.From(props.Public)
					state.Quota = flattenGatewayGatewayResourceRequests(props.ResourceRequests)
					state.Sso = flattenGatewaySsoProperties(props.SsoProperties, model.Sso)
					state.LocalResponseCachePerRoute = flattenGatewayLocalResponseCachePerRouteProperties(props.ResponseCacheProperties)
					state.LocalResponseCachePerInstance = flattenGatewayLocalResponseCachePerInstanceProperties(props.ResponseCacheProperties)
				}

				if sku := resp.Model.Sku; sku != nil {
					state.InstanceCount = pointer.From(sku.Capacity)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
