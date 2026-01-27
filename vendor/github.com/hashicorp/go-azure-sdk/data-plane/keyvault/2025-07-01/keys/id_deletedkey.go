package keys

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &DeletedkeyId{}

// DeletedkeyId is a struct representing the Resource ID for a Deletedkey
type DeletedkeyId struct {
	BaseURI        string
	DeletedkeyName string
}

// NewDeletedkeyID returns a new DeletedkeyId struct
func NewDeletedkeyID(baseURI string, deletedkeyName string) DeletedkeyId {
	return DeletedkeyId{
		BaseURI:        strings.TrimSuffix(baseURI, "/"),
		DeletedkeyName: deletedkeyName,
	}
}

// ParseDeletedkeyID parses 'input' into a DeletedkeyId
func ParseDeletedkeyID(input string) (*DeletedkeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeletedkeyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeletedkeyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDeletedkeyIDInsensitively parses 'input' case-insensitively into a DeletedkeyId
// note: this method should only be used for API response data and not user input
func ParseDeletedkeyIDInsensitively(input string) (*DeletedkeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeletedkeyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeletedkeyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DeletedkeyId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.DeletedkeyName, ok = input.Parsed["deletedkeyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "deletedkeyName", input)
	}

	return nil
}

// ValidateDeletedkeyID checks that 'input' can be parsed as a Deletedkey ID
func ValidateDeletedkeyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeletedkeyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Deletedkey ID
func (id DeletedkeyId) ID() string {
	fmtString := "%s/deletedkeys/%s"
	return fmt.Sprintf(fmtString, id.BaseURI, id.DeletedkeyName)
}

// Path returns the formatted Deletedkey ID without the BaseURI
func (id DeletedkeyId) Path() string {
	fmtString := "/deletedkeys/%s"
	return fmt.Sprintf(fmtString, id.DeletedkeyName)
}

// PathElements returns the values of Deletedkey ID Segments without the BaseURI
func (id DeletedkeyId) PathElements() []any {
	return []any{id.DeletedkeyName}
}

// Segments returns a slice of Resource ID Segments which comprise this Deletedkey ID
func (id DeletedkeyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint_url"),
		resourceids.StaticSegment("staticDeletedkeys", "deletedkeys", "deletedkeys"),
		resourceids.UserSpecifiedSegment("deletedkeyName", "deletedkeyName"),
	}
}

// String returns a human-readable description of this Deletedkey ID
func (id DeletedkeyId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Deletedkey Name: %q", id.DeletedkeyName),
	}
	return fmt.Sprintf("Deletedkey (%s)", strings.Join(components, "\n"))
}
