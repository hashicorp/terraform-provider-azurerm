package networksecuritygroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = NetworkSecurityGroupId{}

// NetworkSecurityGroupId is a struct representing the Resource ID for a Network Security Group
type NetworkSecurityGroupId struct {
	SubscriptionId           string
	ResourceGroupName        string
	NetworkSecurityGroupName string
}

// NewNetworkSecurityGroupID returns a new NetworkSecurityGroupId struct
func NewNetworkSecurityGroupID(subscriptionId string, resourceGroupName string, networkSecurityGroupName string) NetworkSecurityGroupId {
	return NetworkSecurityGroupId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		NetworkSecurityGroupName: networkSecurityGroupName,
	}
}

// ParseNetworkSecurityGroupID parses 'input' into a NetworkSecurityGroupId
func ParseNetworkSecurityGroupID(input string) (*NetworkSecurityGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(NetworkSecurityGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NetworkSecurityGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkSecurityGroupName, ok = parsed.Parsed["networkSecurityGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkSecurityGroupName", *parsed)
	}

	return &id, nil
}

// ParseNetworkSecurityGroupIDInsensitively parses 'input' case-insensitively into a NetworkSecurityGroupId
// note: this method should only be used for API response data and not user input
func ParseNetworkSecurityGroupIDInsensitively(input string) (*NetworkSecurityGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(NetworkSecurityGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NetworkSecurityGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkSecurityGroupName, ok = parsed.Parsed["networkSecurityGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkSecurityGroupName", *parsed)
	}

	return &id, nil
}

// ValidateNetworkSecurityGroupID checks that 'input' can be parsed as a Network Security Group ID
func ValidateNetworkSecurityGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNetworkSecurityGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Network Security Group ID
func (id NetworkSecurityGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkSecurityGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkSecurityGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Network Security Group ID
func (id NetworkSecurityGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkSecurityGroups", "networkSecurityGroups", "networkSecurityGroups"),
		resourceids.UserSpecifiedSegment("networkSecurityGroupName", "networkSecurityGroupValue"),
	}
}

// String returns a human-readable description of this Network Security Group ID
func (id NetworkSecurityGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Security Group Name: %q", id.NetworkSecurityGroupName),
	}
	return fmt.Sprintf("Network Security Group (%s)", strings.Join(components, "\n"))
}
