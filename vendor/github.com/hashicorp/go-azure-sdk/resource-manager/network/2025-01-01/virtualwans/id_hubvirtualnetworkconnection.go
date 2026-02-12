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
	recaser.RegisterResourceId(&HubVirtualNetworkConnectionId{})
}

var _ resourceids.ResourceId = &HubVirtualNetworkConnectionId{}

// HubVirtualNetworkConnectionId is a struct representing the Resource ID for a Hub Virtual Network Connection
type HubVirtualNetworkConnectionId struct {
	SubscriptionId                  string
	ResourceGroupName               string
	VirtualHubName                  string
	HubVirtualNetworkConnectionName string
}

// NewHubVirtualNetworkConnectionID returns a new HubVirtualNetworkConnectionId struct
func NewHubVirtualNetworkConnectionID(subscriptionId string, resourceGroupName string, virtualHubName string, hubVirtualNetworkConnectionName string) HubVirtualNetworkConnectionId {
	return HubVirtualNetworkConnectionId{
		SubscriptionId:                  subscriptionId,
		ResourceGroupName:               resourceGroupName,
		VirtualHubName:                  virtualHubName,
		HubVirtualNetworkConnectionName: hubVirtualNetworkConnectionName,
	}
}

// ParseHubVirtualNetworkConnectionID parses 'input' into a HubVirtualNetworkConnectionId
func ParseHubVirtualNetworkConnectionID(input string) (*HubVirtualNetworkConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HubVirtualNetworkConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HubVirtualNetworkConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseHubVirtualNetworkConnectionIDInsensitively parses 'input' case-insensitively into a HubVirtualNetworkConnectionId
// note: this method should only be used for API response data and not user input
func ParseHubVirtualNetworkConnectionIDInsensitively(input string) (*HubVirtualNetworkConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&HubVirtualNetworkConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := HubVirtualNetworkConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *HubVirtualNetworkConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VirtualHubName, ok = input.Parsed["virtualHubName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualHubName", input)
	}

	if id.HubVirtualNetworkConnectionName, ok = input.Parsed["hubVirtualNetworkConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hubVirtualNetworkConnectionName", input)
	}

	return nil
}

// ValidateHubVirtualNetworkConnectionID checks that 'input' can be parsed as a Hub Virtual Network Connection ID
func ValidateHubVirtualNetworkConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHubVirtualNetworkConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Hub Virtual Network Connection ID
func (id HubVirtualNetworkConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualHubs/%s/hubVirtualNetworkConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName, id.HubVirtualNetworkConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Hub Virtual Network Connection ID
func (id HubVirtualNetworkConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVirtualHubs", "virtualHubs", "virtualHubs"),
		resourceids.UserSpecifiedSegment("virtualHubName", "virtualHubName"),
		resourceids.StaticSegment("staticHubVirtualNetworkConnections", "hubVirtualNetworkConnections", "hubVirtualNetworkConnections"),
		resourceids.UserSpecifiedSegment("hubVirtualNetworkConnectionName", "hubVirtualNetworkConnectionName"),
	}
}

// String returns a human-readable description of this Hub Virtual Network Connection ID
func (id HubVirtualNetworkConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Hub Name: %q", id.VirtualHubName),
		fmt.Sprintf("Hub Virtual Network Connection Name: %q", id.HubVirtualNetworkConnectionName),
	}
	return fmt.Sprintf("Hub Virtual Network Connection (%s)", strings.Join(components, "\n"))
}
