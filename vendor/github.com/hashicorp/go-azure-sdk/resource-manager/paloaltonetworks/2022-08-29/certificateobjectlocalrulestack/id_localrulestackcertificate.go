package certificateobjectlocalrulestack

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LocalRulestackCertificateId{}

// LocalRulestackCertificateId is a struct representing the Resource ID for a Local Rulestack Certificate
type LocalRulestackCertificateId struct {
	SubscriptionId     string
	ResourceGroupName  string
	LocalRulestackName string
	CertificateName    string
}

// NewLocalRulestackCertificateID returns a new LocalRulestackCertificateId struct
func NewLocalRulestackCertificateID(subscriptionId string, resourceGroupName string, localRulestackName string, certificateName string) LocalRulestackCertificateId {
	return LocalRulestackCertificateId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		LocalRulestackName: localRulestackName,
		CertificateName:    certificateName,
	}
}

// ParseLocalRulestackCertificateID parses 'input' into a LocalRulestackCertificateId
func ParseLocalRulestackCertificateID(input string) (*LocalRulestackCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocalRulestackCertificateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocalRulestackCertificateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocalRulestackName, ok = parsed.Parsed["localRulestackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "localRulestackName", *parsed)
	}

	if id.CertificateName, ok = parsed.Parsed["certificateName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "certificateName", *parsed)
	}

	return &id, nil
}

// ParseLocalRulestackCertificateIDInsensitively parses 'input' case-insensitively into a LocalRulestackCertificateId
// note: this method should only be used for API response data and not user input
func ParseLocalRulestackCertificateIDInsensitively(input string) (*LocalRulestackCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocalRulestackCertificateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocalRulestackCertificateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocalRulestackName, ok = parsed.Parsed["localRulestackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "localRulestackName", *parsed)
	}

	if id.CertificateName, ok = parsed.Parsed["certificateName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "certificateName", *parsed)
	}

	return &id, nil
}

// ValidateLocalRulestackCertificateID checks that 'input' can be parsed as a Local Rulestack Certificate ID
func ValidateLocalRulestackCertificateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLocalRulestackCertificateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Local Rulestack Certificate ID
func (id LocalRulestackCertificateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/PaloAltoNetworks.Cloudngfw/localRulestacks/%s/certificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocalRulestackName, id.CertificateName)
}

// Segments returns a slice of Resource ID Segments which comprise this Local Rulestack Certificate ID
func (id LocalRulestackCertificateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticPaloAltoNetworksCloudngfw", "PaloAltoNetworks.Cloudngfw", "PaloAltoNetworks.Cloudngfw"),
		resourceids.StaticSegment("staticLocalRulestacks", "localRulestacks", "localRulestacks"),
		resourceids.UserSpecifiedSegment("localRulestackName", "localRulestackValue"),
		resourceids.StaticSegment("staticCertificates", "certificates", "certificates"),
		resourceids.UserSpecifiedSegment("certificateName", "certificateValue"),
	}
}

// String returns a human-readable description of this Local Rulestack Certificate ID
func (id LocalRulestackCertificateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Local Rulestack Name: %q", id.LocalRulestackName),
		fmt.Sprintf("Certificate Name: %q", id.CertificateName),
	}
	return fmt.Sprintf("Local Rulestack Certificate (%s)", strings.Join(components, "\n"))
}
