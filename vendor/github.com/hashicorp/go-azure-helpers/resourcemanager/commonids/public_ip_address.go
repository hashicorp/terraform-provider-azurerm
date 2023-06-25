// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = PublicIPAddressId{}

// PublicIPAddressId is a struct representing the Resource ID for a Public I P Address
type PublicIPAddressId struct {
	SubscriptionId        string
	ResourceGroupName     string
	PublicIPAddressesName string
}

// NewPublicIPAddressID returns a new PublicIPAddressId struct
func NewPublicIPAddressID(subscriptionId string, resourceGroupName string, publicIPAddressesName string) PublicIPAddressId {
	return PublicIPAddressId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		PublicIPAddressesName: publicIPAddressesName,
	}
}

// ParsePublicIPAddressID parses 'input' into a PublicIPAddressId
func ParsePublicIPAddressID(input string) (*PublicIPAddressId, error) {
	parser := resourceids.NewParserFromResourceIdType(PublicIPAddressId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PublicIPAddressId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PublicIPAddressesName, ok = parsed.Parsed["publicIPAddressesName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "publicIPAddressesName", *parsed)
	}

	return &id, nil
}

// ParsePublicIPAddressIDInsensitively parses 'input' case-insensitively into a PublicIPAddressId
// note: this method should only be used for API response data and not user input
func ParsePublicIPAddressIDInsensitively(input string) (*PublicIPAddressId, error) {
	parser := resourceids.NewParserFromResourceIdType(PublicIPAddressId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PublicIPAddressId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PublicIPAddressesName, ok = parsed.Parsed["publicIPAddressesName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "publicIPAddressesName", *parsed)
	}

	return &id, nil
}

// ValidatePublicIPAddressID checks that 'input' can be parsed as a Public I P Address ID
func ValidatePublicIPAddressID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePublicIPAddressID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Public I P Address ID
func (id PublicIPAddressId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/publicIPAddresses/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PublicIPAddressesName)
}

// Segments returns a slice of Resource ID Segments which comprise this Public I P Address ID
func (id PublicIPAddressId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("publicIPAddresses", "publicIPAddresses", "publicIPAddresses"),
		resourceids.UserSpecifiedSegment("publicIPAddressesName", "publicIPAddressesValue"),
	}
}

// String returns a human-readable description of this Public I P Address ID
func (id PublicIPAddressId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Public I P Addresses Name: %q", id.PublicIPAddressesName),
	}
	return fmt.Sprintf("Public I P Address (%s)", strings.Join(components, "\n"))
}
