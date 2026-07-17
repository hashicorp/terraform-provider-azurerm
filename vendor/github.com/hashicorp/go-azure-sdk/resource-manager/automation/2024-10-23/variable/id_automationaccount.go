package variable

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AutomationAccountId{})
}

var _ resourceids.ResourceId = &AutomationAccountId{}

// AutomationAccountId is a struct representing the Resource ID for a Automation Account
type AutomationAccountId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
}

// NewAutomationAccountID returns a new AutomationAccountId struct
func NewAutomationAccountID(subscriptionId string, resourceGroupName string, automationAccountName string) AutomationAccountId {
	return AutomationAccountId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
	}
}

// ParseAutomationAccountID parses 'input' into a AutomationAccountId
func ParseAutomationAccountID(input string) (*AutomationAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutomationAccountId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutomationAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAutomationAccountIDInsensitively parses 'input' case-insensitively into a AutomationAccountId
// note: this method should only be used for API response data and not user input
func ParseAutomationAccountIDInsensitively(input string) (*AutomationAccountId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutomationAccountId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutomationAccountId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AutomationAccountId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateAutomationAccountID checks that 'input' can be parsed as a Automation Account ID
func ValidateAutomationAccountID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAutomationAccountID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Automation Account ID
func (id AutomationAccountId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Automation Account ID
func (id AutomationAccountId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountName"),
	}
}

// String returns a human-readable description of this Automation Account ID
func (id AutomationAccountId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
	}
	return fmt.Sprintf("Automation Account (%s)", strings.Join(components, "\n"))
}
