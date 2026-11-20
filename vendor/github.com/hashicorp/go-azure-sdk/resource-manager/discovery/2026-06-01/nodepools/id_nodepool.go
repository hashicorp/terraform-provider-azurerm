package nodepools

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NodePoolId{})
}

var _ resourceids.ResourceId = &NodePoolId{}

// NodePoolId is a struct representing the Resource ID for a Node Pool
type NodePoolId struct {
	SubscriptionId    string
	ResourceGroupName string
	SupercomputerName string
	NodePoolName      string
}

// NewNodePoolID returns a new NodePoolId struct
func NewNodePoolID(subscriptionId string, resourceGroupName string, supercomputerName string, nodePoolName string) NodePoolId {
	return NodePoolId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SupercomputerName: supercomputerName,
		NodePoolName:      nodePoolName,
	}
}

// ParseNodePoolID parses 'input' into a NodePoolId
func ParseNodePoolID(input string) (*NodePoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NodePoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NodePoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNodePoolIDInsensitively parses 'input' case-insensitively into a NodePoolId
// note: this method should only be used for API response data and not user input
func ParseNodePoolIDInsensitively(input string) (*NodePoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NodePoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NodePoolId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NodePoolId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SupercomputerName, ok = input.Parsed["supercomputerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "supercomputerName", input)
	}

	if id.NodePoolName, ok = input.Parsed["nodePoolName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "nodePoolName", input)
	}

	return nil
}

// ValidateNodePoolID checks that 'input' can be parsed as a Node Pool ID
func ValidateNodePoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNodePoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Node Pool ID
func (id NodePoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Discovery/supercomputers/%s/nodePools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SupercomputerName, id.NodePoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Node Pool ID
func (id NodePoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDiscovery", "Microsoft.Discovery", "Microsoft.Discovery"),
		resourceids.StaticSegment("staticSupercomputers", "supercomputers", "supercomputers"),
		resourceids.UserSpecifiedSegment("supercomputerName", "supercomputerName"),
		resourceids.StaticSegment("staticNodePools", "nodePools", "nodePools"),
		resourceids.UserSpecifiedSegment("nodePoolName", "nodePoolName"),
	}
}

// String returns a human-readable description of this Node Pool ID
func (id NodePoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Supercomputer Name: %q", id.SupercomputerName),
		fmt.Sprintf("Node Pool Name: %q", id.NodePoolName),
	}
	return fmt.Sprintf("Node Pool (%s)", strings.Join(components, "\n"))
}
