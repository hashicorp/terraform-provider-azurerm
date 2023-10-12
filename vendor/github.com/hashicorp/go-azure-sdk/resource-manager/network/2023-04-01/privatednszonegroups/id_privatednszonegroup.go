package privatednszonegroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PrivateDnsZoneGroupId{}

// PrivateDnsZoneGroupId is a struct representing the Resource ID for a Private Dns Zone Group
type PrivateDnsZoneGroupId struct {
	SubscriptionId          string
	ResourceGroupName       string
	PrivateEndpointName     string
	PrivateDnsZoneGroupName string
}

// NewPrivateDnsZoneGroupID returns a new PrivateDnsZoneGroupId struct
func NewPrivateDnsZoneGroupID(subscriptionId string, resourceGroupName string, privateEndpointName string, privateDnsZoneGroupName string) PrivateDnsZoneGroupId {
	return PrivateDnsZoneGroupId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		PrivateEndpointName:     privateEndpointName,
		PrivateDnsZoneGroupName: privateDnsZoneGroupName,
	}
}

// ParsePrivateDnsZoneGroupID parses 'input' into a PrivateDnsZoneGroupId
func ParsePrivateDnsZoneGroupID(input string) (*PrivateDnsZoneGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateDnsZoneGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateDnsZoneGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PrivateEndpointName, ok = parsed.Parsed["privateEndpointName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateEndpointName", *parsed)
	}

	if id.PrivateDnsZoneGroupName, ok = parsed.Parsed["privateDnsZoneGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateDnsZoneGroupName", *parsed)
	}

	return &id, nil
}

// ParsePrivateDnsZoneGroupIDInsensitively parses 'input' case-insensitively into a PrivateDnsZoneGroupId
// note: this method should only be used for API response data and not user input
func ParsePrivateDnsZoneGroupIDInsensitively(input string) (*PrivateDnsZoneGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateDnsZoneGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateDnsZoneGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PrivateEndpointName, ok = parsed.Parsed["privateEndpointName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateEndpointName", *parsed)
	}

	if id.PrivateDnsZoneGroupName, ok = parsed.Parsed["privateDnsZoneGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateDnsZoneGroupName", *parsed)
	}

	return &id, nil
}

// ValidatePrivateDnsZoneGroupID checks that 'input' can be parsed as a Private Dns Zone Group ID
func ValidatePrivateDnsZoneGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePrivateDnsZoneGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Private Dns Zone Group ID
func (id PrivateDnsZoneGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateEndpoints/%s/privateDnsZoneGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateEndpointName, id.PrivateDnsZoneGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Private Dns Zone Group ID
func (id PrivateDnsZoneGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticPrivateEndpoints", "privateEndpoints", "privateEndpoints"),
		resourceids.UserSpecifiedSegment("privateEndpointName", "privateEndpointValue"),
		resourceids.StaticSegment("staticPrivateDnsZoneGroups", "privateDnsZoneGroups", "privateDnsZoneGroups"),
		resourceids.UserSpecifiedSegment("privateDnsZoneGroupName", "privateDnsZoneGroupValue"),
	}
}

// String returns a human-readable description of this Private Dns Zone Group ID
func (id PrivateDnsZoneGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Endpoint Name: %q", id.PrivateEndpointName),
		fmt.Sprintf("Private Dns Zone Group Name: %q", id.PrivateDnsZoneGroupName),
	}
	return fmt.Sprintf("Private Dns Zone Group (%s)", strings.Join(components, "\n"))
}
