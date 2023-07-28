package datacontainer

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = WorkspaceDataId{}

// WorkspaceDataId is a struct representing the Resource ID for a Workspace Data
type WorkspaceDataId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	DataName          string
}

// NewWorkspaceDataID returns a new WorkspaceDataId struct
func NewWorkspaceDataID(subscriptionId string, resourceGroupName string, workspaceName string, dataName string) WorkspaceDataId {
	return WorkspaceDataId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		DataName:          dataName,
	}
}

// ParseWorkspaceDataID parses 'input' into a WorkspaceDataId
func ParseWorkspaceDataID(input string) (*WorkspaceDataId, error) {
	parser := resourceids.NewParserFromResourceIdType(WorkspaceDataId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := WorkspaceDataId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.DataName, ok = parsed.Parsed["dataName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dataName", *parsed)
	}

	return &id, nil
}

// ParseWorkspaceDataIDInsensitively parses 'input' case-insensitively into a WorkspaceDataId
// note: this method should only be used for API response data and not user input
func ParseWorkspaceDataIDInsensitively(input string) (*WorkspaceDataId, error) {
	parser := resourceids.NewParserFromResourceIdType(WorkspaceDataId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := WorkspaceDataId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.DataName, ok = parsed.Parsed["dataName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dataName", *parsed)
	}

	return &id, nil
}

// ValidateWorkspaceDataID checks that 'input' can be parsed as a Workspace Data ID
func ValidateWorkspaceDataID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkspaceDataID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Workspace Data ID
func (id WorkspaceDataId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MachineLearningServices/workspaces/%s/data/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.DataName)
}

// Segments returns a slice of Resource ID Segments which comprise this Workspace Data ID
func (id WorkspaceDataId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMachineLearningServices", "Microsoft.MachineLearningServices", "Microsoft.MachineLearningServices"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticData", "data", "data"),
		resourceids.UserSpecifiedSegment("dataName", "dataValue"),
	}
}

// String returns a human-readable description of this Workspace Data ID
func (id WorkspaceDataId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Data Name: %q", id.DataName),
	}
	return fmt.Sprintf("Workspace Data (%s)", strings.Join(components, "\n"))
}
