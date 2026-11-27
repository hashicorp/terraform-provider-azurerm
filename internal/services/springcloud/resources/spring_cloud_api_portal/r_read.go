// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_api_portal

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
)

func (s Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := springcloudClient.From(metadata.Client).AppPlatformClient

			id, err := appplatform.ParseApiPortalID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ApiPortalsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			springId := commonids.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroupName, id.SpringName)

			var model SpringCloudAPIPortalModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := SpringCloudAPIPortalModel{
				Name:                 id.ApiPortalName,
				SpringCloudServiceId: springId.ID(),
			}
			if resp.Model != nil {
				if props := resp.Model.Properties; props != nil {
					state.GatewayIds = flattenSpringCloudAPIPortalGatewayIds(props.GatewayIds)
					state.HttpsOnlyEnabled = pointer.From(props.HTTPSOnly)
					state.PublicNetworkAccessEnabled = pointer.From(props.Public)
					state.Sso = flattenAPIPortalSsoProperties(props.SsoProperties, model.Sso)
					state.ApiTryOutEnabled = props.ApiTryOutEnabledState != nil && *props.ApiTryOutEnabledState == appplatform.ApiPortalApiTryOutEnabledStateEnabled
				}

				if sku := resp.Model.Sku; sku != nil {
					state.InstanceCount = pointer.From(sku.Capacity)
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func flattenAPIPortalSsoProperties(input *appplatform.SsoProperties, old []ApiPortalSsoModel) []ApiPortalSsoModel {
	if input == nil {
		return make([]ApiPortalSsoModel, 0)
	}

	oldItems := make(map[string]ApiPortalSsoModel)
	for _, item := range old {
		if item.IssuerUri != "" {
			oldItems[item.IssuerUri] = item
		}
	}

	issuerUri := pointer.From(input.IssuerUri)
	var clientId string
	var clientSecret string
	if oldItem, ok := oldItems[issuerUri]; ok {
		if oldItem.ClientId != "" {
			clientId = oldItem.ClientId
		}
		if oldItem.ClientSecret != "" {
			clientSecret = oldItem.ClientSecret
		}
	}
	return []ApiPortalSsoModel{
		{
			ClientId:     clientId,
			ClientSecret: clientSecret,
			IssuerUri:    issuerUri,
			Scope:        pointer.From(input.Scope),
		},
	}
}

func flattenSpringCloudAPIPortalGatewayIds(ids *[]string) []string {
	if ids == nil || len(*ids) == 0 {
		return []string{}
	}
	out := make([]string, 0)
	for _, id := range *ids {
		gatewayId, err := appplatform.ParseGatewayIDInsensitively(id)
		if err == nil {
			out = append(out, gatewayId.ID())
		}
	}
	return out
}
