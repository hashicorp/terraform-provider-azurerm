package certificateobjectglobalrulestack

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = CertificateId{}

// CertificateId is a struct representing the Resource ID for a Certificate
type CertificateId struct {
	GlobalRulestackName string
	CertificateName     string
}

// NewCertificateID returns a new CertificateId struct
func NewCertificateID(globalRulestackName string, certificateName string) CertificateId {
	return CertificateId{
		GlobalRulestackName: globalRulestackName,
		CertificateName:     certificateName,
	}
}

// ParseCertificateID parses 'input' into a CertificateId
func ParseCertificateID(input string) (*CertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(CertificateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CertificateId{}

	if id.GlobalRulestackName, ok = parsed.Parsed["globalRulestackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "globalRulestackName", *parsed)
	}

	if id.CertificateName, ok = parsed.Parsed["certificateName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "certificateName", *parsed)
	}

	return &id, nil
}

// ParseCertificateIDInsensitively parses 'input' case-insensitively into a CertificateId
// note: this method should only be used for API response data and not user input
func ParseCertificateIDInsensitively(input string) (*CertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(CertificateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CertificateId{}

	if id.GlobalRulestackName, ok = parsed.Parsed["globalRulestackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "globalRulestackName", *parsed)
	}

	if id.CertificateName, ok = parsed.Parsed["certificateName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "certificateName", *parsed)
	}

	return &id, nil
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
	fmtString := "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks/%s/certificates/%s"
	return fmt.Sprintf(fmtString, id.GlobalRulestackName, id.CertificateName)
}

// Segments returns a slice of Resource ID Segments which comprise this Certificate ID
func (id CertificateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticPaloAltoNetworksCloudngfw", "PaloAltoNetworks.Cloudngfw", "PaloAltoNetworks.Cloudngfw"),
		resourceids.StaticSegment("staticGlobalRulestacks", "globalRulestacks", "globalRulestacks"),
		resourceids.UserSpecifiedSegment("globalRulestackName", "globalRulestackValue"),
		resourceids.StaticSegment("staticCertificates", "certificates", "certificates"),
		resourceids.UserSpecifiedSegment("certificateName", "certificateValue"),
	}
}

// String returns a human-readable description of this Certificate ID
func (id CertificateId) String() string {
	components := []string{
		fmt.Sprintf("Global Rulestack Name: %q", id.GlobalRulestackName),
		fmt.Sprintf("Certificate Name: %q", id.CertificateName),
	}
	return fmt.Sprintf("Certificate (%s)", strings.Join(components, "\n"))
}
