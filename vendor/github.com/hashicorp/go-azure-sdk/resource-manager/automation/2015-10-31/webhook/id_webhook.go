package webhook

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = WebHookId{}

// WebHookId is a struct representing the Resource ID for a Web Hook
type WebHookId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
	WebHookName           string
}

// NewWebHookID returns a new WebHookId struct
func NewWebHookID(subscriptionId string, resourceGroupName string, automationAccountName string, webHookName string) WebHookId {
	return WebHookId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
		WebHookName:           webHookName,
	}
}

// ParseWebHookID parses 'input' into a WebHookId
func ParseWebHookID(input string) (*WebHookId, error) {
	parser := resourceids.NewParserFromResourceIdType(WebHookId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := WebHookId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.WebHookName, ok = parsed.Parsed["webHookName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "webHookName", *parsed)
	}

	return &id, nil
}

// ParseWebHookIDInsensitively parses 'input' case-insensitively into a WebHookId
// note: this method should only be used for API response data and not user input
func ParseWebHookIDInsensitively(input string) (*WebHookId, error) {
	parser := resourceids.NewParserFromResourceIdType(WebHookId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := WebHookId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.WebHookName, ok = parsed.Parsed["webHookName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "webHookName", *parsed)
	}

	return &id, nil
}

// ValidateWebHookID checks that 'input' can be parsed as a Web Hook ID
func ValidateWebHookID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWebHookID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Web Hook ID
func (id WebHookId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/webHooks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.WebHookName)
}

// Segments returns a slice of Resource ID Segments which comprise this Web Hook ID
func (id WebHookId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountValue"),
		resourceids.StaticSegment("staticWebHooks", "webHooks", "webHooks"),
		resourceids.UserSpecifiedSegment("webHookName", "webHookValue"),
	}
}

// String returns a human-readable description of this Web Hook ID
func (id WebHookId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Web Hook Name: %q", id.WebHookName),
	}
	return fmt.Sprintf("Web Hook (%s)", strings.Join(components, "\n"))
}
