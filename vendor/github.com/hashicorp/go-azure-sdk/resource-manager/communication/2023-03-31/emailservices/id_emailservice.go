package emailservices

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = EmailServiceId{}

// EmailServiceId is a struct representing the Resource ID for a Email Service
type EmailServiceId struct {
	SubscriptionId    string
	ResourceGroupName string
	EmailServiceName  string
}

// NewEmailServiceID returns a new EmailServiceId struct
func NewEmailServiceID(subscriptionId string, resourceGroupName string, emailServiceName string) EmailServiceId {
	return EmailServiceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		EmailServiceName:  emailServiceName,
	}
}

// ParseEmailServiceID parses 'input' into a EmailServiceId
func ParseEmailServiceID(input string) (*EmailServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(EmailServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EmailServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.EmailServiceName, ok = parsed.Parsed["emailServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "emailServiceName", *parsed)
	}

	return &id, nil
}

// ParseEmailServiceIDInsensitively parses 'input' case-insensitively into a EmailServiceId
// note: this method should only be used for API response data and not user input
func ParseEmailServiceIDInsensitively(input string) (*EmailServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(EmailServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := EmailServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.EmailServiceName, ok = parsed.Parsed["emailServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "emailServiceName", *parsed)
	}

	return &id, nil
}

// ValidateEmailServiceID checks that 'input' can be parsed as a Email Service ID
func ValidateEmailServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEmailServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Email Service ID
func (id EmailServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Communication/emailServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.EmailServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Email Service ID
func (id EmailServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCommunication", "Microsoft.Communication", "Microsoft.Communication"),
		resourceids.StaticSegment("staticEmailServices", "emailServices", "emailServices"),
		resourceids.UserSpecifiedSegment("emailServiceName", "emailServiceValue"),
	}
}

// String returns a human-readable description of this Email Service ID
func (id EmailServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Email Service Name: %q", id.EmailServiceName),
	}
	return fmt.Sprintf("Email Service (%s)", strings.Join(components, "\n"))
}
