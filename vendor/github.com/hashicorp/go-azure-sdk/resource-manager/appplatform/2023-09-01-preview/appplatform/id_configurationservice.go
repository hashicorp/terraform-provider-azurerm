package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ConfigurationServiceId{}

// ConfigurationServiceId is a struct representing the Resource ID for a Configuration Service
type ConfigurationServiceId struct {
	SubscriptionId           string
	ResourceGroupName        string
	SpringName               string
	ConfigurationServiceName string
}

// NewConfigurationServiceID returns a new ConfigurationServiceId struct
func NewConfigurationServiceID(subscriptionId string, resourceGroupName string, springName string, configurationServiceName string) ConfigurationServiceId {
	return ConfigurationServiceId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		SpringName:               springName,
		ConfigurationServiceName: configurationServiceName,
	}
}

// ParseConfigurationServiceID parses 'input' into a ConfigurationServiceId
func ParseConfigurationServiceID(input string) (*ConfigurationServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConfigurationServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConfigurationServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SpringName, ok = parsed.Parsed["springName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "springName", *parsed)
	}

	if id.ConfigurationServiceName, ok = parsed.Parsed["configurationServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "configurationServiceName", *parsed)
	}

	return &id, nil
}

// ParseConfigurationServiceIDInsensitively parses 'input' case-insensitively into a ConfigurationServiceId
// note: this method should only be used for API response data and not user input
func ParseConfigurationServiceIDInsensitively(input string) (*ConfigurationServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConfigurationServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConfigurationServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SpringName, ok = parsed.Parsed["springName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "springName", *parsed)
	}

	if id.ConfigurationServiceName, ok = parsed.Parsed["configurationServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "configurationServiceName", *parsed)
	}

	return &id, nil
}

// ValidateConfigurationServiceID checks that 'input' can be parsed as a Configuration Service ID
func ValidateConfigurationServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConfigurationServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Configuration Service ID
func (id ConfigurationServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/configurationServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.ConfigurationServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Configuration Service ID
func (id ConfigurationServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("springName", "springValue"),
		resourceids.StaticSegment("staticConfigurationServices", "configurationServices", "configurationServices"),
		resourceids.UserSpecifiedSegment("configurationServiceName", "configurationServiceValue"),
	}
}

// String returns a human-readable description of this Configuration Service ID
func (id ConfigurationServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Configuration Service Name: %q", id.ConfigurationServiceName),
	}
	return fmt.Sprintf("Configuration Service (%s)", strings.Join(components, "\n"))
}
