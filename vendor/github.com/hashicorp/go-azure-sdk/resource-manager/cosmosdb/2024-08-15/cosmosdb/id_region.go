package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RegionId{})
}

var _ resourceids.ResourceId = &RegionId{}

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
	parser := resourceids.NewParserFromResourceIdType(&RegionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RegionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRegionIDInsensitively parses 'input' case-insensitively into a RegionId
// note: this method should only be used for API response data and not user input
func ParseRegionIDInsensitively(input string) (*RegionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RegionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RegionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RegionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DatabaseAccountName, ok = input.Parsed["databaseAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", input)
	}

	if id.RegionName, ok = input.Parsed["regionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "regionName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountName"),
		resourceids.StaticSegment("staticRegion", "region", "region"),
		resourceids.UserSpecifiedSegment("regionName", "regionName"),
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
