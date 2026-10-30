// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_provider_registration

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.ResourceProvidersClient
			account := metadata.Client.Account

			var obj ResourceProviderRegistrationModel
			if err := metadata.Decode(&obj); err != nil {
				return err
			}

			id, err := providers.ParseSubscriptionProviderID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			if err := r.checkIfManagedByTerraform(id.ProviderName, account); err != nil {
				return err
			}

			provider, err := client.Get(ctx, *id, providers.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			registrationState := ""
			if model := provider.Model; model != nil && model.RegistrationState != nil {
				registrationState = *model.RegistrationState
			}
			if registrationState == "" {
				return fmt.Errorf("retrieving %s: `registrationState` was nil", *id)
			}

			if !strings.EqualFold(registrationState, "Registered") {
				// Account for API inconsistency before erroring
				pollerType := custompollers.NewResourceProviderExistencePoller(client, *id, 10)
				poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
				if err := poller.PollUntilDone(ctx); err != nil {
					if errors.As(err, &pollers.PollingCancelledError{}) {
						return fmt.Errorf("retrieving %s: expected `registrationState` to be `%s` but got `%s`", *id, Registered, registrationState)
					}
					return fmt.Errorf("polling %s: %+v", id, err)
				}
			}

			if metadata.ResourceData.HasChange("feature") {
				oldFeaturesRaw, newFeaturesRaw := metadata.ResourceData.GetChange("feature")
				if err := r.applyFeatures(ctx, metadata, *id, oldFeaturesRaw.(*pluginsdk.Set).List(), newFeaturesRaw.(*pluginsdk.Set).List()); err != nil {
					return fmt.Errorf("applying features for %s: %+v", *id, err)
				}
			}

			log.Printf("[DEBUG] Registering %s..", *id)
			payload := providers.ProviderRegistrationRequest{}
			if _, err := client.Register(ctx, *id, payload); err != nil {
				return fmt.Errorf("registering %s: %+v", *id, err)
			}

			log.Printf("[DEBUG] Waiting for %s to finish registering..", id)
			pollerType := custompollers.NewResourceProviderRegistrationPollerDefault(client, *id, Registered)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be registered: %s", id, err)
			}
			log.Printf("[DEBUG] Registered Resource Provider %q.", id)

			return nil
		},
		Timeout: 120 * time.Minute,
	}
}
