package dataexport

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DataExportId{}

// DataExportId is a struct representing the Resource ID for a Data Export
type DataExportId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	DataExportName    string
}

// NewDataExportID returns a new DataExportId struct
func NewDataExportID(subscriptionId string, resourceGroupName string, workspaceName string, dataExportName string) DataExportId {
	return DataExportId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		DataExportName:    dataExportName,
	}
}

// ParseDataExportID parses 'input' into a DataExportId
func ParseDataExportID(input string) (*DataExportId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataExportId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataExportId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.DataExportName, ok = parsed.Parsed["dataExportName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dataExportName", *parsed)
	}

	return &id, nil
}

// ParseDataExportIDInsensitively parses 'input' case-insensitively into a DataExportId
// note: this method should only be used for API response data and not user input
func ParseDataExportIDInsensitively(input string) (*DataExportId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataExportId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataExportId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.DataExportName, ok = parsed.Parsed["dataExportName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dataExportName", *parsed)
	}

	return &id, nil
}

// ValidateDataExportID checks that 'input' can be parsed as a Data Export ID
func ValidateDataExportID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataExportID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Data Export ID
func (id DataExportId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/dataExports/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.DataExportName)
}

// Segments returns a slice of Resource ID Segments which comprise this Data Export ID
func (id DataExportId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticDataExports", "dataExports", "dataExports"),
		resourceids.UserSpecifiedSegment("dataExportName", "dataExportValue"),
	}
}

// String returns a human-readable description of this Data Export ID
func (id DataExportId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Data Export Name: %q", id.DataExportName),
	}
	return fmt.Sprintf("Data Export (%s)", strings.Join(components, "\n"))
}
