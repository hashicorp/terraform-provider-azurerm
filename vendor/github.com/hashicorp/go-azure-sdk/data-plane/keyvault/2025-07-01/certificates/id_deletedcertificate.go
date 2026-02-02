package certificates

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &DeletedcertificateId{}

// DeletedcertificateId is a struct representing the Resource ID for a Deletedcertificate
type DeletedcertificateId struct {
	BaseURI                string
	DeletedcertificateName string
}

// NewDeletedcertificateID returns a new DeletedcertificateId struct
func NewDeletedcertificateID(baseURI string, deletedcertificateName string) DeletedcertificateId {
	return DeletedcertificateId{
		BaseURI:                strings.TrimSuffix(baseURI, "/"),
		DeletedcertificateName: deletedcertificateName,
	}
}

// ParseDeletedcertificateID parses 'input' into a DeletedcertificateId
func ParseDeletedcertificateID(input string) (*DeletedcertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeletedcertificateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeletedcertificateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDeletedcertificateIDInsensitively parses 'input' case-insensitively into a DeletedcertificateId
// note: this method should only be used for API response data and not user input
func ParseDeletedcertificateIDInsensitively(input string) (*DeletedcertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DeletedcertificateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DeletedcertificateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DeletedcertificateId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.DeletedcertificateName, ok = input.Parsed["deletedcertificateName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "deletedcertificateName", input)
	}

	return nil
}

// ValidateDeletedcertificateID checks that 'input' can be parsed as a Deletedcertificate ID
func ValidateDeletedcertificateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDeletedcertificateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Deletedcertificate ID
func (id DeletedcertificateId) ID() string {
	fmtString := "%s/deletedcertificates/%s"
	return fmt.Sprintf(fmtString, strings.TrimSuffix(id.BaseURI, "/"), id.DeletedcertificateName)
}

// Path returns the formatted Deletedcertificate ID without the BaseURI
func (id DeletedcertificateId) Path() string {
	fmtString := "/deletedcertificates/%s"
	return fmt.Sprintf(fmtString, id.DeletedcertificateName)
}

// PathElements returns the values of Deletedcertificate ID Segments without the BaseURI
func (id DeletedcertificateId) PathElements() []any {
	return []any{id.DeletedcertificateName}
}

// Segments returns a slice of Resource ID Segments which comprise this Deletedcertificate ID
func (id DeletedcertificateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint-url.example.com"),
		resourceids.StaticSegment("staticDeletedcertificates", "deletedcertificates", "deletedcertificates"),
		resourceids.UserSpecifiedSegment("deletedcertificateName", "deletedcertificateName"),
	}
}

// String returns a human-readable description of this Deletedcertificate ID
func (id DeletedcertificateId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Deletedcertificate Name: %q", id.DeletedcertificateName),
	}
	return fmt.Sprintf("Deletedcertificate (%s)", strings.Join(components, "\n"))
}
