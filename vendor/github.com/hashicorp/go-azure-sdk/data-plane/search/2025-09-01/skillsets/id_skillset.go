package skillsets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &SkillsetId{}

// SkillsetId is a struct representing the Resource ID for a Skillset
type SkillsetId struct {
	BaseURI      string
	SkillsetName string
}

// NewSkillsetID returns a new SkillsetId struct
func NewSkillsetID(baseURI string, skillsetName string) SkillsetId {
	return SkillsetId{
		BaseURI:      strings.TrimSuffix(baseURI, "/"),
		SkillsetName: skillsetName,
	}
}

// ParseSkillsetID parses 'input' into a SkillsetId
func ParseSkillsetID(input string) (*SkillsetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SkillsetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SkillsetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSkillsetIDInsensitively parses 'input' case-insensitively into a SkillsetId
// note: this method should only be used for API response data and not user input
func ParseSkillsetIDInsensitively(input string) (*SkillsetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SkillsetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SkillsetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SkillsetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.SkillsetName, ok = input.Parsed["skillsetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "skillsetName", input)
	}

	return nil
}

// ValidateSkillsetID checks that 'input' can be parsed as a Skillset ID
func ValidateSkillsetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSkillsetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Skillset ID
func (id SkillsetId) ID() string {
	fmtString := "%s/skillsets/%s"
	return fmt.Sprintf(fmtString, strings.TrimSuffix(id.BaseURI, "/"), id.SkillsetName)
}

// Path returns the formatted Skillset ID without the BaseURI
func (id SkillsetId) Path() string {
	fmtString := "/skillsets/%s"
	return fmt.Sprintf(fmtString, id.SkillsetName)
}

// PathElements returns the values of Skillset ID Segments without the BaseURI
func (id SkillsetId) PathElements() []any {
	return []any{id.SkillsetName}
}

// Segments returns a slice of Resource ID Segments which comprise this Skillset ID
func (id SkillsetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint-url.example.com"),
		resourceids.StaticSegment("skillsets", "skillsets", "skillsets"),
		resourceids.UserSpecifiedSegment("skillsetName", "skillsetName"),
	}
}

// String returns a human-readable description of this Skillset ID
func (id SkillsetId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Skillset Name: %q", id.SkillsetName),
	}
	return fmt.Sprintf("Skillset (%s)", strings.Join(components, "\n"))
}
