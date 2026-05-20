package networkconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NetworkConnectionId{})
}

var _ resourceids.ResourceId = &NetworkConnectionId{}

// NetworkConnectionId is a struct representing the Resource ID for a Network Connection
type NetworkConnectionId struct {
	SubscriptionId        string
	ResourceGroupName     string
	NetworkConnectionName string
}

// NewNetworkConnectionID returns a new NetworkConnectionId struct
func NewNetworkConnectionID(subscriptionId string, resourceGroupName string, networkConnectionName string) NetworkConnectionId {
	return NetworkConnectionId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		NetworkConnectionName: networkConnectionName,
	}
}

// ParseNetworkConnectionID parses 'input' into a NetworkConnectionId
func ParseNetworkConnectionID(input string) (*NetworkConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNetworkConnectionIDInsensitively parses 'input' case-insensitively into a NetworkConnectionId
// note: this method should only be used for API response data and not user input
func ParseNetworkConnectionIDInsensitively(input string) (*NetworkConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NetworkConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkConnectionName, ok = input.Parsed["networkConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkConnectionName", input)
	}

	return nil
}

// ValidateNetworkConnectionID checks that 'input' can be parsed as a Network Connection ID
func ValidateNetworkConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNetworkConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Network Connection ID
func (id NetworkConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/networkConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Network Connection ID
func (id NetworkConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticNetworkConnections", "networkConnections", "networkConnections"),
		resourceids.UserSpecifiedSegment("networkConnectionName", "networkConnectionName"),
	}
}

// String returns a human-readable description of this Network Connection ID
func (id NetworkConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Connection Name: %q", id.NetworkConnectionName),
	}
	return fmt.Sprintf("Network Connection (%s)", strings.Join(components, "\n"))
}
