package natgateways

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NatGatewayId{})
}

var _ resourceids.ResourceId = &NatGatewayId{}

// NatGatewayId is a struct representing the Resource ID for a Nat Gateway
type NatGatewayId struct {
	SubscriptionId    string
	ResourceGroupName string
	NatGatewayName    string
}

// NewNatGatewayID returns a new NatGatewayId struct
func NewNatGatewayID(subscriptionId string, resourceGroupName string, natGatewayName string) NatGatewayId {
	return NatGatewayId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NatGatewayName:    natGatewayName,
	}
}

// ParseNatGatewayID parses 'input' into a NatGatewayId
func ParseNatGatewayID(input string) (*NatGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NatGatewayId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NatGatewayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNatGatewayIDInsensitively parses 'input' case-insensitively into a NatGatewayId
// note: this method should only be used for API response data and not user input
func ParseNatGatewayIDInsensitively(input string) (*NatGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NatGatewayId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NatGatewayId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NatGatewayId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NatGatewayName, ok = input.Parsed["natGatewayName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "natGatewayName", input)
	}

	return nil
}

// ValidateNatGatewayID checks that 'input' can be parsed as a Nat Gateway ID
func ValidateNatGatewayID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNatGatewayID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Nat Gateway ID
func (id NatGatewayId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/natGateways/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NatGatewayName)
}

// Segments returns a slice of Resource ID Segments which comprise this Nat Gateway ID
func (id NatGatewayId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNatGateways", "natGateways", "natGateways"),
		resourceids.UserSpecifiedSegment("natGatewayName", "natGatewayName"),
	}
}

// String returns a human-readable description of this Nat Gateway ID
func (id NatGatewayId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Nat Gateway Name: %q", id.NatGatewayName),
	}
	return fmt.Sprintf("Nat Gateway (%s)", strings.Join(components, "\n"))
}
