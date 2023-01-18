package resourceproviders

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
)

// cachedResourceProviders can be (validly) nil - as such this shouldn't be relied on
var cachedResourceProviders *[]string

// CacheSupportedProviders attempts to retrieve the supported Resource Providers from the Resource Manager API
// and caches them, for used in enhanced validation
func CacheSupportedProviders(cacheFunc CacheFunc) error {
	providers, err := cacheFunc()
	if err != nil {
		log.Printf("[DEBUG] error retrieving providers: %s. Enhanced validation will be unavailable", err)
		return err
	}
	providerNames := make([]string, 0)
	for _, provider := range providers {
		if provider.Namespace != nil {
			providerNames = append(providerNames, *provider.Namespace)
		}
	}
	cachedResourceProviders = &providerNames
	return nil
}

// CacheFunc provides an interface to cache resource providers
type CacheFunc func() ([]resources.Provider, error)

// DefaultCacheFunc lists all the available resource providers
func DefaultCacheFunc(ctx context.Context, client *resources.ProvidersClient) CacheFunc {
	return func() ([]resources.Provider, error) {
		return availableResourceProviders(ctx, client)
	}
}

// SavedResourceProvidersCacheFunc returns already saved resource providers
func SavedResourceProvidersCacheFunc(resourceProviders []resources.Provider) CacheFunc {
	return func() ([]resources.Provider, error) {
		return resourceProviders, nil
	}
}
