// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &SubnetId{}

// SubnetId is a struct representing the Resource ID for a Subnet
type SubnetId struct {
	SubscriptionId     string
	ResourceGroupName  string
	VirtualNetworkName string
	SubnetName         string
}

// NewSubnetID returns a new SubnetId struct
func NewSubnetID(subscriptionId string, resourceGroupName string, virtualNetworkName string, subnetName string) SubnetId {
	return SubnetId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		VirtualNetworkName: virtualNetworkName,
		SubnetName:         subnetName,
	}
}

// ParseSubnetID parses 'input' into a SubnetId
func ParseSubnetID(input string) (*SubnetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SubnetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SubnetId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSubnetIDInsensitively parses 'input' case-insensitively into a SubnetId
// note: this method should only be used for API response data and not user input
func ParseSubnetIDInsensitively(input string) (*SubnetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SubnetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SubnetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SubnetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VirtualNetworkName, ok = input.Parsed["virtualNetworkName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualNetworkName", input)
	}

	if id.SubnetName, ok = input.Parsed["subnetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subnetName", input)
	}

	return nil
}

// ValidateSubnetID checks that 'input' can be parsed as a Subnet ID
func ValidateSubnetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSubnetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Subnet ID
func (id SubnetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/subnets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName, id.SubnetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Subnet ID
func (id SubnetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("virtualNetworks", "virtualNetworks", "virtualNetworks"),
		resourceids.UserSpecifiedSegment("virtualNetworkName", "virtualNetworksValue"),
		resourceids.StaticSegment("subnets", "subnets", "subnets"),
		resourceids.UserSpecifiedSegment("subnetName", "subnetValue"),
	}
}

// String returns a human-readable description of this Subnet ID
func (id SubnetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Network Name: %q", id.VirtualNetworkName),
		fmt.Sprintf("Subnet Name: %q", id.SubnetName),
	}
	return fmt.Sprintf("Subnet (%s)", strings.Join(components, "\n"))
}
