// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourceproviders

import (
	"fmt"
	"log"
)

// DetermineWhichRequiredResourceProvidersRequireRegistration determines which Resource Providers require registration to be able to be used
func DetermineWhichRequiredResourceProvidersRequireRegistration(requiredResourceProviders ResourceProviders) (*[]string, error) {
	if registeredResourceProviders == nil || unregisteredResourceProviders == nil {
		return nil, fmt.Errorf("internal-error: the registered/unregistered Resource Provider cache isn't populated")
	}

	requiringRegistration := make([]string, 0)
	for providerName := range requiredResourceProviders {
		if _, isRegistered := registeredResourceProviders[providerName]; isRegistered {
			continue
		}

		if _, isUnregistered := unregisteredResourceProviders[providerName]; !isUnregistered {
			// some RPs may not exist in some non-public clouds, so we'll log a warning here instead of raising an error
			log.Printf("[WARN] The required Resource Provider %q wasn't returned from the Azure API", providerName)
			continue
		}

		requiringRegistration = append(requiringRegistration, providerName)
	}

	return &requiringRegistration, nil
}
