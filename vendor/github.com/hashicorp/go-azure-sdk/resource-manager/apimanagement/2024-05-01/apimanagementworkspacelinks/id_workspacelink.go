package apimanagementworkspacelinks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&WorkspaceLinkId{})
}

var _ resourceids.ResourceId = &WorkspaceLinkId{}

// WorkspaceLinkId is a struct representing the Resource ID for a Workspace Link
type WorkspaceLinkId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	WorkspaceId       string
}

// NewWorkspaceLinkID returns a new WorkspaceLinkId struct
func NewWorkspaceLinkID(subscriptionId string, resourceGroupName string, serviceName string, workspaceId string) WorkspaceLinkId {
	return WorkspaceLinkId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		WorkspaceId:       workspaceId,
	}
}

// ParseWorkspaceLinkID parses 'input' into a WorkspaceLinkId
func ParseWorkspaceLinkID(input string) (*WorkspaceLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceLinkId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWorkspaceLinkIDInsensitively parses 'input' case-insensitively into a WorkspaceLinkId
// note: this method should only be used for API response data and not user input
func ParseWorkspaceLinkIDInsensitively(input string) (*WorkspaceLinkId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceLinkId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceLinkId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WorkspaceLinkId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServiceName, ok = input.Parsed["serviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceName", input)
	}

	if id.WorkspaceId, ok = input.Parsed["workspaceId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workspaceId", input)
	}

	return nil
}

// ValidateWorkspaceLinkID checks that 'input' can be parsed as a Workspace Link ID
func ValidateWorkspaceLinkID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkspaceLinkID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Workspace Link ID
func (id WorkspaceLinkId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/workspaceLinks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId)
}

// Segments returns a slice of Resource ID Segments which comprise this Workspace Link ID
func (id WorkspaceLinkId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticWorkspaceLinks", "workspaceLinks", "workspaceLinks"),
		resourceids.UserSpecifiedSegment("workspaceId", "workspaceId"),
	}
}

// String returns a human-readable description of this Workspace Link ID
func (id WorkspaceLinkId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace: %q", id.WorkspaceId),
	}
	return fmt.Sprintf("Workspace Link (%s)", strings.Join(components, "\n"))
}
