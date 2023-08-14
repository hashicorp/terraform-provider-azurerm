package virtualappliancesites

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = NetworkVirtualApplianceId{}

// NetworkVirtualApplianceId is a struct representing the Resource ID for a Network Virtual Appliance
type NetworkVirtualApplianceId struct {
	SubscriptionId              string
	ResourceGroupName           string
	NetworkVirtualApplianceName string
}

// NewNetworkVirtualApplianceID returns a new NetworkVirtualApplianceId struct
func NewNetworkVirtualApplianceID(subscriptionId string, resourceGroupName string, networkVirtualApplianceName string) NetworkVirtualApplianceId {
	return NetworkVirtualApplianceId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		NetworkVirtualApplianceName: networkVirtualApplianceName,
	}
}

// ParseNetworkVirtualApplianceID parses 'input' into a NetworkVirtualApplianceId
func ParseNetworkVirtualApplianceID(input string) (*NetworkVirtualApplianceId, error) {
	parser := resourceids.NewParserFromResourceIdType(NetworkVirtualApplianceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NetworkVirtualApplianceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkVirtualApplianceName, ok = parsed.Parsed["networkVirtualApplianceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkVirtualApplianceName", *parsed)
	}

	return &id, nil
}

// ParseNetworkVirtualApplianceIDInsensitively parses 'input' case-insensitively into a NetworkVirtualApplianceId
// note: this method should only be used for API response data and not user input
func ParseNetworkVirtualApplianceIDInsensitively(input string) (*NetworkVirtualApplianceId, error) {
	parser := resourceids.NewParserFromResourceIdType(NetworkVirtualApplianceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NetworkVirtualApplianceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkVirtualApplianceName, ok = parsed.Parsed["networkVirtualApplianceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkVirtualApplianceName", *parsed)
	}

	return &id, nil
}

// ValidateNetworkVirtualApplianceID checks that 'input' can be parsed as a Network Virtual Appliance ID
func ValidateNetworkVirtualApplianceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNetworkVirtualApplianceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Network Virtual Appliance ID
func (id NetworkVirtualApplianceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkVirtualAppliances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkVirtualApplianceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Network Virtual Appliance ID
func (id NetworkVirtualApplianceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkVirtualAppliances", "networkVirtualAppliances", "networkVirtualAppliances"),
		resourceids.UserSpecifiedSegment("networkVirtualApplianceName", "networkVirtualApplianceValue"),
	}
}

// String returns a human-readable description of this Network Virtual Appliance ID
func (id NetworkVirtualApplianceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Virtual Appliance Name: %q", id.NetworkVirtualApplianceName),
	}
	return fmt.Sprintf("Network Virtual Appliance (%s)", strings.Join(components, "\n"))
}
