package azure

import (
	"fmt"
	"net/url"
	"strings"
)

// ResourceID represents a parsed long-form Azure Resource Manager ID
// with the Subscription ID, Resource Group and the Provider as top-
// level fields, and other key-value pairs available via a map in the
// Path field.
type ResourceID struct {
	SubscriptionID string
	ResourceGroup  string
	Provider       string
	Path           map[string]string
}

// ParseAzureResourceID converts a long-form Azure Resource Manager ID
// into a ResourceID. We make assumptions about the structure of URLs,
// which is obviously not good, but the best thing available given the
// SDK.
func ParseAzureResourceID(id string) (*ResourceID, error) {
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse Azure ID: %s", err)
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	// We should have an even number of key-value pairs.
	if len(components)%2 != 0 {
		return nil, fmt.Errorf("The number of path segments is not divisible by 2 in %q", path)
	}

	var subscriptionID string

	// Put the constituent key-value pairs into a map
	componentMap := make(map[string]string, len(components)/2)
	for current := 0; current < len(components); current += 2 {
		key := components[current]
		value := components[current+1]

		// Check key/value for empty strings.
		if key == "" || value == "" {
			return nil, fmt.Errorf("Key/Value cannot be empty strings. Key: '%s', Value: '%s'", key, value)
		}

		// Catch the subscriptionID before it can be overwritten by another "subscriptions"
		// value in the ID which is the case for the Service Bus subscription resource
		if key == "subscriptions" && subscriptionID == "" {
			subscriptionID = value
		} else {
			componentMap[key] = value
		}
	}

	// Build up a TargetResourceID from the map
	idObj := &ResourceID{}
	idObj.Path = componentMap

	if subscriptionID != "" {
		idObj.SubscriptionID = subscriptionID
	} else {
		return nil, fmt.Errorf("No subscription ID found in: %q", path)
	}

	if resourceGroup, ok := componentMap["resourceGroups"]; ok {
		idObj.ResourceGroup = resourceGroup
		delete(componentMap, "resourceGroups")
	} else {
		// Some Azure APIs are weird and provide things in lower case...
		// However it's not clear whether the casing of other elements in the URI
		// matter, so we explicitly look for that case here.
		if resourceGroup, ok := componentMap["resourcegroups"]; ok {
			idObj.ResourceGroup = resourceGroup
			delete(componentMap, "resourcegroups")
		}
	}

	// It is OK not to have a provider in the case of a resource group
	if provider, ok := componentMap["providers"]; ok {
		idObj.Provider = provider
		delete(componentMap, "providers")
	}

	return idObj, nil
}

// PopSegment retrieves a segment from the Path and returns it
// if found it removes it from the Path then return the value
// if not found, this returns nil
func (id *ResourceID) PopSegment(name string) (string, error) {
	val, ok := id.Path[name]
	if !ok {
		return "", fmt.Errorf("ID was missing the `%s` element", name)
	}

	delete(id.Path, name)
	return val, nil
}

// ValidateNoEmptySegments validates ...
func (id *ResourceID) ValidateNoEmptySegments(sourceId string) error {
	if len(id.Path) == 0 {
		return nil
	}

	return fmt.Errorf("ID contained more segments than required: %q", sourceId)
}
