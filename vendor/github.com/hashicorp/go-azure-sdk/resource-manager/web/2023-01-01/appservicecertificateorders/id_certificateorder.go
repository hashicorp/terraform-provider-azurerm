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
	recaser.RegisterResourceId(&CertificateOrderId{})
}

var _ resourceids.ResourceId = &CertificateOrderId{}

// CertificateOrderId is a struct representing the Resource ID for a Certificate Order
type CertificateOrderId struct {
	SubscriptionId       string
	ResourceGroupName    string
	CertificateOrderName string
}

// NewCertificateOrderID returns a new CertificateOrderId struct
func NewCertificateOrderID(subscriptionId string, resourceGroupName string, certificateOrderName string) CertificateOrderId {
	return CertificateOrderId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		CertificateOrderName: certificateOrderName,
	}
}

// ParseCertificateOrderID parses 'input' into a CertificateOrderId
func ParseCertificateOrderID(input string) (*CertificateOrderId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CertificateOrderId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CertificateOrderId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCertificateOrderIDInsensitively parses 'input' case-insensitively into a CertificateOrderId
// note: this method should only be used for API response data and not user input
func ParseCertificateOrderIDInsensitively(input string) (*CertificateOrderId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CertificateOrderId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CertificateOrderId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CertificateOrderId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateCertificateOrderID checks that 'input' can be parsed as a Certificate Order ID
func ValidateCertificateOrderID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCertificateOrderID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Certificate Order ID
func (id CertificateOrderId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CertificateRegistration/certificateOrders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CertificateOrderName)
}

// Segments returns a slice of Resource ID Segments which comprise this Certificate Order ID
func (id CertificateOrderId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCertificateRegistration", "Microsoft.CertificateRegistration", "Microsoft.CertificateRegistration"),
		resourceids.StaticSegment("staticCertificateOrders", "certificateOrders", "certificateOrders"),
		resourceids.UserSpecifiedSegment("certificateOrderName", "certificateOrderValue"),
	}
}

// String returns a human-readable description of this Certificate Order ID
func (id CertificateOrderId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Certificate Order Name: %q", id.CertificateOrderName),
	}
	return fmt.Sprintf("Certificate Order (%s)", strings.Join(components, "\n"))
}
