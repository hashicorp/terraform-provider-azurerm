package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RegionId{}

// RegionId is a struct representing the Resource ID for a Region
type RegionId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DatabaseAccountName string
	RegionName          string
}

// NewRegionID returns a new RegionId struct
func NewRegionID(subscriptionId string, resourceGroupName string, databaseAccountName string, regionName string) RegionId {
	return RegionId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DatabaseAccountName: databaseAccountName,
		RegionName:          regionName,
	}
}

// ParseRegionID parses 'input' into a RegionId
func ParseRegionID(input string) (*RegionId, error) {
	parser := resourceids.NewParserFromResourceIdType(RegionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RegionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.RegionName, ok = parsed.Parsed["regionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "regionName", *parsed)
	}

	return &id, nil
}

// ParseRegionIDInsensitively parses 'input' case-insensitively into a RegionId
// note: this method should only be used for API response data and not user input
func ParseRegionIDInsensitively(input string) (*RegionId, error) {
	parser := resourceids.NewParserFromResourceIdType(RegionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RegionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.RegionName, ok = parsed.Parsed["regionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "regionName", *parsed)
	}

	return &id, nil
}

// ValidateRegionID checks that 'input' can be parsed as a Region ID
func ValidateRegionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRegionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Region ID
func (id RegionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/region/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.RegionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Region ID
func (id RegionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountValue"),
		resourceids.StaticSegment("staticRegion", "region", "region"),
		resourceids.UserSpecifiedSegment("regionName", "regionValue"),
	}
}

// String returns a human-readable description of this Region ID
func (id RegionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Region Name: %q", id.RegionName),
	}
	return fmt.Sprintf("Region (%s)", strings.Join(components, "\n"))
}
