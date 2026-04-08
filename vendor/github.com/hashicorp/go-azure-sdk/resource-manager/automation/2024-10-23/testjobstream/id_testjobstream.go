package testjobstream

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TestJobStreamId{})
}

var _ resourceids.ResourceId = &TestJobStreamId{}

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
	parser := resourceids.NewParserFromResourceIdType(&TestJobStreamId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TestJobStreamId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTestJobStreamIDInsensitively parses 'input' case-insensitively into a TestJobStreamId
// note: this method should only be used for API response data and not user input
func ParseTestJobStreamIDInsensitively(input string) (*TestJobStreamId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TestJobStreamId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TestJobStreamId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TestJobStreamId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.RunbookName, ok = input.Parsed["runbookName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "runbookName", input)
	}

	if id.JobStreamId, ok = input.Parsed["jobStreamId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "jobStreamId", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountName"),
		resourceids.StaticSegment("staticRunbooks", "runbooks", "runbooks"),
		resourceids.UserSpecifiedSegment("runbookName", "runbookName"),
		resourceids.StaticSegment("staticDraft", "draft", "draft"),
		resourceids.StaticSegment("staticTestJob", "testJob", "testJob"),
		resourceids.StaticSegment("staticStreams", "streams", "streams"),
		resourceids.UserSpecifiedSegment("jobStreamId", "jobStreamId"),
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
