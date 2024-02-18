// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &KubernetesFleetId{}

// KubernetesFleetId is a struct representing the Resource ID for a KubernetesFleet
type KubernetesFleetId struct {
	SubscriptionId    string
	ResourceGroupName string
	FleetName         string
}

// NewKubernetesFleetID returns a new KubernetesFleetId struct
func NewKubernetesFleetID(subscriptionId string, resourceGroupName string, fleetName string) KubernetesFleetId {
	return KubernetesFleetId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		FleetName:         fleetName,
	}
}

// ParseKubernetesFleetID parses 'input' into a KubernetesFleetId
func ParseKubernetesFleetID(input string) (*KubernetesFleetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&KubernetesFleetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := KubernetesFleetId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseKubernetesFleetIDInsensitively parses 'input' case-insensitively into a KubernetesFleetId
// note: this method should only be used for API response data and not user input
func ParseKubernetesFleetIDInsensitively(input string) (*KubernetesFleetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&KubernetesFleetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := KubernetesFleetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *KubernetesFleetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FleetName, ok = input.Parsed["fleetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "fleetName", input)
	}

	return nil
}

// ValidateKubernetesFleetID checks that 'input' can be parsed as a KubernetesFleet ID
func ValidateKubernetesFleetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseKubernetesFleetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted KubernetesFleet ID
func (id KubernetesFleetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/fleets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FleetName)
}

// Segments returns a slice of Resource ID Segments which comprise this KubernetesFleet ID
func (id KubernetesFleetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerService", "Microsoft.ContainerService", "Microsoft.ContainerService"),
		resourceids.StaticSegment("staticKubernetesFleets", "fleets", "fleets"),
		resourceids.UserSpecifiedSegment("fleetName", "fleetValue"),
	}
}

// String returns a human-readable description of this KubernetesFleet ID
func (id KubernetesFleetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("KubernetesFleet Name: %q", id.FleetName),
	}
	return fmt.Sprintf("KubernetesFleet (%s)", strings.Join(components, "\n"))
}
