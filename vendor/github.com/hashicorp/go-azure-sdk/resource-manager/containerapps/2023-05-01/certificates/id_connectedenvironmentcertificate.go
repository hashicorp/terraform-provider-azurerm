package certificates

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ConnectedEnvironmentCertificateId{}

// ConnectedEnvironmentCertificateId is a struct representing the Resource ID for a Connected Environment Certificate
type ConnectedEnvironmentCertificateId struct {
	SubscriptionId           string
	ResourceGroupName        string
	ConnectedEnvironmentName string
	CertificateName          string
}

// NewConnectedEnvironmentCertificateID returns a new ConnectedEnvironmentCertificateId struct
func NewConnectedEnvironmentCertificateID(subscriptionId string, resourceGroupName string, connectedEnvironmentName string, certificateName string) ConnectedEnvironmentCertificateId {
	return ConnectedEnvironmentCertificateId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		ConnectedEnvironmentName: connectedEnvironmentName,
		CertificateName:          certificateName,
	}
}

// ParseConnectedEnvironmentCertificateID parses 'input' into a ConnectedEnvironmentCertificateId
func ParseConnectedEnvironmentCertificateID(input string) (*ConnectedEnvironmentCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConnectedEnvironmentCertificateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConnectedEnvironmentCertificateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ConnectedEnvironmentName, ok = parsed.Parsed["connectedEnvironmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "connectedEnvironmentName", *parsed)
	}

	if id.CertificateName, ok = parsed.Parsed["certificateName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "certificateName", *parsed)
	}

	return &id, nil
}

// ParseConnectedEnvironmentCertificateIDInsensitively parses 'input' case-insensitively into a ConnectedEnvironmentCertificateId
// note: this method should only be used for API response data and not user input
func ParseConnectedEnvironmentCertificateIDInsensitively(input string) (*ConnectedEnvironmentCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConnectedEnvironmentCertificateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConnectedEnvironmentCertificateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ConnectedEnvironmentName, ok = parsed.Parsed["connectedEnvironmentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "connectedEnvironmentName", *parsed)
	}

	if id.CertificateName, ok = parsed.Parsed["certificateName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "certificateName", *parsed)
	}

	return &id, nil
}

// ValidateConnectedEnvironmentCertificateID checks that 'input' can be parsed as a Connected Environment Certificate ID
func ValidateConnectedEnvironmentCertificateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConnectedEnvironmentCertificateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Connected Environment Certificate ID
func (id ConnectedEnvironmentCertificateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/connectedEnvironments/%s/certificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ConnectedEnvironmentName, id.CertificateName)
}

// Segments returns a slice of Resource ID Segments which comprise this Connected Environment Certificate ID
func (id ConnectedEnvironmentCertificateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticConnectedEnvironments", "connectedEnvironments", "connectedEnvironments"),
		resourceids.UserSpecifiedSegment("connectedEnvironmentName", "connectedEnvironmentValue"),
		resourceids.StaticSegment("staticCertificates", "certificates", "certificates"),
		resourceids.UserSpecifiedSegment("certificateName", "certificateValue"),
	}
}

// String returns a human-readable description of this Connected Environment Certificate ID
func (id ConnectedEnvironmentCertificateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Connected Environment Name: %q", id.ConnectedEnvironmentName),
		fmt.Sprintf("Certificate Name: %q", id.CertificateName),
	}
	return fmt.Sprintf("Connected Environment Certificate (%s)", strings.Join(components, "\n"))
}
