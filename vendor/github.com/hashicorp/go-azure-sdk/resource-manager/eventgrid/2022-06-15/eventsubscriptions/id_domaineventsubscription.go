package eventsubscriptions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DomainEventSubscriptionId{}

// DomainEventSubscriptionId is a struct representing the Resource ID for a Domain Event Subscription
type DomainEventSubscriptionId struct {
	SubscriptionId        string
	ResourceGroupName     string
	DomainName            string
	EventSubscriptionName string
}

// NewDomainEventSubscriptionID returns a new DomainEventSubscriptionId struct
func NewDomainEventSubscriptionID(subscriptionId string, resourceGroupName string, domainName string, eventSubscriptionName string) DomainEventSubscriptionId {
	return DomainEventSubscriptionId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		DomainName:            domainName,
		EventSubscriptionName: eventSubscriptionName,
	}
}

// ParseDomainEventSubscriptionID parses 'input' into a DomainEventSubscriptionId
func ParseDomainEventSubscriptionID(input string) (*DomainEventSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(DomainEventSubscriptionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DomainEventSubscriptionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DomainName, ok = parsed.Parsed["domainName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "domainName", *parsed)
	}

	if id.EventSubscriptionName, ok = parsed.Parsed["eventSubscriptionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "eventSubscriptionName", *parsed)
	}

	return &id, nil
}

// ParseDomainEventSubscriptionIDInsensitively parses 'input' case-insensitively into a DomainEventSubscriptionId
// note: this method should only be used for API response data and not user input
func ParseDomainEventSubscriptionIDInsensitively(input string) (*DomainEventSubscriptionId, error) {
	parser := resourceids.NewParserFromResourceIdType(DomainEventSubscriptionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DomainEventSubscriptionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DomainName, ok = parsed.Parsed["domainName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "domainName", *parsed)
	}

	if id.EventSubscriptionName, ok = parsed.Parsed["eventSubscriptionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "eventSubscriptionName", *parsed)
	}

	return &id, nil
}

// ValidateDomainEventSubscriptionID checks that 'input' can be parsed as a Domain Event Subscription ID
func ValidateDomainEventSubscriptionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDomainEventSubscriptionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Domain Event Subscription ID
func (id DomainEventSubscriptionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/domains/%s/eventSubscriptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DomainName, id.EventSubscriptionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Domain Event Subscription ID
func (id DomainEventSubscriptionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticDomains", "domains", "domains"),
		resourceids.UserSpecifiedSegment("domainName", "domainValue"),
		resourceids.StaticSegment("staticEventSubscriptions", "eventSubscriptions", "eventSubscriptions"),
		resourceids.UserSpecifiedSegment("eventSubscriptionName", "eventSubscriptionValue"),
	}
}

// String returns a human-readable description of this Domain Event Subscription ID
func (id DomainEventSubscriptionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Domain Name: %q", id.DomainName),
		fmt.Sprintf("Event Subscription Name: %q", id.EventSubscriptionName),
	}
	return fmt.Sprintf("Domain Event Subscription (%s)", strings.Join(components, "\n"))
}
