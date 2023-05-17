package variable

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = VariableId{}

// VariableId is a struct representing the Resource ID for a Variable
type VariableId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
	VariableName          string
}

// NewVariableID returns a new VariableId struct
func NewVariableID(subscriptionId string, resourceGroupName string, automationAccountName string, variableName string) VariableId {
	return VariableId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
		VariableName:          variableName,
	}
}

// ParseVariableID parses 'input' into a VariableId
func ParseVariableID(input string) (*VariableId, error) {
	parser := resourceids.NewParserFromResourceIdType(VariableId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VariableId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.VariableName, ok = parsed.Parsed["variableName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "variableName", *parsed)
	}

	return &id, nil
}

// ParseVariableIDInsensitively parses 'input' case-insensitively into a VariableId
// note: this method should only be used for API response data and not user input
func ParseVariableIDInsensitively(input string) (*VariableId, error) {
	parser := resourceids.NewParserFromResourceIdType(VariableId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VariableId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.VariableName, ok = parsed.Parsed["variableName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "variableName", *parsed)
	}

	return &id, nil
}

// ValidateVariableID checks that 'input' can be parsed as a Variable ID
func ValidateVariableID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVariableID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Variable ID
func (id VariableId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/variables/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.VariableName)
}

// Segments returns a slice of Resource ID Segments which comprise this Variable ID
func (id VariableId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountValue"),
		resourceids.StaticSegment("staticVariables", "variables", "variables"),
		resourceids.UserSpecifiedSegment("variableName", "variableValue"),
	}
}

// String returns a human-readable description of this Variable ID
func (id VariableId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Variable Name: %q", id.VariableName),
	}
	return fmt.Sprintf("Variable (%s)", strings.Join(components, "\n"))
}
