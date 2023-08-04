package networkprofiles

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = NetworkProfileId{}

// NetworkProfileId is a struct representing the Resource ID for a Network Profile
type NetworkProfileId struct {
	SubscriptionId     string
	ResourceGroupName  string
	NetworkProfileName string
}

// NewNetworkProfileID returns a new NetworkProfileId struct
func NewNetworkProfileID(subscriptionId string, resourceGroupName string, networkProfileName string) NetworkProfileId {
	return NetworkProfileId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		NetworkProfileName: networkProfileName,
	}
}

// ParseNetworkProfileID parses 'input' into a NetworkProfileId
func ParseNetworkProfileID(input string) (*NetworkProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(NetworkProfileId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NetworkProfileId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkProfileName, ok = parsed.Parsed["networkProfileName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkProfileName", *parsed)
	}

	return &id, nil
}

// ParseNetworkProfileIDInsensitively parses 'input' case-insensitively into a NetworkProfileId
// note: this method should only be used for API response data and not user input
func ParseNetworkProfileIDInsensitively(input string) (*NetworkProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(NetworkProfileId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NetworkProfileId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkProfileName, ok = parsed.Parsed["networkProfileName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkProfileName", *parsed)
	}

	return &id, nil
}

// ValidateNetworkProfileID checks that 'input' can be parsed as a Network Profile ID
func ValidateNetworkProfileID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNetworkProfileID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Network Profile ID
func (id NetworkProfileId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkProfiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkProfileName)
}

// Segments returns a slice of Resource ID Segments which comprise this Network Profile ID
func (id NetworkProfileId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkProfiles", "networkProfiles", "networkProfiles"),
		resourceids.UserSpecifiedSegment("networkProfileName", "networkProfileValue"),
	}
}

// String returns a human-readable description of this Network Profile ID
func (id NetworkProfileId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Profile Name: %q", id.NetworkProfileName),
	}
	return fmt.Sprintf("Network Profile (%s)", strings.Join(components, "\n"))
}
