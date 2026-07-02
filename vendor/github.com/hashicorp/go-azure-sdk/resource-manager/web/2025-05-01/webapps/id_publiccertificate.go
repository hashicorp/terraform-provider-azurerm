package webapps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PublicCertificateId{})
}

var _ resourceids.ResourceId = &PublicCertificateId{}

// PublicCertificateId is a struct representing the Resource ID for a Public Certificate
type PublicCertificateId struct {
	SubscriptionId        string
	ResourceGroupName     string
	SiteName              string
	PublicCertificateName string
}

// NewPublicCertificateID returns a new PublicCertificateId struct
func NewPublicCertificateID(subscriptionId string, resourceGroupName string, siteName string, publicCertificateName string) PublicCertificateId {
	return PublicCertificateId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		SiteName:              siteName,
		PublicCertificateName: publicCertificateName,
	}
}

// ParsePublicCertificateID parses 'input' into a PublicCertificateId
func ParsePublicCertificateID(input string) (*PublicCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PublicCertificateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PublicCertificateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePublicCertificateIDInsensitively parses 'input' case-insensitively into a PublicCertificateId
// note: this method should only be used for API response data and not user input
func ParsePublicCertificateIDInsensitively(input string) (*PublicCertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PublicCertificateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PublicCertificateId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PublicCertificateId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SiteName, ok = input.Parsed["siteName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "siteName", input)
	}

	if id.PublicCertificateName, ok = input.Parsed["publicCertificateName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "publicCertificateName", input)
	}

	return nil
}

// ValidatePublicCertificateID checks that 'input' can be parsed as a Public Certificate ID
func ValidatePublicCertificateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePublicCertificateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Public Certificate ID
func (id PublicCertificateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/publicCertificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SiteName, id.PublicCertificateName)
}

// Segments returns a slice of Resource ID Segments which comprise this Public Certificate ID
func (id PublicCertificateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticSites", "sites", "sites"),
		resourceids.UserSpecifiedSegment("siteName", "siteName"),
		resourceids.StaticSegment("staticPublicCertificates", "publicCertificates", "publicCertificates"),
		resourceids.UserSpecifiedSegment("publicCertificateName", "publicCertificateName"),
	}
}

// String returns a human-readable description of this Public Certificate ID
func (id PublicCertificateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Site Name: %q", id.SiteName),
		fmt.Sprintf("Public Certificate Name: %q", id.PublicCertificateName),
	}
	return fmt.Sprintf("Public Certificate (%s)", strings.Join(components, "\n"))
}
