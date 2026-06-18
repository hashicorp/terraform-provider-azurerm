package integrationruntime

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NodeId{})
}

var _ resourceids.ResourceId = &NodeId{}

// NodeId is a struct representing the Resource ID for a Node
type NodeId struct {
	SubscriptionId         string
	ResourceGroupName      string
	WorkspaceName          string
	IntegrationRuntimeName string
	NodeName               string
}

// NewNodeID returns a new NodeId struct
func NewNodeID(subscriptionId string, resourceGroupName string, workspaceName string, integrationRuntimeName string, nodeName string) NodeId {
	return NodeId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		WorkspaceName:          workspaceName,
		IntegrationRuntimeName: integrationRuntimeName,
		NodeName:               nodeName,
	}
}

// ParseNodeID parses 'input' into a NodeId
func ParseNodeID(input string) (*NodeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NodeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NodeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNodeIDInsensitively parses 'input' case-insensitively into a NodeId
// note: this method should only be used for API response data and not user input
func ParseNodeIDInsensitively(input string) (*NodeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NodeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NodeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NodeId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.WorkspaceName, ok = input.Parsed["workspaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", input)
	}

	if id.IntegrationRuntimeName, ok = input.Parsed["integrationRuntimeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "integrationRuntimeName", input)
	}

	if id.NodeName, ok = input.Parsed["nodeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "nodeName", input)
	}

	return nil
}

// ValidateNodeID checks that 'input' can be parsed as a Node ID
func ValidateNodeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNodeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Node ID
func (id NodeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Synapse/workspaces/%s/integrationRuntimes/%s/nodes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.IntegrationRuntimeName, id.NodeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Node ID
func (id NodeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSynapse", "Microsoft.Synapse", "Microsoft.Synapse"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceName"),
		resourceids.StaticSegment("staticIntegrationRuntimes", "integrationRuntimes", "integrationRuntimes"),
		resourceids.UserSpecifiedSegment("integrationRuntimeName", "integrationRuntimeName"),
		resourceids.StaticSegment("staticNodes", "nodes", "nodes"),
		resourceids.UserSpecifiedSegment("nodeName", "nodeName"),
	}
}

// String returns a human-readable description of this Node ID
func (id NodeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Integration Runtime Name: %q", id.IntegrationRuntimeName),
		fmt.Sprintf("Node Name: %q", id.NodeName),
	}
	return fmt.Sprintf("Node (%s)", strings.Join(components, "\n"))
}
