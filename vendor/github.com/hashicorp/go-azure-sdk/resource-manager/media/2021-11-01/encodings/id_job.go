package encodings

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = JobId{}

// JobId is a struct representing the Resource ID for a Job
type JobId struct {
	SubscriptionId    string
	ResourceGroupName string
	MediaServiceName  string
	TransformName     string
	JobName           string
}

// NewJobID returns a new JobId struct
func NewJobID(subscriptionId string, resourceGroupName string, mediaServiceName string, transformName string, jobName string) JobId {
	return JobId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MediaServiceName:  mediaServiceName,
		TransformName:     transformName,
		JobName:           jobName,
	}
}

// ParseJobID parses 'input' into a JobId
func ParseJobID(input string) (*JobId, error) {
	parser := resourceids.NewParserFromResourceIdType(JobId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := JobId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MediaServiceName, ok = parsed.Parsed["mediaServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mediaServiceName", *parsed)
	}

	if id.TransformName, ok = parsed.Parsed["transformName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "transformName", *parsed)
	}

	if id.JobName, ok = parsed.Parsed["jobName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "jobName", *parsed)
	}

	return &id, nil
}

// ParseJobIDInsensitively parses 'input' case-insensitively into a JobId
// note: this method should only be used for API response data and not user input
func ParseJobIDInsensitively(input string) (*JobId, error) {
	parser := resourceids.NewParserFromResourceIdType(JobId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := JobId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.MediaServiceName, ok = parsed.Parsed["mediaServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mediaServiceName", *parsed)
	}

	if id.TransformName, ok = parsed.Parsed["transformName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "transformName", *parsed)
	}

	if id.JobName, ok = parsed.Parsed["jobName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "jobName", *parsed)
	}

	return &id, nil
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
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaServices/%s/transforms/%s/jobs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MediaServiceName, id.TransformName, id.JobName)
}

// Segments returns a slice of Resource ID Segments which comprise this Job ID
func (id JobId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMedia", "Microsoft.Media", "Microsoft.Media"),
		resourceids.StaticSegment("staticMediaServices", "mediaServices", "mediaServices"),
		resourceids.UserSpecifiedSegment("mediaServiceName", "mediaServiceValue"),
		resourceids.StaticSegment("staticTransforms", "transforms", "transforms"),
		resourceids.UserSpecifiedSegment("transformName", "transformValue"),
		resourceids.StaticSegment("staticJobs", "jobs", "jobs"),
		resourceids.UserSpecifiedSegment("jobName", "jobValue"),
	}
}

// String returns a human-readable description of this Job ID
func (id JobId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Media Service Name: %q", id.MediaServiceName),
		fmt.Sprintf("Transform Name: %q", id.TransformName),
		fmt.Sprintf("Job Name: %q", id.JobName),
	}
	return fmt.Sprintf("Job (%s)", strings.Join(components, "\n"))
}
