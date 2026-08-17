package networksecurityperimeterlinks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NetworkSecurityPerimeterLinkId{})
}

var _ resourceids.ResourceId = &NetworkSecurityPerimeterLinkId{}

// NetworkSecurityPerimeterLinkId is a struct representing the Resource ID for a Network Security Perimeter Link
type NetworkSecurityPerimeterLinkId struct {
	SubscriptionId               string
	ResourceGroupName            string
	NetworkSecurityPerimeterName string
	LinkName                     string
}

// NewNetworkSecurityPerimeterLinkID returns a new NetworkSecurityPerimeterLinkId struct
func NewNetworkSecurityPerimeterLinkID(subscriptionId string, resourceGroupName string, networkSecurityPerimeterName string, linkName string) NetworkSecurityPerimeterLinkId {
	return NetworkSecurityPerimeterLinkId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		NetworkSecurityPerimeterName: networkSecurityPerimeterName,
		LinkName:                     linkName,
	}
}

// ParseNetworkSecurityPerimeterLinkID parses 'input' into a NetworkSecurityPerimeterLinkId
func ParseNetworkSecurityPerimeterLinkID(input string) (*NetworkSecurityPerimeterLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkSecurityPerimeterLinkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkSecurityPerimeterLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNetworkSecurityPerimeterLinkIDInsensitively parses 'input' case-insensitively into a NetworkSecurityPerimeterLinkId
// note: this method should only be used for API response data and not user input
func ParseNetworkSecurityPerimeterLinkIDInsensitively(input string) (*NetworkSecurityPerimeterLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkSecurityPerimeterLinkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkSecurityPerimeterLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NetworkSecurityPerimeterLinkId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.LinkName, ok = input.Parsed["linkName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "linkName", input)
	}

	return nil
}

// ValidateNetworkSecurityPerimeterLinkID checks that 'input' can be parsed as a Network Security Perimeter Link ID
func ValidateNetworkSecurityPerimeterLinkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNetworkSecurityPerimeterLinkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Network Security Perimeter Link ID
func (id NetworkSecurityPerimeterLinkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkSecurityPerimeters/%s/links/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkSecurityPerimeterName, id.LinkName)
}

// Segments returns a slice of Resource ID Segments which comprise this Network Security Perimeter Link ID
func (id NetworkSecurityPerimeterLinkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkSecurityPerimeters", "networkSecurityPerimeters", "networkSecurityPerimeters"),
		resourceids.UserSpecifiedSegment("networkSecurityPerimeterName", "networkSecurityPerimeterName"),
		resourceids.StaticSegment("staticLinks", "links", "links"),
		resourceids.UserSpecifiedSegment("linkName", "linkName"),
	}
}

// String returns a human-readable description of this Network Security Perimeter Link ID
func (id NetworkSecurityPerimeterLinkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Security Perimeter Name: %q", id.NetworkSecurityPerimeterName),
		fmt.Sprintf("Link Name: %q", id.LinkName),
	}
	return fmt.Sprintf("Network Security Perimeter Link (%s)", strings.Join(components, "\n"))
}
