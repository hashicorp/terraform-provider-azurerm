package deletedstorage

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &DeletedstorageId{}

// DeletedstorageId is a struct representing the Resource ID for a Deletedstorage
type DeletedstorageId struct {
	BaseURI            string
	DeletedstorageName string
}

// NewDeletedstorageID returns a new DeletedstorageId struct
func NewDeletedstorageID(baseURI string, deletedstorageName string) DeletedstorageId {
	return DeletedstorageId{
		BaseURI:            strings.TrimSuffix(baseURI, "/"),
		DeletedstorageName: deletedstorageName,
	}
}

// ParseDeletedstorageID parses 'input' into a DeletedstorageId
func ParseDeletedstorageID(input string) (*DeletedstorageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeletedstorageId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeletedstorageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDeletedstorageIDInsensitively parses 'input' case-insensitively into a DeletedstorageId
// note: this method should only be used for API response data and not user input
func ParseDeletedstorageIDInsensitively(input string) (*DeletedstorageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeletedstorageId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeletedstorageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DeletedstorageId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.DeletedstorageName, ok = input.Parsed["deletedstorageName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "deletedstorageName", input)
	}

	return nil
}

// ValidateDeletedstorageID checks that 'input' can be parsed as a Deletedstorage ID
func ValidateDeletedstorageID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeletedstorageID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Deletedstorage ID
func (id DeletedstorageId) ID() string {
	fmtString := "%s/deletedstorage/%s"
	return fmt.Sprintf(fmtString, id.BaseURI, id.DeletedstorageName)
}

// Path returns the formatted Deletedstorage ID without the BaseURI
func (id DeletedstorageId) Path() string {
	fmtString := "/deletedstorage/%s"
	return fmt.Sprintf(fmtString, id.DeletedstorageName)
}

// PathElements returns the values of Deletedstorage ID Segments without the BaseURI
func (id DeletedstorageId) PathElements() []any {
	return []any{id.DeletedstorageName}
}

// Segments returns a slice of Resource ID Segments which comprise this Deletedstorage ID
func (id DeletedstorageId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint_url"),
		resourceids.StaticSegment("staticDeletedstorage", "deletedstorage", "deletedstorage"),
		resourceids.UserSpecifiedSegment("deletedstorageName", "deletedstorageName"),
	}
}

// String returns a human-readable description of this Deletedstorage ID
func (id DeletedstorageId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Deletedstorage Name: %q", id.DeletedstorageName),
	}
	return fmt.Sprintf("Deletedstorage (%s)", strings.Join(components, "\n"))
}
