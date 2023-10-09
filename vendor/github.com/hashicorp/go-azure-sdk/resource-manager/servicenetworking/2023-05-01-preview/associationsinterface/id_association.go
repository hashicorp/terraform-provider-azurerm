package associationsinterface

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = AssociationId{}

// AssociationId is a struct representing the Resource ID for a Association
type AssociationId struct {
	SubscriptionId        string
	ResourceGroupName     string
	TrafficControllerName string
	AssociationName       string
}

// NewAssociationID returns a new AssociationId struct
func NewAssociationID(subscriptionId string, resourceGroupName string, trafficControllerName string, associationName string) AssociationId {
	return AssociationId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		TrafficControllerName: trafficControllerName,
		AssociationName:       associationName,
	}
}

// ParseAssociationID parses 'input' into a AssociationId
func ParseAssociationID(input string) (*AssociationId, error) {
	parser := resourceids.NewParserFromResourceIdType(AssociationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AssociationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.TrafficControllerName, ok = parsed.Parsed["trafficControllerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "trafficControllerName", *parsed)
	}

	if id.AssociationName, ok = parsed.Parsed["associationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "associationName", *parsed)
	}

	return &id, nil
}

// ParseAssociationIDInsensitively parses 'input' case-insensitively into a AssociationId
// note: this method should only be used for API response data and not user input
func ParseAssociationIDInsensitively(input string) (*AssociationId, error) {
	parser := resourceids.NewParserFromResourceIdType(AssociationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AssociationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.TrafficControllerName, ok = parsed.Parsed["trafficControllerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "trafficControllerName", *parsed)
	}

	if id.AssociationName, ok = parsed.Parsed["associationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "associationName", *parsed)
	}

	return &id, nil
}

// ValidateAssociationID checks that 'input' can be parsed as a Association ID
func ValidateAssociationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAssociationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Association ID
func (id AssociationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceNetworking/trafficControllers/%s/associations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.TrafficControllerName, id.AssociationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Association ID
func (id AssociationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftServiceNetworking", "Microsoft.ServiceNetworking", "Microsoft.ServiceNetworking"),
		resourceids.StaticSegment("staticTrafficControllers", "trafficControllers", "trafficControllers"),
		resourceids.UserSpecifiedSegment("trafficControllerName", "trafficControllerValue"),
		resourceids.StaticSegment("staticAssociations", "associations", "associations"),
		resourceids.UserSpecifiedSegment("associationName", "associationValue"),
	}
}

// String returns a human-readable description of this Association ID
func (id AssociationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Traffic Controller Name: %q", id.TrafficControllerName),
		fmt.Sprintf("Association Name: %q", id.AssociationName),
	}
	return fmt.Sprintf("Association (%s)", strings.Join(components, "\n"))
}
