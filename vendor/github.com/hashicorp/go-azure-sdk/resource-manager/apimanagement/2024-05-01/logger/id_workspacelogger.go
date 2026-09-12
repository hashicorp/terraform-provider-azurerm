package logger

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&WorkspaceLoggerId{})
}

var _ resourceids.ResourceId = &WorkspaceLoggerId{}

// WorkspaceLoggerId is a struct representing the Resource ID for a Workspace Logger
type WorkspaceLoggerId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	WorkspaceId       string
	LoggerId          string
}

// NewWorkspaceLoggerID returns a new WorkspaceLoggerId struct
func NewWorkspaceLoggerID(subscriptionId string, resourceGroupName string, serviceName string, workspaceId string, loggerId string) WorkspaceLoggerId {
	return WorkspaceLoggerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		WorkspaceId:       workspaceId,
		LoggerId:          loggerId,
	}
}

// ParseWorkspaceLoggerID parses 'input' into a WorkspaceLoggerId
func ParseWorkspaceLoggerID(input string) (*WorkspaceLoggerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceLoggerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceLoggerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWorkspaceLoggerIDInsensitively parses 'input' case-insensitively into a WorkspaceLoggerId
// note: this method should only be used for API response data and not user input
func ParseWorkspaceLoggerIDInsensitively(input string) (*WorkspaceLoggerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceLoggerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceLoggerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WorkspaceLoggerId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.LoggerId, ok = input.Parsed["loggerId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "loggerId", input)
	}

	return nil
}

// ValidateWorkspaceLoggerID checks that 'input' can be parsed as a Workspace Logger ID
func ValidateWorkspaceLoggerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkspaceLoggerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Workspace Logger ID
func (id WorkspaceLoggerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/workspaces/%s/loggers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId, id.LoggerId)
}

// Segments returns a slice of Resource ID Segments which comprise this Workspace Logger ID
func (id WorkspaceLoggerId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticLoggers", "loggers", "loggers"),
		resourceids.UserSpecifiedSegment("loggerId", "loggerId"),
	}
}

// String returns a human-readable description of this Workspace Logger ID
func (id WorkspaceLoggerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace: %q", id.WorkspaceId),
		fmt.Sprintf("Logger: %q", id.LoggerId),
	}
	return fmt.Sprintf("Workspace Logger (%s)", strings.Join(components, "\n"))
}
