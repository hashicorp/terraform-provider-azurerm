package storage

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &StorageId{}

// StorageId is a struct representing the Resource ID for a Storage
type StorageId struct {
	BaseURI     string
	StorageName string
}

// NewStorageID returns a new StorageId struct
func NewStorageID(baseURI string, storageName string) StorageId {
	return StorageId{
		BaseURI:     strings.TrimSuffix(baseURI, "/"),
		StorageName: storageName,
	}
}

// ParseStorageID parses 'input' into a StorageId
func ParseStorageID(input string) (*StorageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStorageIDInsensitively parses 'input' case-insensitively into a StorageId
// note: this method should only be used for API response data and not user input
func ParseStorageIDInsensitively(input string) (*StorageId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StorageId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.StorageName, ok = input.Parsed["storageName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageName", input)
	}

	return nil
}

// ValidateStorageID checks that 'input' can be parsed as a Storage ID
func ValidateStorageID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStorageID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Storage ID
func (id StorageId) ID() string {
	fmtString := "%s/storage/%s"
	return fmt.Sprintf(fmtString, id.BaseURI, id.StorageName)
}

// Path returns the formatted Storage ID without the BaseURI
func (id StorageId) Path() string {
	fmtString := "/storage/%s"
	return fmt.Sprintf(fmtString, id.StorageName)
}

// PathElements returns the values of Storage ID Segments without the BaseURI
func (id StorageId) PathElements() []any {
	return []any{id.StorageName}
}

// Segments returns a slice of Resource ID Segments which comprise this Storage ID
func (id StorageId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint_url"),
		resourceids.StaticSegment("staticStorage", "storage", "storage"),
		resourceids.UserSpecifiedSegment("storageName", "storageName"),
	}
}

// String returns a human-readable description of this Storage ID
func (id StorageId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Storage Name: %q", id.StorageName),
	}
	return fmt.Sprintf("Storage (%s)", strings.Join(components, "\n"))
}
