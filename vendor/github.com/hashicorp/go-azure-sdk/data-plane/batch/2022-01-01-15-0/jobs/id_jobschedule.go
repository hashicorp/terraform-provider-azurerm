package jobs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &JobscheduleId{}

// JobscheduleId is a struct representing the Resource ID for a Jobschedule
type JobscheduleId struct {
	BaseURI       string
	JobScheduleId string
}

// NewJobscheduleID returns a new JobscheduleId struct
func NewJobscheduleID(baseURI string, jobScheduleId string) JobscheduleId {
	return JobscheduleId{
		BaseURI:       strings.TrimSuffix(baseURI, "/"),
		JobScheduleId: jobScheduleId,
	}
}

// ParseJobscheduleID parses 'input' into a JobscheduleId
func ParseJobscheduleID(input string) (*JobscheduleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&JobscheduleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := JobscheduleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseJobscheduleIDInsensitively parses 'input' case-insensitively into a JobscheduleId
// note: this method should only be used for API response data and not user input
func ParseJobscheduleIDInsensitively(input string) (*JobscheduleId, error) {
	parser := resourceids.NewParserFromResourceIdType(&JobscheduleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := JobscheduleId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *JobscheduleId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.JobScheduleId, ok = input.Parsed["jobScheduleId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "jobScheduleId", input)
	}

	return nil
}

// ValidateJobscheduleID checks that 'input' can be parsed as a Jobschedule ID
func ValidateJobscheduleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseJobscheduleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Jobschedule ID
func (id JobscheduleId) ID() string {
	fmtString := "%s/jobschedules/%s"
	return fmt.Sprintf(fmtString, id.BaseURI, id.JobScheduleId)
}

// Path returns the formatted Jobschedule ID without the BaseURI
func (id JobscheduleId) Path() string {
	fmtString := "/jobschedules/%s"
	return fmt.Sprintf(fmtString, id.JobScheduleId)
}

// PathElements returns the values of Jobschedule ID Segments without the BaseURI
func (id JobscheduleId) PathElements() []any {
	return []any{id.JobScheduleId}
}

// Segments returns a slice of Resource ID Segments which comprise this Jobschedule ID
func (id JobscheduleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint_url"),
		resourceids.StaticSegment("staticJobschedules", "jobschedules", "jobschedules"),
		resourceids.UserSpecifiedSegment("jobScheduleId", "jobScheduleId"),
	}
}

// String returns a human-readable description of this Jobschedule ID
func (id JobscheduleId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Job Schedule: %q", id.JobScheduleId),
	}
	return fmt.Sprintf("Jobschedule (%s)", strings.Join(components, "\n"))
}
