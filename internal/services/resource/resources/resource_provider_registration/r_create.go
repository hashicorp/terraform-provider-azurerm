// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_provider_registration

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Resource.ResourceProvidersClient
			account := metadata.Client.Account

			var obj ResourceProviderRegistrationModel
			if err := metadata.Decode(&obj); err != nil {
				return err
			}

			resourceId := providers.NewSubscriptionProviderID(account.SubscriptionId, obj.Name)
			if err := r.checkIfManagedByTerraform(resourceId.ProviderName, account); err != nil {
				return err
			}

			provider, err := client.Get(ctx, resourceId, providers.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(provider.HttpResponse) {
					return fmt.Errorf("%s was not found", resourceId)
				}

				return fmt.Errorf("retrieving %q: %+v", resourceId, err)
			}

			registrationState := ""
			if model := provider.Model; model != nil && model.RegistrationState != nil {
				registrationState = *model.RegistrationState
			}
			if registrationState == "" {
				return fmt.Errorf("retrieving %s: `registrationState` was nil", resourceId)
			}

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				if strings.EqualFold(registrationState, "Registered") {
					return metadata.ResourceRequiresImport(r.ResourceType(), resourceId)
				}
			}

			if metadata.ResourceData.HasChange("feature") {
				oldFeaturesRaw, newFeaturesRaw := metadata.ResourceData.GetChange("feature")
				if err := r.applyFeatures(ctx, metadata, resourceId, oldFeaturesRaw.(*pluginsdk.Set).List(), newFeaturesRaw.(*pluginsdk.Set).List()); err != nil {
					return fmt.Errorf("applying features for %q: %+v", resourceId, err)
				}
			}

			log.Printf("[DEBUG] Registering %s..", resourceId)
			if _, err := client.Register(ctx, resourceId, providers.ProviderRegistrationRequest{}); err != nil {
				return fmt.Errorf("registering %s: %+v", resourceId, err)
			}
			metadata.SetID(resourceId)

			log.Printf("[DEBUG] Waiting for %s to finish registering..", resourceId)
			pollerType := custompollers.NewResourceProviderRegistrationPollerDefault(client, resourceId, Registered)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for %s to be registered: %s", resourceId, err)
			}
			log.Printf("[DEBUG] Registered Resource Provider %q.", resourceId)

			return nil
		},

		Timeout: 120 * time.Minute,
	}
}
