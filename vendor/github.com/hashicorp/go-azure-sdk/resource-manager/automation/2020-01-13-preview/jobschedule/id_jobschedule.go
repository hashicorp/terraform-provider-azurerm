package jobschedule

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = JobScheduleId{}

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
	parser := resourceids.NewParserFromResourceIdType(JobScheduleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := JobScheduleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.JobScheduleId, ok = parsed.Parsed["jobScheduleId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "jobScheduleId", *parsed)
	}

	return &id, nil
}

// ParseJobScheduleIDInsensitively parses 'input' case-insensitively into a JobScheduleId
// note: this method should only be used for API response data and not user input
func ParseJobScheduleIDInsensitively(input string) (*JobScheduleId, error) {
	parser := resourceids.NewParserFromResourceIdType(JobScheduleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := JobScheduleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.AutomationAccountName, ok = parsed.Parsed["automationAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "automationAccountName", *parsed)
	}

	if id.JobScheduleId, ok = parsed.Parsed["jobScheduleId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "jobScheduleId", *parsed)
	}

	return &id, nil
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
		resourceids.UserSpecifiedSegment("automationAccountName", "automationAccountValue"),
		resourceids.StaticSegment("staticJobSchedules", "jobSchedules", "jobSchedules"),
		resourceids.UserSpecifiedSegment("jobScheduleId", "jobScheduleIdValue"),
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
