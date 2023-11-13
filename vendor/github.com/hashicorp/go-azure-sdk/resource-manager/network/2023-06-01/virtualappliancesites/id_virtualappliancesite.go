package virtualappliancesites

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = VirtualApplianceSiteId{}

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
	parser := resourceids.NewParserFromResourceIdType(VirtualApplianceSiteId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualApplianceSiteId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkVirtualApplianceName, ok = parsed.Parsed["networkVirtualApplianceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkVirtualApplianceName", *parsed)
	}

	if id.VirtualApplianceSiteName, ok = parsed.Parsed["virtualApplianceSiteName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualApplianceSiteName", *parsed)
	}

	return &id, nil
}

// ParseVirtualApplianceSiteIDInsensitively parses 'input' case-insensitively into a VirtualApplianceSiteId
// note: this method should only be used for API response data and not user input
func ParseVirtualApplianceSiteIDInsensitively(input string) (*VirtualApplianceSiteId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualApplianceSiteId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualApplianceSiteId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkVirtualApplianceName, ok = parsed.Parsed["networkVirtualApplianceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkVirtualApplianceName", *parsed)
	}

	if id.VirtualApplianceSiteName, ok = parsed.Parsed["virtualApplianceSiteName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "virtualApplianceSiteName", *parsed)
	}

	return &id, nil
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
		resourceids.UserSpecifiedSegment("networkVirtualApplianceName", "networkVirtualApplianceValue"),
		resourceids.StaticSegment("staticVirtualApplianceSites", "virtualApplianceSites", "virtualApplianceSites"),
		resourceids.UserSpecifiedSegment("virtualApplianceSiteName", "virtualApplianceSiteValue"),
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
