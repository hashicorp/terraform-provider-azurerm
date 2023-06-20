package nginxcertificate

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
	SubscriptionId      string
	ResourceGroupName   string
	NginxDeploymentName string
	CertificateName     string
}

// NewCertificateID returns a new CertificateId struct
func NewCertificateID(subscriptionId string, resourceGroupName string, nginxDeploymentName string, certificateName string) CertificateId {
	return CertificateId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		NginxDeploymentName: nginxDeploymentName,
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

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NginxDeploymentName, ok = parsed.Parsed["nginxDeploymentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "nginxDeploymentName", *parsed)
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

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NginxDeploymentName, ok = parsed.Parsed["nginxDeploymentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "nginxDeploymentName", *parsed)
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
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Nginx.NginxPlus/nginxDeployments/%s/certificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NginxDeploymentName, id.CertificateName)
}

// Segments returns a slice of Resource ID Segments which comprise this Certificate ID
func (id CertificateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticNginxNginxPlus", "Nginx.NginxPlus", "Nginx.NginxPlus"),
		resourceids.StaticSegment("staticNginxDeployments", "nginxDeployments", "nginxDeployments"),
		resourceids.UserSpecifiedSegment("nginxDeploymentName", "nginxDeploymentValue"),
		resourceids.StaticSegment("staticCertificates", "certificates", "certificates"),
		resourceids.UserSpecifiedSegment("certificateName", "certificateValue"),
	}
}

// String returns a human-readable description of this Certificate ID
func (id CertificateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Nginx Deployment Name: %q", id.NginxDeploymentName),
		fmt.Sprintf("Certificate Name: %q", id.CertificateName),
	}
	return fmt.Sprintf("Certificate (%s)", strings.Join(components, "\n"))
}
