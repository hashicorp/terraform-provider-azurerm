package dedicatedhsms

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DedicatedHSMId{})
}

var _ resourceids.ResourceId = &DedicatedHSMId{}

// DedicatedHSMId is a struct representing the Resource ID for a Dedicated H S M
type DedicatedHSMId struct {
	SubscriptionId    string
	ResourceGroupName string
	DedicatedHSMName  string
}

// NewDedicatedHSMID returns a new DedicatedHSMId struct
func NewDedicatedHSMID(subscriptionId string, resourceGroupName string, dedicatedHSMName string) DedicatedHSMId {
	return DedicatedHSMId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DedicatedHSMName:  dedicatedHSMName,
	}
}

// ParseDedicatedHSMID parses 'input' into a DedicatedHSMId
func ParseDedicatedHSMID(input string) (*DedicatedHSMId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DedicatedHSMId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DedicatedHSMId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDedicatedHSMIDInsensitively parses 'input' case-insensitively into a DedicatedHSMId
// note: this method should only be used for API response data and not user input
func ParseDedicatedHSMIDInsensitively(input string) (*DedicatedHSMId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DedicatedHSMId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DedicatedHSMId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DedicatedHSMId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DedicatedHSMName, ok = input.Parsed["dedicatedHSMName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dedicatedHSMName", input)
	}

	return nil
}

// ValidateDedicatedHSMID checks that 'input' can be parsed as a Dedicated H S M ID
func ValidateDedicatedHSMID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDedicatedHSMID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dedicated H S M ID
func (id DedicatedHSMId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HardwareSecurityModules/dedicatedHSMs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DedicatedHSMName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dedicated H S M ID
func (id DedicatedHSMId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftHardwareSecurityModules", "Microsoft.HardwareSecurityModules", "Microsoft.HardwareSecurityModules"),
		resourceids.StaticSegment("staticDedicatedHSMs", "dedicatedHSMs", "dedicatedHSMs"),
		resourceids.UserSpecifiedSegment("dedicatedHSMName", "dedicatedHSMName"),
	}
}

// String returns a human-readable description of this Dedicated H S M ID
func (id DedicatedHSMId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dedicated H S M Name: %q", id.DedicatedHSMName),
	}
	return fmt.Sprintf("Dedicated H S M (%s)", strings.Join(components, "\n"))
}
