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
)

func (s Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudAPIPortalModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppPlatform.AppPlatformClient
			springId, err := commonids.ParseSpringCloudServiceID(model.SpringCloudServiceId)
			if err != nil {
				return err
			}
			id := appplatform.NewApiPortalID(springId.SubscriptionId, springId.ResourceGroupName, springId.ServiceName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ApiPortalsGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(s.ResourceType(), id)
				}
			}

			service, err := client.ServicesGet(ctx, *springId)
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", springId, err)
			}
			if service.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", springId)
			}
			if service.Model.Sku == nil || service.Model.Sku.Name == nil || service.Model.Sku.Tier == nil {
				return fmt.Errorf("invalid `sku` for %s", springId)
			}

			apiTryOutEnabledState := appplatform.ApiPortalApiTryOutEnabledStateDisabled
			if model.ApiTryOutEnabled {
				apiTryOutEnabledState = appplatform.ApiPortalApiTryOutEnabledStateEnabled
			}

			apiPortalResource := appplatform.ApiPortalResource{
				Properties: &appplatform.ApiPortalProperties{
					GatewayIds:            pointer.To(model.GatewayIds),
					HTTPSOnly:             pointer.To(model.HttpsOnlyEnabled),
					Public:                pointer.To(model.PublicNetworkAccessEnabled),
					SsoProperties:         expandAPIPortalSsoProperties(model.Sso),
					ApiTryOutEnabledState: pointer.To(apiTryOutEnabledState),
				},
				Sku: &appplatform.Sku{
					Name:     service.Model.Sku.Name,
					Tier:     service.Model.Sku.Tier,
					Capacity: pointer.To(model.InstanceCount),
				},
			}

			if err := client.ApiPortalsCreateOrUpdateCallbackThenPoll(ctx, id, apiPortalResource, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}

func expandAPIPortalSsoProperties(input []ApiPortalSsoModel) *appplatform.SsoProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	return &appplatform.SsoProperties{
		Scope:        pointer.To(v.Scope),
		ClientId:     pointer.To(v.ClientId),
		ClientSecret: pointer.To(v.ClientSecret),
		IssuerUri:    pointer.To(v.IssuerUri),
	}
}
