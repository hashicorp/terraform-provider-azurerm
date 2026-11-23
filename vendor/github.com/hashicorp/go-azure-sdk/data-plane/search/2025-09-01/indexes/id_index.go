package indexes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &IndexId{}

// IndexId is a struct representing the Resource ID for a Index
type IndexId struct {
	BaseURI   string
	IndexName string
}

// NewIndexID returns a new IndexId struct
func NewIndexID(baseURI string, indexName string) IndexId {
	return IndexId{
		BaseURI:   strings.TrimSuffix(baseURI, "/"),
		IndexName: indexName,
	}
}

// ParseIndexID parses 'input' into a IndexId
func ParseIndexID(input string) (*IndexId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IndexId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IndexId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseIndexIDInsensitively parses 'input' case-insensitively into a IndexId
// note: this method should only be used for API response data and not user input
func ParseIndexIDInsensitively(input string) (*IndexId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IndexId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IndexId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *IndexId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.IndexName, ok = input.Parsed["indexName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "indexName", input)
	}

	return nil
}

// ValidateIndexID checks that 'input' can be parsed as a Index ID
func ValidateIndexID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseIndexID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Index ID
func (id IndexId) ID() string {
	fmtString := "%s/indexes/%s"
	return fmt.Sprintf(fmtString, strings.TrimSuffix(id.BaseURI, "/"), id.IndexName)
}

// Path returns the formatted Index ID without the BaseURI
func (id IndexId) Path() string {
	fmtString := "/indexes/%s"
	return fmt.Sprintf(fmtString, id.IndexName)
}

// PathElements returns the values of Index ID Segments without the BaseURI
func (id IndexId) PathElements() []any {
	return []any{id.IndexName}
}

// Segments returns a slice of Resource ID Segments which comprise this Index ID
func (id IndexId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint-url.example.com"),
		resourceids.StaticSegment("indexes", "indexes", "indexes"),
		resourceids.UserSpecifiedSegment("indexName", "indexName"),
	}
}

// String returns a human-readable description of this Index ID
func (id IndexId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Index Name: %q", id.IndexName),
	}
	return fmt.Sprintf("Index (%s)", strings.Join(components, "\n"))
}
