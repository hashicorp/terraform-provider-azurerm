// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourceproviders

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-09-01/providers"
)

// cachedResourceProviders can be (validly) nil - as such this shouldn't be relied on
var cachedResourceProviders *[]string
var registeredResourceProviders *map[string]struct{}
var unregisteredResourceProviders *map[string]struct{}

var cacheLock = &sync.Mutex{}

// CacheSupportedProviders attempts to retrieve the supported Resource Providers from the Resource Manager API
// and caches them, for used in enhanced validation
func CacheSupportedProviders(ctx context.Context, client *providers.ProvidersClient, subscriptionId commonids.SubscriptionId) error {
	// already populated
	if cachedResourceProviders != nil {
		return nil
	}

	if err := populateCache(ctx, client, subscriptionId); err != nil {
		return fmt.Errorf("populating cache: %+v", err)
	}

	return nil
}

func ClearCache() {
	cacheLock.Lock()
	cachedResourceProviders = nil
	registeredResourceProviders = nil
	unregisteredResourceProviders = nil
	cacheLock.Unlock()
}

func populateCache(ctx context.Context, client *providers.ProvidersClient, subscriptionId commonids.SubscriptionId) error {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	providers, err := client.ListComplete(ctx, subscriptionId, providers.DefaultListOperationOptions())
	if err != nil {
		return fmt.Errorf("listing Resource Providers: %+v", err)
	}

	providerNames := make([]string, 0)
	registeredProviders := make(map[string]struct{}, 0)
	unregisteredProviders := make(map[string]struct{}, 0)
	for _, provider := range providers.Items {
		if provider.Namespace == nil {
			continue
		}

		providerNames = append(providerNames, *provider.Namespace)
		registered := provider.RegistrationState != nil && strings.EqualFold(*provider.RegistrationState, "registered")
		if registered {
			registeredProviders[*provider.Namespace] = struct{}{}
		} else {
			unregisteredProviders[*provider.Namespace] = struct{}{}
		}
	}

	cachedResourceProviders = &providerNames
	registeredResourceProviders = &registeredProviders
	unregisteredResourceProviders = &unregisteredProviders
	return nil
}
