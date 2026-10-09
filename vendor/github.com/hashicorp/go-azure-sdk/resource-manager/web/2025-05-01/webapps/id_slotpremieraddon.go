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
	recaser.RegisterResourceId(&SlotPremierAddonId{})
}

var _ resourceids.ResourceId = &SlotPremierAddonId{}

// SlotPremierAddonId is a struct representing the Resource ID for a Slot Premier Addon
type SlotPremierAddonId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	SlotName          string
	PremierAddonName  string
}

// NewSlotPremierAddonID returns a new SlotPremierAddonId struct
func NewSlotPremierAddonID(subscriptionId string, resourceGroupName string, siteName string, slotName string, premierAddonName string) SlotPremierAddonId {
	return SlotPremierAddonId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		SlotName:          slotName,
		PremierAddonName:  premierAddonName,
	}
}

// ParseSlotPremierAddonID parses 'input' into a SlotPremierAddonId
func ParseSlotPremierAddonID(input string) (*SlotPremierAddonId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotPremierAddonId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotPremierAddonId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSlotPremierAddonIDInsensitively parses 'input' case-insensitively into a SlotPremierAddonId
// note: this method should only be used for API response data and not user input
func ParseSlotPremierAddonIDInsensitively(input string) (*SlotPremierAddonId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SlotPremierAddonId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SlotPremierAddonId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SlotPremierAddonId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.PremierAddonName, ok = input.Parsed["premierAddonName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "premierAddonName", input)
	}

	return nil
}

// ValidateSlotPremierAddonID checks that 'input' can be parsed as a Slot Premier Addon ID
func ValidateSlotPremierAddonID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSlotPremierAddonID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Slot Premier Addon ID
func (id SlotPremierAddonId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s/premierAddons/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.SlotName, id.PremierAddonName)
}

// Segments returns a slice of Resource ID Segments which comprise this Slot Premier Addon ID
func (id SlotPremierAddonId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticPremierAddons", "premierAddons", "premierAddons"),
		resourceids.UserSpecifiedSegment("premierAddonName", "premierAddonName"),
	}
}

// String returns a human-readable description of this Slot Premier Addon ID
func (id SlotPremierAddonId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Slot Name: %q", id.SlotName),
		fmt.Sprintf("Premier Addon Name: %q", id.PremierAddonName),
	}
	return fmt.Sprintf("Slot Premier Addon (%s)", strings.Join(components, "\n"))
}
