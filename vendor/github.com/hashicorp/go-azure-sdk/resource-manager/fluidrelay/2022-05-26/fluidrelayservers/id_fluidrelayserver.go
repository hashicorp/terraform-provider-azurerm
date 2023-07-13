package fluidrelayservers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = FluidRelayServerId{}

// FluidRelayServerId is a struct representing the Resource ID for a Fluid Relay Server
type FluidRelayServerId struct {
	SubscriptionId       string
	ResourceGroup        string
	FluidRelayServerName string
}

// NewFluidRelayServerID returns a new FluidRelayServerId struct
func NewFluidRelayServerID(subscriptionId string, resourceGroup string, fluidRelayServerName string) FluidRelayServerId {
	return FluidRelayServerId{
		SubscriptionId:       subscriptionId,
		ResourceGroup:        resourceGroup,
		FluidRelayServerName: fluidRelayServerName,
	}
}

// ParseFluidRelayServerID parses 'input' into a FluidRelayServerId
func ParseFluidRelayServerID(input string) (*FluidRelayServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(FluidRelayServerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FluidRelayServerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroup, ok = parsed.Parsed["resourceGroup"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroup", *parsed)
	}

	if id.FluidRelayServerName, ok = parsed.Parsed["fluidRelayServerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fluidRelayServerName", *parsed)
	}

	return &id, nil
}

// ParseFluidRelayServerIDInsensitively parses 'input' case-insensitively into a FluidRelayServerId
// note: this method should only be used for API response data and not user input
func ParseFluidRelayServerIDInsensitively(input string) (*FluidRelayServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(FluidRelayServerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FluidRelayServerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroup, ok = parsed.Parsed["resourceGroup"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroup", *parsed)
	}

	if id.FluidRelayServerName, ok = parsed.Parsed["fluidRelayServerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fluidRelayServerName", *parsed)
	}

	return &id, nil
}

// ValidateFluidRelayServerID checks that 'input' can be parsed as a Fluid Relay Server ID
func ValidateFluidRelayServerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFluidRelayServerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Fluid Relay Server ID
func (id FluidRelayServerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.FluidRelay/fluidRelayServers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FluidRelayServerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Fluid Relay Server ID
func (id FluidRelayServerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroup", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftFluidRelay", "Microsoft.FluidRelay", "Microsoft.FluidRelay"),
		resourceids.StaticSegment("staticFluidRelayServers", "fluidRelayServers", "fluidRelayServers"),
		resourceids.UserSpecifiedSegment("fluidRelayServerName", "fluidRelayServerValue"),
	}
}

// String returns a human-readable description of this Fluid Relay Server ID
func (id FluidRelayServerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group: %q", id.ResourceGroup),
		fmt.Sprintf("Fluid Relay Server Name: %q", id.FluidRelayServerName),
	}
	return fmt.Sprintf("Fluid Relay Server (%s)", strings.Join(components, "\n"))
}
