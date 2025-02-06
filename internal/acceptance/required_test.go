// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acceptance

import (
	"context"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders"
)

// since this depends on GetAuthConfig which lives in this package
// unfortunately this has to live in a different package to the other func

func TestAccEnsureRequiredResourceProvidersAreRegistered(t *testing.T) {
	config := GetAuthConfig(t)
	if config == nil {
		return
	}

	builder := clients.ClientBuilder{
		AuthConfig:                  config,
		DisableCorrelationRequestID: true,
		DisableTerraformPartnerID:   false,
		PartnerID:                   "",
		SubscriptionID:              os.Getenv("ARM_SUBSCRIPTION_ID"),
		TerraformVersion:            "0.0.0",
	}
	armClient, err := clients.Build(context.Background(), builder)
	if err != nil {
		t.Fatalf("Error building ARM Client: %+v", err)
	}

	client := armClient.Resource.ResourceProvidersClient
	ctx := armClient.StopContext

	requiredResourceProviders := resourceproviders.Legacy()
	subscriptionId := commonids.NewSubscriptionID(armClient.Account.SubscriptionId)

	if err = resourceproviders.EnsureRegistered(ctx, client, subscriptionId, requiredResourceProviders); err != nil {
		t.Fatalf("Error registering Resource Providers: %+v", err)
	}

	// refresh the cache now things have been re-registered
	resourceproviders.ClearCache()
	if err := resourceproviders.CacheSupportedProviders(ctx, client, subscriptionId); err != nil {
		t.Fatalf("re-caching Resource Providers: %+v", err)
	}

	stillRequiringRegistration, err := resourceproviders.DetermineWhichRequiredResourceProvidersRequireRegistration(requiredResourceProviders)
	if err != nil {
		t.Fatalf("determining which Resource Providers still require Registration: %+v", err)
	}
	if len(*stillRequiringRegistration) > 0 {
		t.Fatalf("'%d' Resource Providers are still Pending Registration: %s", len(*stillRequiringRegistration), spew.Sprint(stillRequiringRegistration))
	}
}
