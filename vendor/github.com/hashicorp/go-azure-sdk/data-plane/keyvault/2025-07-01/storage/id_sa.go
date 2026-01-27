package storage

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &SaId{}

// SaId is a struct representing the Resource ID for a Sa
type SaId struct {
	BaseURI     string
	StorageName string
	SaName      string
}

// NewSaID returns a new SaId struct
func NewSaID(baseURI string, storageName string, saName string) SaId {
	return SaId{
		BaseURI:     strings.TrimSuffix(baseURI, "/"),
		StorageName: storageName,
		SaName:      saName,
	}
}

// ParseSaID parses 'input' into a SaId
func ParseSaID(input string) (*SaId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SaId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SaId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSaIDInsensitively parses 'input' case-insensitively into a SaId
// note: this method should only be used for API response data and not user input
func ParseSaIDInsensitively(input string) (*SaId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SaId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SaId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SaId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.StorageName, ok = input.Parsed["storageName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageName", input)
	}

	if id.SaName, ok = input.Parsed["saName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "saName", input)
	}

	return nil
}

// ValidateSaID checks that 'input' can be parsed as a Sa ID
func ValidateSaID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSaID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Sa ID
func (id SaId) ID() string {
	fmtString := "%s/storage/%s/sas/%s"
	return fmt.Sprintf(fmtString, id.BaseURI, id.StorageName, id.SaName)
}

// Path returns the formatted Sa ID without the BaseURI
func (id SaId) Path() string {
	fmtString := "/storage/%s/sas/%s"
	return fmt.Sprintf(fmtString, id.StorageName, id.SaName)
}

// PathElements returns the values of Sa ID Segments without the BaseURI
func (id SaId) PathElements() []any {
	return []any{id.StorageName, id.SaName}
}

// Segments returns a slice of Resource ID Segments which comprise this Sa ID
func (id SaId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint_url"),
		resourceids.StaticSegment("staticStorage", "storage", "storage"),
		resourceids.UserSpecifiedSegment("storageName", "storageName"),
		resourceids.StaticSegment("staticSas", "sas", "sas"),
		resourceids.UserSpecifiedSegment("saName", "saName"),
	}
}

// String returns a human-readable description of this Sa ID
func (id SaId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Storage Name: %q", id.StorageName),
		fmt.Sprintf("Sa Name: %q", id.SaName),
	}
	return fmt.Sprintf("Sa (%s)", strings.Join(components, "\n"))
}
