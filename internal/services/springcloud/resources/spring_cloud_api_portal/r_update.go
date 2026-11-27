// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_api_portal

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
)

func (s Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudAPIPortalModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := springcloudClient.From(metadata.Client).AppPlatformClient

			id, err := appplatform.ParseApiPortalID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ApiPortalsGet(ctx, *id)
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

			if metadata.ResourceData.HasChange("gateway_ids") {
				properties.GatewayIds = pointer.To(model.GatewayIds)
			}

			if metadata.ResourceData.HasChange("https_only_enabled") {
				properties.HTTPSOnly = pointer.To(model.HttpsOnlyEnabled)
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				properties.Public = pointer.To(model.PublicNetworkAccessEnabled)
			}

			if metadata.ResourceData.HasChange("sso") {
				properties.SsoProperties = expandAPIPortalSsoProperties(model.Sso)
			}

			if metadata.ResourceData.HasChange("instance_count") {
				sku.Capacity = pointer.To(model.InstanceCount)
			}

			if metadata.ResourceData.HasChange("api_try_out_enabled") {
				apiTryOutEnabledState := appplatform.ApiPortalApiTryOutEnabledStateDisabled
				if model.ApiTryOutEnabled {
					apiTryOutEnabledState = appplatform.ApiPortalApiTryOutEnabledStateEnabled
				}
				properties.ApiTryOutEnabledState = pointer.To(apiTryOutEnabledState)
			}

			apiPortalResource := appplatform.ApiPortalResource{
				Properties: properties,
				Sku:        sku,
			}
			if err = client.ApiPortalsCreateOrUpdateThenPoll(ctx, *id, apiPortalResource); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
