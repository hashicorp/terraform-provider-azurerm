package testlines

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TestLineId{})
}

var _ resourceids.ResourceId = &TestLineId{}

// TestLineId is a struct representing the Resource ID for a Test Line
type TestLineId struct {
	SubscriptionId            string
	ResourceGroupName         string
	CommunicationsGatewayName string
	TestLineName              string
}

// NewTestLineID returns a new TestLineId struct
func NewTestLineID(subscriptionId string, resourceGroupName string, communicationsGatewayName string, testLineName string) TestLineId {
	return TestLineId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		CommunicationsGatewayName: communicationsGatewayName,
		TestLineName:              testLineName,
	}
}

// ParseTestLineID parses 'input' into a TestLineId
func ParseTestLineID(input string) (*TestLineId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TestLineId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TestLineId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTestLineIDInsensitively parses 'input' case-insensitively into a TestLineId
// note: this method should only be used for API response data and not user input
func ParseTestLineIDInsensitively(input string) (*TestLineId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TestLineId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TestLineId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TestLineId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CommunicationsGatewayName, ok = input.Parsed["communicationsGatewayName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "communicationsGatewayName", input)
	}

	if id.TestLineName, ok = input.Parsed["testLineName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "testLineName", input)
	}

	return nil
}

// ValidateTestLineID checks that 'input' can be parsed as a Test Line ID
func ValidateTestLineID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTestLineID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Test Line ID
func (id TestLineId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.VoiceServices/communicationsGateways/%s/testLines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CommunicationsGatewayName, id.TestLineName)
}

// Segments returns a slice of Resource ID Segments which comprise this Test Line ID
func (id TestLineId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftVoiceServices", "Microsoft.VoiceServices", "Microsoft.VoiceServices"),
		resourceids.StaticSegment("staticCommunicationsGateways", "communicationsGateways", "communicationsGateways"),
		resourceids.UserSpecifiedSegment("communicationsGatewayName", "communicationsGatewayName"),
		resourceids.StaticSegment("staticTestLines", "testLines", "testLines"),
		resourceids.UserSpecifiedSegment("testLineName", "testLineName"),
	}
}

// String returns a human-readable description of this Test Line ID
func (id TestLineId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Communications Gateway Name: %q", id.CommunicationsGatewayName),
		fmt.Sprintf("Test Line Name: %q", id.TestLineName),
	}
	return fmt.Sprintf("Test Line (%s)", strings.Join(components, "\n"))
}
