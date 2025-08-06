package automations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AutomationId{})
}

var _ resourceids.ResourceId = &AutomationId{}

// AutomationId is a struct representing the Resource ID for a Automation
type AutomationId struct {
	SubscriptionId    string
	ResourceGroupName string
	AutomationName    string
}

// NewAutomationID returns a new AutomationId struct
func NewAutomationID(subscriptionId string, resourceGroupName string, automationName string) AutomationId {
	return AutomationId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AutomationName:    automationName,
	}
}

// ParseAutomationID parses 'input' into a AutomationId
func ParseAutomationID(input string) (*AutomationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutomationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutomationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAutomationIDInsensitively parses 'input' case-insensitively into a AutomationId
// note: this method should only be used for API response data and not user input
func ParseAutomationIDInsensitively(input string) (*AutomationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutomationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutomationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AutomationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AutomationName, ok = input.Parsed["automationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "automationName", input)
	}

	return nil
}

// ValidateAutomationID checks that 'input' can be parsed as a Automation ID
func ValidateAutomationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAutomationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Automation ID
func (id AutomationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Security/automations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Automation ID
func (id AutomationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSecurity", "Microsoft.Security", "Microsoft.Security"),
		resourceids.StaticSegment("staticAutomations", "automations", "automations"),
		resourceids.UserSpecifiedSegment("automationName", "automationName"),
	}
}

// String returns a human-readable description of this Automation ID
func (id AutomationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Name: %q", id.AutomationName),
	}
	return fmt.Sprintf("Automation (%s)", strings.Join(components, "\n"))
}
