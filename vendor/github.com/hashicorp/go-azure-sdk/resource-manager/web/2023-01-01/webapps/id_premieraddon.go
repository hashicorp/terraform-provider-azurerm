package webapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &PremierAddonId{}

// PremierAddonId is a struct representing the Resource ID for a Premier Addon
type PremierAddonId struct {
	SubscriptionId    string
	ResourceGroupName string
	SiteName          string
	PremierAddonName  string
}

// NewPremierAddonID returns a new PremierAddonId struct
func NewPremierAddonID(subscriptionId string, resourceGroupName string, siteName string, premierAddonName string) PremierAddonId {
	return PremierAddonId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SiteName:          siteName,
		PremierAddonName:  premierAddonName,
	}
}

// ParsePremierAddonID parses 'input' into a PremierAddonId
func ParsePremierAddonID(input string) (*PremierAddonId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PremierAddonId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PremierAddonId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePremierAddonIDInsensitively parses 'input' case-insensitively into a PremierAddonId
// note: this method should only be used for API response data and not user input
func ParsePremierAddonIDInsensitively(input string) (*PremierAddonId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PremierAddonId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PremierAddonId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PremierAddonId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.PremierAddonName, ok = input.Parsed["premierAddonName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "premierAddonName", input)
	}

	return nil
}

// ValidatePremierAddonID checks that 'input' can be parsed as a Premier Addon ID
func ValidatePremierAddonID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePremierAddonID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Premier Addon ID
func (id PremierAddonId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/premierAddons/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.PremierAddonName)
}

// Segments returns a slice of Resource ID Segments which comprise this Premier Addon ID
func (id PremierAddonId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteValue"),
		resourceids.StaticSegment("staticPremierAddons", "premierAddons", "premierAddons"),
		resourceids.UserSpecifiedSegment("premierAddonName", "premierAddonValue"),
	}
}

// String returns a human-readable description of this Premier Addon ID
func (id PremierAddonId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Premier Addon Name: %q", id.PremierAddonName),
	}
	return fmt.Sprintf("Premier Addon (%s)", strings.Join(components, "\n"))
}
