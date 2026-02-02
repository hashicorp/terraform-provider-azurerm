package backuprestores

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &RestoreId{}

// RestoreId is a struct representing the Resource ID for a Restore
type RestoreId struct {
	BaseURI string
	JobId   string
}

// NewRestoreID returns a new RestoreId struct
func NewRestoreID(baseURI string, jobId string) RestoreId {
	return RestoreId{
		BaseURI: strings.TrimSuffix(baseURI, "/"),
		JobId:   jobId,
	}
}

// ParseRestoreID parses 'input' into a RestoreId
func ParseRestoreID(input string) (*RestoreId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RestoreId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RestoreId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRestoreIDInsensitively parses 'input' case-insensitively into a RestoreId
// note: this method should only be used for API response data and not user input
func ParseRestoreIDInsensitively(input string) (*RestoreId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RestoreId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RestoreId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RestoreId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.JobId, ok = input.Parsed["jobId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "jobId", input)
	}

	return nil
}

// ValidateRestoreID checks that 'input' can be parsed as a Restore ID
func ValidateRestoreID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRestoreID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Restore ID
func (id RestoreId) ID() string {
	fmtString := "%s/restore/%s"
	return fmt.Sprintf(fmtString, strings.TrimSuffix(id.BaseURI, "/"), id.JobId)
}

// Path returns the formatted Restore ID without the BaseURI
func (id RestoreId) Path() string {
	fmtString := "/restore/%s"
	return fmt.Sprintf(fmtString, id.JobId)
}

// PathElements returns the values of Restore ID Segments without the BaseURI
func (id RestoreId) PathElements() []any {
	return []any{id.JobId}
}

// Segments returns a slice of Resource ID Segments which comprise this Restore ID
func (id RestoreId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint-url.example.com"),
		resourceids.StaticSegment("staticRestore", "restore", "restore"),
		resourceids.UserSpecifiedSegment("jobId", "jobId"),
	}
}

// String returns a human-readable description of this Restore ID
func (id RestoreId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Job: %q", id.JobId),
	}
	return fmt.Sprintf("Restore (%s)", strings.Join(components, "\n"))
}
