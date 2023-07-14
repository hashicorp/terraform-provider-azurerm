// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

// specific parser to prevent decoding of the ID
func parseAzureResourceID(id string) (*azure.ResourceID, error) {
	id = strings.TrimPrefix(id, "/")
	id = strings.TrimSuffix(id, "/")

	components := strings.Split(id, "/")

	// We should have an even number of key-value pairs.
	if len(components)%2 != 0 {
		return nil, fmt.Errorf("the number of path segments is not divisible by 2 in %q", id)
	}

	var subscriptionID string
	var provider string

	// Put the constituent key-value pairs into a map
	componentMap := make(map[string]string, len(components)/2)
	for current := 0; current < len(components); current += 2 {
		key := components[current]
		value, err := url.QueryUnescape(components[current+1])
		if err != nil {
			return nil, fmt.Errorf("parsing value of id components %q: %+v", value, err)
		}

		// Check key/value for empty strings.
		if key == "" || value == "" {
			return nil, fmt.Errorf("Key/Value cannot be empty strings. Key: '%s', Value: '%s'", key, value)
		}

		switch {
		case key == "subscriptions" && subscriptionID == "":
			// Catch the subscriptionID before it can be overwritten by another "subscriptions"
			// value in the ID which is the case for the Service Bus subscription resource
			subscriptionID = value
		case key == "providers" && provider == "":
			// Catch the provider before it can be overwritten by another "providers"
			// value in the ID which can be the case for the Role Assignment resource
			provider = value
		default:
			componentMap[key] = value
		}
	}

	// Build up a TargetResourceID from the map
	idObj := &azure.ResourceID{}
	idObj.Path = componentMap

	if subscriptionID != "" {
		idObj.SubscriptionID = subscriptionID
	} else {
		return nil, fmt.Errorf("no subscription ID found in: %q", id)
	}

	if resourceGroup, ok := componentMap["resourceGroups"]; ok {
		idObj.ResourceGroup = resourceGroup
		delete(componentMap, "resourceGroups")
	} else if resourceGroup, ok := componentMap["resourcegroups"]; ok {
		// Some Azure APIs are weird and provide things in lower case...
		// However it's not clear whether the casing of other elements in the URI
		// matter, so we explicitly look for that case here.
		idObj.ResourceGroup = resourceGroup
		delete(componentMap, "resourcegroups")
	}

	if provider != "" {
		idObj.Provider = provider
	}

	if secondaryProvider := componentMap["providers"]; secondaryProvider != "" {
		idObj.SecondaryProvider = secondaryProvider
		delete(componentMap, "providers")
	}

	return idObj, nil
}
