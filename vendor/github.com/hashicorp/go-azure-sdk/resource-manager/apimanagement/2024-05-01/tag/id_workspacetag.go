package tag

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&WorkspaceTagId{})
}

var _ resourceids.ResourceId = &WorkspaceTagId{}

// WorkspaceTagId is a struct representing the Resource ID for a Workspace Tag
type WorkspaceTagId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	WorkspaceId       string
	TagId             string
}

// NewWorkspaceTagID returns a new WorkspaceTagId struct
func NewWorkspaceTagID(subscriptionId string, resourceGroupName string, serviceName string, workspaceId string, tagId string) WorkspaceTagId {
	return WorkspaceTagId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		WorkspaceId:       workspaceId,
		TagId:             tagId,
	}
}

// ParseWorkspaceTagID parses 'input' into a WorkspaceTagId
func ParseWorkspaceTagID(input string) (*WorkspaceTagId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceTagId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceTagId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWorkspaceTagIDInsensitively parses 'input' case-insensitively into a WorkspaceTagId
// note: this method should only be used for API response data and not user input
func ParseWorkspaceTagIDInsensitively(input string) (*WorkspaceTagId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceTagId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceTagId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WorkspaceTagId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.TagId, ok = input.Parsed["tagId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "tagId", input)
	}

	return nil
}

// ValidateWorkspaceTagID checks that 'input' can be parsed as a Workspace Tag ID
func ValidateWorkspaceTagID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkspaceTagID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Workspace Tag ID
func (id WorkspaceTagId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/workspaces/%s/tags/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId, id.TagId)
}

// Segments returns a slice of Resource ID Segments which comprise this Workspace Tag ID
func (id WorkspaceTagId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceName"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceId", "workspaceId"),
		resourceids.StaticSegment("staticTags", "tags", "tags"),
		resourceids.UserSpecifiedSegment("tagId", "tagId"),
	}
}

// String returns a human-readable description of this Workspace Tag ID
func (id WorkspaceTagId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace: %q", id.WorkspaceId),
		fmt.Sprintf("Tag: %q", id.TagId),
	}
	return fmt.Sprintf("Workspace Tag (%s)", strings.Join(components, "\n"))
}
