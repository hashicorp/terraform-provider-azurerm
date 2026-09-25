package indexers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &IndexerId{}

// IndexerId is a struct representing the Resource ID for a Indexer
type IndexerId struct {
	BaseURI     string
	IndexerName string
}

// NewIndexerID returns a new IndexerId struct
func NewIndexerID(baseURI string, indexerName string) IndexerId {
	return IndexerId{
		BaseURI:     strings.TrimSuffix(baseURI, "/"),
		IndexerName: indexerName,
	}
}

// ParseIndexerID parses 'input' into a IndexerId
func ParseIndexerID(input string) (*IndexerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IndexerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IndexerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseIndexerIDInsensitively parses 'input' case-insensitively into a IndexerId
// note: this method should only be used for API response data and not user input
func ParseIndexerIDInsensitively(input string) (*IndexerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IndexerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IndexerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *IndexerId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.IndexerName, ok = input.Parsed["indexerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "indexerName", input)
	}

	return nil
}

// ValidateIndexerID checks that 'input' can be parsed as a Indexer ID
func ValidateIndexerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseIndexerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Indexer ID
func (id IndexerId) ID() string {
	fmtString := "%s/indexers/%s"
	return fmt.Sprintf(fmtString, strings.TrimSuffix(id.BaseURI, "/"), id.IndexerName)
}

// Path returns the formatted Indexer ID without the BaseURI
func (id IndexerId) Path() string {
	fmtString := "/indexers/%s"
	return fmt.Sprintf(fmtString, id.IndexerName)
}

// PathElements returns the values of Indexer ID Segments without the BaseURI
func (id IndexerId) PathElements() []any {
	return []any{id.IndexerName}
}

// Segments returns a slice of Resource ID Segments which comprise this Indexer ID
func (id IndexerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint-url.example.com"),
		resourceids.StaticSegment("indexers", "indexers", "indexers"),
		resourceids.UserSpecifiedSegment("indexerName", "indexerName"),
	}
}

// String returns a human-readable description of this Indexer ID
func (id IndexerId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Indexer Name: %q", id.IndexerName),
	}
	return fmt.Sprintf("Indexer (%s)", strings.Join(components, "\n"))
}
