package configurations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NodeConfigurationId{})
}

var _ resourceids.ResourceId = &NodeConfigurationId{}

// NodeConfigurationId is a struct representing the Resource ID for a Node Configuration
type NodeConfigurationId struct {
	SubscriptionId        string
	ResourceGroupName     string
	ServerGroupsv2Name    string
	NodeConfigurationName string
}

// NewNodeConfigurationID returns a new NodeConfigurationId struct
func NewNodeConfigurationID(subscriptionId string, resourceGroupName string, serverGroupsv2Name string, nodeConfigurationName string) NodeConfigurationId {
	return NodeConfigurationId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		ServerGroupsv2Name:    serverGroupsv2Name,
		NodeConfigurationName: nodeConfigurationName,
	}
}

// ParseNodeConfigurationID parses 'input' into a NodeConfigurationId
func ParseNodeConfigurationID(input string) (*NodeConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NodeConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NodeConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNodeConfigurationIDInsensitively parses 'input' case-insensitively into a NodeConfigurationId
// note: this method should only be used for API response data and not user input
func ParseNodeConfigurationIDInsensitively(input string) (*NodeConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NodeConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NodeConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NodeConfigurationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServerGroupsv2Name, ok = input.Parsed["serverGroupsv2Name"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serverGroupsv2Name", input)
	}

	if id.NodeConfigurationName, ok = input.Parsed["nodeConfigurationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "nodeConfigurationName", input)
	}

	return nil
}

// ValidateNodeConfigurationID checks that 'input' can be parsed as a Node Configuration ID
func ValidateNodeConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNodeConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Node Configuration ID
func (id NodeConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforPostgreSQL/serverGroupsv2/%s/nodeConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServerGroupsv2Name, id.NodeConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Node Configuration ID
func (id NodeConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDBforPostgreSQL", "Microsoft.DBforPostgreSQL", "Microsoft.DBforPostgreSQL"),
		resourceids.StaticSegment("staticServerGroupsv2", "serverGroupsv2", "serverGroupsv2"),
		resourceids.UserSpecifiedSegment("serverGroupsv2Name", "serverGroupsv2Name"),
		resourceids.StaticSegment("staticNodeConfigurations", "nodeConfigurations", "nodeConfigurations"),
		resourceids.UserSpecifiedSegment("nodeConfigurationName", "nodeConfigurationName"),
	}
}

// String returns a human-readable description of this Node Configuration ID
func (id NodeConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Server Groupsv 2 Name: %q", id.ServerGroupsv2Name),
		fmt.Sprintf("Node Configuration Name: %q", id.NodeConfigurationName),
	}
	return fmt.Sprintf("Node Configuration (%s)", strings.Join(components, "\n"))
}
