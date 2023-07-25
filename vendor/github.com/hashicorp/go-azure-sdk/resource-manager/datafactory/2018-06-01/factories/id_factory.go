package factories

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = FactoryId{}

// FactoryId is a struct representing the Resource ID for a Factory
type FactoryId struct {
	SubscriptionId    string
	ResourceGroupName string
	FactoryName       string
}

// NewFactoryID returns a new FactoryId struct
func NewFactoryID(subscriptionId string, resourceGroupName string, factoryName string) FactoryId {
	return FactoryId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		FactoryName:       factoryName,
	}
}

// ParseFactoryID parses 'input' into a FactoryId
func ParseFactoryID(input string) (*FactoryId, error) {
	parser := resourceids.NewParserFromResourceIdType(FactoryId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FactoryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.FactoryName, ok = parsed.Parsed["factoryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "factoryName", *parsed)
	}

	return &id, nil
}

// ParseFactoryIDInsensitively parses 'input' case-insensitively into a FactoryId
// note: this method should only be used for API response data and not user input
func ParseFactoryIDInsensitively(input string) (*FactoryId, error) {
	parser := resourceids.NewParserFromResourceIdType(FactoryId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FactoryId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.FactoryName, ok = parsed.Parsed["factoryName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "factoryName", *parsed)
	}

	return &id, nil
}

// ValidateFactoryID checks that 'input' can be parsed as a Factory ID
func ValidateFactoryID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFactoryID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Factory ID
func (id FactoryId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataFactory/factories/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FactoryName)
}

// Segments returns a slice of Resource ID Segments which comprise this Factory ID
func (id FactoryId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataFactory", "Microsoft.DataFactory", "Microsoft.DataFactory"),
		resourceids.StaticSegment("staticFactories", "factories", "factories"),
		resourceids.UserSpecifiedSegment("factoryName", "factoryValue"),
	}
}

// String returns a human-readable description of this Factory ID
func (id FactoryId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Factory Name: %q", id.FactoryName),
	}
	return fmt.Sprintf("Factory (%s)", strings.Join(components, "\n"))
}
