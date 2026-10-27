// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_provider_registration

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.ResourceProvidersClient
			account := metadata.Client.Account

			id, err := providers.ParseSubscriptionProviderID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := r.checkIfManagedByTerraform(id.ProviderName, account); err != nil {
				return err
			}

			if err = r.applyFeatures(ctx, metadata, *id, metadata.ResourceData.Get("feature").(*pluginsdk.Set).List(), make([]any, 0)); err != nil {
				return fmt.Errorf("applying features for %s: %+v", *id, err)
			}

			if _, err := client.Unregister(ctx, *id); err != nil {
				return fmt.Errorf("unregistering Resource Provider %q: %+v", *id, err)
			}

			log.Printf("[DEBUG] Waiting for %s to finish unregistering..", *id)
			pollerType := custompollers.NewResourceProviderRegistrationPollerDefault(client, *id, Unregistered)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become unregistered: %+v", *id, err)
			}
			log.Printf("[DEBUG] Unregistered Resource Provider %q.", *id)

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
