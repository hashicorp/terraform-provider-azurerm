package apidiagnostic

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&WorkspaceApiDiagnosticId{})
}

var _ resourceids.ResourceId = &WorkspaceApiDiagnosticId{}

// WorkspaceApiDiagnosticId is a struct representing the Resource ID for a Workspace Api Diagnostic
type WorkspaceApiDiagnosticId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	WorkspaceId       string
	ApiId             string
	DiagnosticId      string
}

// NewWorkspaceApiDiagnosticID returns a new WorkspaceApiDiagnosticId struct
func NewWorkspaceApiDiagnosticID(subscriptionId string, resourceGroupName string, serviceName string, workspaceId string, apiId string, diagnosticId string) WorkspaceApiDiagnosticId {
	return WorkspaceApiDiagnosticId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		WorkspaceId:       workspaceId,
		ApiId:             apiId,
		DiagnosticId:      diagnosticId,
	}
}

// ParseWorkspaceApiDiagnosticID parses 'input' into a WorkspaceApiDiagnosticId
func ParseWorkspaceApiDiagnosticID(input string) (*WorkspaceApiDiagnosticId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceApiDiagnosticId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceApiDiagnosticId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWorkspaceApiDiagnosticIDInsensitively parses 'input' case-insensitively into a WorkspaceApiDiagnosticId
// note: this method should only be used for API response data and not user input
func ParseWorkspaceApiDiagnosticIDInsensitively(input string) (*WorkspaceApiDiagnosticId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceApiDiagnosticId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceApiDiagnosticId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WorkspaceApiDiagnosticId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ApiId, ok = input.Parsed["apiId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "apiId", input)
	}

	if id.DiagnosticId, ok = input.Parsed["diagnosticId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "diagnosticId", input)
	}

	return nil
}

// ValidateWorkspaceApiDiagnosticID checks that 'input' can be parsed as a Workspace Api Diagnostic ID
func ValidateWorkspaceApiDiagnosticID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkspaceApiDiagnosticID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Workspace Api Diagnostic ID
func (id WorkspaceApiDiagnosticId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/workspaces/%s/apis/%s/diagnostics/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId, id.ApiId, id.DiagnosticId)
}

// Segments returns a slice of Resource ID Segments which comprise this Workspace Api Diagnostic ID
func (id WorkspaceApiDiagnosticId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticApis", "apis", "apis"),
		resourceids.UserSpecifiedSegment("apiId", "apiId"),
		resourceids.StaticSegment("staticDiagnostics", "diagnostics", "diagnostics"),
		resourceids.UserSpecifiedSegment("diagnosticId", "diagnosticId"),
	}
}

// String returns a human-readable description of this Workspace Api Diagnostic ID
func (id WorkspaceApiDiagnosticId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace: %q", id.WorkspaceId),
		fmt.Sprintf("Api: %q", id.ApiId),
		fmt.Sprintf("Diagnostic: %q", id.DiagnosticId),
	}
	return fmt.Sprintf("Workspace Api Diagnostic (%s)", strings.Join(components, "\n"))
}
