package namedvalue

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&WorkspaceNamedValueId{})
}

var _ resourceids.ResourceId = &WorkspaceNamedValueId{}

// WorkspaceNamedValueId is a struct representing the Resource ID for a Workspace Named Value
type WorkspaceNamedValueId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	WorkspaceId       string
	NamedValueId      string
}

// NewWorkspaceNamedValueID returns a new WorkspaceNamedValueId struct
func NewWorkspaceNamedValueID(subscriptionId string, resourceGroupName string, serviceName string, workspaceId string, namedValueId string) WorkspaceNamedValueId {
	return WorkspaceNamedValueId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		WorkspaceId:       workspaceId,
		NamedValueId:      namedValueId,
	}
}

// ParseWorkspaceNamedValueID parses 'input' into a WorkspaceNamedValueId
func ParseWorkspaceNamedValueID(input string) (*WorkspaceNamedValueId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceNamedValueId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceNamedValueId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWorkspaceNamedValueIDInsensitively parses 'input' case-insensitively into a WorkspaceNamedValueId
// note: this method should only be used for API response data and not user input
func ParseWorkspaceNamedValueIDInsensitively(input string) (*WorkspaceNamedValueId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceNamedValueId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceNamedValueId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WorkspaceNamedValueId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.NamedValueId, ok = input.Parsed["namedValueId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "namedValueId", input)
	}

	return nil
}

// ValidateWorkspaceNamedValueID checks that 'input' can be parsed as a Workspace Named Value ID
func ValidateWorkspaceNamedValueID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkspaceNamedValueID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Workspace Named Value ID
func (id WorkspaceNamedValueId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/workspaces/%s/namedValues/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId, id.NamedValueId)
}

// Segments returns a slice of Resource ID Segments which comprise this Workspace Named Value ID
func (id WorkspaceNamedValueId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceId", "workspaceIdValue"),
		resourceids.StaticSegment("staticNamedValues", "namedValues", "namedValues"),
		resourceids.UserSpecifiedSegment("namedValueId", "namedValueIdValue"),
	}
}

// String returns a human-readable description of this Workspace Named Value ID
func (id WorkspaceNamedValueId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace: %q", id.WorkspaceId),
		fmt.Sprintf("Named Value: %q", id.NamedValueId),
	}
	return fmt.Sprintf("Workspace Named Value (%s)", strings.Join(components, "\n"))
}
