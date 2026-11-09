package virtualapplianceskus

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NetworkVirtualApplianceSkuId{})
}

var _ resourceids.ResourceId = &NetworkVirtualApplianceSkuId{}

// NetworkVirtualApplianceSkuId is a struct representing the Resource ID for a Network Virtual Appliance Sku
type NetworkVirtualApplianceSkuId struct {
	SubscriptionId                 string
	NetworkVirtualApplianceSkuName string
}

// NewNetworkVirtualApplianceSkuID returns a new NetworkVirtualApplianceSkuId struct
func NewNetworkVirtualApplianceSkuID(subscriptionId string, networkVirtualApplianceSkuName string) NetworkVirtualApplianceSkuId {
	return NetworkVirtualApplianceSkuId{
		SubscriptionId:                 subscriptionId,
		NetworkVirtualApplianceSkuName: networkVirtualApplianceSkuName,
	}
}

// ParseNetworkVirtualApplianceSkuID parses 'input' into a NetworkVirtualApplianceSkuId
func ParseNetworkVirtualApplianceSkuID(input string) (*NetworkVirtualApplianceSkuId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkVirtualApplianceSkuId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkVirtualApplianceSkuId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNetworkVirtualApplianceSkuIDInsensitively parses 'input' case-insensitively into a NetworkVirtualApplianceSkuId
// note: this method should only be used for API response data and not user input
func ParseNetworkVirtualApplianceSkuIDInsensitively(input string) (*NetworkVirtualApplianceSkuId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkVirtualApplianceSkuId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkVirtualApplianceSkuId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NetworkVirtualApplianceSkuId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.NetworkVirtualApplianceSkuName, ok = input.Parsed["networkVirtualApplianceSkuName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkVirtualApplianceSkuName", input)
	}

	return nil
}

// ValidateNetworkVirtualApplianceSkuID checks that 'input' can be parsed as a Network Virtual Appliance Sku ID
func ValidateNetworkVirtualApplianceSkuID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNetworkVirtualApplianceSkuID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Network Virtual Appliance Sku ID
func (id NetworkVirtualApplianceSkuId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Network/networkVirtualApplianceSkus/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.NetworkVirtualApplianceSkuName)
}

// Segments returns a slice of Resource ID Segments which comprise this Network Virtual Appliance Sku ID
func (id NetworkVirtualApplianceSkuId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkVirtualApplianceSkus", "networkVirtualApplianceSkus", "networkVirtualApplianceSkus"),
		resourceids.UserSpecifiedSegment("networkVirtualApplianceSkuName", "networkVirtualApplianceSkuName"),
	}
}

// String returns a human-readable description of this Network Virtual Appliance Sku ID
func (id NetworkVirtualApplianceSkuId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Network Virtual Appliance Sku Name: %q", id.NetworkVirtualApplianceSkuName),
	}
	return fmt.Sprintf("Network Virtual Appliance Sku (%s)", strings.Join(components, "\n"))
}
