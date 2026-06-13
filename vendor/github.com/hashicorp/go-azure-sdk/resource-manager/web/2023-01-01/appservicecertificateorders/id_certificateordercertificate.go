package appservicecertificateorders

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CertificateOrderCertificateId{})
}

var _ resourceids.ResourceId = &CertificateOrderCertificateId{}

// CertificateOrderCertificateId is a struct representing the Resource ID for a Certificate Order Certificate
type CertificateOrderCertificateId struct {
	SubscriptionId       string
	ResourceGroupName    string
	CertificateOrderName string
	CertificateName      string
}

// NewCertificateOrderCertificateID returns a new CertificateOrderCertificateId struct
func NewCertificateOrderCertificateID(subscriptionId string, resourceGroupName string, certificateOrderName string, certificateName string) CertificateOrderCertificateId {
	return CertificateOrderCertificateId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		CertificateOrderName: certificateOrderName,
		CertificateName:      certificateName,
	}
}

// ParseCertificateOrderCertificateID parses 'input' into a CertificateOrderCertificateId
func ParseCertificateOrderCertificateID(input string) (*CertificateOrderCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CertificateOrderCertificateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CertificateOrderCertificateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCertificateOrderCertificateIDInsensitively parses 'input' case-insensitively into a CertificateOrderCertificateId
// note: this method should only be used for API response data and not user input
func ParseCertificateOrderCertificateIDInsensitively(input string) (*CertificateOrderCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CertificateOrderCertificateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CertificateOrderCertificateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CertificateOrderCertificateId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CertificateOrderName, ok = input.Parsed["certificateOrderName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "certificateOrderName", input)
	}

	if id.CertificateName, ok = input.Parsed["certificateName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "certificateName", input)
	}

	return nil
}

// ValidateCertificateOrderCertificateID checks that 'input' can be parsed as a Certificate Order Certificate ID
func ValidateCertificateOrderCertificateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCertificateOrderCertificateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Certificate Order Certificate ID
func (id CertificateOrderCertificateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CertificateRegistration/certificateOrders/%s/certificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CertificateOrderName, id.CertificateName)
}

// Segments returns a slice of Resource ID Segments which comprise this Certificate Order Certificate ID
func (id CertificateOrderCertificateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCertificateRegistration", "Microsoft.CertificateRegistration", "Microsoft.CertificateRegistration"),
		resourceids.StaticSegment("staticCertificateOrders", "certificateOrders", "certificateOrders"),
		resourceids.UserSpecifiedSegment("certificateOrderName", "certificateOrderName"),
		resourceids.StaticSegment("staticCertificates", "certificates", "certificates"),
		resourceids.UserSpecifiedSegment("certificateName", "certificateName"),
	}
}

// String returns a human-readable description of this Certificate Order Certificate ID
func (id CertificateOrderCertificateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Certificate Order Name: %q", id.CertificateOrderName),
		fmt.Sprintf("Certificate Name: %q", id.CertificateName),
	}
	return fmt.Sprintf("Certificate Order Certificate (%s)", strings.Join(components, "\n"))
}
