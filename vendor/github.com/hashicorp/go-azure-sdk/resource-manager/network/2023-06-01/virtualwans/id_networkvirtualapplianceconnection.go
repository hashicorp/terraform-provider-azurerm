package virtualwans

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = NetworkVirtualApplianceConnectionId{}

// NetworkVirtualApplianceConnectionId is a struct representing the Resource ID for a Network Virtual Appliance Connection
type NetworkVirtualApplianceConnectionId struct {
	SubscriptionId                        string
	ResourceGroupName                     string
	NetworkVirtualApplianceName           string
	NetworkVirtualApplianceConnectionName string
}

// NewNetworkVirtualApplianceConnectionID returns a new NetworkVirtualApplianceConnectionId struct
func NewNetworkVirtualApplianceConnectionID(subscriptionId string, resourceGroupName string, networkVirtualApplianceName string, networkVirtualApplianceConnectionName string) NetworkVirtualApplianceConnectionId {
	return NetworkVirtualApplianceConnectionId{
		SubscriptionId:                        subscriptionId,
		ResourceGroupName:                     resourceGroupName,
		NetworkVirtualApplianceName:           networkVirtualApplianceName,
		NetworkVirtualApplianceConnectionName: networkVirtualApplianceConnectionName,
	}
}

// ParseNetworkVirtualApplianceConnectionID parses 'input' into a NetworkVirtualApplianceConnectionId
func ParseNetworkVirtualApplianceConnectionID(input string) (*NetworkVirtualApplianceConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(NetworkVirtualApplianceConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NetworkVirtualApplianceConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkVirtualApplianceName, ok = parsed.Parsed["networkVirtualApplianceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkVirtualApplianceName", *parsed)
	}

	if id.NetworkVirtualApplianceConnectionName, ok = parsed.Parsed["networkVirtualApplianceConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkVirtualApplianceConnectionName", *parsed)
	}

	return &id, nil
}

// ParseNetworkVirtualApplianceConnectionIDInsensitively parses 'input' case-insensitively into a NetworkVirtualApplianceConnectionId
// note: this method should only be used for API response data and not user input
func ParseNetworkVirtualApplianceConnectionIDInsensitively(input string) (*NetworkVirtualApplianceConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(NetworkVirtualApplianceConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NetworkVirtualApplianceConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkVirtualApplianceName, ok = parsed.Parsed["networkVirtualApplianceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkVirtualApplianceName", *parsed)
	}

	if id.NetworkVirtualApplianceConnectionName, ok = parsed.Parsed["networkVirtualApplianceConnectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkVirtualApplianceConnectionName", *parsed)
	}

	return &id, nil
}

// ValidateNetworkVirtualApplianceConnectionID checks that 'input' can be parsed as a Network Virtual Appliance Connection ID
func ValidateNetworkVirtualApplianceConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNetworkVirtualApplianceConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Network Virtual Appliance Connection ID
func (id NetworkVirtualApplianceConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkVirtualAppliances/%s/networkVirtualApplianceConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkVirtualApplianceName, id.NetworkVirtualApplianceConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Network Virtual Appliance Connection ID
func (id NetworkVirtualApplianceConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkVirtualAppliances", "networkVirtualAppliances", "networkVirtualAppliances"),
		resourceids.UserSpecifiedSegment("networkVirtualApplianceName", "networkVirtualApplianceValue"),
		resourceids.StaticSegment("staticNetworkVirtualApplianceConnections", "networkVirtualApplianceConnections", "networkVirtualApplianceConnections"),
		resourceids.UserSpecifiedSegment("networkVirtualApplianceConnectionName", "networkVirtualApplianceConnectionValue"),
	}
}

// String returns a human-readable description of this Network Virtual Appliance Connection ID
func (id NetworkVirtualApplianceConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Virtual Appliance Name: %q", id.NetworkVirtualApplianceName),
		fmt.Sprintf("Network Virtual Appliance Connection Name: %q", id.NetworkVirtualApplianceConnectionName),
	}
	return fmt.Sprintf("Network Virtual Appliance Connection (%s)", strings.Join(components, "\n"))
}
