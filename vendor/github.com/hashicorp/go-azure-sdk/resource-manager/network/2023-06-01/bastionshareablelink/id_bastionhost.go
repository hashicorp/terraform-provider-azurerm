package bastionshareablelink

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = BastionHostId{}

// BastionHostId is a struct representing the Resource ID for a Bastion Host
type BastionHostId struct {
	SubscriptionId    string
	ResourceGroupName string
	BastionHostName   string
}

// NewBastionHostID returns a new BastionHostId struct
func NewBastionHostID(subscriptionId string, resourceGroupName string, bastionHostName string) BastionHostId {
	return BastionHostId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		BastionHostName:   bastionHostName,
	}
}

// ParseBastionHostID parses 'input' into a BastionHostId
func ParseBastionHostID(input string) (*BastionHostId, error) {
	parser := resourceids.NewParserFromResourceIdType(BastionHostId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BastionHostId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.BastionHostName, ok = parsed.Parsed["bastionHostName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "bastionHostName", *parsed)
	}

	return &id, nil
}

// ParseBastionHostIDInsensitively parses 'input' case-insensitively into a BastionHostId
// note: this method should only be used for API response data and not user input
func ParseBastionHostIDInsensitively(input string) (*BastionHostId, error) {
	parser := resourceids.NewParserFromResourceIdType(BastionHostId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BastionHostId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.BastionHostName, ok = parsed.Parsed["bastionHostName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "bastionHostName", *parsed)
	}

	return &id, nil
}

// ValidateBastionHostID checks that 'input' can be parsed as a Bastion Host ID
func ValidateBastionHostID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBastionHostID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Bastion Host ID
func (id BastionHostId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/bastionHosts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.BastionHostName)
}

// Segments returns a slice of Resource ID Segments which comprise this Bastion Host ID
func (id BastionHostId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticBastionHosts", "bastionHosts", "bastionHosts"),
		resourceids.UserSpecifiedSegment("bastionHostName", "bastionHostValue"),
	}
}

// String returns a human-readable description of this Bastion Host ID
func (id BastionHostId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Bastion Host Name: %q", id.BastionHostName),
	}
	return fmt.Sprintf("Bastion Host (%s)", strings.Join(components, "\n"))
}
