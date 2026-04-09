package nodetype

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NodeTypeId{})
}

var _ resourceids.ResourceId = &NodeTypeId{}

// NodeTypeId is a struct representing the Resource ID for a Node Type
type NodeTypeId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ManagedClusterName string
	NodeTypeName       string
}

// NewNodeTypeID returns a new NodeTypeId struct
func NewNodeTypeID(subscriptionId string, resourceGroupName string, managedClusterName string, nodeTypeName string) NodeTypeId {
	return NodeTypeId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ManagedClusterName: managedClusterName,
		NodeTypeName:       nodeTypeName,
	}
}

// ParseNodeTypeID parses 'input' into a NodeTypeId
func ParseNodeTypeID(input string) (*NodeTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NodeTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NodeTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNodeTypeIDInsensitively parses 'input' case-insensitively into a NodeTypeId
// note: this method should only be used for API response data and not user input
func ParseNodeTypeIDInsensitively(input string) (*NodeTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NodeTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NodeTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NodeTypeId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ManagedClusterName, ok = input.Parsed["managedClusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "managedClusterName", input)
	}

	if id.NodeTypeName, ok = input.Parsed["nodeTypeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "nodeTypeName", input)
	}

	return nil
}

// ValidateNodeTypeID checks that 'input' can be parsed as a Node Type ID
func ValidateNodeTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNodeTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Node Type ID
func (id NodeTypeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ServiceFabric/managedClusters/%s/nodeTypes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, id.NodeTypeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Node Type ID
func (id NodeTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftServiceFabric", "Microsoft.ServiceFabric", "Microsoft.ServiceFabric"),
		resourceids.StaticSegment("staticManagedClusters", "managedClusters", "managedClusters"),
		resourceids.UserSpecifiedSegment("managedClusterName", "managedClusterName"),
		resourceids.StaticSegment("staticNodeTypes", "nodeTypes", "nodeTypes"),
		resourceids.UserSpecifiedSegment("nodeTypeName", "nodeTypeName"),
	}
}

// String returns a human-readable description of this Node Type ID
func (id NodeTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Cluster Name: %q", id.ManagedClusterName),
		fmt.Sprintf("Node Type Name: %q", id.NodeTypeName),
	}
	return fmt.Sprintf("Node Type (%s)", strings.Join(components, "\n"))
}
