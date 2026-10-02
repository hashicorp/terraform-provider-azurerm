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
	recaser.RegisterResourceId(&VersionStepId{})
}

var _ resourceids.ResourceId = &VersionStepId{}

// VersionStepId is a struct representing the Resource ID for a Version Step
type VersionStepId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServerName        string
	JobAgentName      string
	JobName           string
	VersionName       string
	StepName          string
}

// NewVersionStepID returns a new VersionStepId struct
func NewVersionStepID(subscriptionId string, resourceGroupName string, serverName string, jobAgentName string, jobName string, versionName string, stepName string) VersionStepId {
	return VersionStepId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServerName:        serverName,
		JobAgentName:      jobAgentName,
		JobName:           jobName,
		VersionName:       versionName,
		StepName:          stepName,
	}
}

// ParseVersionStepID parses 'input' into a VersionStepId
func ParseVersionStepID(input string) (*VersionStepId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VersionStepId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VersionStepId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseVersionStepIDInsensitively parses 'input' case-insensitively into a VersionStepId
// note: this method should only be used for API response data and not user input
func ParseVersionStepIDInsensitively(input string) (*VersionStepId, error) {
	parser := resourceids.NewParserFromResourceIdType(&VersionStepId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := VersionStepId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *VersionStepId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.VersionName, ok = input.Parsed["versionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "versionName", input)
	}

	if id.StepName, ok = input.Parsed["stepName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "stepName", input)
	}

	return nil
}

// ValidateVersionStepID checks that 'input' can be parsed as a Version Step ID
func ValidateVersionStepID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVersionStepID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Version Step ID
func (id VersionStepId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/jobAgents/%s/jobs/%s/versions/%s/steps/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.JobAgentName, id.JobName, id.VersionName, id.StepName)
}

// Segments returns a slice of Resource ID Segments which comprise this Version Step ID
func (id VersionStepId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticVersions", "versions", "versions"),
		resourceids.UserSpecifiedSegment("versionName", "versionName"),
		resourceids.StaticSegment("staticSteps", "steps", "steps"),
		resourceids.UserSpecifiedSegment("stepName", "stepName"),
	}
}

// String returns a human-readable description of this Version Step ID
func (id VersionStepId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Server Name: %q", id.ServerName),
		fmt.Sprintf("Job Agent Name: %q", id.JobAgentName),
		fmt.Sprintf("Job Name: %q", id.JobName),
		fmt.Sprintf("Version Name: %q", id.VersionName),
		fmt.Sprintf("Step Name: %q", id.StepName),
	}
	return fmt.Sprintf("Version Step (%s)", strings.Join(components, "\n"))
}
