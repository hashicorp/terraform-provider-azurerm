package loadtest

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LoadTestId{}

// LoadTestId is a struct representing the Resource ID for a Load Test
type LoadTestId struct {
	SubscriptionId    string
	ResourceGroupName string
	LoadTestName      string
}

// NewLoadTestID returns a new LoadTestId struct
func NewLoadTestID(subscriptionId string, resourceGroupName string, loadTestName string) LoadTestId {
	return LoadTestId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		LoadTestName:      loadTestName,
	}
}

// ParseLoadTestID parses 'input' into a LoadTestId
func ParseLoadTestID(input string) (*LoadTestId, error) {
	parser := resourceids.NewParserFromResourceIdType(LoadTestId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LoadTestId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LoadTestName, ok = parsed.Parsed["loadTestName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadTestName", *parsed)
	}

	return &id, nil
}

// ParseLoadTestIDInsensitively parses 'input' case-insensitively into a LoadTestId
// note: this method should only be used for API response data and not user input
func ParseLoadTestIDInsensitively(input string) (*LoadTestId, error) {
	parser := resourceids.NewParserFromResourceIdType(LoadTestId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LoadTestId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LoadTestName, ok = parsed.Parsed["loadTestName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loadTestName", *parsed)
	}

	return &id, nil
}

// ValidateLoadTestID checks that 'input' can be parsed as a Load Test ID
func ValidateLoadTestID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLoadTestID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Load Test ID
func (id LoadTestId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.LoadTestService/loadTests/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LoadTestName)
}

// Segments returns a slice of Resource ID Segments which comprise this Load Test ID
func (id LoadTestId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLoadTestService", "Microsoft.LoadTestService", "Microsoft.LoadTestService"),
		resourceids.StaticSegment("staticLoadTests", "loadTests", "loadTests"),
		resourceids.UserSpecifiedSegment("loadTestName", "loadTestValue"),
	}
}

// String returns a human-readable description of this Load Test ID
func (id LoadTestId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Load Test Name: %q", id.LoadTestName),
	}
	return fmt.Sprintf("Load Test (%s)", strings.Join(components, "\n"))
}
