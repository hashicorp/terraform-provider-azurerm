// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_provider_registration

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.ResourceWithUpdate         = Resource{}
	_ sdk.ResourceWithCustomImporter = Resource{}
)

type Resource struct{}

func (r Resource) ResourceType() string {
	return "azurerm_resource_provider_registration"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return providers.ValidateSubscriptionProviderID
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.Resource.ResourceProvidersClient
		account := metadata.Client.Account

		id, err := providers.ParseSubscriptionProviderID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		provider, err := client.Get(ctx, *id, providers.DefaultGetOperationOptions())
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		namespace := ""
		registrationState := ""
		if model := provider.Model; model != nil {
			if model.Namespace != nil {
				namespace = *model.Namespace
			}
			if model.RegistrationState != nil {
				registrationState = *model.RegistrationState
			}
		}
		if namespace != id.ProviderName {
			return fmt.Errorf("importing %s: expected %q but got %q", *id, id.ProviderName, namespace)
		}

		if !strings.EqualFold(registrationState, "Registered") {
			// Account for API inconsistency before erroring
			pollerType := custompollers.NewResourceProviderExistencePoller(client, *id, 10)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				if errors.As(err, &pollers.PollingCancelledError{}) {
					return fmt.Errorf("importing %s: Resource Provider must be registered to be imported", id.ProviderName)
				}
				return fmt.Errorf("polling %s: %+v", id, err)
			}
		}

		if err := r.checkIfManagedByTerraform(id.ProviderName, account); err != nil {
			return fmt.Errorf("importing %s: %+v", *id, err)
		}

		return nil
	}
}
