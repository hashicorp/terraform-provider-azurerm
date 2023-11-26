package testjobstream

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = TestJobStreamId{}

// TestJobStreamId is a struct representing the Resource ID for a Test Job Stream
type TestJobStreamId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
	RunbookName           string
	JobStreamId           string
}

// NewTestJobStreamID returns a new TestJobStreamId struct
func NewTestJobStreamID(subscriptionId string, resourceGroupName string, automationAccountName string, runbookName string, jobStreamId string) TestJobStreamId {
	return TestJobStreamId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
		RunbookName:           runbookName,
		JobStreamId:           jobStreamId,
	}
}

// ParseTestJobStreamID parses 'input' into a TestJobStreamId
func ParseTestJobStreamID(input string) (*TestJobStreamId, error) {
	parser := resourceids.NewParserFromResourceIdType(TestJobStreamId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TestJobStreamId{}

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

	if id.JobStreamId, ok = parsed.Parsed["jobStreamId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "jobStreamId", *parsed)
	}

	return &id, nil
}

// ParseTestJobStreamIDInsensitively parses 'input' case-insensitively into a TestJobStreamId
// note: this method should only be used for API response data and not user input
func ParseTestJobStreamIDInsensitively(input string) (*TestJobStreamId, error) {
	parser := resourceids.NewParserFromResourceIdType(TestJobStreamId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := TestJobStreamId{}

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

	if id.JobStreamId, ok = parsed.Parsed["jobStreamId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "jobStreamId", *parsed)
	}

	return &id, nil
}

// ValidateTestJobStreamID checks that 'input' can be parsed as a Test Job Stream ID
func ValidateTestJobStreamID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTestJobStreamID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Test Job Stream ID
func (id TestJobStreamId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/runbooks/%s/draft/testJob/streams/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.RunbookName, id.JobStreamId)
}

// Segments returns a slice of Resource ID Segments which comprise this Test Job Stream ID
func (id TestJobStreamId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticDraft", "draft", "draft"),
		resourceids.StaticSegment("staticTestJob", "testJob", "testJob"),
		resourceids.StaticSegment("staticStreams", "streams", "streams"),
		resourceids.UserSpecifiedSegment("jobStreamId", "jobStreamIdValue"),
	}
}

// String returns a human-readable description of this Test Job Stream ID
func (id TestJobStreamId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Runbook Name: %q", id.RunbookName),
		fmt.Sprintf("Job Stream: %q", id.JobStreamId),
	}
	return fmt.Sprintf("Test Job Stream (%s)", strings.Join(components, "\n"))
}
