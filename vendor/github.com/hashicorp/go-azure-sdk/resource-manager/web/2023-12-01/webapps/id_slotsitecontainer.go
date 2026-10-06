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
	recaser.RegisterResourceId(&SlotSitecontainerId{})
}

var _ resourceids.ResourceId = &SlotSitecontainerId{}

// SlotSitecontainerId is a struct representing the Resource ID for a Slot Sitecontainer
type SlotSitecontainerId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	SlotName          string
	SitecontainerName string
}

// NewSlotSitecontainerID returns a new SlotSitecontainerId struct
func NewSlotSitecontainerID(subscriptionId string, resourceGroupName string, siteName string, slotName string, sitecontainerName string) SlotSitecontainerId {
	return SlotSitecontainerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		SlotName:          slotName,
		SitecontainerName: sitecontainerName,
	}
}

// ParseSlotSitecontainerID parses 'input' into a SlotSitecontainerId
func ParseSlotSitecontainerID(input string) (*SlotSitecontainerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotSitecontainerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotSitecontainerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSlotSitecontainerIDInsensitively parses 'input' case-insensitively into a SlotSitecontainerId
// note: this method should only be used for API response data and not user input
func ParseSlotSitecontainerIDInsensitively(input string) (*SlotSitecontainerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotSitecontainerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotSitecontainerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SlotSitecontainerId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.SitecontainerName, ok = input.Parsed["sitecontainerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "sitecontainerName", input)
	}

	return nil
}

// ValidateSlotSitecontainerID checks that 'input' can be parsed as a Slot Sitecontainer ID
func ValidateSlotSitecontainerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSlotSitecontainerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Slot Sitecontainer ID
func (id SlotSitecontainerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/sitecontainers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.SitecontainerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Slot Sitecontainer ID
func (id SlotSitecontainerId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticSitecontainers", "sitecontainers", "sitecontainers"),
		resourceids.UserSpecifiedSegment("sitecontainerName", "sitecontainerName"),
	}
}

// String returns a human-readable description of this Slot Sitecontainer ID
func (id SlotSitecontainerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Sitecontainer Name: %q", id.SitecontainerName),
	}
	return fmt.Sprintf("Slot Sitecontainer (%s)", strings.Join(components, "\n"))
}
