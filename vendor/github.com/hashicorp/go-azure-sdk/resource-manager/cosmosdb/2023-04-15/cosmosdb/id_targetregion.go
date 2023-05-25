package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = TargetRegionId{}

// TargetRegionId is a struct representing the Resource ID for a Target Region
type TargetRegionId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DatabaseAccountName string
	TargetRegionName    string
}

// NewTargetRegionID returns a new TargetRegionId struct
func NewTargetRegionID(subscriptionId string, resourceGroupName string, databaseAccountName string, targetRegionName string) TargetRegionId {
	return TargetRegionId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DatabaseAccountName: databaseAccountName,
		TargetRegionName:    targetRegionName,
	}
}

// ParseTargetRegionID parses 'input' into a TargetRegionId
func ParseTargetRegionID(input string) (*TargetRegionId, error) {
	parser := resourceids.NewParserFromResourceIdType(TargetRegionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TargetRegionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.TargetRegionName, ok = parsed.Parsed["targetRegionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "targetRegionName", *parsed)
	}

	return &id, nil
}

// ParseTargetRegionIDInsensitively parses 'input' case-insensitively into a TargetRegionId
// note: this method should only be used for API response data and not user input
func ParseTargetRegionIDInsensitively(input string) (*TargetRegionId, error) {
	parser := resourceids.NewParserFromResourceIdType(TargetRegionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TargetRegionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.TargetRegionName, ok = parsed.Parsed["targetRegionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "targetRegionName", *parsed)
	}

	return &id, nil
}

// ValidateTargetRegionID checks that 'input' can be parsed as a Target Region ID
func ValidateTargetRegionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTargetRegionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Target Region ID
func (id TargetRegionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/targetRegion/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.TargetRegionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Target Region ID
func (id TargetRegionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountValue"),
		resourceids.StaticSegment("staticTargetRegion", "targetRegion", "targetRegion"),
		resourceids.UserSpecifiedSegment("targetRegionName", "targetRegionValue"),
	}
}

// String returns a human-readable description of this Target Region ID
func (id TargetRegionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Target Region Name: %q", id.TargetRegionName),
	}
	return fmt.Sprintf("Target Region (%s)", strings.Join(components, "\n"))
}
