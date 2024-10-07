package loadtests

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&LoadTestId{})
}

var _ resourceids.ResourceId = &LoadTestId{}

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
	parser := resourceids.NewParserFromResourceIdType(&LoadTestId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LoadTestId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLoadTestIDInsensitively parses 'input' case-insensitively into a LoadTestId
// note: this method should only be used for API response data and not user input
func ParseLoadTestIDInsensitively(input string) (*LoadTestId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LoadTestId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LoadTestId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LoadTestId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.LoadTestName, ok = input.Parsed["loadTestName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "loadTestName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("loadTestName", "loadTestName"),
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
