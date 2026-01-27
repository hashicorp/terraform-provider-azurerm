package certificates

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &CertificateId{}

// CertificateId is a struct representing the Resource ID for a Certificate
type CertificateId struct {
	BaseURI         string
	CertificateName string
}

// NewCertificateID returns a new CertificateId struct
func NewCertificateID(baseURI string, certificateName string) CertificateId {
	return CertificateId{
		BaseURI:         strings.TrimSuffix(baseURI, "/"),
		CertificateName: certificateName,
	}
}

// ParseCertificateID parses 'input' into a CertificateId
func ParseCertificateID(input string) (*CertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CertificateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CertificateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCertificateIDInsensitively parses 'input' case-insensitively into a CertificateId
// note: this method should only be used for API response data and not user input
func ParseCertificateIDInsensitively(input string) (*CertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CertificateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CertificateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CertificateId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.CertificateName, ok = input.Parsed["certificateName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "certificateName", input)
	}

	return nil
}

// ValidateCertificateID checks that 'input' can be parsed as a Certificate ID
func ValidateCertificateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCertificateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Certificate ID
func (id CertificateId) ID() string {
	fmtString := "%s/certificates/%s"
	return fmt.Sprintf(fmtString, id.BaseURI, id.CertificateName)
}

// Path returns the formatted Certificate ID without the BaseURI
func (id CertificateId) Path() string {
	fmtString := "/certificates/%s"
	return fmt.Sprintf(fmtString, id.CertificateName)
}

// PathElements returns the values of Certificate ID Segments without the BaseURI
func (id CertificateId) PathElements() []any {
	return []any{id.CertificateName}
}

// Segments returns a slice of Resource ID Segments which comprise this Certificate ID
func (id CertificateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint_url"),
		resourceids.StaticSegment("staticCertificates", "certificates", "certificates"),
		resourceids.UserSpecifiedSegment("certificateName", "certificateName"),
	}
}

// String returns a human-readable description of this Certificate ID
func (id CertificateId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Certificate Name: %q", id.CertificateName),
	}
	return fmt.Sprintf("Certificate (%s)", strings.Join(components, "\n"))
}
