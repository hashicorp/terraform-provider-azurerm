package certificates

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &CertificateversionId{}

// CertificateversionId is a struct representing the Resource ID for a Certificateversion
type CertificateversionId struct {
	BaseURI            string
	CertificateName    string
	Certificateversion string
}

// NewCertificateversionID returns a new CertificateversionId struct
func NewCertificateversionID(baseURI string, certificateName string, certificateversion string) CertificateversionId {
	return CertificateversionId{
		BaseURI:            strings.TrimSuffix(baseURI, "/"),
		CertificateName:    certificateName,
		Certificateversion: certificateversion,
	}
}

// ParseCertificateversionID parses 'input' into a CertificateversionId
func ParseCertificateversionID(input string) (*CertificateversionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CertificateversionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CertificateversionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCertificateversionIDInsensitively parses 'input' case-insensitively into a CertificateversionId
// note: this method should only be used for API response data and not user input
func ParseCertificateversionIDInsensitively(input string) (*CertificateversionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CertificateversionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CertificateversionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CertificateversionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.BaseURI, ok = input.Parsed["baseURI"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "baseURI", input)
	}

	if id.CertificateName, ok = input.Parsed["certificateName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "certificateName", input)
	}

	if id.Certificateversion, ok = input.Parsed["certificateversion"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "certificateversion", input)
	}

	return nil
}

// ValidateCertificateversionID checks that 'input' can be parsed as a Certificateversion ID
func ValidateCertificateversionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCertificateversionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Certificateversion ID
func (id CertificateversionId) ID() string {
	fmtString := "%s/certificates/%s/%s"
	return fmt.Sprintf(fmtString, id.BaseURI, id.CertificateName, id.Certificateversion)
}

// Path returns the formatted Certificateversion ID without the BaseURI
func (id CertificateversionId) Path() string {
	fmtString := "/certificates/%s/%s"
	return fmt.Sprintf(fmtString, id.CertificateName, id.Certificateversion)
}

// PathElements returns the values of Certificateversion ID Segments without the BaseURI
func (id CertificateversionId) PathElements() []any {
	return []any{id.CertificateName, id.Certificateversion}
}

// Segments returns a slice of Resource ID Segments which comprise this Certificateversion ID
func (id CertificateversionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.DataPlaneBaseURISegment("baseURI", "https://endpoint_url"),
		resourceids.StaticSegment("staticCertificates", "certificates", "certificates"),
		resourceids.UserSpecifiedSegment("certificateName", "certificateName"),
		resourceids.UserSpecifiedSegment("certificateversion", "certificateversion"),
	}
}

// String returns a human-readable description of this Certificateversion ID
func (id CertificateversionId) String() string {
	components := []string{
		fmt.Sprintf("Base URI: %q", id.BaseURI),
		fmt.Sprintf("Certificate Name: %q", id.CertificateName),
		fmt.Sprintf("Certificateversion: %q", id.Certificateversion),
	}
	return fmt.Sprintf("Certificateversion (%s)", strings.Join(components, "\n"))
}
