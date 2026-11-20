// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_dev_tool_portal

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func (s Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.DevToolPortalClient

			id, err := parse.SpringCloudDevToolPortalID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.DevToolPortalName)
			if err != nil {
				if response.WasNotFound(resp.Response.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			state := SpringCloudDevToolPortalModel{
				Name:                 id.DevToolPortalName,
				SpringCloudServiceId: commonids.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroup, id.SpringName).ID(),
			}

			var model SpringCloudDevToolPortalModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if props := resp.Properties; props != nil {
				if props.Public != nil {
					state.PublicNetworkAccessEnabled = *props.Public
				}
				if props.SsoProperties != nil {
					state.Sso = flattenSpringCloudDevToolPortalSsoProperties(props.SsoProperties, model)
				}
				if props.Features != nil {
					if props.Features.ApplicationAccelerator != nil && props.Features.ApplicationAccelerator.State == appplatform.DevToolPortalFeatureStateEnabled {
						state.ApplicationAcceleratorEnabled = true
					}
					if props.Features.ApplicationLiveView != nil && props.Features.ApplicationLiveView.State == appplatform.DevToolPortalFeatureStateEnabled {
						state.ApplicationLiveViewEnabled = true
					}
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func flattenSpringCloudDevToolPortalSsoProperties(properties *appplatform.DevToolPortalSsoProperties, model SpringCloudDevToolPortalModel) []SsoModel {
	if properties == nil {
		return []SsoModel{}
	}

	clientSecret := ""
	if len(model.Sso) != 0 {
		clientSecret = model.Sso[0].ClientSecret
	}

	scopes := make([]string, 0)
	if properties.Scopes != nil {
		scopes = *properties.Scopes
	}

	return []SsoModel{
		{
			ClientId:     pointer.From(properties.ClientID),
			ClientSecret: clientSecret,
			MetadataUrl:  pointer.From(properties.MetadataURL),
			Scope:        scopes,
		},
	}
}
