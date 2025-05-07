package customizationtasks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DevCenterCatalogId{})
}

var _ resourceids.ResourceId = &DevCenterCatalogId{}

// DevCenterCatalogId is a struct representing the Resource ID for a Dev Center Catalog
type DevCenterCatalogId struct {
	SubscriptionId    string
	ResourceGroupName string
	DevCenterName     string
	CatalogName       string
}

// NewDevCenterCatalogID returns a new DevCenterCatalogId struct
func NewDevCenterCatalogID(subscriptionId string, resourceGroupName string, devCenterName string, catalogName string) DevCenterCatalogId {
	return DevCenterCatalogId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DevCenterName:     devCenterName,
		CatalogName:       catalogName,
	}
}

// ParseDevCenterCatalogID parses 'input' into a DevCenterCatalogId
func ParseDevCenterCatalogID(input string) (*DevCenterCatalogId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DevCenterCatalogId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DevCenterCatalogId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDevCenterCatalogIDInsensitively parses 'input' case-insensitively into a DevCenterCatalogId
// note: this method should only be used for API response data and not user input
func ParseDevCenterCatalogIDInsensitively(input string) (*DevCenterCatalogId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DevCenterCatalogId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DevCenterCatalogId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DevCenterCatalogId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DevCenterName, ok = input.Parsed["devCenterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "devCenterName", input)
	}

	if id.CatalogName, ok = input.Parsed["catalogName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "catalogName", input)
	}

	return nil
}

// ValidateDevCenterCatalogID checks that 'input' can be parsed as a Dev Center Catalog ID
func ValidateDevCenterCatalogID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDevCenterCatalogID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dev Center Catalog ID
func (id DevCenterCatalogId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevCenter/devCenters/%s/catalogs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DevCenterName, id.CatalogName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dev Center Catalog ID
func (id DevCenterCatalogId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevCenter", "Microsoft.DevCenter", "Microsoft.DevCenter"),
		resourceids.StaticSegment("staticDevCenters", "devCenters", "devCenters"),
		resourceids.UserSpecifiedSegment("devCenterName", "devCenterName"),
		resourceids.StaticSegment("staticCatalogs", "catalogs", "catalogs"),
		resourceids.UserSpecifiedSegment("catalogName", "catalogName"),
	}
}

// String returns a human-readable description of this Dev Center Catalog ID
func (id DevCenterCatalogId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dev Center Name: %q", id.DevCenterName),
		fmt.Sprintf("Catalog Name: %q", id.CatalogName),
	}
	return fmt.Sprintf("Dev Center Catalog (%s)", strings.Join(components, "\n"))
}
