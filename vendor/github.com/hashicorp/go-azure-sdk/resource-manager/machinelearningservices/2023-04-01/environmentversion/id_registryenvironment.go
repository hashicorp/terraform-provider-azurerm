package environmentversion

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RegistryEnvironmentId{}

// RegistryEnvironmentId is a struct representing the Resource ID for a Registry Environment
type RegistryEnvironmentId struct {
	SubscriptionId    string
	ResourceGroupName string
	RegistryName      string
	EnvironmentName   string
}

// NewRegistryEnvironmentID returns a new RegistryEnvironmentId struct
func NewRegistryEnvironmentID(subscriptionId string, resourceGroupName string, registryName string, environmentName string) RegistryEnvironmentId {
	return RegistryEnvironmentId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RegistryName:      registryName,
		EnvironmentName:   environmentName,
	}
}

// ParseRegistryEnvironmentID parses 'input' into a RegistryEnvironmentId
func ParseRegistryEnvironmentID(input string) (*RegistryEnvironmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(RegistryEnvironmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RegistryEnvironmentId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.EnvironmentName, ok = parsed.Parsed["environmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "environmentName", *parsed)
	}

	return &id, nil
}

// ParseRegistryEnvironmentIDInsensitively parses 'input' case-insensitively into a RegistryEnvironmentId
// note: this method should only be used for API response data and not user input
func ParseRegistryEnvironmentIDInsensitively(input string) (*RegistryEnvironmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(RegistryEnvironmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RegistryEnvironmentId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.RegistryName, ok = parsed.Parsed["registryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "registryName", *parsed)
	}

	if id.EnvironmentName, ok = parsed.Parsed["environmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "environmentName", *parsed)
	}

	return &id, nil
}

// ValidateRegistryEnvironmentID checks that 'input' can be parsed as a Registry Environment ID
func ValidateRegistryEnvironmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRegistryEnvironmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Registry Environment ID
func (id RegistryEnvironmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MachineLearningServices/registries/%s/environments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RegistryName, id.EnvironmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Registry Environment ID
func (id RegistryEnvironmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMachineLearningServices", "Microsoft.MachineLearningServices", "Microsoft.MachineLearningServices"),
		resourceids.StaticSegment("staticRegistries", "registries", "registries"),
		resourceids.UserSpecifiedSegment("registryName", "registryValue"),
		resourceids.StaticSegment("staticEnvironments", "environments", "environments"),
		resourceids.UserSpecifiedSegment("environmentName", "environmentValue"),
	}
}

// String returns a human-readable description of this Registry Environment ID
func (id RegistryEnvironmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Registry Name: %q", id.RegistryName),
		fmt.Sprintf("Environment Name: %q", id.EnvironmentName),
	}
	return fmt.Sprintf("Registry Environment (%s)", strings.Join(components, "\n"))
}
