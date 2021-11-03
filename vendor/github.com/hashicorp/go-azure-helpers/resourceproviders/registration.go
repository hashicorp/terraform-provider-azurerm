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
	providers := make(map[string]struct{})

	// filter out any providers already registered and not in the required list.
	for _, p := range availableResourceProviders {
		// Skip it if it's not in the required list.
		if _, ok := requiredResourceProviders[*p.Namespace]; !ok {
			continue
		}

		// If it's in the required list but not registered.
		if strings.ToLower(*p.RegistrationState) != "registered" {
			log.Printf("[DEBUG] Adding provider registration for namespace %s\n", *p.Namespace)
			providers[*p.Namespace] = requiredResourceProviders[*p.Namespace]
		}
	}

	return providers
}

// RegisterForSubscription registers the specified Resource Providers in the current Subscription
func RegisterForSubscription(ctx context.Context, client resources.ProvidersClient, providersToRegister map[string]struct{}) error {
	var err error
	var failedProviders []string
	var wg sync.WaitGroup
	wg.Add(len(providersToRegister))

	for providerName := range providersToRegister {
		go func(p string) {
			defer wg.Done()
			log.Printf("[DEBUG] Registering Resource Provider %q with namespace", p)
			if innerErr := registerWithSubscription(ctx, p, client); innerErr != nil {
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
		err = fmt.Errorf("Cannnot register providers: %s. Errors were: %s", strings.Join(failedProviders, ", "), err)
	}
	return err
}

func registerWithSubscription(ctx context.Context, providerName string, client resources.ProvidersClient) error {
	if _, err := client.Register(ctx, providerName); err != nil {
		return fmt.Errorf("Cannot register provider %s with Azure Resource Manager: %s.", providerName, err)
	}

	return nil
}
