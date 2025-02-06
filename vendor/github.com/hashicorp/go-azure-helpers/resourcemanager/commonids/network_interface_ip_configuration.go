// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &NetworkInterfaceIPConfigurationId{}

// NetworkInterfaceIPConfigurationId is a struct representing the Resource ID for a Network Interface I P Configuration
type NetworkInterfaceIPConfigurationId struct {
	SubscriptionId       string
	ResourceGroupName    string
	NetworkInterfaceName string
	IpConfigurationName  string
}

// NewNetworkInterfaceIPConfigurationID returns a new NetworkInterfaceIPConfigurationId struct
func NewNetworkInterfaceIPConfigurationID(subscriptionId string, resourceGroupName string, networkInterfaceName string, ipConfigurationName string) NetworkInterfaceIPConfigurationId {
	return NetworkInterfaceIPConfigurationId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		NetworkInterfaceName: networkInterfaceName,
		IpConfigurationName:  ipConfigurationName,
	}
}

// ParseNetworkInterfaceIPConfigurationID parses 'input' into a NetworkInterfaceIPConfigurationId
func ParseNetworkInterfaceIPConfigurationID(input string) (*NetworkInterfaceIPConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkInterfaceIPConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkInterfaceIPConfigurationId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNetworkInterfaceIPConfigurationIDInsensitively parses 'input' case-insensitively into a NetworkInterfaceIPConfigurationId
// note: this method should only be used for API response data and not user input
func ParseNetworkInterfaceIPConfigurationIDInsensitively(input string) (*NetworkInterfaceIPConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkInterfaceIPConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkInterfaceIPConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NetworkInterfaceIPConfigurationId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.IpConfigurationName, ok = input.Parsed["ipConfigurationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "ipConfigurationName", input)
	}

	return nil
}

// ValidateNetworkInterfaceIPConfigurationID checks that 'input' can be parsed as a Network Interface I P Configuration ID
func ValidateNetworkInterfaceIPConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNetworkInterfaceIPConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Network Interface I P Configuration ID
func (id NetworkInterfaceIPConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkInterfaces/%s/ipConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkInterfaceName, id.IpConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Network Interface I P Configuration ID
func (id NetworkInterfaceIPConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("networkInterfaces", "networkInterfaces", "networkInterfaces"),
		resourceids.UserSpecifiedSegment("networkInterfaceName", "networkInterfaceValue"),
		resourceids.StaticSegment("ipConfigurations", "ipConfigurations", "ipConfigurations"),
		resourceids.UserSpecifiedSegment("ipConfigurationName", "ipConfigurationValue"),
	}
}

// String returns a human-readable description of this Network Interface I P Configuration ID
func (id NetworkInterfaceIPConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Interface Name: %q", id.NetworkInterfaceName),
		fmt.Sprintf("Ip Configuration Name: %q", id.IpConfigurationName),
	}
	return fmt.Sprintf("Network Interface I P Configuration (%s)", strings.Join(components, "\n"))
}
