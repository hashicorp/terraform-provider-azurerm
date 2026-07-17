package networkanchors

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NetworkAnchorId{})
}

var _ resourceids.ResourceId = &NetworkAnchorId{}

// NetworkAnchorId is a struct representing the Resource ID for a Network Anchor
type NetworkAnchorId struct {
	SubscriptionId    string
	ResourceGroupName string
	NetworkAnchorName string
}

// NewNetworkAnchorID returns a new NetworkAnchorId struct
func NewNetworkAnchorID(subscriptionId string, resourceGroupName string, networkAnchorName string) NetworkAnchorId {
	return NetworkAnchorId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NetworkAnchorName: networkAnchorName,
	}
}

// ParseNetworkAnchorID parses 'input' into a NetworkAnchorId
func ParseNetworkAnchorID(input string) (*NetworkAnchorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkAnchorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkAnchorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNetworkAnchorIDInsensitively parses 'input' case-insensitively into a NetworkAnchorId
// note: this method should only be used for API response data and not user input
func ParseNetworkAnchorIDInsensitively(input string) (*NetworkAnchorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkAnchorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkAnchorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NetworkAnchorId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkAnchorName, ok = input.Parsed["networkAnchorName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkAnchorName", input)
	}

	return nil
}

// ValidateNetworkAnchorID checks that 'input' can be parsed as a Network Anchor ID
func ValidateNetworkAnchorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNetworkAnchorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Network Anchor ID
func (id NetworkAnchorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/networkAnchors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkAnchorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Network Anchor ID
func (id NetworkAnchorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticNetworkAnchors", "networkAnchors", "networkAnchors"),
		resourceids.UserSpecifiedSegment("networkAnchorName", "networkAnchorName"),
	}
}

// String returns a human-readable description of this Network Anchor ID
func (id NetworkAnchorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Anchor Name: %q", id.NetworkAnchorName),
	}
	return fmt.Sprintf("Network Anchor (%s)", strings.Join(components, "\n"))
}
