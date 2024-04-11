package webpubsub

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &CustomDomainId{}

// CustomDomainId is a struct representing the Resource ID for a Custom Domain
type CustomDomainId struct {
	SubscriptionId    string
	ResourceGroupName string
	WebPubSubName     string
	CustomDomainName  string
}

// NewCustomDomainID returns a new CustomDomainId struct
func NewCustomDomainID(subscriptionId string, resourceGroupName string, webPubSubName string, customDomainName string) CustomDomainId {
	return CustomDomainId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WebPubSubName:     webPubSubName,
		CustomDomainName:  customDomainName,
	}
}

// ParseCustomDomainID parses 'input' into a CustomDomainId
func ParseCustomDomainID(input string) (*CustomDomainId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CustomDomainId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CustomDomainId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCustomDomainIDInsensitively parses 'input' case-insensitively into a CustomDomainId
// note: this method should only be used for API response data and not user input
func ParseCustomDomainIDInsensitively(input string) (*CustomDomainId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CustomDomainId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CustomDomainId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CustomDomainId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.WebPubSubName, ok = input.Parsed["webPubSubName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "webPubSubName", input)
	}

	if id.CustomDomainName, ok = input.Parsed["customDomainName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "customDomainName", input)
	}

	return nil
}

// ValidateCustomDomainID checks that 'input' can be parsed as a Custom Domain ID
func ValidateCustomDomainID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCustomDomainID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Custom Domain ID
func (id CustomDomainId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SignalRService/webPubSub/%s/customDomains/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WebPubSubName, id.CustomDomainName)
}

// Segments returns a slice of Resource ID Segments which comprise this Custom Domain ID
func (id CustomDomainId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSignalRService", "Microsoft.SignalRService", "Microsoft.SignalRService"),
		resourceids.StaticSegment("staticWebPubSub", "webPubSub", "webPubSub"),
		resourceids.UserSpecifiedSegment("webPubSubName", "webPubSubValue"),
		resourceids.StaticSegment("staticCustomDomains", "customDomains", "customDomains"),
		resourceids.UserSpecifiedSegment("customDomainName", "customDomainValue"),
	}
}

// String returns a human-readable description of this Custom Domain ID
func (id CustomDomainId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Web Pub Sub Name: %q", id.WebPubSubName),
		fmt.Sprintf("Custom Domain Name: %q", id.CustomDomainName),
	}
	return fmt.Sprintf("Custom Domain (%s)", strings.Join(components, "\n"))
}
