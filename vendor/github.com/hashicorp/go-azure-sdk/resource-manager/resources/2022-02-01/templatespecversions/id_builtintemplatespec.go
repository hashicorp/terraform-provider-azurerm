package templatespecversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&BuiltInTemplateSpecId{})
}

var _ resourceids.ResourceId = &BuiltInTemplateSpecId{}

// BuiltInTemplateSpecId is a struct representing the Resource ID for a Built In Template Spec
type BuiltInTemplateSpecId struct {
	BuiltInTemplateSpecName string
}

// NewBuiltInTemplateSpecID returns a new BuiltInTemplateSpecId struct
func NewBuiltInTemplateSpecID(builtInTemplateSpecName string) BuiltInTemplateSpecId {
	return BuiltInTemplateSpecId{
		BuiltInTemplateSpecName: builtInTemplateSpecName,
	}
}

// ParseBuiltInTemplateSpecID parses 'input' into a BuiltInTemplateSpecId
func ParseBuiltInTemplateSpecID(input string) (*BuiltInTemplateSpecId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BuiltInTemplateSpecId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BuiltInTemplateSpecId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseBuiltInTemplateSpecIDInsensitively parses 'input' case-insensitively into a BuiltInTemplateSpecId
// note: this method should only be used for API response data and not user input
func ParseBuiltInTemplateSpecIDInsensitively(input string) (*BuiltInTemplateSpecId, error) {
	parser := resourceids.NewParserFromResourceIdType(&BuiltInTemplateSpecId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := BuiltInTemplateSpecId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *BuiltInTemplateSpecId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BuiltInTemplateSpecName, ok = input.Parsed["builtInTemplateSpecName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "builtInTemplateSpecName", input)
	}

	return nil
}

// ValidateBuiltInTemplateSpecID checks that 'input' can be parsed as a Built In Template Spec ID
func ValidateBuiltInTemplateSpecID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBuiltInTemplateSpecID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Built In Template Spec ID
func (id BuiltInTemplateSpecId) ID() string {
	fmtString := "/providers/Microsoft.Resources/builtInTemplateSpecs/%s"
	return fmt.Sprintf(fmtString, id.BuiltInTemplateSpecName)
}

// Segments returns a slice of Resource ID Segments which comprise this Built In Template Spec ID
func (id BuiltInTemplateSpecId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftResources", "Microsoft.Resources", "Microsoft.Resources"),
		resourceids.StaticSegment("staticBuiltInTemplateSpecs", "builtInTemplateSpecs", "builtInTemplateSpecs"),
		resourceids.UserSpecifiedSegment("builtInTemplateSpecName", "builtInTemplateSpecName"),
	}
}

// String returns a human-readable description of this Built In Template Spec ID
func (id BuiltInTemplateSpecId) String() string {
	components := []string{
		fmt.Sprintf("Built In Template Spec Name: %q", id.BuiltInTemplateSpecName),
	}
	return fmt.Sprintf("Built In Template Spec (%s)", strings.Join(components, "\n"))
}
