package virtualwans

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NetworkVirtualApplianceConnectionId{})
}

var _ resourceids.ResourceId = &NetworkVirtualApplianceConnectionId{}

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
	parser := resourceids.NewParserFromResourceIdType(&NetworkVirtualApplianceConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkVirtualApplianceConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNetworkVirtualApplianceConnectionIDInsensitively parses 'input' case-insensitively into a NetworkVirtualApplianceConnectionId
// note: this method should only be used for API response data and not user input
func ParseNetworkVirtualApplianceConnectionIDInsensitively(input string) (*NetworkVirtualApplianceConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkVirtualApplianceConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkVirtualApplianceConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NetworkVirtualApplianceConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkVirtualApplianceName, ok = input.Parsed["networkVirtualApplianceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkVirtualApplianceName", input)
	}

	if id.NetworkVirtualApplianceConnectionName, ok = input.Parsed["networkVirtualApplianceConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkVirtualApplianceConnectionName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("networkVirtualApplianceName", "networkVirtualApplianceName"),
		resourceids.StaticSegment("staticNetworkVirtualApplianceConnections", "networkVirtualApplianceConnections", "networkVirtualApplianceConnections"),
		resourceids.UserSpecifiedSegment("networkVirtualApplianceConnectionName", "networkVirtualApplianceConnectionName"),
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
