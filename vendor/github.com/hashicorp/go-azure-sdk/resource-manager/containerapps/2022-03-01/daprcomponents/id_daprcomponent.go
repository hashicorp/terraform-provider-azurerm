package daprcomponents

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DaprComponentId{}

// DaprComponentId is a struct representing the Resource ID for a Dapr Component
type DaprComponentId struct {
	SubscriptionId         string
	ResourceGroupName      string
	ManagedEnvironmentName string
	DaprComponentName      string
}

// NewDaprComponentID returns a new DaprComponentId struct
func NewDaprComponentID(subscriptionId string, resourceGroupName string, managedEnvironmentName string, daprComponentName string) DaprComponentId {
	return DaprComponentId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		ManagedEnvironmentName: managedEnvironmentName,
		DaprComponentName:      daprComponentName,
	}
}

// ParseDaprComponentID parses 'input' into a DaprComponentId
func ParseDaprComponentID(input string) (*DaprComponentId, error) {
	parser := resourceids.NewParserFromResourceIdType(DaprComponentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DaprComponentId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ManagedEnvironmentName, ok = parsed.Parsed["managedEnvironmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedEnvironmentName", *parsed)
	}

	if id.DaprComponentName, ok = parsed.Parsed["daprComponentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "daprComponentName", *parsed)
	}

	return &id, nil
}

// ParseDaprComponentIDInsensitively parses 'input' case-insensitively into a DaprComponentId
// note: this method should only be used for API response data and not user input
func ParseDaprComponentIDInsensitively(input string) (*DaprComponentId, error) {
	parser := resourceids.NewParserFromResourceIdType(DaprComponentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DaprComponentId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ManagedEnvironmentName, ok = parsed.Parsed["managedEnvironmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "managedEnvironmentName", *parsed)
	}

	if id.DaprComponentName, ok = parsed.Parsed["daprComponentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "daprComponentName", *parsed)
	}

	return &id, nil
}

// ValidateDaprComponentID checks that 'input' can be parsed as a Dapr Component ID
func ValidateDaprComponentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDaprComponentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dapr Component ID
func (id DaprComponentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/managedEnvironments/%s/daprComponents/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName, id.DaprComponentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dapr Component ID
func (id DaprComponentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticManagedEnvironments", "managedEnvironments", "managedEnvironments"),
		resourceids.UserSpecifiedSegment("managedEnvironmentName", "managedEnvironmentValue"),
		resourceids.StaticSegment("staticDaprComponents", "daprComponents", "daprComponents"),
		resourceids.UserSpecifiedSegment("daprComponentName", "daprComponentValue"),
	}
}

// String returns a human-readable description of this Dapr Component ID
func (id DaprComponentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Environment Name: %q", id.ManagedEnvironmentName),
		fmt.Sprintf("Dapr Component Name: %q", id.DaprComponentName),
	}
	return fmt.Sprintf("Dapr Component (%s)", strings.Join(components, "\n"))
}
