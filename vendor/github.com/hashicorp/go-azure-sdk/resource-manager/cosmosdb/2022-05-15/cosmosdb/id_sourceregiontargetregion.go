package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SourceRegionTargetRegionId{}

// SourceRegionTargetRegionId is a struct representing the Resource ID for a Source Region Target Region
type SourceRegionTargetRegionId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DatabaseAccountName string
	SourceRegionName    string
	TargetRegionName    string
}

// NewSourceRegionTargetRegionID returns a new SourceRegionTargetRegionId struct
func NewSourceRegionTargetRegionID(subscriptionId string, resourceGroupName string, databaseAccountName string, sourceRegionName string, targetRegionName string) SourceRegionTargetRegionId {
	return SourceRegionTargetRegionId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DatabaseAccountName: databaseAccountName,
		SourceRegionName:    sourceRegionName,
		TargetRegionName:    targetRegionName,
	}
}

// ParseSourceRegionTargetRegionID parses 'input' into a SourceRegionTargetRegionId
func ParseSourceRegionTargetRegionID(input string) (*SourceRegionTargetRegionId, error) {
	parser := resourceids.NewParserFromResourceIdType(SourceRegionTargetRegionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SourceRegionTargetRegionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.SourceRegionName, ok = parsed.Parsed["sourceRegionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sourceRegionName", *parsed)
	}

	if id.TargetRegionName, ok = parsed.Parsed["targetRegionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "targetRegionName", *parsed)
	}

	return &id, nil
}

// ParseSourceRegionTargetRegionIDInsensitively parses 'input' case-insensitively into a SourceRegionTargetRegionId
// note: this method should only be used for API response data and not user input
func ParseSourceRegionTargetRegionIDInsensitively(input string) (*SourceRegionTargetRegionId, error) {
	parser := resourceids.NewParserFromResourceIdType(SourceRegionTargetRegionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SourceRegionTargetRegionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.SourceRegionName, ok = parsed.Parsed["sourceRegionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sourceRegionName", *parsed)
	}

	if id.TargetRegionName, ok = parsed.Parsed["targetRegionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "targetRegionName", *parsed)
	}

	return &id, nil
}

// ValidateSourceRegionTargetRegionID checks that 'input' can be parsed as a Source Region Target Region ID
func ValidateSourceRegionTargetRegionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSourceRegionTargetRegionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Source Region Target Region ID
func (id SourceRegionTargetRegionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/sourceRegion/%s/targetRegion/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.SourceRegionName, id.TargetRegionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Source Region Target Region ID
func (id SourceRegionTargetRegionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountValue"),
		resourceids.StaticSegment("staticSourceRegion", "sourceRegion", "sourceRegion"),
		resourceids.UserSpecifiedSegment("sourceRegionName", "sourceRegionValue"),
		resourceids.StaticSegment("staticTargetRegion", "targetRegion", "targetRegion"),
		resourceids.UserSpecifiedSegment("targetRegionName", "targetRegionValue"),
	}
}

// String returns a human-readable description of this Source Region Target Region ID
func (id SourceRegionTargetRegionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Source Region Name: %q", id.SourceRegionName),
		fmt.Sprintf("Target Region Name: %q", id.TargetRegionName),
	}
	return fmt.Sprintf("Source Region Target Region (%s)", strings.Join(components, "\n"))
}
