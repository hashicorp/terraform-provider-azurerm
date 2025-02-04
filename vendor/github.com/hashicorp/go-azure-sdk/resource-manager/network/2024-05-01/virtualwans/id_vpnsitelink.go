package virtualwans

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VpnSiteLinkId{})
}

var _ resourceids.ResourceId = &VpnSiteLinkId{}

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
	parser := resourceids.NewParserFromResourceIdType(&VpnSiteLinkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VpnSiteLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVpnSiteLinkIDInsensitively parses 'input' case-insensitively into a VpnSiteLinkId
// note: this method should only be used for API response data and not user input
func ParseVpnSiteLinkIDInsensitively(input string) (*VpnSiteLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VpnSiteLinkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VpnSiteLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VpnSiteLinkId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VpnSiteName, ok = input.Parsed["vpnSiteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "vpnSiteName", input)
	}

	if id.VpnSiteLinkName, ok = input.Parsed["vpnSiteLinkName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "vpnSiteLinkName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("vpnSiteName", "vpnSiteName"),
		resourceids.StaticSegment("staticVpnSiteLinks", "vpnSiteLinks", "vpnSiteLinks"),
		resourceids.UserSpecifiedSegment("vpnSiteLinkName", "vpnSiteLinkName"),
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
