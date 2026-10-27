// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_provider_feature_registration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2021-07-01/features"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/custompollers"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.FeaturesClient

			featureId, err := features.ParseFeatureID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Unregister(ctx, *featureId); err != nil {
				return fmt.Errorf("unregistering %s: %+v", featureId, err)
			}

			pollerType := custompollers.NewResourceProviderFeatureRegistrationPoller(client, *featureId, Unregistered)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for registration of %s: %+v", featureId, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
