package deletedstorage

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &DeletedstorageSaId{}

// DeletedstorageSaId is a struct representing the Resource ID for a Deletedstorage Sa
type DeletedstorageSaId struct {
	BaseURI            string
	DeletedstorageName string
	SaName             string
}

// NewDeletedstorageSaID returns a new DeletedstorageSaId struct
func NewDeletedstorageSaID(baseURI string, deletedstorageName string, saName string) DeletedstorageSaId {
	return DeletedstorageSaId{
		BaseURI:            strings.TrimSuffix(baseURI, "/"),
		DeletedstorageName: deletedstorageName,
		SaName:             saName,
	}
}

// ParseDeletedstorageSaID parses 'input' into a DeletedstorageSaId
func ParseDeletedstorageSaID(input string) (*DeletedstorageSaId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeletedstorageSaId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeletedstorageSaId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDeletedstorageSaIDInsensitively parses 'input' case-insensitively into a DeletedstorageSaId
// note: this method should only be used for API response data and not user input
func ParseDeletedstorageSaIDInsensitively(input string) (*DeletedstorageSaId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeletedstorageSaId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeletedstorageSaId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DeletedstorageSaId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.DeletedstorageName, ok = input.Parsed["deletedstorageName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "deletedstorageName", input)
	}

	if id.SaName, ok = input.Parsed["saName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "saName", input)
	}

	return nil
}

// ValidateDeletedstorageSaID checks that 'input' can be parsed as a Deletedstorage Sa ID
func ValidateDeletedstorageSaID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeletedstorageSaID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Deletedstorage Sa ID
func (id DeletedstorageSaId) ID() string {
	fmtString := "%s/deletedstorage/%s/sas/%s"
	return fmt.Sprintf(fmtString, strings.TrimSuffix(id.BaseURI, "/"), id.DeletedstorageName, id.SaName)
}

// Path returns the formatted Deletedstorage Sa ID without the BaseURI
func (id DeletedstorageSaId) Path() string {
	fmtString := "/deletedstorage/%s/sas/%s"
	return fmt.Sprintf(fmtString, id.DeletedstorageName, id.SaName)
}

// PathElements returns the values of Deletedstorage Sa ID Segments without the BaseURI
func (id DeletedstorageSaId) PathElements() []any {
	return []any{id.DeletedstorageName, id.SaName}
}

// Segments returns a slice of Resource ID Segments which comprise this Deletedstorage Sa ID
func (id DeletedstorageSaId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint-url.example.com"),
		resourceids.StaticSegment("staticDeletedstorage", "deletedstorage", "deletedstorage"),
		resourceids.UserSpecifiedSegment("deletedstorageName", "deletedstorageName"),
		resourceids.StaticSegment("staticSas", "sas", "sas"),
		resourceids.UserSpecifiedSegment("saName", "saName"),
	}
}

// String returns a human-readable description of this Deletedstorage Sa ID
func (id DeletedstorageSaId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Deletedstorage Name: %q", id.DeletedstorageName),
		fmt.Sprintf("Sa Name: %q", id.SaName),
	}
	return fmt.Sprintf("Deletedstorage Sa (%s)", strings.Join(components, "\n"))
}
