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

func (s Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudDevToolPortalModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppPlatform.DevToolPortalClient
			// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
			// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
			springId, err := commonids.ParseSpringCloudServiceIDInsensitively(model.SpringCloudServiceId)
			if err != nil {
				return fmt.Errorf("parsing spring service ID: %+v", err)
			}
			id := parse.NewSpringCloudDevToolPortalID(springId.SubscriptionId, springId.ResourceGroupName, springId.ServiceName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.DevToolPortalName)
				if err != nil && !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.Response.Response) {
					return metadata.ResourceRequiresImport(s.ResourceType(), id)
				}
			}

			DevToolPortalResource := appplatform.DevToolPortalResource{
				Properties: &appplatform.DevToolPortalProperties{
					Public:        pointer.To(model.PublicNetworkAccessEnabled),
					SsoProperties: expandSpringCloudDevToolPortalSsoProperties(model.Sso),
					Features:      expandSpringCloudDevToolPortalFeatures(model.ApplicationAcceleratorEnabled, model.ApplicationLiveViewEnabled),
				},
			}
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.DevToolPortalName, DevToolPortalResource)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandSpringCloudDevToolPortalFeatures(applicationAcceleratorEnabled, applicationLiveViewEnabled bool) *appplatform.DevToolPortalFeatureSettings {
	applicationAcceleratorState := appplatform.DevToolPortalFeatureStateDisabled
	if applicationAcceleratorEnabled {
		applicationAcceleratorState = appplatform.DevToolPortalFeatureStateEnabled
	}

	applicationLiveViewState := appplatform.DevToolPortalFeatureStateDisabled
	if applicationLiveViewEnabled {
		applicationLiveViewState = appplatform.DevToolPortalFeatureStateEnabled
	}

	return &appplatform.DevToolPortalFeatureSettings{
		ApplicationAccelerator: &appplatform.DevToolPortalFeatureDetail{
			State: applicationAcceleratorState,
		},
		ApplicationLiveView: &appplatform.DevToolPortalFeatureDetail{
			State: applicationLiveViewState,
		},
	}
}

func expandSpringCloudDevToolPortalSsoProperties(input []SsoModel) *appplatform.DevToolPortalSsoProperties {
	if len(input) == 0 {
		return nil
	}
	sso := input[0]

	return &appplatform.DevToolPortalSsoProperties{
		Scopes:       &sso.Scope,
		ClientID:     pointer.To(sso.ClientId),
		ClientSecret: pointer.To(sso.ClientSecret),
		MetadataURL:  pointer.To(sso.MetadataUrl),
	}
}
