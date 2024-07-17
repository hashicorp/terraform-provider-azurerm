package policyfragment

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&WorkspacePolicyFragmentId{})
}

var _ resourceids.ResourceId = &WorkspacePolicyFragmentId{}

// WorkspacePolicyFragmentId is a struct representing the Resource ID for a Workspace Policy Fragment
type WorkspacePolicyFragmentId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ServiceName        string
	WorkspaceId        string
	PolicyFragmentName string
}

// NewWorkspacePolicyFragmentID returns a new WorkspacePolicyFragmentId struct
func NewWorkspacePolicyFragmentID(subscriptionId string, resourceGroupName string, serviceName string, workspaceId string, policyFragmentName string) WorkspacePolicyFragmentId {
	return WorkspacePolicyFragmentId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ServiceName:        serviceName,
		WorkspaceId:        workspaceId,
		PolicyFragmentName: policyFragmentName,
	}
}

// ParseWorkspacePolicyFragmentID parses 'input' into a WorkspacePolicyFragmentId
func ParseWorkspacePolicyFragmentID(input string) (*WorkspacePolicyFragmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspacePolicyFragmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspacePolicyFragmentId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseWorkspacePolicyFragmentIDInsensitively parses 'input' case-insensitively into a WorkspacePolicyFragmentId
// note: this method should only be used for API response data and not user input
func ParseWorkspacePolicyFragmentIDInsensitively(input string) (*WorkspacePolicyFragmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&WorkspacePolicyFragmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := WorkspacePolicyFragmentId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *WorkspacePolicyFragmentId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.PolicyFragmentName, ok = input.Parsed["policyFragmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "policyFragmentName", input)
	}

	return nil
}

// ValidateWorkspacePolicyFragmentID checks that 'input' can be parsed as a Workspace Policy Fragment ID
func ValidateWorkspacePolicyFragmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWorkspacePolicyFragmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Workspace Policy Fragment ID
func (id WorkspacePolicyFragmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/workspaces/%s/policyFragments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId, id.PolicyFragmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Workspace Policy Fragment ID
func (id WorkspacePolicyFragmentId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticPolicyFragments", "policyFragments", "policyFragments"),
		resourceids.UserSpecifiedSegment("policyFragmentName", "policyFragmentValue"),
	}
}

// String returns a human-readable description of this Workspace Policy Fragment ID
func (id WorkspacePolicyFragmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Workspace: %q", id.WorkspaceId),
		fmt.Sprintf("Policy Fragment Name: %q", id.PolicyFragmentName),
	}
	return fmt.Sprintf("Workspace Policy Fragment (%s)", strings.Join(components, "\n"))
}
