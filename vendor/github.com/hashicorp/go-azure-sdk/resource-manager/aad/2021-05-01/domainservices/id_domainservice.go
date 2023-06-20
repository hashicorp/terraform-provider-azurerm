package domainservices

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DomainServiceId{}

// DomainServiceId is a struct representing the Resource ID for a Domain Service
type DomainServiceId struct {
	SubscriptionId    string
	ResourceGroupName string
	DomainServiceName string
}

// NewDomainServiceID returns a new DomainServiceId struct
func NewDomainServiceID(subscriptionId string, resourceGroupName string, domainServiceName string) DomainServiceId {
	return DomainServiceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DomainServiceName: domainServiceName,
	}
}

// ParseDomainServiceID parses 'input' into a DomainServiceId
func ParseDomainServiceID(input string) (*DomainServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(DomainServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DomainServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DomainServiceName, ok = parsed.Parsed["domainServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "domainServiceName", *parsed)
	}

	return &id, nil
}

// ParseDomainServiceIDInsensitively parses 'input' case-insensitively into a DomainServiceId
// note: this method should only be used for API response data and not user input
func ParseDomainServiceIDInsensitively(input string) (*DomainServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(DomainServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DomainServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DomainServiceName, ok = parsed.Parsed["domainServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "domainServiceName", *parsed)
	}

	return &id, nil
}

// ValidateDomainServiceID checks that 'input' can be parsed as a Domain Service ID
func ValidateDomainServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDomainServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Domain Service ID
func (id DomainServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AAD/domainServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DomainServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Domain Service ID
func (id DomainServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAAD", "Microsoft.AAD", "Microsoft.AAD"),
		resourceids.StaticSegment("staticDomainServices", "domainServices", "domainServices"),
		resourceids.UserSpecifiedSegment("domainServiceName", "domainServiceValue"),
	}
}

// String returns a human-readable description of this Domain Service ID
func (id DomainServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Domain Service Name: %q", id.DomainServiceName),
	}
	return fmt.Sprintf("Domain Service (%s)", strings.Join(components, "\n"))
}
