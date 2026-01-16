package jobdefinitions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&JobDefinitionId{})
}

var _ resourceids.ResourceId = &JobDefinitionId{}

// JobDefinitionId is a struct representing the Resource ID for a Job Definition
type JobDefinitionId struct {
	SubscriptionId    string
	ResourceGroupName string
	StorageMoverName  string
	ProjectName       string
	JobDefinitionName string
}

// NewJobDefinitionID returns a new JobDefinitionId struct
func NewJobDefinitionID(subscriptionId string, resourceGroupName string, storageMoverName string, projectName string, jobDefinitionName string) JobDefinitionId {
	return JobDefinitionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		StorageMoverName:  storageMoverName,
		ProjectName:       projectName,
		JobDefinitionName: jobDefinitionName,
	}
}

// ParseJobDefinitionID parses 'input' into a JobDefinitionId
func ParseJobDefinitionID(input string) (*JobDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&JobDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := JobDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseJobDefinitionIDInsensitively parses 'input' case-insensitively into a JobDefinitionId
// note: this method should only be used for API response data and not user input
func ParseJobDefinitionIDInsensitively(input string) (*JobDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&JobDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := JobDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *JobDefinitionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.StorageMoverName, ok = input.Parsed["storageMoverName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageMoverName", input)
	}

	if id.ProjectName, ok = input.Parsed["projectName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "projectName", input)
	}

	if id.JobDefinitionName, ok = input.Parsed["jobDefinitionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "jobDefinitionName", input)
	}

	return nil
}

// ValidateJobDefinitionID checks that 'input' can be parsed as a Job Definition ID
func ValidateJobDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseJobDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Job Definition ID
func (id JobDefinitionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StorageMover/storageMovers/%s/projects/%s/jobDefinitions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StorageMoverName, id.ProjectName, id.JobDefinitionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Job Definition ID
func (id JobDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStorageMover", "Microsoft.StorageMover", "Microsoft.StorageMover"),
		resourceids.StaticSegment("staticStorageMovers", "storageMovers", "storageMovers"),
		resourceids.UserSpecifiedSegment("storageMoverName", "storageMoverName"),
		resourceids.StaticSegment("staticProjects", "projects", "projects"),
		resourceids.UserSpecifiedSegment("projectName", "projectName"),
		resourceids.StaticSegment("staticJobDefinitions", "jobDefinitions", "jobDefinitions"),
		resourceids.UserSpecifiedSegment("jobDefinitionName", "jobDefinitionName"),
	}
}

// String returns a human-readable description of this Job Definition ID
func (id JobDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Storage Mover Name: %q", id.StorageMoverName),
		fmt.Sprintf("Project Name: %q", id.ProjectName),
		fmt.Sprintf("Job Definition Name: %q", id.JobDefinitionName),
	}
	return fmt.Sprintf("Job Definition (%s)", strings.Join(components, "\n"))
}
