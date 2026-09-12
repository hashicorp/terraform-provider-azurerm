package synonymmaps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &SynonymmapId{}

// SynonymmapId is a struct representing the Resource ID for a Synonymmap
type SynonymmapId struct {
	BaseURI        string
	SynonymmapName string
}

// NewSynonymmapID returns a new SynonymmapId struct
func NewSynonymmapID(baseURI string, synonymmapName string) SynonymmapId {
	return SynonymmapId{
		BaseURI:        strings.TrimSuffix(baseURI, "/"),
		SynonymmapName: synonymmapName,
	}
}

// ParseSynonymmapID parses 'input' into a SynonymmapId
func ParseSynonymmapID(input string) (*SynonymmapId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SynonymmapId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SynonymmapId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSynonymmapIDInsensitively parses 'input' case-insensitively into a SynonymmapId
// note: this method should only be used for API response data and not user input
func ParseSynonymmapIDInsensitively(input string) (*SynonymmapId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SynonymmapId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SynonymmapId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SynonymmapId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.SynonymmapName, ok = input.Parsed["synonymmapName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "synonymmapName", input)
	}

	return nil
}

// ValidateSynonymmapID checks that 'input' can be parsed as a Synonymmap ID
func ValidateSynonymmapID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSynonymmapID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Synonymmap ID
func (id SynonymmapId) ID() string {
	fmtString := "%s/synonymmaps/%s"
	return fmt.Sprintf(fmtString, strings.TrimSuffix(id.BaseURI, "/"), id.SynonymmapName)
}

// Path returns the formatted Synonymmap ID without the BaseURI
func (id SynonymmapId) Path() string {
	fmtString := "/synonymmaps/%s"
	return fmt.Sprintf(fmtString, id.SynonymmapName)
}

// PathElements returns the values of Synonymmap ID Segments without the BaseURI
func (id SynonymmapId) PathElements() []any {
	return []any{id.SynonymmapName}
}

// Segments returns a slice of Resource ID Segments which comprise this Synonymmap ID
func (id SynonymmapId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint-url.example.com"),
		resourceids.StaticSegment("synonymmaps", "synonymmaps", "synonymmaps"),
		resourceids.UserSpecifiedSegment("synonymmapName", "synonymmapName"),
	}
}

// String returns a human-readable description of this Synonymmap ID
func (id SynonymmapId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Synonymmap Name: %q", id.SynonymmapName),
	}
	return fmt.Sprintf("Synonymmap (%s)", strings.Join(components, "\n"))
}
