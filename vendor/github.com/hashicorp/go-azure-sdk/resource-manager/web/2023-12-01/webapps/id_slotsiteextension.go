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
	recaser.RegisterResourceId(&SlotSiteExtensionId{})
}

var _ resourceids.ResourceId = &SlotSiteExtensionId{}

// SlotSiteExtensionId is a struct representing the Resource ID for a Slot Site Extension
type SlotSiteExtensionId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	SlotName          string
	SiteExtensionId   string
}

// NewSlotSiteExtensionID returns a new SlotSiteExtensionId struct
func NewSlotSiteExtensionID(subscriptionId string, resourceGroupName string, siteName string, slotName string, siteExtensionId string) SlotSiteExtensionId {
	return SlotSiteExtensionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		SlotName:          slotName,
		SiteExtensionId:   siteExtensionId,
	}
}

// ParseSlotSiteExtensionID parses 'input' into a SlotSiteExtensionId
func ParseSlotSiteExtensionID(input string) (*SlotSiteExtensionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotSiteExtensionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotSiteExtensionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSlotSiteExtensionIDInsensitively parses 'input' case-insensitively into a SlotSiteExtensionId
// note: this method should only be used for API response data and not user input
func ParseSlotSiteExtensionIDInsensitively(input string) (*SlotSiteExtensionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotSiteExtensionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotSiteExtensionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SlotSiteExtensionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.SlotName, ok = input.Parsed["slotName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "slotName", input)
	}

	if id.SiteExtensionId, ok = input.Parsed["siteExtensionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "siteExtensionId", input)
	}

	return nil
}

// ValidateSlotSiteExtensionID checks that 'input' can be parsed as a Slot Site Extension ID
func ValidateSlotSiteExtensionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSlotSiteExtensionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Slot Site Extension ID
func (id SlotSiteExtensionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/siteExtensions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.SiteExtensionId)
}

// Segments returns a slice of Resource ID Segments which comprise this Slot Site Extension ID
func (id SlotSiteExtensionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticSlots", "slots", "slots"),
		resourceids.UserSpecifiedSegment("slotName", "slotName"),
		resourceids.StaticSegment("staticSiteExtensions", "siteExtensions", "siteExtensions"),
		resourceids.UserSpecifiedSegment("siteExtensionId", "siteExtensionId"),
	}
}

// String returns a human-readable description of this Slot Site Extension ID
func (id SlotSiteExtensionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Site Extension: %q", id.SiteExtensionId),
	}
	return fmt.Sprintf("Slot Site Extension (%s)", strings.Join(components, "\n"))
}
