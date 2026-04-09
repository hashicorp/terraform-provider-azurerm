package automationrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AutomationRuleId{})
}

var _ resourceids.ResourceId = &AutomationRuleId{}

// AutomationRuleId is a struct representing the Resource ID for a Automation Rule
type AutomationRuleId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	AutomationRuleId  string
}

// NewAutomationRuleID returns a new AutomationRuleId struct
func NewAutomationRuleID(subscriptionId string, resourceGroupName string, workspaceName string, automationRuleId string) AutomationRuleId {
	return AutomationRuleId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		AutomationRuleId:  automationRuleId,
	}
}

// ParseAutomationRuleID parses 'input' into a AutomationRuleId
func ParseAutomationRuleID(input string) (*AutomationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutomationRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutomationRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAutomationRuleIDInsensitively parses 'input' case-insensitively into a AutomationRuleId
// note: this method should only be used for API response data and not user input
func ParseAutomationRuleIDInsensitively(input string) (*AutomationRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutomationRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutomationRuleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AutomationRuleId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.AutomationRuleId, ok = input.Parsed["automationRuleId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "automationRuleId", input)
	}

	return nil
}

// ValidateAutomationRuleID checks that 'input' can be parsed as a Automation Rule ID
func ValidateAutomationRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAutomationRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Automation Rule ID
func (id AutomationRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/automationRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.AutomationRuleId)
}

// Segments returns a slice of Resource ID Segments which comprise this Automation Rule ID
func (id AutomationRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceName"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSecurityInsights", "Microsoft.SecurityInsights", "Microsoft.SecurityInsights"),
		resourceids.StaticSegment("staticAutomationRules", "automationRules", "automationRules"),
		resourceids.UserSpecifiedSegment("automationRuleId", "automationRuleId"),
	}
}

// String returns a human-readable description of this Automation Rule ID
func (id AutomationRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Automation Rule: %q", id.AutomationRuleId),
	}
	return fmt.Sprintf("Automation Rule (%s)", strings.Join(components, "\n"))
}
