package recordsets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PrivateDnsZoneId{}

// PrivateDnsZoneId is a struct representing the Resource ID for a Private Dns Zone
type PrivateDnsZoneId struct {
	SubscriptionId     string
	ResourceGroupName  string
	PrivateDnsZoneName string
}

// NewPrivateDnsZoneID returns a new PrivateDnsZoneId struct
func NewPrivateDnsZoneID(subscriptionId string, resourceGroupName string, privateDnsZoneName string) PrivateDnsZoneId {
	return PrivateDnsZoneId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		PrivateDnsZoneName: privateDnsZoneName,
	}
}

// ParsePrivateDnsZoneID parses 'input' into a PrivateDnsZoneId
func ParsePrivateDnsZoneID(input string) (*PrivateDnsZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateDnsZoneId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateDnsZoneId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PrivateDnsZoneName, ok = parsed.Parsed["privateDnsZoneName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateDnsZoneName", *parsed)
	}

	return &id, nil
}

// ParsePrivateDnsZoneIDInsensitively parses 'input' case-insensitively into a PrivateDnsZoneId
// note: this method should only be used for API response data and not user input
func ParsePrivateDnsZoneIDInsensitively(input string) (*PrivateDnsZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateDnsZoneId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateDnsZoneId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PrivateDnsZoneName, ok = parsed.Parsed["privateDnsZoneName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateDnsZoneName", *parsed)
	}

	return &id, nil
}

// ValidatePrivateDnsZoneID checks that 'input' can be parsed as a Private Dns Zone ID
func ValidatePrivateDnsZoneID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePrivateDnsZoneID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Private Dns Zone ID
func (id PrivateDnsZoneId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZones/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateDnsZoneName)
}

// Segments returns a slice of Resource ID Segments which comprise this Private Dns Zone ID
func (id PrivateDnsZoneId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticPrivateDnsZones", "privateDnsZones", "privateDnsZones"),
		resourceids.UserSpecifiedSegment("privateDnsZoneName", "privateDnsZoneValue"),
	}
}

// String returns a human-readable description of this Private Dns Zone ID
func (id PrivateDnsZoneId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Dns Zone Name: %q", id.PrivateDnsZoneName),
	}
	return fmt.Sprintf("Private Dns Zone (%s)", strings.Join(components, "\n"))
}
