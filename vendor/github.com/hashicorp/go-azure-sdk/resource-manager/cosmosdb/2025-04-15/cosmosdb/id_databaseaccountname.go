package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DatabaseAccountNameId{})
}

var _ resourceids.ResourceId = &DatabaseAccountNameId{}

// DatabaseAccountNameId is a struct representing the Resource ID for a Database Account Name
type DatabaseAccountNameId struct {
	DatabaseAccountName string
}

// NewDatabaseAccountNameID returns a new DatabaseAccountNameId struct
func NewDatabaseAccountNameID(databaseAccountName string) DatabaseAccountNameId {
	return DatabaseAccountNameId{
		DatabaseAccountName: databaseAccountName,
	}
}

// ParseDatabaseAccountNameID parses 'input' into a DatabaseAccountNameId
func ParseDatabaseAccountNameID(input string) (*DatabaseAccountNameId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatabaseAccountNameId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatabaseAccountNameId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDatabaseAccountNameIDInsensitively parses 'input' case-insensitively into a DatabaseAccountNameId
// note: this method should only be used for API response data and not user input
func ParseDatabaseAccountNameIDInsensitively(input string) (*DatabaseAccountNameId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatabaseAccountNameId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatabaseAccountNameId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DatabaseAccountNameId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.DatabaseAccountName, ok = input.Parsed["databaseAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", input)
	}

	return nil
}

// ValidateDatabaseAccountNameID checks that 'input' can be parsed as a Database Account Name ID
func ValidateDatabaseAccountNameID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDatabaseAccountNameID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Database Account Name ID
func (id DatabaseAccountNameId) ID() string {
	fmtString := "/providers/Microsoft.DocumentDB/databaseAccountNames/%s"
	return fmt.Sprintf(fmtString, id.DatabaseAccountName)
}

// Segments returns a slice of Resource ID Segments which comprise this Database Account Name ID
func (id DatabaseAccountNameId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccountNames", "databaseAccountNames", "databaseAccountNames"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "accountName"),
	}
}

// String returns a human-readable description of this Database Account Name ID
func (id DatabaseAccountNameId) String() string {
	components := []string{
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
	}
	return fmt.Sprintf("Database Account Name (%s)", strings.Join(components, "\n"))
}
