package diagnostic

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&WorkspaceDiagnosticId{})
}

var _ resourceids.ResourceId = &WorkspaceDiagnosticId{}

// WorkspaceDiagnosticId is a struct representing the Resource ID for a Workspace Diagnostic
type WorkspaceDiagnosticId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	WorkspaceId       string
	DiagnosticId      string
}

// NewWorkspaceDiagnosticID returns a new WorkspaceDiagnosticId struct
func NewWorkspaceDiagnosticID(subscriptionId string, resourceGroupName string, serviceName string, workspaceId string, diagnosticId string) WorkspaceDiagnosticId {
	return WorkspaceDiagnosticId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		WorkspaceId:       workspaceId,
		DiagnosticId:      diagnosticId,
	}
}

// ParseWorkspaceDiagnosticID parses 'input' into a WorkspaceDiagnosticId
func ParseWorkspaceDiagnosticID(input string) (*WorkspaceDiagnosticId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceDiagnosticId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceDiagnosticId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWorkspaceDiagnosticIDInsensitively parses 'input' case-insensitively into a WorkspaceDiagnosticId
// note: this method should only be used for API response data and not user input
func ParseWorkspaceDiagnosticIDInsensitively(input string) (*WorkspaceDiagnosticId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceDiagnosticId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceDiagnosticId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WorkspaceDiagnosticId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DiagnosticId, ok = input.Parsed["diagnosticId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "diagnosticId", input)
	}

	return nil
}

// ValidateWorkspaceDiagnosticID checks that 'input' can be parsed as a Workspace Diagnostic ID
func ValidateWorkspaceDiagnosticID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkspaceDiagnosticID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Workspace Diagnostic ID
func (id WorkspaceDiagnosticId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/workspaces/%s/diagnostics/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId, id.DiagnosticId)
}

// Segments returns a slice of Resource ID Segments which comprise this Workspace Diagnostic ID
func (id WorkspaceDiagnosticId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticDiagnostics", "diagnostics", "diagnostics"),
		resourceids.UserSpecifiedSegment("diagnosticId", "diagnosticId"),
	}
}

// String returns a human-readable description of this Workspace Diagnostic ID
func (id WorkspaceDiagnosticId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace: %q", id.WorkspaceId),
		fmt.Sprintf("Diagnostic: %q", id.DiagnosticId),
	}
	return fmt.Sprintf("Workspace Diagnostic (%s)", strings.Join(components, "\n"))
}
