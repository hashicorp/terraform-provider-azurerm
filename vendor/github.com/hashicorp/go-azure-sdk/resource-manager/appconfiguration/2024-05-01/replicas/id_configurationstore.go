package replicas

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ConfigurationStoreId{})
}

var _ resourceids.ResourceId = &ConfigurationStoreId{}

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
	parser := resourceids.NewParserFromResourceIdType(&ConfigurationStoreId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConfigurationStoreId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConfigurationStoreIDInsensitively parses 'input' case-insensitively into a ConfigurationStoreId
// note: this method should only be used for API response data and not user input
func ParseConfigurationStoreIDInsensitively(input string) (*ConfigurationStoreId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConfigurationStoreId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConfigurationStoreId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConfigurationStoreId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ConfigurationStoreName, ok = input.Parsed["configurationStoreName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "configurationStoreName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("configurationStoreName", "configurationStoreName"),
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
