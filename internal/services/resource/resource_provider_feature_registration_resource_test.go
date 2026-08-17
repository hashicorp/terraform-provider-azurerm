// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2021-07-01/features"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/testclient"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ResourceProviderFeatureRegistrationResource struct{}

func TestAccResourceProviderFeatureRegistration_basic(t *testing.T) {
	if os.Getenv("ARM_SUBSCRIPTION_ID_ALT") == "" {
		t.Skip("Skipping as `ARM_SUBSCRIPTION_ID_ALT` was not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_resource_provider_feature_registration", "test")
	r := ResourceProviderFeatureRegistrationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			PreConfig: func() {
				// Ensure the feature is not currently registered
				if err := r.unRegisterFeature(data, "EncryptionAtHost", "Microsoft.Compute"); err != nil {
					t.Fatalf("unregistering feature: %+v", err)
				}
			},
			Config: r.basic(data, "EncryptionAtHost", "Microsoft.Compute"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceProviderFeatureRegistration_requiresImport(t *testing.T) {
	if os.Getenv("ARM_SUBSCRIPTION_ID_ALT") == "" {
		t.Skip("Skipping as `ARM_SUBSCRIPTION_ID_ALT` was not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_resource_provider_feature_registration", "test")
	r := ResourceProviderFeatureRegistrationResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			PreConfig: func() {
				// Ensure the feature is not currently registered
				if err := r.unRegisterFeature(data, "InGuestScheduledPatchVMPreview", "Microsoft.Compute"); err != nil {
					t.Fatalf("unregistering feature: %+v", err)
				}
			},
			Config: r.basic(data, "InGuestScheduledPatchVMPreview", "Microsoft.Compute"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport(data, "InGuestScheduledPatchVMPreview", "Microsoft.Compute")
		}),
	})
}

func (ResourceProviderFeatureRegistrationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := features.ParseFeatureID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Resource.FeaturesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	isRegistered := false
	if model := resp.Model; model != nil && model.Properties != nil && model.Properties.State != nil {
		isRegistered = strings.EqualFold(*model.Properties.State, resource.Registered)
	}

	return pointer.To(isRegistered), nil
}

func (ResourceProviderFeatureRegistrationResource) basic(data acceptance.TestData, name, providerName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}

  resource_provider_registrations = "none"
  resource_providers_to_register  = [%[2]q]

  // Run this test in the alternate subscription to prevent unregistering features required by other tests
  subscription_id = %[3]q
}

resource "azurerm_resource_provider_feature_registration" "test" {
  name          = %[1]q
  provider_name = %[2]q
}
`, name, providerName, data.Subscriptions.Secondary)
}

func (r ResourceProviderFeatureRegistrationResource) basicForResourceIdentity(data acceptance.TestData) string {
	return r.basic(data, "AutomaticZoneRebalancing", "Microsoft.Compute")
}

func (r ResourceProviderFeatureRegistrationResource) requiresImport(data acceptance.TestData, name, providerName string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_resource_provider_feature_registration" "import" {
  name          = azurerm_resource_provider_feature_registration.test.name
  provider_name = azurerm_resource_provider_feature_registration.test.provider_name
}
`, r.basic(data, name, providerName))
}

func (r ResourceProviderFeatureRegistrationResource) unRegisterFeature(data acceptance.TestData, name, providerName string) error {
	client, err := testclient.Build()
	if err != nil {
		return fmt.Errorf("building client: %+v", err)
	}

	ctx, cancel := context.WithDeadline(client.StopContext, time.Now().Add(30*time.Minute))
	defer cancel()

	featuresClient := client.Resource.FeaturesClient
	id := features.NewFeatureID(data.Subscriptions.Secondary, providerName, name)
	feature, err := featuresClient.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(feature.HttpResponse) {
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	registrationState := ""
	if model := feature.Model; model != nil && feature.Model != nil && feature.Model.Properties != nil && feature.Model.Properties.State != nil {
		registrationState = *feature.Model.Properties.State
	}

	if registrationState == "" {
		return fmt.Errorf("retrieving %s: unable to retrieve registration state", id)
	}

	if !strings.EqualFold(registrationState, resource.Registered) {
		return nil
	}

	if _, err := featuresClient.Unregister(ctx, id); err != nil {
		return fmt.Errorf("unregistering %s: %+v", id, err)
	}

	pollerType := custompollers.NewResourceProviderFeatureRegistrationPoller(featuresClient, id, resource.Unregistered)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for registration state of %s: %+v", id, err)
	}

	return nil
}
