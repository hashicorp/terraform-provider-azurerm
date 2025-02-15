package jobs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&JobId{})
}

var _ resourceids.ResourceId = &JobId{}

// JobId is a struct representing the Resource ID for a Job
type JobId struct {
	SubscriptionId    string
	ResourceGroupName string
	JobName           string
}

// NewJobID returns a new JobId struct
func NewJobID(subscriptionId string, resourceGroupName string, jobName string) JobId {
	return JobId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		JobName:           jobName,
	}
}

// ParseJobID parses 'input' into a JobId
func ParseJobID(input string) (*JobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&JobId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := JobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseJobIDInsensitively parses 'input' case-insensitively into a JobId
// note: this method should only be used for API response data and not user input
func ParseJobIDInsensitively(input string) (*JobId, error) {
	parser := resourceids.NewParserFromResourceIdType(&JobId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := JobId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *JobId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.JobName, ok = input.Parsed["jobName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "jobName", input)
	}

	return nil
}

// ValidateJobID checks that 'input' can be parsed as a Job ID
func ValidateJobID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseJobID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Job ID
func (id JobId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/jobs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.JobName)
}

// Segments returns a slice of Resource ID Segments which comprise this Job ID
func (id JobId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticJobs", "jobs", "jobs"),
		resourceids.UserSpecifiedSegment("jobName", "jobName"),
	}
}

// String returns a human-readable description of this Job ID
func (id JobId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Job Name: %q", id.JobName),
	}
	return fmt.Sprintf("Job (%s)", strings.Join(components, "\n"))
}
