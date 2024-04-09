package schema

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &WorkspaceSchemaId{}

// WorkspaceSchemaId is a struct representing the Resource ID for a Workspace Schema
type WorkspaceSchemaId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	WorkspaceId       string
	SchemaId          string
}

// NewWorkspaceSchemaID returns a new WorkspaceSchemaId struct
func NewWorkspaceSchemaID(subscriptionId string, resourceGroupName string, serviceName string, workspaceId string, schemaId string) WorkspaceSchemaId {
	return WorkspaceSchemaId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		WorkspaceId:       workspaceId,
		SchemaId:          schemaId,
	}
}

// ParseWorkspaceSchemaID parses 'input' into a WorkspaceSchemaId
func ParseWorkspaceSchemaID(input string) (*WorkspaceSchemaId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceSchemaId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceSchemaId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWorkspaceSchemaIDInsensitively parses 'input' case-insensitively into a WorkspaceSchemaId
// note: this method should only be used for API response data and not user input
func ParseWorkspaceSchemaIDInsensitively(input string) (*WorkspaceSchemaId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspaceSchemaId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspaceSchemaId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WorkspaceSchemaId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.SchemaId, ok = input.Parsed["schemaId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "schemaId", input)
	}

	return nil
}

// ValidateWorkspaceSchemaID checks that 'input' can be parsed as a Workspace Schema ID
func ValidateWorkspaceSchemaID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkspaceSchemaID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Workspace Schema ID
func (id WorkspaceSchemaId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/workspaces/%s/schemas/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId, id.SchemaId)
}

// Segments returns a slice of Resource ID Segments which comprise this Workspace Schema ID
func (id WorkspaceSchemaId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticSchemas", "schemas", "schemas"),
		resourceids.UserSpecifiedSegment("schemaId", "schemaIdValue"),
	}
}

// String returns a human-readable description of this Workspace Schema ID
func (id WorkspaceSchemaId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace: %q", id.WorkspaceId),
		fmt.Sprintf("Schema: %q", id.SchemaId),
	}
	return fmt.Sprintf("Workspace Schema (%s)", strings.Join(components, "\n"))
}
