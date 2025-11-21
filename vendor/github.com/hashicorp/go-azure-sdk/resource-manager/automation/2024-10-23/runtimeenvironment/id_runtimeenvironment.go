package runtimeenvironment

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RuntimeEnvironmentId{})
}

var _ resourceids.ResourceId = &RuntimeEnvironmentId{}

// RuntimeEnvironmentId is a struct representing the Resource ID for a Runtime Environment
type RuntimeEnvironmentId struct {
	SubscriptionId         string
	ResourceGroupName      string
	AutomationAccountName  string
	RuntimeEnvironmentName string
}

// NewRuntimeEnvironmentID returns a new RuntimeEnvironmentId struct
func NewRuntimeEnvironmentID(subscriptionId string, resourceGroupName string, automationAccountName string, runtimeEnvironmentName string) RuntimeEnvironmentId {
	return RuntimeEnvironmentId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		AutomationAccountName:  automationAccountName,
		RuntimeEnvironmentName: runtimeEnvironmentName,
	}
}

// ParseRuntimeEnvironmentID parses 'input' into a RuntimeEnvironmentId
func ParseRuntimeEnvironmentID(input string) (*RuntimeEnvironmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RuntimeEnvironmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RuntimeEnvironmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRuntimeEnvironmentIDInsensitively parses 'input' case-insensitively into a RuntimeEnvironmentId
// note: this method should only be used for API response data and not user input
func ParseRuntimeEnvironmentIDInsensitively(input string) (*RuntimeEnvironmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RuntimeEnvironmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RuntimeEnvironmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RuntimeEnvironmentId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.RuntimeEnvironmentName, ok = input.Parsed["runtimeEnvironmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "runtimeEnvironmentName", input)
	}

	return nil
}

// ValidateRuntimeEnvironmentID checks that 'input' can be parsed as a Runtime Environment ID
func ValidateRuntimeEnvironmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRuntimeEnvironmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Runtime Environment ID
func (id RuntimeEnvironmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/runtimeEnvironments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.RuntimeEnvironmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Runtime Environment ID
func (id RuntimeEnvironmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountName"),
		resourceids.StaticSegment("staticRuntimeEnvironments", "runtimeEnvironments", "runtimeEnvironments"),
		resourceids.UserSpecifiedSegment("runtimeEnvironmentName", "runtimeEnvironmentName"),
	}
}

// String returns a human-readable description of this Runtime Environment ID
func (id RuntimeEnvironmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Runtime Environment Name: %q", id.RuntimeEnvironmentName),
	}
	return fmt.Sprintf("Runtime Environment (%s)", strings.Join(components, "\n"))
}
