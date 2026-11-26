// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_provider_feature_registration

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2021-07-01/features"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.FeaturesClient

			var obj ResourceProviderFeatureRegistrationModel
			if err := metadata.Decode(&obj); err != nil {
				return err
			}

			featureId := features.NewFeatureID(metadata.Client.Account.SubscriptionId, obj.ProviderName, obj.Name)

			provider, err := client.Get(ctx, featureId)
			if err != nil {
				if response.WasNotFound(provider.HttpResponse) {
					return fmt.Errorf("%s was not found", featureId)
				}

				return fmt.Errorf("retrieving %s: %+v", featureId, err)
			}

			registrationState := ""
			if model := provider.Model; model != nil && model.Properties != nil && model.Properties.State != nil {
				registrationState = *model.Properties.State
			}

			if registrationState == "" {
				return fmt.Errorf("retrieving %s: `registrationState` was nil", featureId)
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				if strings.EqualFold(registrationState, Registered) {
					return metadata.ResourceRequiresImport(r.ResourceType(), featureId)
				}
			}

			if strings.EqualFold(registrationState, Pending) {
				return fmt.Errorf("%s which requires manual approval can not be managed by Terraform", featureId)
			}

			resp, err := client.Register(ctx, featureId)
			if err != nil {
				return fmt.Errorf("registering %s: %+v", featureId, err)
			}

			if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.State != nil {
				if strings.EqualFold(*resp.Model.Properties.State, Pending) {
					return fmt.Errorf("%s which requires manual approval can not be managed by terraform", featureId)
				}
			}

			metadata.SetID(featureId)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &featureId); err != nil {
				return err
			}

			pollerType := custompollers.NewResourceProviderFeatureRegistrationPoller(client, featureId, Registered)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for registration of %s: %+v", featureId, err)
			}

			return nil
		},

		Timeout: 30 * time.Minute,
	}
}
