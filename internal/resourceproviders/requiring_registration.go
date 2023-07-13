// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resourceproviders

import (
	"fmt"
)

// DetermineWhichRequiredResourceProvidersRequireRegistration determines which Resource Providers require registration to be able to be used
func DetermineWhichRequiredResourceProvidersRequireRegistration(requiredResourceProviders map[string]struct{}) (*[]string, error) {
	if registeredResourceProviders == nil || unregisteredResourceProviders == nil {
		return nil, fmt.Errorf("internal-error: the registered/unregistered Resource Provider cache isn't populated")
	}

	requiringRegistration := make([]string, 0)
	for providerName := range requiredResourceProviders {
		if _, isRegistered := (*registeredResourceProviders)[providerName]; isRegistered {
			continue
		}

		if _, isUnregistered := (*unregisteredResourceProviders)[providerName]; !isUnregistered {
			// this is likely a typo in the Required Resource Providers list, so we should surface this
			return nil, fmt.Errorf("the required Resource Provider %q wasn't returned from the Azure API", providerName)
		}

		requiringRegistration = append(requiringRegistration, providerName)
	}

	return &requiringRegistration, nil
}
