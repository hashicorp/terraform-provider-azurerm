package apiversionset

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&WorkspaceApiVersionSetId{})
}

var _ resourceids.ResourceId = &WorkspaceApiVersionSetId{}

// WorkspaceApiVersionSetId is a struct representing the Resource ID for a Workspace Api Version Set
type WorkspaceApiVersionSetId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	WorkspaceId       string
	VersionSetId      string
}

// NewWorkspaceApiVersionSetID returns a new WorkspaceApiVersionSetId struct
func NewWorkspaceApiVersionSetID(subscriptionId string, resourceGroupName string, serviceName string, workspaceId string, versionSetId string) WorkspaceApiVersionSetId {
	return WorkspaceApiVersionSetId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		WorkspaceId:       workspaceId,
		VersionSetId:      versionSetId,
	}
}

// ParseWorkspaceApiVersionSetID parses 'input' into a WorkspaceApiVersionSetId
func ParseWorkspaceApiVersionSetID(input string) (*WorkspaceApiVersionSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceApiVersionSetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceApiVersionSetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWorkspaceApiVersionSetIDInsensitively parses 'input' case-insensitively into a WorkspaceApiVersionSetId
// note: this method should only be used for API response data and not user input
func ParseWorkspaceApiVersionSetIDInsensitively(input string) (*WorkspaceApiVersionSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceApiVersionSetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceApiVersionSetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WorkspaceApiVersionSetId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.VersionSetId, ok = input.Parsed["versionSetId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "versionSetId", input)
	}

	return nil
}

// ValidateWorkspaceApiVersionSetID checks that 'input' can be parsed as a Workspace Api Version Set ID
func ValidateWorkspaceApiVersionSetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkspaceApiVersionSetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Workspace Api Version Set ID
func (id WorkspaceApiVersionSetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/workspaces/%s/apiVersionSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId, id.VersionSetId)
}

// Segments returns a slice of Resource ID Segments which comprise this Workspace Api Version Set ID
func (id WorkspaceApiVersionSetId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticApiVersionSets", "apiVersionSets", "apiVersionSets"),
		resourceids.UserSpecifiedSegment("versionSetId", "versionSetId"),
	}
}

// String returns a human-readable description of this Workspace Api Version Set ID
func (id WorkspaceApiVersionSetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace: %q", id.WorkspaceId),
		fmt.Sprintf("Version Set: %q", id.VersionSetId),
	}
	return fmt.Sprintf("Workspace Api Version Set (%s)", strings.Join(components, "\n"))
}
