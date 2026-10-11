package certificates

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &IssuernameId{}

// IssuernameId is a struct representing the Resource ID for a Issuername
type IssuernameId struct {
	BaseURI    string
	IssuerName string
}

// NewIssuernameID returns a new IssuernameId struct
func NewIssuernameID(baseURI string, issuerName string) IssuernameId {
	return IssuernameId{
		BaseURI:    strings.TrimSuffix(baseURI, "/"),
		IssuerName: issuerName,
	}
}

// ParseIssuernameID parses 'input' into a IssuernameId
func ParseIssuernameID(input string) (*IssuernameId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IssuernameId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IssuernameId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseIssuernameIDInsensitively parses 'input' case-insensitively into a IssuernameId
// note: this method should only be used for API response data and not user input
func ParseIssuernameIDInsensitively(input string) (*IssuernameId, error) {
	parser := resourceids.NewParserFromResourceIdType(&IssuernameId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := IssuernameId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *IssuernameId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.IssuerName, ok = input.Parsed["issuerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "issuerName", input)
	}

	return nil
}

// ValidateIssuernameID checks that 'input' can be parsed as a Issuername ID
func ValidateIssuernameID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseIssuernameID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Issuername ID
func (id IssuernameId) ID() string {
	fmtString := "%s/certificates/issuers/%s"
	return fmt.Sprintf(fmtString, strings.TrimSuffix(id.BaseURI, "/"), id.IssuerName)
}

// Path returns the formatted Issuername ID without the BaseURI
func (id IssuernameId) Path() string {
	fmtString := "/certificates/issuers/%s"
	return fmt.Sprintf(fmtString, id.IssuerName)
}

// PathElements returns the values of Issuername ID Segments without the BaseURI
func (id IssuernameId) PathElements() []any {
	return []any{id.IssuerName}
}

// Segments returns a slice of Resource ID Segments which comprise this Issuername ID
func (id IssuernameId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint-url.example.com"),
		resourceids.StaticSegment("staticCertificates", "certificates", "certificates"),
		resourceids.StaticSegment("staticIssuers", "issuers", "issuers"),
		resourceids.UserSpecifiedSegment("issuerName", "issuerName"),
	}
}

// String returns a human-readable description of this Issuername ID
func (id IssuernameId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Issuer Name: %q", id.IssuerName),
	}
	return fmt.Sprintf("Issuername (%s)", strings.Join(components, "\n"))
}
