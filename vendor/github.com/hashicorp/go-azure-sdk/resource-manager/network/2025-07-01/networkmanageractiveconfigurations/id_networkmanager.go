package networkmanageractiveconfigurations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NetworkManagerId{})
}

var _ resourceids.ResourceId = &NetworkManagerId{}

// NetworkManagerId is a struct representing the Resource ID for a Network Manager
type NetworkManagerId struct {
	SubscriptionId     string
	ResourceGroupName  string
	NetworkManagerName string
}

// NewNetworkManagerID returns a new NetworkManagerId struct
func NewNetworkManagerID(subscriptionId string, resourceGroupName string, networkManagerName string) NetworkManagerId {
	return NetworkManagerId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		NetworkManagerName: networkManagerName,
	}
}

// ParseNetworkManagerID parses 'input' into a NetworkManagerId
func ParseNetworkManagerID(input string) (*NetworkManagerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkManagerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkManagerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNetworkManagerIDInsensitively parses 'input' case-insensitively into a NetworkManagerId
// note: this method should only be used for API response data and not user input
func ParseNetworkManagerIDInsensitively(input string) (*NetworkManagerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkManagerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkManagerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NetworkManagerId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkManagerName, ok = input.Parsed["networkManagerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkManagerName", input)
	}

	return nil
}

// ValidateNetworkManagerID checks that 'input' can be parsed as a Network Manager ID
func ValidateNetworkManagerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNetworkManagerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Network Manager ID
func (id NetworkManagerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Network Manager ID
func (id NetworkManagerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkManagers", "networkManagers", "networkManagers"),
		resourceids.UserSpecifiedSegment("networkManagerName", "networkManagerName"),
	}
}

// String returns a human-readable description of this Network Manager ID
func (id NetworkManagerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Manager Name: %q", id.NetworkManagerName),
	}
	return fmt.Sprintf("Network Manager (%s)", strings.Join(components, "\n"))
}
