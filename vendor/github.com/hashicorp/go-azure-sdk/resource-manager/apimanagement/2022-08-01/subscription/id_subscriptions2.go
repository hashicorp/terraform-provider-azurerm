package subscription

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &Subscriptions2Id{}

// Subscriptions2Id is a struct representing the Resource ID for a Subscriptions 2
type Subscriptions2Id struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	SubscriptionName  string
}

// NewSubscriptions2ID returns a new Subscriptions2Id struct
func NewSubscriptions2ID(subscriptionId string, resourceGroupName string, serviceName string, subscriptionName string) Subscriptions2Id {
	return Subscriptions2Id{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		SubscriptionName:  subscriptionName,
	}
}

// ParseSubscriptions2ID parses 'input' into a Subscriptions2Id
func ParseSubscriptions2ID(input string) (*Subscriptions2Id, error) {
	parser := resourceids.NewParserFromResourceIdType(&Subscriptions2Id{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := Subscriptions2Id{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSubscriptions2IDInsensitively parses 'input' case-insensitively into a Subscriptions2Id
// note: this method should only be used for API response data and not user input
func ParseSubscriptions2IDInsensitively(input string) (*Subscriptions2Id, error) {
	parser := resourceids.NewParserFromResourceIdType(&Subscriptions2Id{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := Subscriptions2Id{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *Subscriptions2Id) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServiceName, ok = input.Parsed["serviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceName", input)
	}

	if id.SubscriptionName, ok = input.Parsed["subscriptionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionName", input)
	}

	return nil
}

// ValidateSubscriptions2ID checks that 'input' can be parsed as a Subscriptions 2 ID
func ValidateSubscriptions2ID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSubscriptions2ID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Subscriptions 2 ID
func (id Subscriptions2Id) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/subscriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.SubscriptionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Subscriptions 2 ID
func (id Subscriptions2Id) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
		resourceids.StaticSegment("staticSubscriptions2", "subscriptions", "subscriptions"),
		resourceids.UserSpecifiedSegment("subscriptionName", "subscriptionValue"),
	}
}

// String returns a human-readable description of this Subscriptions 2 ID
func (id Subscriptions2Id) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Subscription Name: %q", id.SubscriptionName),
	}
	return fmt.Sprintf("Subscriptions 2 (%s)", strings.Join(components, "\n"))
}
