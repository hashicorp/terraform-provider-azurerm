package gatewaycertificateauthority

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &CertificateAuthorityId{}

// CertificateAuthorityId is a struct representing the Resource ID for a Certificate Authority
type CertificateAuthorityId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	GatewayId         string
	CertificateId     string
}

// NewCertificateAuthorityID returns a new CertificateAuthorityId struct
func NewCertificateAuthorityID(subscriptionId string, resourceGroupName string, serviceName string, gatewayId string, certificateId string) CertificateAuthorityId {
	return CertificateAuthorityId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		GatewayId:         gatewayId,
		CertificateId:     certificateId,
	}
}

// ParseCertificateAuthorityID parses 'input' into a CertificateAuthorityId
func ParseCertificateAuthorityID(input string) (*CertificateAuthorityId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CertificateAuthorityId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CertificateAuthorityId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCertificateAuthorityIDInsensitively parses 'input' case-insensitively into a CertificateAuthorityId
// note: this method should only be used for API response data and not user input
func ParseCertificateAuthorityIDInsensitively(input string) (*CertificateAuthorityId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CertificateAuthorityId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CertificateAuthorityId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CertificateAuthorityId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServiceName, ok = input.Parsed["serviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceName", input)
	}

	if id.GatewayId, ok = input.Parsed["gatewayId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "gatewayId", input)
	}

	if id.CertificateId, ok = input.Parsed["certificateId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "certificateId", input)
	}

	return nil
}

// ValidateCertificateAuthorityID checks that 'input' can be parsed as a Certificate Authority ID
func ValidateCertificateAuthorityID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCertificateAuthorityID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Certificate Authority ID
func (id CertificateAuthorityId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/gateways/%s/certificateAuthorities/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.GatewayId, id.CertificateId)
}

// Segments returns a slice of Resource ID Segments which comprise this Certificate Authority ID
func (id CertificateAuthorityId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
		resourceids.StaticSegment("staticGateways", "gateways", "gateways"),
		resourceids.UserSpecifiedSegment("gatewayId", "gatewayIdValue"),
		resourceids.StaticSegment("staticCertificateAuthorities", "certificateAuthorities", "certificateAuthorities"),
		resourceids.UserSpecifiedSegment("certificateId", "certificateIdValue"),
	}
}

// String returns a human-readable description of this Certificate Authority ID
func (id CertificateAuthorityId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Gateway: %q", id.GatewayId),
		fmt.Sprintf("Certificate: %q", id.CertificateId),
	}
	return fmt.Sprintf("Certificate Authority (%s)", strings.Join(components, "\n"))
}
