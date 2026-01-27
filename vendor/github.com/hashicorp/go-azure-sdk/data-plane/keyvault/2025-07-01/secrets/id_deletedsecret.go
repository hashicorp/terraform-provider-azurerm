package secrets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &DeletedsecretId{}

// DeletedsecretId is a struct representing the Resource ID for a Deletedsecret
type DeletedsecretId struct {
	BaseURI           string
	DeletedsecretName string
}

// NewDeletedsecretID returns a new DeletedsecretId struct
func NewDeletedsecretID(baseURI string, deletedsecretName string) DeletedsecretId {
	return DeletedsecretId{
		BaseURI:           strings.TrimSuffix(baseURI, "/"),
		DeletedsecretName: deletedsecretName,
	}
}

// ParseDeletedsecretID parses 'input' into a DeletedsecretId
func ParseDeletedsecretID(input string) (*DeletedsecretId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeletedsecretId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeletedsecretId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDeletedsecretIDInsensitively parses 'input' case-insensitively into a DeletedsecretId
// note: this method should only be used for API response data and not user input
func ParseDeletedsecretIDInsensitively(input string) (*DeletedsecretId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeletedsecretId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeletedsecretId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DeletedsecretId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.DeletedsecretName, ok = input.Parsed["deletedsecretName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "deletedsecretName", input)
	}

	return nil
}

// ValidateDeletedsecretID checks that 'input' can be parsed as a Deletedsecret ID
func ValidateDeletedsecretID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeletedsecretID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Deletedsecret ID
func (id DeletedsecretId) ID() string {
	fmtString := "%s/deletedsecrets/%s"
	return fmt.Sprintf(fmtString, id.BaseURI, id.DeletedsecretName)
}

// Path returns the formatted Deletedsecret ID without the BaseURI
func (id DeletedsecretId) Path() string {
	fmtString := "/deletedsecrets/%s"
	return fmt.Sprintf(fmtString, id.DeletedsecretName)
}

// PathElements returns the values of Deletedsecret ID Segments without the BaseURI
func (id DeletedsecretId) PathElements() []any {
	return []any{id.DeletedsecretName}
}

// Segments returns a slice of Resource ID Segments which comprise this Deletedsecret ID
func (id DeletedsecretId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint_url"),
		resourceids.StaticSegment("staticDeletedsecrets", "deletedsecrets", "deletedsecrets"),
		resourceids.UserSpecifiedSegment("deletedsecretName", "deletedsecretName"),
	}
}

// String returns a human-readable description of this Deletedsecret ID
func (id DeletedsecretId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Deletedsecret Name: %q", id.DeletedsecretName),
	}
	return fmt.Sprintf("Deletedsecret (%s)", strings.Join(components, "\n"))
}
