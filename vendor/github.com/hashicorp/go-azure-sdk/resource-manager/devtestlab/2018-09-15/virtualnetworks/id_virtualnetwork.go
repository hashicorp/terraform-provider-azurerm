package virtualnetworks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VirtualNetworkId{})
}

var _ resourceids.ResourceId = &VirtualNetworkId{}

// VirtualNetworkId is a struct representing the Resource ID for a Virtual Network
type VirtualNetworkId struct {
	SubscriptionId     string
	ResourceGroupName  string
	LabName            string
	VirtualNetworkName string
}

// NewVirtualNetworkID returns a new VirtualNetworkId struct
func NewVirtualNetworkID(subscriptionId string, resourceGroupName string, labName string, virtualNetworkName string) VirtualNetworkId {
	return VirtualNetworkId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		LabName:            labName,
		VirtualNetworkName: virtualNetworkName,
	}
}

// ParseVirtualNetworkID parses 'input' into a VirtualNetworkId
func ParseVirtualNetworkID(input string) (*VirtualNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualNetworkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualNetworkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualNetworkIDInsensitively parses 'input' case-insensitively into a VirtualNetworkId
// note: this method should only be used for API response data and not user input
func ParseVirtualNetworkIDInsensitively(input string) (*VirtualNetworkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualNetworkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualNetworkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualNetworkId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.LabName, ok = input.Parsed["labName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "labName", input)
	}

	if id.VirtualNetworkName, ok = input.Parsed["virtualNetworkName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualNetworkName", input)
	}

	return nil
}

// ValidateVirtualNetworkID checks that 'input' can be parsed as a Virtual Network ID
func ValidateVirtualNetworkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualNetworkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Network ID
func (id VirtualNetworkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevTestLab/labs/%s/virtualNetworks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LabName, id.VirtualNetworkName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Network ID
func (id VirtualNetworkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevTestLab", "Microsoft.DevTestLab", "Microsoft.DevTestLab"),
		resourceids.StaticSegment("staticLabs", "labs", "labs"),
		resourceids.UserSpecifiedSegment("labName", "labName"),
		resourceids.StaticSegment("staticVirtualNetworks", "virtualNetworks", "virtualNetworks"),
		resourceids.UserSpecifiedSegment("virtualNetworkName", "virtualNetworkName"),
	}
}

// String returns a human-readable description of this Virtual Network ID
func (id VirtualNetworkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Lab Name: %q", id.LabName),
		fmt.Sprintf("Virtual Network Name: %q", id.VirtualNetworkName),
	}
	return fmt.Sprintf("Virtual Network (%s)", strings.Join(components, "\n"))
}
