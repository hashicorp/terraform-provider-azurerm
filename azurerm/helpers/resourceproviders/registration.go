package resourceproviders

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
)

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

func RegisterForSubscription(ctx context.Context, client resources.ProvidersClient, providersToRegister map[string]struct{}) error {
	var err error
	var wg sync.WaitGroup
	wg.Add(len(providersToRegister))

	for providerName := range providersToRegister {
		go func(p string) {
			defer wg.Done()
			log.Printf("[DEBUG] Registering Resource Provider %q with namespace", p)
			if innerErr := registerWithSubscription(ctx, p, client); err != nil {
				err = innerErr
			}
		}(providerName)
	}

	wg.Wait()

	return err
}

func registerWithSubscription(ctx context.Context, providerName string, client resources.ProvidersClient) error {
	_, err := client.Register(ctx, providerName)
	if err != nil {
		return fmt.Errorf("Cannot register provider %s with Azure Resource Manager: %s.", providerName, err)
	}

	return nil
}
