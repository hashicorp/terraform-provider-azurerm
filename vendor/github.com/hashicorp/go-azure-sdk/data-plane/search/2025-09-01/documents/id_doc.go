package documents

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &DocId{}

// DocId is a struct representing the Resource ID for a Doc
type DocId struct {
	BaseURI string
	DocName string
}

// NewDocID returns a new DocId struct
func NewDocID(baseURI string, docName string) DocId {
	return DocId{
		BaseURI: strings.TrimSuffix(baseURI, "/"),
		DocName: docName,
	}
}

// ParseDocID parses 'input' into a DocId
func ParseDocID(input string) (*DocId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DocId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DocId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDocIDInsensitively parses 'input' case-insensitively into a DocId
// note: this method should only be used for API response data and not user input
func ParseDocIDInsensitively(input string) (*DocId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DocId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DocId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DocId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.DocName, ok = input.Parsed["docName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "docName", input)
	}

	return nil
}

// ValidateDocID checks that 'input' can be parsed as a Doc ID
func ValidateDocID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDocID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Doc ID
func (id DocId) ID() string {
	fmtString := "%s/docs/%s"
	return fmt.Sprintf(fmtString, strings.TrimSuffix(id.BaseURI, "/"), id.DocName)
}

// Path returns the formatted Doc ID without the BaseURI
func (id DocId) Path() string {
	fmtString := "/docs/%s"
	return fmt.Sprintf(fmtString, id.DocName)
}

// PathElements returns the values of Doc ID Segments without the BaseURI
func (id DocId) PathElements() []any {
	return []any{id.DocName}
}

// Segments returns a slice of Resource ID Segments which comprise this Doc ID
func (id DocId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint-url.example.com"),
		resourceids.StaticSegment("docs", "docs", "docs"),
		resourceids.UserSpecifiedSegment("docName", "docName"),
	}
}

// String returns a human-readable description of this Doc ID
func (id DocId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Doc Name: %q", id.DocName),
	}
	return fmt.Sprintf("Doc (%s)", strings.Join(components, "\n"))
}
