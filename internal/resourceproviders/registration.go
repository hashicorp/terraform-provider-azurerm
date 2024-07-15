// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourceproviders

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceproviders/custompollers"
)

// EnsureRegistered tries to determine whether all requiredRPs are registered in the subscription, and attempts to
// register them if it appears they are not. Note that this may fail if a resource provider is not available in the
// current cloud environment (a warning message will be logged to indicate when a resource provider is not listed).
func EnsureRegistered(ctx context.Context, client *providers.ProvidersClient, subscriptionId commonids.SubscriptionId, requiredRPs ResourceProviders) error {
	// Cache supported resource providers if RP registration and enhanced validation are not both disabled
	if len(requiredRPs) == 0 && !features.EnhancedValidationEnabled() {
		log.Printf("[DEBUG] Skipping populating the resource provider cache, since resource provider registration and enhanced validation are both disabled")
		return nil
	}

	if cachedResourceProviders == nil || registeredResourceProviders == nil || unregisteredResourceProviders == nil {
		if err := populateCache(ctx, client, subscriptionId); err != nil {
			return fmt.Errorf("populating Resource Provider cache: %+v", err)
		}
	}

	log.Printf("[DEBUG] Determining which Resource Providers require Registration")
	providersToRegister, err := DetermineWhichRequiredResourceProvidersRequireRegistration(requiredRPs)
	if err != nil {
		return fmt.Errorf("determining which Resource Providers require registration: %+v", err)
	}

	if len(*providersToRegister) == 0 {
		log.Printf("[DEBUG] All required Resource Providers are registered")
		return nil
	}

	log.Printf("[DEBUG] Registering %d Resource Providers", len(*providersToRegister))
	if err = registerForSubscription(ctx, client, subscriptionId, *providersToRegister); err != nil {
		return userError(err)
	}

	return nil
}

// registerForSubscription registers the specified Resource Providers in the current Subscription
func registerForSubscription(ctx context.Context, client *providers.ProvidersClient, subscriptionId commonids.SubscriptionId, providersToRegister []string) error {
	errs := &registrationErrors{}
	var wg sync.WaitGroup
	wg.Add(len(providersToRegister))

	for _, providerName := range providersToRegister {
		go func(p string) {
			defer wg.Done()
			log.Printf("[DEBUG] Registering Resource Provider %q with namespace", p)
			if err := registerWithSubscription(ctx, client, subscriptionId, p); err != nil {
				errs.append(err)
			}
		}(providerName)
	}

	wg.Wait()

	if errs.hasErr() {
		return errs
	}

	return nil
}

func registerWithSubscription(ctx context.Context, client *providers.ProvidersClient, subscriptionId commonids.SubscriptionId, providerName string) error {
	providerId := providers.NewSubscriptionProviderID(subscriptionId.SubscriptionId, providerName)
	log.Printf("[DEBUG] Registering %s..", providerId)
	if resp, err := client.Register(ctx, providerId, providers.ProviderRegistrationRequest{}); err != nil {
		msg := fmt.Sprintf("registering resource provider %q: %s", providerName, err)

		if response.WasForbidden(resp.HttpResponse) {
			// a 403 response was received, so wrap ErrNoAuthorization in order to expose messaging for this
			return fmt.Errorf("%w: %s", ErrNoAuthorization, msg)
		}

		return errors.New(msg)
	}

	log.Printf("[DEBUG] Waiting for %s to finish registering..", providerId)
	pollerType := custompollers.NewResourceProviderRegistrationPoller(client, providerId)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be registered: %s", providerId, err)
	}

	log.Printf("[DEBUG] %s is registered.", providerId)

	return nil
}
