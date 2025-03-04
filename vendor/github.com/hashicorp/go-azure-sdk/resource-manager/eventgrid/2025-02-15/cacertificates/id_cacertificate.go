package cacertificates

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CaCertificateId{})
}

var _ resourceids.ResourceId = &CaCertificateId{}

// CaCertificateId is a struct representing the Resource ID for a Ca Certificate
type CaCertificateId struct {
	SubscriptionId    string
	ResourceGroupName string
	NamespaceName     string
	CaCertificateName string
}

// NewCaCertificateID returns a new CaCertificateId struct
func NewCaCertificateID(subscriptionId string, resourceGroupName string, namespaceName string, caCertificateName string) CaCertificateId {
	return CaCertificateId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NamespaceName:     namespaceName,
		CaCertificateName: caCertificateName,
	}
}

// ParseCaCertificateID parses 'input' into a CaCertificateId
func ParseCaCertificateID(input string) (*CaCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CaCertificateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CaCertificateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCaCertificateIDInsensitively parses 'input' case-insensitively into a CaCertificateId
// note: this method should only be used for API response data and not user input
func ParseCaCertificateIDInsensitively(input string) (*CaCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CaCertificateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CaCertificateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CaCertificateId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NamespaceName, ok = input.Parsed["namespaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "namespaceName", input)
	}

	if id.CaCertificateName, ok = input.Parsed["caCertificateName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "caCertificateName", input)
	}

	return nil
}

// ValidateCaCertificateID checks that 'input' can be parsed as a Ca Certificate ID
func ValidateCaCertificateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCaCertificateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Ca Certificate ID
func (id CaCertificateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.EventGrid/namespaces/%s/caCertificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NamespaceName, id.CaCertificateName)
}

// Segments returns a slice of Resource ID Segments which comprise this Ca Certificate ID
func (id CaCertificateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftEventGrid", "Microsoft.EventGrid", "Microsoft.EventGrid"),
		resourceids.StaticSegment("staticNamespaces", "namespaces", "namespaces"),
		resourceids.UserSpecifiedSegment("namespaceName", "namespaceName"),
		resourceids.StaticSegment("staticCaCertificates", "caCertificates", "caCertificates"),
		resourceids.UserSpecifiedSegment("caCertificateName", "caCertificateName"),
	}
}

// String returns a human-readable description of this Ca Certificate ID
func (id CaCertificateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Namespace Name: %q", id.NamespaceName),
		fmt.Sprintf("Ca Certificate Name: %q", id.CaCertificateName),
	}
	return fmt.Sprintf("Ca Certificate (%s)", strings.Join(components, "\n"))
}
