package webtestsapis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = WebTestId{}

// WebTestId is a struct representing the Resource ID for a Web Test
type WebTestId struct {
	SubscriptionId    string
	ResourceGroupName string
	WebTestName       string
}

// NewWebTestID returns a new WebTestId struct
func NewWebTestID(subscriptionId string, resourceGroupName string, webTestName string) WebTestId {
	return WebTestId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WebTestName:       webTestName,
	}
}

// ParseWebTestID parses 'input' into a WebTestId
func ParseWebTestID(input string) (*WebTestId, error) {
	parser := resourceids.NewParserFromResourceIdType(WebTestId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := WebTestId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WebTestName, ok = parsed.Parsed["webTestName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "webTestName", *parsed)
	}

	return &id, nil
}

// ParseWebTestIDInsensitively parses 'input' case-insensitively into a WebTestId
// note: this method should only be used for API response data and not user input
func ParseWebTestIDInsensitively(input string) (*WebTestId, error) {
	parser := resourceids.NewParserFromResourceIdType(WebTestId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := WebTestId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WebTestName, ok = parsed.Parsed["webTestName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "webTestName", *parsed)
	}

	return &id, nil
}

// ValidateWebTestID checks that 'input' can be parsed as a Web Test ID
func ValidateWebTestID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseWebTestID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Web Test ID
func (id WebTestId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/webTests/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WebTestName)
}

// Segments returns a slice of Resource ID Segments which comprise this Web Test ID
func (id WebTestId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticWebTests", "webTests", "webTests"),
		resourceids.UserSpecifiedSegment("webTestName", "webTestValue"),
	}
}

// String returns a human-readable description of this Web Test ID
func (id WebTestId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Web Test Name: %q", id.WebTestName),
	}
	return fmt.Sprintf("Web Test (%s)", strings.Join(components, "\n"))
}
