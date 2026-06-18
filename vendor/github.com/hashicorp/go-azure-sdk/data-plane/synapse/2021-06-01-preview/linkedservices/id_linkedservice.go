package linkedservices

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &LinkedServiceId{}

// LinkedServiceId is a struct representing the Resource ID for a Linked Service
type LinkedServiceId struct {
	BaseURI           string
	LinkedServiceName string
}

// NewLinkedServiceID returns a new LinkedServiceId struct
func NewLinkedServiceID(baseURI string, linkedServiceName string) LinkedServiceId {
	return LinkedServiceId{
		BaseURI:           strings.TrimSuffix(baseURI, "/"),
		LinkedServiceName: linkedServiceName,
	}
}

// ParseLinkedServiceID parses 'input' into a LinkedServiceId
func ParseLinkedServiceID(input string) (*LinkedServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LinkedServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LinkedServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLinkedServiceIDInsensitively parses 'input' case-insensitively into a LinkedServiceId
// note: this method should only be used for API response data and not user input
func ParseLinkedServiceIDInsensitively(input string) (*LinkedServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LinkedServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LinkedServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LinkedServiceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.LinkedServiceName, ok = input.Parsed["linkedServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "linkedServiceName", input)
	}

	return nil
}

// ValidateLinkedServiceID checks that 'input' can be parsed as a Linked Service ID
func ValidateLinkedServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLinkedServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Linked Service ID
func (id LinkedServiceId) ID() string {
	fmtString := "%s/linkedServices/%s"
	return fmt.Sprintf(fmtString, strings.TrimSuffix(id.BaseURI, "/"), id.LinkedServiceName)
}

// Path returns the formatted Linked Service ID without the BaseURI
func (id LinkedServiceId) Path() string {
	fmtString := "/linkedServices/%s"
	return fmt.Sprintf(fmtString, id.LinkedServiceName)
}

// PathElements returns the values of Linked Service ID Segments without the BaseURI
func (id LinkedServiceId) PathElements() []any {
	return []any{id.LinkedServiceName}
}

// Segments returns a slice of Resource ID Segments which comprise this Linked Service ID
func (id LinkedServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint-url.example.com"),
		resourceids.StaticSegment("staticLinkedServices", "linkedServices", "linkedServices"),
		resourceids.UserSpecifiedSegment("linkedServiceName", "linkedServiceName"),
	}
}

// String returns a human-readable description of this Linked Service ID
func (id LinkedServiceId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Linked Service Name: %q", id.LinkedServiceName),
	}
	return fmt.Sprintf("Linked Service (%s)", strings.Join(components, "\n"))
}
