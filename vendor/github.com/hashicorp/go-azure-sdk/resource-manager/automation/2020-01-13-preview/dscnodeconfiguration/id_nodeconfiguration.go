package dscnodeconfiguration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = NodeConfigurationId{}

// NodeConfigurationId is a struct representing the Resource ID for a Node Configuration
type NodeConfigurationId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
	NodeConfigurationName string
}

// NewNodeConfigurationID returns a new NodeConfigurationId struct
func NewNodeConfigurationID(subscriptionId string, resourceGroupName string, automationAccountName string, nodeConfigurationName string) NodeConfigurationId {
	return NodeConfigurationId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
		NodeConfigurationName: nodeConfigurationName,
	}
}

// ParseNodeConfigurationID parses 'input' into a NodeConfigurationId
func ParseNodeConfigurationID(input string) (*NodeConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(NodeConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NodeConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.NodeConfigurationName, ok = parsed.Parsed["nodeConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "nodeConfigurationName", *parsed)
	}

	return &id, nil
}

// ParseNodeConfigurationIDInsensitively parses 'input' case-insensitively into a NodeConfigurationId
// note: this method should only be used for API response data and not user input
func ParseNodeConfigurationIDInsensitively(input string) (*NodeConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(NodeConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := NodeConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.NodeConfigurationName, ok = parsed.Parsed["nodeConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "nodeConfigurationName", *parsed)
	}

	return &id, nil
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
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/nodeConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.NodeConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Node Configuration ID
func (id NodeConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountValue"),
		resourceids.StaticSegment("staticNodeConfigurations", "nodeConfigurations", "nodeConfigurations"),
		resourceids.UserSpecifiedSegment("nodeConfigurationName", "nodeConfigurationValue"),
	}
}

// String returns a human-readable description of this Node Configuration ID
func (id NodeConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Node Configuration Name: %q", id.NodeConfigurationName),
	}
	return fmt.Sprintf("Node Configuration (%s)", strings.Join(components, "\n"))
}
