package site

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SiteId{}

// SiteId is a struct representing the Resource ID for a Site
type SiteId struct {
	SubscriptionId    string
	ResourceGroupName string
	MobileNetworkName string
	SiteName          string
}

// NewSiteID returns a new SiteId struct
func NewSiteID(subscriptionId string, resourceGroupName string, mobileNetworkName string, siteName string) SiteId {
	return SiteId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MobileNetworkName: mobileNetworkName,
		SiteName:          siteName,
	}
}

// ParseSiteID parses 'input' into a SiteId
func ParseSiteID(input string) (*SiteId, error) {
	parser := resourceids.NewParserFromResourceIdType(SiteId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SiteId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MobileNetworkName, ok = parsed.Parsed["mobileNetworkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mobileNetworkName", *parsed)
	}

	if id.SiteName, ok = parsed.Parsed["siteName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "siteName", *parsed)
	}

	return &id, nil
}

// ParseSiteIDInsensitively parses 'input' case-insensitively into a SiteId
// note: this method should only be used for API response data and not user input
func ParseSiteIDInsensitively(input string) (*SiteId, error) {
	parser := resourceids.NewParserFromResourceIdType(SiteId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SiteId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MobileNetworkName, ok = parsed.Parsed["mobileNetworkName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mobileNetworkName", *parsed)
	}

	if id.SiteName, ok = parsed.Parsed["siteName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "siteName", *parsed)
	}

	return &id, nil
}

// ValidateSiteID checks that 'input' can be parsed as a Site ID
func ValidateSiteID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSiteID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Site ID
func (id SiteId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MobileNetwork/mobileNetworks/%s/sites/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MobileNetworkName, id.SiteName)
}

// Segments returns a slice of Resource ID Segments which comprise this Site ID
func (id SiteId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMobileNetwork", "Microsoft.MobileNetwork", "Microsoft.MobileNetwork"),
		resourceids.StaticSegment("staticMobileNetworks", "mobileNetworks", "mobileNetworks"),
		resourceids.UserSpecifiedSegment("mobileNetworkName", "mobileNetworkValue"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteValue"),
	}
}

// String returns a human-readable description of this Site ID
func (id SiteId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Mobile Network Name: %q", id.MobileNetworkName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
	}
	return fmt.Sprintf("Site (%s)", strings.Join(components, "\n"))
}
