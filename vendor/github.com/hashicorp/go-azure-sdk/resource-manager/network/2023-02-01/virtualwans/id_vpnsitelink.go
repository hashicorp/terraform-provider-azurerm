package virtualwans

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = VpnSiteLinkId{}

// VpnSiteLinkId is a struct representing the Resource ID for a Vpn Site Link
type VpnSiteLinkId struct {
	SubscriptionId    string
	ResourceGroupName string
	VpnSiteName       string
	VpnSiteLinkName   string
}

// NewVpnSiteLinkID returns a new VpnSiteLinkId struct
func NewVpnSiteLinkID(subscriptionId string, resourceGroupName string, vpnSiteName string, vpnSiteLinkName string) VpnSiteLinkId {
	return VpnSiteLinkId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VpnSiteName:       vpnSiteName,
		VpnSiteLinkName:   vpnSiteLinkName,
	}
}

// ParseVpnSiteLinkID parses 'input' into a VpnSiteLinkId
func ParseVpnSiteLinkID(input string) (*VpnSiteLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(VpnSiteLinkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VpnSiteLinkId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VpnSiteName, ok = parsed.Parsed["vpnSiteName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnSiteName", *parsed)
	}

	if id.VpnSiteLinkName, ok = parsed.Parsed["vpnSiteLinkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnSiteLinkName", *parsed)
	}

	return &id, nil
}

// ParseVpnSiteLinkIDInsensitively parses 'input' case-insensitively into a VpnSiteLinkId
// note: this method should only be used for API response data and not user input
func ParseVpnSiteLinkIDInsensitively(input string) (*VpnSiteLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(VpnSiteLinkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VpnSiteLinkId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VpnSiteName, ok = parsed.Parsed["vpnSiteName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnSiteName", *parsed)
	}

	if id.VpnSiteLinkName, ok = parsed.Parsed["vpnSiteLinkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnSiteLinkName", *parsed)
	}

	return &id, nil
}

// ValidateVpnSiteLinkID checks that 'input' can be parsed as a Vpn Site Link ID
func ValidateVpnSiteLinkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVpnSiteLinkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Vpn Site Link ID
func (id VpnSiteLinkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnSites/%s/vpnSiteLinks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VpnSiteName, id.VpnSiteLinkName)
}

// Segments returns a slice of Resource ID Segments which comprise this Vpn Site Link ID
func (id VpnSiteLinkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVpnSites", "vpnSites", "vpnSites"),
		resourceids.UserSpecifiedSegment("vpnSiteName", "vpnSiteValue"),
		resourceids.StaticSegment("staticVpnSiteLinks", "vpnSiteLinks", "vpnSiteLinks"),
		resourceids.UserSpecifiedSegment("vpnSiteLinkName", "vpnSiteLinkValue"),
	}
}

// String returns a human-readable description of this Vpn Site Link ID
func (id VpnSiteLinkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vpn Site Name: %q", id.VpnSiteName),
		fmt.Sprintf("Vpn Site Link Name: %q", id.VpnSiteLinkName),
	}
	return fmt.Sprintf("Vpn Site Link (%s)", strings.Join(components, "\n"))
}
