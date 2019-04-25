package resourceproviders

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
)

// DetermineResourceProvidersRequiringRegistration determines which Resource Providers require registration to be able to be used
func DetermineResourceProvidersRequiringRegistration(availableResourceProviders []resources.Provider, requiredResourceProviders map[string]struct{}) map[string]struct{} {
	providers := requiredResourceProviders

	// filter out any providers already registered
	for _, p := range availableResourceProviders {
		if _, ok := providers[*p.Namespace]; !ok {
			continue
		}

		if strings.ToLower(*p.RegistrationState) == "registered" {
			log.Printf("[DEBUG] Skipping provider registration for namespace %s\n", *p.Namespace)
			delete(providers, *p.Namespace)
		}
	}

	return providers
}

// RegisterForSubscription registers the specified Resource Providers in the current Subscription
func RegisterForSubscription(ctx context.Context, client resources.ProvidersClient, providersToRegister map[string]struct{}) error {
	var err error
	var wg sync.WaitGroup
	wg.Add(len(providersToRegister))

	for providerName := range providersToRegister {
		go func(p string) {
			defer wg.Done()
			log.Printf("[DEBUG] Registering Resource Provider %q with namespace", p)
			if innerErr := registerWithSubscription(ctx, p, client); innerErr != nil {
				err = innerErr
			}
		}(providerName)
	}

	wg.Wait()

	return err
}

func registerWithSubscription(ctx context.Context, providerName string, client resources.ProvidersClient) error {
	if _, err := client.Register(ctx, providerName); err != nil {
		return fmt.Errorf("Cannot register provider %s with Azure Resource Manager: %s.", providerName, err)
	}

	return nil
}
