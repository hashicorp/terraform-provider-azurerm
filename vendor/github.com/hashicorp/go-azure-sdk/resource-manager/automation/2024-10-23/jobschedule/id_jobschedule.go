package jobschedule

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&JobScheduleId{})
}

var _ resourceids.ResourceId = &JobScheduleId{}

// JobScheduleId is a struct representing the Resource ID for a Job Schedule
type JobScheduleId struct {
	SubscriptionId        string
	ResourceGroupName     string
	AutomationAccountName string
	JobScheduleId         string
}

// NewJobScheduleID returns a new JobScheduleId struct
func NewJobScheduleID(subscriptionId string, resourceGroupName string, automationAccountName string, jobScheduleId string) JobScheduleId {
	return JobScheduleId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		AutomationAccountName: automationAccountName,
		JobScheduleId:         jobScheduleId,
	}
}

// ParseJobScheduleID parses 'input' into a JobScheduleId
func ParseJobScheduleID(input string) (*JobScheduleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&JobScheduleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := JobScheduleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseJobScheduleIDInsensitively parses 'input' case-insensitively into a JobScheduleId
// note: this method should only be used for API response data and not user input
func ParseJobScheduleIDInsensitively(input string) (*JobScheduleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&JobScheduleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := JobScheduleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *JobScheduleId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.JobScheduleId, ok = input.Parsed["jobScheduleId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "jobScheduleId", input)
	}

	return nil
}

// ValidateJobScheduleID checks that 'input' can be parsed as a Job Schedule ID
func ValidateJobScheduleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseJobScheduleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Job Schedule ID
func (id JobScheduleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/jobSchedules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, id.JobScheduleId)
}

// Segments returns a slice of Resource ID Segments which comprise this Job Schedule ID
func (id JobScheduleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutomation", "Microsoft.Automation", "Microsoft.Automation"),
		resourceids.StaticSegment("staticAutomationAccounts", "automationAccounts", "automationAccounts"),
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountName"),
		resourceids.StaticSegment("staticJobSchedules", "jobSchedules", "jobSchedules"),
		resourceids.UserSpecifiedSegment("jobScheduleId", "jobScheduleId"),
	}
}

// String returns a human-readable description of this Job Schedule ID
func (id JobScheduleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Automation Account Name: %q", id.AutomationAccountName),
		fmt.Sprintf("Job Schedule: %q", id.JobScheduleId),
	}
	return fmt.Sprintf("Job Schedule (%s)", strings.Join(components, "\n"))
}
