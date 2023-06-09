package vpnsites

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = VpnSiteId{}

// VpnSiteId is a struct representing the Resource ID for a Vpn Site
type VpnSiteId struct {
	SubscriptionId    string
	ResourceGroupName string
	VpnSiteName       string
}

// NewVpnSiteID returns a new VpnSiteId struct
func NewVpnSiteID(subscriptionId string, resourceGroupName string, vpnSiteName string) VpnSiteId {
	return VpnSiteId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		VpnSiteName:       vpnSiteName,
	}
}

// ParseVpnSiteID parses 'input' into a VpnSiteId
func ParseVpnSiteID(input string) (*VpnSiteId, error) {
	parser := resourceids.NewParserFromResourceIdType(VpnSiteId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VpnSiteId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VpnSiteName, ok = parsed.Parsed["vpnSiteName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnSiteName", *parsed)
	}

	return &id, nil
}

// ParseVpnSiteIDInsensitively parses 'input' case-insensitively into a VpnSiteId
// note: this method should only be used for API response data and not user input
func ParseVpnSiteIDInsensitively(input string) (*VpnSiteId, error) {
	parser := resourceids.NewParserFromResourceIdType(VpnSiteId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VpnSiteId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VpnSiteName, ok = parsed.Parsed["vpnSiteName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnSiteName", *parsed)
	}

	return &id, nil
}

// ValidateVpnSiteID checks that 'input' can be parsed as a Vpn Site ID
func ValidateVpnSiteID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVpnSiteID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Vpn Site ID
func (id VpnSiteId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnSites/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VpnSiteName)
}

// Segments returns a slice of Resource ID Segments which comprise this Vpn Site ID
func (id VpnSiteId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVpnSites", "vpnSites", "vpnSites"),
		resourceids.UserSpecifiedSegment("vpnSiteName", "vpnSiteValue"),
	}
}

// String returns a human-readable description of this Vpn Site ID
func (id VpnSiteId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vpn Site Name: %q", id.VpnSiteName),
	}
	return fmt.Sprintf("Vpn Site (%s)", strings.Join(components, "\n"))
}
