package testjobstream

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RunbookId{}

// RunbookId is a struct representing the Resource ID for a Runbook
type RunbookId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
	RunbookName           string
}

// NewRunbookID returns a new RunbookId struct
func NewRunbookID(subscriptionId string, resourceGroupName string, automationAccountName string, runbookName string) RunbookId {
	return RunbookId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
		RunbookName:           runbookName,
	}
}

// ParseRunbookID parses 'input' into a RunbookId
func ParseRunbookID(input string) (*RunbookId, error) {
	parser := resourceids.NewParserFromResourceIdType(RunbookId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RunbookId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.RunbookName, ok = parsed.Parsed["runbookName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "runbookName", *parsed)
	}

	return &id, nil
}

// ParseRunbookIDInsensitively parses 'input' case-insensitively into a RunbookId
// note: this method should only be used for API response data and not user input
func ParseRunbookIDInsensitively(input string) (*RunbookId, error) {
	parser := resourceids.NewParserFromResourceIdType(RunbookId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RunbookId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.RunbookName, ok = parsed.Parsed["runbookName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "runbookName", *parsed)
	}

	return &id, nil
}

// ValidateRunbookID checks that 'input' can be parsed as a Runbook ID
func ValidateRunbookID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRunbookID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Runbook ID
func (id RunbookId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/runbooks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.RunbookName)
}

// Segments returns a slice of Resource ID Segments which comprise this Runbook ID
func (id RunbookId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountValue"),
		resourceids.StaticSegment("staticRunbooks", "runbooks", "runbooks"),
		resourceids.UserSpecifiedSegment("runbookName", "runbookValue"),
	}
}

// String returns a human-readable description of this Runbook ID
func (id RunbookId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Runbook Name: %q", id.RunbookName),
	}
	return fmt.Sprintf("Runbook (%s)", strings.Join(components, "\n"))
}
