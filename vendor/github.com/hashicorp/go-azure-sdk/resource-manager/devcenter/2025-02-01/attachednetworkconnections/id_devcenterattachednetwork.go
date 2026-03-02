package attachednetworkconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DevCenterAttachedNetworkId{})
}

var _ resourceids.ResourceId = &DevCenterAttachedNetworkId{}

// DevCenterAttachedNetworkId is a struct representing the Resource ID for a Dev Center Attached Network
type DevCenterAttachedNetworkId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DevCenterName       string
	AttachedNetworkName string
}

// NewDevCenterAttachedNetworkID returns a new DevCenterAttachedNetworkId struct
func NewDevCenterAttachedNetworkID(subscriptionId string, resourceGroupName string, devCenterName string, attachedNetworkName string) DevCenterAttachedNetworkId {
	return DevCenterAttachedNetworkId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DevCenterName:       devCenterName,
		AttachedNetworkName: attachedNetworkName,
	}
}

// ParseDevCenterAttachedNetworkID parses 'input' into a DevCenterAttachedNetworkId
func ParseDevCenterAttachedNetworkID(input string) (*DevCenterAttachedNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DevCenterAttachedNetworkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DevCenterAttachedNetworkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDevCenterAttachedNetworkIDInsensitively parses 'input' case-insensitively into a DevCenterAttachedNetworkId
// note: this method should only be used for API response data and not user input
func ParseDevCenterAttachedNetworkIDInsensitively(input string) (*DevCenterAttachedNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DevCenterAttachedNetworkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DevCenterAttachedNetworkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DevCenterAttachedNetworkId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DevCenterName, ok = input.Parsed["devCenterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "devCenterName", input)
	}

	if id.AttachedNetworkName, ok = input.Parsed["attachedNetworkName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "attachedNetworkName", input)
	}

	return nil
}

// ValidateDevCenterAttachedNetworkID checks that 'input' can be parsed as a Dev Center Attached Network ID
func ValidateDevCenterAttachedNetworkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDevCenterAttachedNetworkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dev Center Attached Network ID
func (id DevCenterAttachedNetworkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/devCenters/%s/attachedNetworks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DevCenterName, id.AttachedNetworkName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dev Center Attached Network ID
func (id DevCenterAttachedNetworkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticDevCenters", "devCenters", "devCenters"),
		resourceids.UserSpecifiedSegment("devCenterName", "devCenterName"),
		resourceids.StaticSegment("staticAttachedNetworks", "attachedNetworks", "attachedNetworks"),
		resourceids.UserSpecifiedSegment("attachedNetworkName", "attachedNetworkName"),
	}
}

// String returns a human-readable description of this Dev Center Attached Network ID
func (id DevCenterAttachedNetworkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dev Center Name: %q", id.DevCenterName),
		fmt.Sprintf("Attached Network Name: %q", id.AttachedNetworkName),
	}
	return fmt.Sprintf("Dev Center Attached Network (%s)", strings.Join(components, "\n"))
}
