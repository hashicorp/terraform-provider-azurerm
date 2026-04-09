package resourceproviders

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SourceControlId{})
}

var _ resourceids.ResourceId = &SourceControlId{}

// SourceControlId is a struct representing the Resource ID for a Source Control
type SourceControlId struct {
	SourceControlName string
}

// NewSourceControlID returns a new SourceControlId struct
func NewSourceControlID(sourceControlName string) SourceControlId {
	return SourceControlId{
		SourceControlName: sourceControlName,
	}
}

// ParseSourceControlID parses 'input' into a SourceControlId
func ParseSourceControlID(input string) (*SourceControlId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SourceControlId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SourceControlId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSourceControlIDInsensitively parses 'input' case-insensitively into a SourceControlId
// note: this method should only be used for API response data and not user input
func ParseSourceControlIDInsensitively(input string) (*SourceControlId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SourceControlId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SourceControlId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SourceControlId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SourceControlName, ok = input.Parsed["sourceControlName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "sourceControlName", input)
	}

	return nil
}

// ValidateSourceControlID checks that 'input' can be parsed as a Source Control ID
func ValidateSourceControlID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSourceControlID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Source Control ID
func (id SourceControlId) ID() string {
	fmtString := "/providers/Microsoft.Web/sourceControls/%s"
	return fmt.Sprintf(fmtString, id.SourceControlName)
}

// Segments returns a slice of Resource ID Segments which comprise this Source Control ID
func (id SourceControlId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSourceControls", "sourceControls", "sourceControls"),
		resourceids.UserSpecifiedSegment("sourceControlName", "sourceControlName"),
	}
}

// String returns a human-readable description of this Source Control ID
func (id SourceControlId) String() string {
	components := []string{
		fmt.Sprintf("Source Control Name: %q", id.SourceControlName),
	}
	return fmt.Sprintf("Source Control (%s)", strings.Join(components, "\n"))
}
