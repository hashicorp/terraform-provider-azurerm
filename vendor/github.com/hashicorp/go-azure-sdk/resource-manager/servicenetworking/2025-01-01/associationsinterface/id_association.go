package associationsinterface

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AssociationId{})
}

var _ resourceids.ResourceId = &AssociationId{}

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
	parser := resourceids.NewParserFromResourceIdType(&AssociationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AssociationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAssociationIDInsensitively parses 'input' case-insensitively into a AssociationId
// note: this method should only be used for API response data and not user input
func ParseAssociationIDInsensitively(input string) (*AssociationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AssociationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AssociationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AssociationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.TrafficControllerName, ok = input.Parsed["trafficControllerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "trafficControllerName", input)
	}

	if id.AssociationName, ok = input.Parsed["associationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "associationName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("trafficControllerName", "trafficControllerName"),
		resourceids.StaticSegment("staticAssociations", "associations", "associations"),
		resourceids.UserSpecifiedSegment("associationName", "associationName"),
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
