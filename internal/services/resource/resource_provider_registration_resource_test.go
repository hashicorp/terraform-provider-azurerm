// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/testclient"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// NOTE: this can be moved up a level when all the others are

type ResourceProviderRegistrationResource struct{}

func TestAccResourceProviderRegistration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_provider_registration", "test")
	r := ResourceProviderRegistrationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic("Microsoft.BlockchainTokens"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceProviderRegistration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_provider_registration", "test")
	r := ResourceProviderRegistrationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic("Microsoft.ApiCenter"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport("Microsoft.ApiCenter")
		}),
	})
}

func TestAccResourceProviderRegistration_feature(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_provider_registration", "test")
	r := ResourceProviderRegistrationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			PreConfig: func() {
				// Last error may cause resource provider still in `Registered` status.Need to unregister it before a new test.
				if err := r.unRegisterProviders("Microsoft.ApiSecurity"); err != nil {
					t.Fatalf("Failed to reset feature registration with error: %+v", err)
				}
			},
			Config: r.multiFeature(true, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiFeature(true, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiFeature(false, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiFeature(false, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ResourceProviderRegistrationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := providers.ParseSubscriptionProviderID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Resource.ResourceProvidersClient.Get(ctx, *id, providers.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	isRegistered := false
	if model := resp.Model; model != nil && model.RegistrationState != nil {
		isRegistered = strings.EqualFold(*model.RegistrationState, "Registered")
	}
	return pointer.To(isRegistered), nil
}

func (ResourceProviderRegistrationResource) basic(name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  skip_provider_registration = true
}

resource "azurerm_resource_provider_registration" "test" {
  name = %q
  lifecycle {
    ignore_changes = [feature]
  }
}
`, name)
}

func (r ResourceProviderRegistrationResource) requiresImport(name string) string {
	template := r.basic(name)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_provider_registration" "import" {
  name = azurerm_resource_provider_registration.test.name
}
`, template)
}

func (r ResourceProviderRegistrationResource) unRegisterProviders(resourceProviders ...string) error {
	client, err := testclient.Build()
	if err != nil {
		return fmt.Errorf("building client: %+v", err)
	}

	for _, rp := range resourceProviders {
		if err = r.unRegisterProvider(client, rp); err != nil {
			return err
		}
	}

	return nil
}

func (r ResourceProviderRegistrationResource) unRegisterProvider(client *clients.Client, resourceProvider string) error {
	ctx, cancel := context.WithDeadline(client.StopContext, time.Now().Add(30*time.Minute))
	defer cancel()

	providersClient := client.Resource.ResourceProvidersClient
	id := providers.NewSubscriptionProviderID(client.Account.SubscriptionId, resourceProvider)
	provider, err := providersClient.Get(ctx, id, providers.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	registrationState := ""
	if model := provider.Model; model != nil && model.RegistrationState != nil {
		registrationState = *model.RegistrationState
	}
	if registrationState == "" {
		return fmt.Errorf("retrieving Resource Provider %q: `registrationState` was nil", resourceProvider)
	}

	if !strings.EqualFold(registrationState, "Registered") {
		return nil
	}

	if _, err := providersClient.Unregister(ctx, id); err != nil {
		return fmt.Errorf("unregistering %s: %+v", id, err)
	}

	pollerType := custompollers.NewResourceProviderUnregistrationPoller(providersClient, id)
	poller := pollers.NewPoller(pollerType, 10*time.Minute, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become unregistered: %+v", id, err)
	}

	return nil
}

func (ResourceProviderRegistrationResource) multiFeature(registered1 bool, registered2 bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  skip_provider_registration = true
}

resource "azurerm_resource_provider_registration" "test" {
  name = "Microsoft.ApiSecurity"
  feature {
    name       = "PP2CanaryAccessDEV"
    registered = %t
  }
  feature {
    name       = "PP3CanaryAccessDEV"
    registered = %t
  }
}
`, registered1, registered2)
}
