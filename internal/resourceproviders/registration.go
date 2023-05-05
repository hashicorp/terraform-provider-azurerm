package resourceproviders

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
)

func EnsureRegistered(ctx context.Context, client *providers.ProvidersClient, subscriptionId commonids.SubscriptionId, requiredRPs map[string]struct{}) error {
	if registeredResourceProviders == nil || unregisteredResourceProviders == nil {
		if err := populateCache(ctx, client, subscriptionId); err != nil {
			return fmt.Errorf("populating Resource Provider cache: %+v", err)
		}
	}

	log.Printf("[DEBUG] Determining which Resource Providers require Registration")
	providersToRegister, err := DetermineWhichRequiredResourceProvidersRequireRegistration(requiredRPs)
	if err != nil {
		return fmt.Errorf("determining which Required Resource Providers require registration: %+v", err)
	}

	if len(*providersToRegister) > 0 {
		log.Printf("[DEBUG] Registering %d Resource Providers", len(*providersToRegister))
		if err := registerForSubscription(ctx, client, subscriptionId, *providersToRegister); err != nil {
			return err
		}
	} else {
		log.Printf("[DEBUG] All required Resource Providers are registered")
	}

	return nil
}

// registerForSubscription registers the specified Resource Providers in the current Subscription
func registerForSubscription(ctx context.Context, client *providers.ProvidersClient, subscriptionId commonids.SubscriptionId, providersToRegister []string) error {
	var err error
	var failedProviders []string
	var wg sync.WaitGroup
	wg.Add(len(providersToRegister))

	for _, providerName := range providersToRegister {
		go func(p string) {
			defer wg.Done()
			log.Printf("[DEBUG] Registering Resource Provider %q with namespace", p)
			if innerErr := registerWithSubscription(ctx, client, subscriptionId, p); innerErr != nil {
				failedProviders = append(failedProviders, p)
				if err == nil {
					err = innerErr
				} else {
					err = fmt.Errorf("%s\n%s", err, innerErr)
				}
			}
		}(providerName)
	}

	wg.Wait()

	if len(failedProviders) > 0 {
		err = fmt.Errorf("Cannot register providers: %s. Errors were: %s", strings.Join(failedProviders, ", "), err)
	}
	return err
}

func registerWithSubscription(ctx context.Context, client *providers.ProvidersClient, subscriptionId commonids.SubscriptionId, providerName string) error {
	providerId := providers.NewSubscriptionProviderID(subscriptionId.SubscriptionId, providerName)
	if _, err := client.Register(ctx, providerId, providers.ProviderRegistrationRequest{}); err != nil {
		return fmt.Errorf("Cannot register provider %s with Azure Resource Manager: %s.", providerName, err)
	}

	return nil
}
