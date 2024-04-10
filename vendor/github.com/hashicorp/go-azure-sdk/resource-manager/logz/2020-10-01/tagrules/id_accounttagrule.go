package tagrules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &AccountTagRuleId{}

// AccountTagRuleId is a struct representing the Resource ID for a Account Tag Rule
type AccountTagRuleId struct {
	SubscriptionId    string
	ResourceGroupName string
	MonitorName       string
	AccountName       string
	TagRuleName       string
}

// NewAccountTagRuleID returns a new AccountTagRuleId struct
func NewAccountTagRuleID(subscriptionId string, resourceGroupName string, monitorName string, accountName string, tagRuleName string) AccountTagRuleId {
	return AccountTagRuleId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MonitorName:       monitorName,
		AccountName:       accountName,
		TagRuleName:       tagRuleName,
	}
}

// ParseAccountTagRuleID parses 'input' into a AccountTagRuleId
func ParseAccountTagRuleID(input string) (*AccountTagRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AccountTagRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AccountTagRuleId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAccountTagRuleIDInsensitively parses 'input' case-insensitively into a AccountTagRuleId
// note: this method should only be used for API response data and not user input
func ParseAccountTagRuleIDInsensitively(input string) (*AccountTagRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AccountTagRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AccountTagRuleId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AccountTagRuleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.MonitorName, ok = input.Parsed["monitorName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "monitorName", input)
	}

	if id.AccountName, ok = input.Parsed["accountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "accountName", input)
	}

	if id.TagRuleName, ok = input.Parsed["tagRuleName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "tagRuleName", input)
	}

	return nil
}

// ValidateAccountTagRuleID checks that 'input' can be parsed as a Account Tag Rule ID
func ValidateAccountTagRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAccountTagRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Account Tag Rule ID
func (id AccountTagRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logz/monitors/%s/accounts/%s/tagRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MonitorName, id.AccountName, id.TagRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Account Tag Rule ID
func (id AccountTagRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLogz", "Microsoft.Logz", "Microsoft.Logz"),
		resourceids.StaticSegment("staticMonitors", "monitors", "monitors"),
		resourceids.UserSpecifiedSegment("monitorName", "monitorValue"),
		resourceids.StaticSegment("staticAccounts", "accounts", "accounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticTagRules", "tagRules", "tagRules"),
		resourceids.UserSpecifiedSegment("tagRuleName", "tagRuleValue"),
	}
}

// String returns a human-readable description of this Account Tag Rule ID
func (id AccountTagRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Monitor Name: %q", id.MonitorName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Tag Rule Name: %q", id.TagRuleName),
	}
	return fmt.Sprintf("Account Tag Rule (%s)", strings.Join(components, "\n"))
}
