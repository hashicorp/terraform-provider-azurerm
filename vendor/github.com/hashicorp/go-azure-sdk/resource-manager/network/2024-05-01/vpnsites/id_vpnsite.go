package vpnsites

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&VpnSiteId{})
}

var _ resourceids.ResourceId = &VpnSiteId{}

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
	parser := resourceids.NewParserFromResourceIdType(&VpnSiteId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VpnSiteId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVpnSiteIDInsensitively parses 'input' case-insensitively into a VpnSiteId
// note: this method should only be used for API response data and not user input
func ParseVpnSiteIDInsensitively(input string) (*VpnSiteId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VpnSiteId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VpnSiteId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VpnSiteId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
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
		resourceids.UserSpecifiedSegment("vpnSiteName", "vpnSiteName"),
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
