package jobsteps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&StepId{})
}

var _ resourceids.ResourceId = &StepId{}

// StepId is a struct representing the Resource ID for a Step
type StepId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServerName        string
	JobAgentName      string
	JobName           string
	StepName          string
}

// NewStepID returns a new StepId struct
func NewStepID(subscriptionId string, resourceGroupName string, serverName string, jobAgentName string, jobName string, stepName string) StepId {
	return StepId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServerName:        serverName,
		JobAgentName:      jobAgentName,
		JobName:           jobName,
		StepName:          stepName,
	}
}

// ParseStepID parses 'input' into a StepId
func ParseStepID(input string) (*StepId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StepId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StepId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStepIDInsensitively parses 'input' case-insensitively into a StepId
// note: this method should only be used for API response data and not user input
func ParseStepIDInsensitively(input string) (*StepId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StepId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StepId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StepId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServerName, ok = input.Parsed["serverName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serverName", input)
	}

	if id.JobAgentName, ok = input.Parsed["jobAgentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "jobAgentName", input)
	}

	if id.JobName, ok = input.Parsed["jobName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "jobName", input)
	}

	if id.StepName, ok = input.Parsed["stepName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "stepName", input)
	}

	return nil
}

// ValidateStepID checks that 'input' can be parsed as a Step ID
func ValidateStepID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStepID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Step ID
func (id StepId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/jobAgents/%s/jobs/%s/steps/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.JobAgentName, id.JobName, id.StepName)
}

// Segments returns a slice of Resource ID Segments which comprise this Step ID
func (id StepId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSql", "Microsoft.Sql", "Microsoft.Sql"),
		resourceids.StaticSegment("staticServers", "servers", "servers"),
		resourceids.UserSpecifiedSegment("serverName", "serverName"),
		resourceids.StaticSegment("staticJobAgents", "jobAgents", "jobAgents"),
		resourceids.UserSpecifiedSegment("jobAgentName", "jobAgentName"),
		resourceids.StaticSegment("staticJobs", "jobs", "jobs"),
		resourceids.UserSpecifiedSegment("jobName", "jobName"),
		resourceids.StaticSegment("staticSteps", "steps", "steps"),
		resourceids.UserSpecifiedSegment("stepName", "stepName"),
	}
}

// String returns a human-readable description of this Step ID
func (id StepId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Server Name: %q", id.ServerName),
		fmt.Sprintf("Job Agent Name: %q", id.JobAgentName),
		fmt.Sprintf("Job Name: %q", id.JobName),
		fmt.Sprintf("Step Name: %q", id.StepName),
	}
	return fmt.Sprintf("Step (%s)", strings.Join(components, "\n"))
}
