package fluidrelaycontainers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = FluidRelayContainerId{}

// FluidRelayContainerId is a struct representing the Resource ID for a Fluid Relay Container
type FluidRelayContainerId struct {
	SubscriptionId          string
	ResourceGroup           string
	FluidRelayServerName    string
	FluidRelayContainerName string
}

// NewFluidRelayContainerID returns a new FluidRelayContainerId struct
func NewFluidRelayContainerID(subscriptionId string, resourceGroup string, fluidRelayServerName string, fluidRelayContainerName string) FluidRelayContainerId {
	return FluidRelayContainerId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		FluidRelayServerName:    fluidRelayServerName,
		FluidRelayContainerName: fluidRelayContainerName,
	}
}

// ParseFluidRelayContainerID parses 'input' into a FluidRelayContainerId
func ParseFluidRelayContainerID(input string) (*FluidRelayContainerId, error) {
	parser := resourceids.NewParserFromResourceIdType(FluidRelayContainerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FluidRelayContainerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroup, ok = parsed.Parsed["resourceGroup"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroup", *parsed)
	}

	if id.FluidRelayServerName, ok = parsed.Parsed["fluidRelayServerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fluidRelayServerName", *parsed)
	}

	if id.FluidRelayContainerName, ok = parsed.Parsed["fluidRelayContainerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fluidRelayContainerName", *parsed)
	}

	return &id, nil
}

// ParseFluidRelayContainerIDInsensitively parses 'input' case-insensitively into a FluidRelayContainerId
// note: this method should only be used for API response data and not user input
func ParseFluidRelayContainerIDInsensitively(input string) (*FluidRelayContainerId, error) {
	parser := resourceids.NewParserFromResourceIdType(FluidRelayContainerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FluidRelayContainerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroup, ok = parsed.Parsed["resourceGroup"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroup", *parsed)
	}

	if id.FluidRelayServerName, ok = parsed.Parsed["fluidRelayServerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fluidRelayServerName", *parsed)
	}

	if id.FluidRelayContainerName, ok = parsed.Parsed["fluidRelayContainerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "fluidRelayContainerName", *parsed)
	}

	return &id, nil
}

// ValidateFluidRelayContainerID checks that 'input' can be parsed as a Fluid Relay Container ID
func ValidateFluidRelayContainerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFluidRelayContainerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Fluid Relay Container ID
func (id FluidRelayContainerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.FluidRelay/fluidRelayServers/%s/fluidRelayContainers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FluidRelayServerName, id.FluidRelayContainerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Fluid Relay Container ID
func (id FluidRelayContainerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroup", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftFluidRelay", "Microsoft.FluidRelay", "Microsoft.FluidRelay"),
		resourceids.StaticSegment("staticFluidRelayServers", "fluidRelayServers", "fluidRelayServers"),
		resourceids.UserSpecifiedSegment("fluidRelayServerName", "fluidRelayServerValue"),
		resourceids.StaticSegment("staticFluidRelayContainers", "fluidRelayContainers", "fluidRelayContainers"),
		resourceids.UserSpecifiedSegment("fluidRelayContainerName", "fluidRelayContainerValue"),
	}
}

// String returns a human-readable description of this Fluid Relay Container ID
func (id FluidRelayContainerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group: %q", id.ResourceGroup),
		fmt.Sprintf("Fluid Relay Server Name: %q", id.FluidRelayServerName),
		fmt.Sprintf("Fluid Relay Container Name: %q", id.FluidRelayContainerName),
	}
	return fmt.Sprintf("Fluid Relay Container (%s)", strings.Join(components, "\n"))
}
