package resourceproviders

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2017-03-09/resources/mgmt/resources"
)

// cachedResourceProviders can be (validly) nil - as such this shouldn't be relied on
var cachedResourceProviders *[]resources.Provider

// CachedSupportedProviders returns cached the supported Resource Providers, if not cached it attempts to retrieve from the Resource Manager API
// and caches them, for used in enhanced validation
func CachedSupportedProviders(ctx context.Context, client *resources.ProvidersClient) (*[]resources.Provider, error) {
	if cachedResourceProviders != nil {
		return cachedResourceProviders, nil
	}
	providers, err := availableResourceProviders(ctx, client)
	if err != nil {
		return nil, err
	}
	cached := make([]resources.Provider, 0)
	for _, provider := range providers {
		cached = append(cached, resources.Provider{
			Namespace:         provider.Namespace,
			RegistrationState: provider.RegistrationState,
		})
	}
	cachedResourceProviders = &cached
	return cachedResourceProviders, nil
}
