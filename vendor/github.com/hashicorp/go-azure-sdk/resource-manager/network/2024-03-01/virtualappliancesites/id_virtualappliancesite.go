package virtualappliancesites

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VirtualApplianceSiteId{})
}

var _ resourceids.ResourceId = &VirtualApplianceSiteId{}

// VirtualApplianceSiteId is a struct representing the Resource ID for a Virtual Appliance Site
type VirtualApplianceSiteId struct {
	SubscriptionId              string
	ResourceGroupName           string
	NetworkVirtualApplianceName string
	VirtualApplianceSiteName    string
}

// NewVirtualApplianceSiteID returns a new VirtualApplianceSiteId struct
func NewVirtualApplianceSiteID(subscriptionId string, resourceGroupName string, networkVirtualApplianceName string, virtualApplianceSiteName string) VirtualApplianceSiteId {
	return VirtualApplianceSiteId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		NetworkVirtualApplianceName: networkVirtualApplianceName,
		VirtualApplianceSiteName:    virtualApplianceSiteName,
	}
}

// ParseVirtualApplianceSiteID parses 'input' into a VirtualApplianceSiteId
func ParseVirtualApplianceSiteID(input string) (*VirtualApplianceSiteId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualApplianceSiteId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualApplianceSiteId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVirtualApplianceSiteIDInsensitively parses 'input' case-insensitively into a VirtualApplianceSiteId
// note: this method should only be used for API response data and not user input
func ParseVirtualApplianceSiteIDInsensitively(input string) (*VirtualApplianceSiteId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VirtualApplianceSiteId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VirtualApplianceSiteId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VirtualApplianceSiteId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkVirtualApplianceName, ok = input.Parsed["networkVirtualApplianceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkVirtualApplianceName", input)
	}

	if id.VirtualApplianceSiteName, ok = input.Parsed["virtualApplianceSiteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "virtualApplianceSiteName", input)
	}

	return nil
}

// ValidateVirtualApplianceSiteID checks that 'input' can be parsed as a Virtual Appliance Site ID
func ValidateVirtualApplianceSiteID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualApplianceSiteID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Appliance Site ID
func (id VirtualApplianceSiteId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkVirtualAppliances/%s/virtualApplianceSites/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkVirtualApplianceName, id.VirtualApplianceSiteName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Appliance Site ID
func (id VirtualApplianceSiteId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkVirtualAppliances", "networkVirtualAppliances", "networkVirtualAppliances"),
		resourceids.UserSpecifiedSegment("networkVirtualApplianceName", "networkVirtualApplianceName"),
		resourceids.StaticSegment("staticVirtualApplianceSites", "virtualApplianceSites", "virtualApplianceSites"),
		resourceids.UserSpecifiedSegment("virtualApplianceSiteName", "virtualApplianceSiteName"),
	}
}

// String returns a human-readable description of this Virtual Appliance Site ID
func (id VirtualApplianceSiteId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Virtual Appliance Name: %q", id.NetworkVirtualApplianceName),
		fmt.Sprintf("Virtual Appliance Site Name: %q", id.VirtualApplianceSiteName),
	}
	return fmt.Sprintf("Virtual Appliance Site (%s)", strings.Join(components, "\n"))
}
