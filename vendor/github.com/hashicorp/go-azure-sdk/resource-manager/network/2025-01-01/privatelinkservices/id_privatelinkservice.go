package privatelinkservices

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PrivateLinkServiceId{})
}

var _ resourceids.ResourceId = &PrivateLinkServiceId{}

// PrivateLinkServiceId is a struct representing the Resource ID for a Private Link Service
type PrivateLinkServiceId struct {
	SubscriptionId         string
	ResourceGroupName      string
	PrivateLinkServiceName string
}

// NewPrivateLinkServiceID returns a new PrivateLinkServiceId struct
func NewPrivateLinkServiceID(subscriptionId string, resourceGroupName string, privateLinkServiceName string) PrivateLinkServiceId {
	return PrivateLinkServiceId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		PrivateLinkServiceName: privateLinkServiceName,
	}
}

// ParsePrivateLinkServiceID parses 'input' into a PrivateLinkServiceId
func ParsePrivateLinkServiceID(input string) (*PrivateLinkServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrivateLinkServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateLinkServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePrivateLinkServiceIDInsensitively parses 'input' case-insensitively into a PrivateLinkServiceId
// note: this method should only be used for API response data and not user input
func ParsePrivateLinkServiceIDInsensitively(input string) (*PrivateLinkServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrivateLinkServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateLinkServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PrivateLinkServiceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.PrivateLinkServiceName, ok = input.Parsed["privateLinkServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "privateLinkServiceName", input)
	}

	return nil
}

// ValidatePrivateLinkServiceID checks that 'input' can be parsed as a Private Link Service ID
func ValidatePrivateLinkServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePrivateLinkServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Private Link Service ID
func (id PrivateLinkServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateLinkServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateLinkServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Private Link Service ID
func (id PrivateLinkServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticPrivateLinkServices", "privateLinkServices", "privateLinkServices"),
		resourceids.UserSpecifiedSegment("privateLinkServiceName", "privateLinkServiceName"),
	}
}

// String returns a human-readable description of this Private Link Service ID
func (id PrivateLinkServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Link Service Name: %q", id.PrivateLinkServiceName),
	}
	return fmt.Sprintf("Private Link Service (%s)", strings.Join(components, "\n"))
}
