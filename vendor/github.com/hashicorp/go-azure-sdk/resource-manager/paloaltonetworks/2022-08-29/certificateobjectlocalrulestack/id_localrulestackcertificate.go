package certificateobjectlocalrulestack

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LocalRuleStackCertificateId{}

// LocalRuleStackCertificateId is a struct representing the Resource ID for a Local Rule Stack Certificate
type LocalRuleStackCertificateId struct {
	SubscriptionId     string
	ResourceGroupName  string
	LocalRuleStackName string
	CertificateName    string
}

// NewLocalRuleStackCertificateID returns a new LocalRuleStackCertificateId struct
func NewLocalRuleStackCertificateID(subscriptionId string, resourceGroupName string, localRuleStackName string, certificateName string) LocalRuleStackCertificateId {
	return LocalRuleStackCertificateId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		LocalRuleStackName: localRuleStackName,
		CertificateName:    certificateName,
	}
}

// ParseLocalRuleStackCertificateID parses 'input' into a LocalRuleStackCertificateId
func ParseLocalRuleStackCertificateID(input string) (*LocalRuleStackCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocalRuleStackCertificateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocalRuleStackCertificateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocalRuleStackName, ok = parsed.Parsed["localRuleStackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "localRuleStackName", *parsed)
	}

	if id.CertificateName, ok = parsed.Parsed["certificateName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "certificateName", *parsed)
	}

	return &id, nil
}

// ParseLocalRuleStackCertificateIDInsensitively parses 'input' case-insensitively into a LocalRuleStackCertificateId
// note: this method should only be used for API response data and not user input
func ParseLocalRuleStackCertificateIDInsensitively(input string) (*LocalRuleStackCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocalRuleStackCertificateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocalRuleStackCertificateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocalRuleStackName, ok = parsed.Parsed["localRuleStackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "localRuleStackName", *parsed)
	}

	if id.CertificateName, ok = parsed.Parsed["certificateName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "certificateName", *parsed)
	}

	return &id, nil
}

// ValidateLocalRuleStackCertificateID checks that 'input' can be parsed as a Local Rule Stack Certificate ID
func ValidateLocalRuleStackCertificateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLocalRuleStackCertificateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Local Rule Stack Certificate ID
func (id LocalRuleStackCertificateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/PaloAltoNetworks.CloudNGFW/localRuleStacks/%s/certificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocalRuleStackName, id.CertificateName)
}

// Segments returns a slice of Resource ID Segments which comprise this Local Rule Stack Certificate ID
func (id LocalRuleStackCertificateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticPaloAltoNetworksCloudNGFW", "PaloAltoNetworks.CloudNGFW", "PaloAltoNetworks.CloudNGFW"),
		resourceids.StaticSegment("staticLocalRuleStacks", "localRuleStacks", "localRuleStacks"),
		resourceids.UserSpecifiedSegment("localRuleStackName", "localRuleStackValue"),
		resourceids.StaticSegment("staticCertificates", "certificates", "certificates"),
		resourceids.UserSpecifiedSegment("certificateName", "certificateValue"),
	}
}

// String returns a human-readable description of this Local Rule Stack Certificate ID
func (id LocalRuleStackCertificateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Local Rule Stack Name: %q", id.LocalRuleStackName),
		fmt.Sprintf("Certificate Name: %q", id.CertificateName),
	}
	return fmt.Sprintf("Local Rule Stack Certificate (%s)", strings.Join(components, "\n"))
}
