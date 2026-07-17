package backend

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&WorkspaceBackendId{})
}

var _ resourceids.ResourceId = &WorkspaceBackendId{}

// WorkspaceBackendId is a struct representing the Resource ID for a Workspace Backend
type WorkspaceBackendId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	WorkspaceId       string
	BackendId         string
}

// NewWorkspaceBackendID returns a new WorkspaceBackendId struct
func NewWorkspaceBackendID(subscriptionId string, resourceGroupName string, serviceName string, workspaceId string, backendId string) WorkspaceBackendId {
	return WorkspaceBackendId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		WorkspaceId:       workspaceId,
		BackendId:         backendId,
	}
}

// ParseWorkspaceBackendID parses 'input' into a WorkspaceBackendId
func ParseWorkspaceBackendID(input string) (*WorkspaceBackendId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceBackendId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceBackendId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWorkspaceBackendIDInsensitively parses 'input' case-insensitively into a WorkspaceBackendId
// note: this method should only be used for API response data and not user input
func ParseWorkspaceBackendIDInsensitively(input string) (*WorkspaceBackendId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceBackendId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceBackendId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WorkspaceBackendId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.BackendId, ok = input.Parsed["backendId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "backendId", input)
	}

	return nil
}

// ValidateWorkspaceBackendID checks that 'input' can be parsed as a Workspace Backend ID
func ValidateWorkspaceBackendID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkspaceBackendID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Workspace Backend ID
func (id WorkspaceBackendId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/workspaces/%s/backends/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId, id.BackendId)
}

// Segments returns a slice of Resource ID Segments which comprise this Workspace Backend ID
func (id WorkspaceBackendId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticBackends", "backends", "backends"),
		resourceids.UserSpecifiedSegment("backendId", "backendId"),
	}
}

// String returns a human-readable description of this Workspace Backend ID
func (id WorkspaceBackendId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace: %q", id.WorkspaceId),
		fmt.Sprintf("Backend: %q", id.BackendId),
	}
	return fmt.Sprintf("Workspace Backend (%s)", strings.Join(components, "\n"))
}
