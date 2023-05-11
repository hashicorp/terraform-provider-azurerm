package configurationstores

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ConfigurationStoreId{}

// ConfigurationStoreId is a struct representing the Resource ID for a Configuration Store
type ConfigurationStoreId struct {
	SubscriptionId         string
	ResourceGroupName      string
	ConfigurationStoreName string
}

// NewConfigurationStoreID returns a new ConfigurationStoreId struct
func NewConfigurationStoreID(subscriptionId string, resourceGroupName string, configurationStoreName string) ConfigurationStoreId {
	return ConfigurationStoreId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		ConfigurationStoreName: configurationStoreName,
	}
}

// ParseConfigurationStoreID parses 'input' into a ConfigurationStoreId
func ParseConfigurationStoreID(input string) (*ConfigurationStoreId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConfigurationStoreId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConfigurationStoreId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ConfigurationStoreName, ok = parsed.Parsed["configurationStoreName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "configurationStoreName", *parsed)
	}

	return &id, nil
}

// ParseConfigurationStoreIDInsensitively parses 'input' case-insensitively into a ConfigurationStoreId
// note: this method should only be used for API response data and not user input
func ParseConfigurationStoreIDInsensitively(input string) (*ConfigurationStoreId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConfigurationStoreId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConfigurationStoreId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ConfigurationStoreName, ok = parsed.Parsed["configurationStoreName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "configurationStoreName", *parsed)
	}

	return &id, nil
}

// ValidateConfigurationStoreID checks that 'input' can be parsed as a Configuration Store ID
func ValidateConfigurationStoreID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConfigurationStoreID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Configuration Store ID
func (id ConfigurationStoreId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppConfiguration/configurationStores/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ConfigurationStoreName)
}

// Segments returns a slice of Resource ID Segments which comprise this Configuration Store ID
func (id ConfigurationStoreId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppConfiguration", "Microsoft.AppConfiguration", "Microsoft.AppConfiguration"),
		resourceids.StaticSegment("staticConfigurationStores", "configurationStores", "configurationStores"),
		resourceids.UserSpecifiedSegment("configurationStoreName", "configurationStoreValue"),
	}
}

// String returns a human-readable description of this Configuration Store ID
func (id ConfigurationStoreId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Configuration Store Name: %q", id.ConfigurationStoreName),
	}
	return fmt.Sprintf("Configuration Store (%s)", strings.Join(components, "\n"))
}
