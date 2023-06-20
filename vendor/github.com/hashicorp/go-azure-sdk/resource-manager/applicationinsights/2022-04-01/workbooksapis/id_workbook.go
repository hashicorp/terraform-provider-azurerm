package workbooksapis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = WorkbookId{}

// WorkbookId is a struct representing the Resource ID for a Workbook
type WorkbookId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkbookName      string
}

// NewWorkbookID returns a new WorkbookId struct
func NewWorkbookID(subscriptionId string, resourceGroupName string, workbookName string) WorkbookId {
	return WorkbookId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkbookName:      workbookName,
	}
}

// ParseWorkbookID parses 'input' into a WorkbookId
func ParseWorkbookID(input string) (*WorkbookId, error) {
	parser := resourceids.NewParserFromResourceIdType(WorkbookId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := WorkbookId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkbookName, ok = parsed.Parsed["workbookName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workbookName", *parsed)
	}

	return &id, nil
}

// ParseWorkbookIDInsensitively parses 'input' case-insensitively into a WorkbookId
// note: this method should only be used for API response data and not user input
func ParseWorkbookIDInsensitively(input string) (*WorkbookId, error) {
	parser := resourceids.NewParserFromResourceIdType(WorkbookId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := WorkbookId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkbookName, ok = parsed.Parsed["workbookName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workbookName", *parsed)
	}

	return &id, nil
}

// ValidateWorkbookID checks that 'input' can be parsed as a Workbook ID
func ValidateWorkbookID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkbookID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Workbook ID
func (id WorkbookId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/workbooks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkbookName)
}

// Segments returns a slice of Resource ID Segments which comprise this Workbook ID
func (id WorkbookId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticWorkbooks", "workbooks", "workbooks"),
		resourceids.UserSpecifiedSegment("workbookName", "workbookValue"),
	}
}

// String returns a human-readable description of this Workbook ID
func (id WorkbookId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workbook Name: %q", id.WorkbookName),
	}
	return fmt.Sprintf("Workbook (%s)", strings.Join(components, "\n"))
}
