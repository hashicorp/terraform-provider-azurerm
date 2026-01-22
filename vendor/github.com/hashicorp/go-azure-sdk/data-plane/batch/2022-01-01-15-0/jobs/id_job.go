package jobs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &JobId{}

// JobId is a struct representing the Resource ID for a Job
type JobId struct {
	BaseURI string
	JobId   string
}

// NewJobID returns a new JobId struct
func NewJobID(baseURI string, jobId string) JobId {
	return JobId{
		BaseURI: baseURI,
		JobId:   jobId,
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

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.JobId, ok = input.Parsed["jobId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "jobId", input)
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
	fmtString := "%s/jobs/%s"
	return fmt.Sprintf(fmtString, id.BaseURI, id.JobId)
}

// Path returns the formatted Job ID without the Scope / BaseURI
func (id JobId) Path() string {
	fmtString := "/jobs/%s"
	return fmt.Sprintf(fmtString, id.JobId)
}

// Segments returns a slice of Resource ID Segments which comprise this Job ID
func (id JobId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint_url"),
		resourceids.StaticSegment("staticJobs", "jobs", "jobs"),
		resourceids.UserSpecifiedSegment("jobId", "jobId"),
	}
}

// String returns a human-readable description of this Job ID
func (id JobId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Job: %q", id.JobId),
	}
	return fmt.Sprintf("Job (%s)", strings.Join(components, "\n"))
}
