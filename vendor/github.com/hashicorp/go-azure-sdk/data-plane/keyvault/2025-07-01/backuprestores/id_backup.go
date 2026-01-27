package backuprestores

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &BackupId{}

// BackupId is a struct representing the Resource ID for a Backup
type BackupId struct {
	BaseURI string
	JobId   string
}

// NewBackupID returns a new BackupId struct
func NewBackupID(baseURI string, jobId string) BackupId {
	return BackupId{
		BaseURI: strings.TrimSuffix(baseURI, "/"),
		JobId:   jobId,
	}
}

// ParseBackupID parses 'input' into a BackupId
func ParseBackupID(input string) (*BackupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BackupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBackupIDInsensitively parses 'input' case-insensitively into a BackupId
// note: this method should only be used for API response data and not user input
func ParseBackupIDInsensitively(input string) (*BackupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BackupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BackupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BackupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.JobId, ok = input.Parsed["jobId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "jobId", input)
	}

	return nil
}

// ValidateBackupID checks that 'input' can be parsed as a Backup ID
func ValidateBackupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBackupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Backup ID
func (id BackupId) ID() string {
	fmtString := "%s/backup/%s"
	return fmt.Sprintf(fmtString, id.BaseURI, id.JobId)
}

// Path returns the formatted Backup ID without the BaseURI
func (id BackupId) Path() string {
	fmtString := "/backup/%s"
	return fmt.Sprintf(fmtString, id.JobId)
}

// PathElements returns the values of Backup ID Segments without the BaseURI
func (id BackupId) PathElements() []any {
	return []any{id.JobId}
}

// Segments returns a slice of Resource ID Segments which comprise this Backup ID
func (id BackupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint_url"),
		resourceids.StaticSegment("staticBackup", "backup", "backup"),
		resourceids.UserSpecifiedSegment("jobId", "jobId"),
	}
}

// String returns a human-readable description of this Backup ID
func (id BackupId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Job: %q", id.JobId),
	}
	return fmt.Sprintf("Backup (%s)", strings.Join(components, "\n"))
}
