package certificateprofiles

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CertificateProfileId{})
}

var _ resourceids.ResourceId = &CertificateProfileId{}

// CertificateProfileId is a struct representing the Resource ID for a Certificate Profile
type CertificateProfileId struct {
	SubscriptionId         string
	ResourceGroupName      string
	CodeSigningAccountName string
	CertificateProfileName string
}

// NewCertificateProfileID returns a new CertificateProfileId struct
func NewCertificateProfileID(subscriptionId string, resourceGroupName string, codeSigningAccountName string, certificateProfileName string) CertificateProfileId {
	return CertificateProfileId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		CodeSigningAccountName: codeSigningAccountName,
		CertificateProfileName: certificateProfileName,
	}
}

// ParseCertificateProfileID parses 'input' into a CertificateProfileId
func ParseCertificateProfileID(input string) (*CertificateProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CertificateProfileId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CertificateProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCertificateProfileIDInsensitively parses 'input' case-insensitively into a CertificateProfileId
// note: this method should only be used for API response data and not user input
func ParseCertificateProfileIDInsensitively(input string) (*CertificateProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CertificateProfileId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CertificateProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CertificateProfileId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CodeSigningAccountName, ok = input.Parsed["codeSigningAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "codeSigningAccountName", input)
	}

	if id.CertificateProfileName, ok = input.Parsed["certificateProfileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "certificateProfileName", input)
	}

	return nil
}

// ValidateCertificateProfileID checks that 'input' can be parsed as a Certificate Profile ID
func ValidateCertificateProfileID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCertificateProfileID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Certificate Profile ID
func (id CertificateProfileId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CodeSigning/codeSigningAccounts/%s/certificateProfiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CodeSigningAccountName, id.CertificateProfileName)
}

// Segments returns a slice of Resource ID Segments which comprise this Certificate Profile ID
func (id CertificateProfileId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCodeSigning", "Microsoft.CodeSigning", "Microsoft.CodeSigning"),
		resourceids.StaticSegment("staticCodeSigningAccounts", "codeSigningAccounts", "codeSigningAccounts"),
		resourceids.UserSpecifiedSegment("codeSigningAccountName", "codeSigningAccountName"),
		resourceids.StaticSegment("staticCertificateProfiles", "certificateProfiles", "certificateProfiles"),
		resourceids.UserSpecifiedSegment("certificateProfileName", "certificateProfileName"),
	}
}

// String returns a human-readable description of this Certificate Profile ID
func (id CertificateProfileId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Code Signing Account Name: %q", id.CodeSigningAccountName),
		fmt.Sprintf("Certificate Profile Name: %q", id.CertificateProfileName),
	}
	return fmt.Sprintf("Certificate Profile (%s)", strings.Join(components, "\n"))
}
