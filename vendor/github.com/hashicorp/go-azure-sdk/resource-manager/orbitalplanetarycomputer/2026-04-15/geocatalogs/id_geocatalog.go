package geocatalogs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&GeoCatalogId{})
}

var _ resourceids.ResourceId = &GeoCatalogId{}

// GeoCatalogId is a struct representing the Resource ID for a Geo Catalog
type GeoCatalogId struct {
	SubscriptionId    string
	ResourceGroupName string
	GeoCatalogName    string
}

// NewGeoCatalogID returns a new GeoCatalogId struct
func NewGeoCatalogID(subscriptionId string, resourceGroupName string, geoCatalogName string) GeoCatalogId {
	return GeoCatalogId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		GeoCatalogName:    geoCatalogName,
	}
}

// ParseGeoCatalogID parses 'input' into a GeoCatalogId
func ParseGeoCatalogID(input string) (*GeoCatalogId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GeoCatalogId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GeoCatalogId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseGeoCatalogIDInsensitively parses 'input' case-insensitively into a GeoCatalogId
// note: this method should only be used for API response data and not user input
func ParseGeoCatalogIDInsensitively(input string) (*GeoCatalogId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GeoCatalogId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GeoCatalogId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *GeoCatalogId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.GeoCatalogName, ok = input.Parsed["geoCatalogName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "geoCatalogName", input)
	}

	return nil
}

// ValidateGeoCatalogID checks that 'input' can be parsed as a Geo Catalog ID
func ValidateGeoCatalogID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGeoCatalogID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Geo Catalog ID
func (id GeoCatalogId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Orbital/geoCatalogs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.GeoCatalogName)
}

// Segments returns a slice of Resource ID Segments which comprise this Geo Catalog ID
func (id GeoCatalogId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOrbital", "Microsoft.Orbital", "Microsoft.Orbital"),
		resourceids.StaticSegment("staticGeoCatalogs", "geoCatalogs", "geoCatalogs"),
		resourceids.UserSpecifiedSegment("geoCatalogName", "geoCatalogName"),
	}
}

// String returns a human-readable description of this Geo Catalog ID
func (id GeoCatalogId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Geo Catalog Name: %q", id.GeoCatalogName),
	}
	return fmt.Sprintf("Geo Catalog (%s)", strings.Join(components, "\n"))
}
