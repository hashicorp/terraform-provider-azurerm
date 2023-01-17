package resourceproviders

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
)

// cachedResourceProviders can be (validly) nil - as such this shouldn't be relied on
var cachedResourceProviders *[]resources.Provider

// CacheSupportedProviders attempts to retrieve the supported Resource Providers from the Resource Manager API
// and caches them, for used in enhanced validation
func CacheSupportedProviders(ctx context.Context, client *resources.ProvidersClient) ([]resources.Provider, error) {
	if cachedResourceProviders != nil {
		return *cachedResourceProviders, nil
	}
	providers, err := availableResourceProviders(ctx, client)
	if err != nil {
		log.Printf("[DEBUG] error retrieving providers: %s. Enhanced validation will be unavailable", err)
		return nil, err
	}
	cachedResourceProviders = &providers
	return providers, nil
}
