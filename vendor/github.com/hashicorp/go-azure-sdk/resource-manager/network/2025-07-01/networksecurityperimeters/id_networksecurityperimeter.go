package networksecurityperimeters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NetworkSecurityPerimeterId{})
}

var _ resourceids.ResourceId = &NetworkSecurityPerimeterId{}

// NetworkSecurityPerimeterId is a struct representing the Resource ID for a Network Security Perimeter
type NetworkSecurityPerimeterId struct {
	SubscriptionId               string
	ResourceGroupName            string
	NetworkSecurityPerimeterName string
}

// NewNetworkSecurityPerimeterID returns a new NetworkSecurityPerimeterId struct
func NewNetworkSecurityPerimeterID(subscriptionId string, resourceGroupName string, networkSecurityPerimeterName string) NetworkSecurityPerimeterId {
	return NetworkSecurityPerimeterId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		NetworkSecurityPerimeterName: networkSecurityPerimeterName,
	}
}

// ParseNetworkSecurityPerimeterID parses 'input' into a NetworkSecurityPerimeterId
func ParseNetworkSecurityPerimeterID(input string) (*NetworkSecurityPerimeterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkSecurityPerimeterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkSecurityPerimeterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNetworkSecurityPerimeterIDInsensitively parses 'input' case-insensitively into a NetworkSecurityPerimeterId
// note: this method should only be used for API response data and not user input
func ParseNetworkSecurityPerimeterIDInsensitively(input string) (*NetworkSecurityPerimeterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkSecurityPerimeterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkSecurityPerimeterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NetworkSecurityPerimeterId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkSecurityPerimeterName, ok = input.Parsed["networkSecurityPerimeterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkSecurityPerimeterName", input)
	}

	return nil
}

// ValidateNetworkSecurityPerimeterID checks that 'input' can be parsed as a Network Security Perimeter ID
func ValidateNetworkSecurityPerimeterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNetworkSecurityPerimeterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Network Security Perimeter ID
func (id NetworkSecurityPerimeterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkSecurityPerimeters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkSecurityPerimeterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Network Security Perimeter ID
func (id NetworkSecurityPerimeterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkSecurityPerimeters", "networkSecurityPerimeters", "networkSecurityPerimeters"),
		resourceids.UserSpecifiedSegment("networkSecurityPerimeterName", "networkSecurityPerimeterName"),
	}
}

// String returns a human-readable description of this Network Security Perimeter ID
func (id NetworkSecurityPerimeterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Security Perimeter Name: %q", id.NetworkSecurityPerimeterName),
	}
	return fmt.Sprintf("Network Security Perimeter (%s)", strings.Join(components, "\n"))
}
