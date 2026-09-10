package nodereports

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
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
	NodeId                string
}

// NewNodeID returns a new NodeId struct
func NewNodeID(subscriptionId string, resourceGroupName string, automationAccountName string, nodeId string) NodeId {
	return NodeId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
		NodeId:                nodeId,
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

	if id.AutomationAccountName, ok = input.Parsed["automationAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", input)
	}

	if id.NodeId, ok = input.Parsed["nodeId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "nodeId", input)
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
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/nodes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.NodeId)
}

// Segments returns a slice of Resource ID Segments which comprise this Node ID
func (id NodeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountName"),
		resourceids.StaticSegment("staticNodes", "nodes", "nodes"),
		resourceids.UserSpecifiedSegment("nodeId", "nodeId"),
	}
}

// String returns a human-readable description of this Node ID
func (id NodeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Node: %q", id.NodeId),
	}
	return fmt.Sprintf("Node (%s)", strings.Join(components, "\n"))
}
