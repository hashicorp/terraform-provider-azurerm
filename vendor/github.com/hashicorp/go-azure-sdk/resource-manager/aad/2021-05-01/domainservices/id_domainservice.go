package domainservices

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DomainServiceId{})
}

var _ resourceids.ResourceId = &DomainServiceId{}

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
	parser := resourceids.NewParserFromResourceIdType(&DomainServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DomainServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDomainServiceIDInsensitively parses 'input' case-insensitively into a DomainServiceId
// note: this method should only be used for API response data and not user input
func ParseDomainServiceIDInsensitively(input string) (*DomainServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DomainServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DomainServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DomainServiceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DomainServiceName, ok = input.Parsed["domainServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "domainServiceName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("domainServiceName", "domainServiceName"),
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
