package networkprofiles

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NetworkProfileId{})
}

var _ resourceids.ResourceId = &NetworkProfileId{}

// NetworkProfileId is a struct representing the Resource ID for a Network Profile
type NetworkProfileId struct {
	SubscriptionId     string
	ResourceGroupName  string
	NetworkProfileName string
}

// NewNetworkProfileID returns a new NetworkProfileId struct
func NewNetworkProfileID(subscriptionId string, resourceGroupName string, networkProfileName string) NetworkProfileId {
	return NetworkProfileId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		NetworkProfileName: networkProfileName,
	}
}

// ParseNetworkProfileID parses 'input' into a NetworkProfileId
func ParseNetworkProfileID(input string) (*NetworkProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkProfileId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNetworkProfileIDInsensitively parses 'input' case-insensitively into a NetworkProfileId
// note: this method should only be used for API response data and not user input
func ParseNetworkProfileIDInsensitively(input string) (*NetworkProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkProfileId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NetworkProfileId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkProfileName, ok = input.Parsed["networkProfileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkProfileName", input)
	}

	return nil
}

// ValidateNetworkProfileID checks that 'input' can be parsed as a Network Profile ID
func ValidateNetworkProfileID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNetworkProfileID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Network Profile ID
func (id NetworkProfileId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkProfiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkProfileName)
}

// Segments returns a slice of Resource ID Segments which comprise this Network Profile ID
func (id NetworkProfileId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkProfiles", "networkProfiles", "networkProfiles"),
		resourceids.UserSpecifiedSegment("networkProfileName", "networkProfileName"),
	}
}

// String returns a human-readable description of this Network Profile ID
func (id NetworkProfileId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Profile Name: %q", id.NetworkProfileName),
	}
	return fmt.Sprintf("Network Profile (%s)", strings.Join(components, "\n"))
}
