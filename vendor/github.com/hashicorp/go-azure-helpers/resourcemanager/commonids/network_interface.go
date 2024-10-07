// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &NetworkInterfaceId{}

// NetworkInterfaceId is a struct representing the Resource ID for a Network Interface
type NetworkInterfaceId struct {
	SubscriptionId       string
	ResourceGroupName    string
	NetworkInterfaceName string
}

// NewNetworkInterfaceID returns a new NetworkInterfaceId struct
func NewNetworkInterfaceID(subscriptionId string, resourceGroupName string, networkInterfaceName string) NetworkInterfaceId {
	return NetworkInterfaceId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		NetworkInterfaceName: networkInterfaceName,
	}
}

// ParseNetworkInterfaceID parses 'input' into a NetworkInterfaceId
func ParseNetworkInterfaceID(input string) (*NetworkInterfaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkInterfaceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkInterfaceId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNetworkInterfaceIDInsensitively parses 'input' case-insensitively into a NetworkInterfaceId
// note: this method should only be used for API response data and not user input
func ParseNetworkInterfaceIDInsensitively(input string) (*NetworkInterfaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkInterfaceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkInterfaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NetworkInterfaceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkInterfaceName, ok = input.Parsed["networkInterfaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkInterfaceName", input)
	}

	return nil
}

// ValidateNetworkInterfaceID checks that 'input' can be parsed as a Network Interface ID
func ValidateNetworkInterfaceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNetworkInterfaceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Network Interface ID
func (id NetworkInterfaceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkInterfaces/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkInterfaceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Network Interface ID
func (id NetworkInterfaceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("networkInterfaces", "networkInterfaces", "networkInterfaces"),
		resourceids.UserSpecifiedSegment("networkInterfaceName", "networkInterfaceValue"),
	}
}

// String returns a human-readable description of this Network Interface ID
func (id NetworkInterfaceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Interface Name: %q", id.NetworkInterfaceName),
	}
	return fmt.Sprintf("Network Interface (%s)", strings.Join(components, "\n"))
}
