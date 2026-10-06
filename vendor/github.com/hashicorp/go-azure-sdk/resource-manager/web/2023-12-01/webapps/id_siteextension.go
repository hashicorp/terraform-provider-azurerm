package webapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SiteExtensionId{})
}

var _ resourceids.ResourceId = &SiteExtensionId{}

// SiteExtensionId is a struct representing the Resource ID for a Site Extension
type SiteExtensionId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	SiteExtensionId   string
}

// NewSiteExtensionID returns a new SiteExtensionId struct
func NewSiteExtensionID(subscriptionId string, resourceGroupName string, siteName string, siteExtensionId string) SiteExtensionId {
	return SiteExtensionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		SiteExtensionId:   siteExtensionId,
	}
}

// ParseSiteExtensionID parses 'input' into a SiteExtensionId
func ParseSiteExtensionID(input string) (*SiteExtensionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SiteExtensionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SiteExtensionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSiteExtensionIDInsensitively parses 'input' case-insensitively into a SiteExtensionId
// note: this method should only be used for API response data and not user input
func ParseSiteExtensionIDInsensitively(input string) (*SiteExtensionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SiteExtensionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SiteExtensionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SiteExtensionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SiteName, ok = input.Parsed["siteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "siteName", input)
	}

	if id.SiteExtensionId, ok = input.Parsed["siteExtensionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "siteExtensionId", input)
	}

	return nil
}

// ValidateSiteExtensionID checks that 'input' can be parsed as a Site Extension ID
func ValidateSiteExtensionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSiteExtensionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Site Extension ID
func (id SiteExtensionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/siteExtensions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SiteExtensionId)
}

// Segments returns a slice of Resource ID Segments which comprise this Site Extension ID
func (id SiteExtensionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticSiteExtensions", "siteExtensions", "siteExtensions"),
		resourceids.UserSpecifiedSegment("siteExtensionId", "siteExtensionId"),
	}
}

// String returns a human-readable description of this Site Extension ID
func (id SiteExtensionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Site Extension: %q", id.SiteExtensionId),
	}
	return fmt.Sprintf("Site Extension (%s)", strings.Join(components, "\n"))
}
