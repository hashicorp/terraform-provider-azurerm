package staticmembers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = NetworkGroupId{}

// NetworkGroupId is a struct representing the Resource ID for a Network Group
type NetworkGroupId struct {
	SubscriptionId     string
	ResourceGroupName  string
	NetworkManagerName string
	NetworkGroupName   string
}

// NewNetworkGroupID returns a new NetworkGroupId struct
func NewNetworkGroupID(subscriptionId string, resourceGroupName string, networkManagerName string, networkGroupName string) NetworkGroupId {
	return NetworkGroupId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		NetworkManagerName: networkManagerName,
		NetworkGroupName:   networkGroupName,
	}
}

// ParseNetworkGroupID parses 'input' into a NetworkGroupId
func ParseNetworkGroupID(input string) (*NetworkGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(NetworkGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NetworkGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkManagerName, ok = parsed.Parsed["networkManagerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkManagerName", *parsed)
	}

	if id.NetworkGroupName, ok = parsed.Parsed["networkGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkGroupName", *parsed)
	}

	return &id, nil
}

// ParseNetworkGroupIDInsensitively parses 'input' case-insensitively into a NetworkGroupId
// note: this method should only be used for API response data and not user input
func ParseNetworkGroupIDInsensitively(input string) (*NetworkGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(NetworkGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NetworkGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkManagerName, ok = parsed.Parsed["networkManagerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkManagerName", *parsed)
	}

	if id.NetworkGroupName, ok = parsed.Parsed["networkGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkGroupName", *parsed)
	}

	return &id, nil
}

// ValidateNetworkGroupID checks that 'input' can be parsed as a Network Group ID
func ValidateNetworkGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNetworkGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Network Group ID
func (id NetworkGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/networkGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.NetworkGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Network Group ID
func (id NetworkGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkManagers", "networkManagers", "networkManagers"),
		resourceids.UserSpecifiedSegment("networkManagerName", "networkManagerValue"),
		resourceids.StaticSegment("staticNetworkGroups", "networkGroups", "networkGroups"),
		resourceids.UserSpecifiedSegment("networkGroupName", "networkGroupValue"),
	}
}

// String returns a human-readable description of this Network Group ID
func (id NetworkGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Manager Name: %q", id.NetworkManagerName),
		fmt.Sprintf("Network Group Name: %q", id.NetworkGroupName),
	}
	return fmt.Sprintf("Network Group (%s)", strings.Join(components, "\n"))
}
