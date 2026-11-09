package expressroutelags

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ExpressRouteLagId{})
}

var _ resourceids.ResourceId = &ExpressRouteLagId{}

// ExpressRouteLagId is a struct representing the Resource ID for a Express Route Lag
type ExpressRouteLagId struct {
	SubscriptionId      string
	ResourceGroupName   string
	ExpressRouteLagName string
}

// NewExpressRouteLagID returns a new ExpressRouteLagId struct
func NewExpressRouteLagID(subscriptionId string, resourceGroupName string, expressRouteLagName string) ExpressRouteLagId {
	return ExpressRouteLagId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		ExpressRouteLagName: expressRouteLagName,
	}
}

// ParseExpressRouteLagID parses 'input' into a ExpressRouteLagId
func ParseExpressRouteLagID(input string) (*ExpressRouteLagId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExpressRouteLagId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExpressRouteLagId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseExpressRouteLagIDInsensitively parses 'input' case-insensitively into a ExpressRouteLagId
// note: this method should only be used for API response data and not user input
func ParseExpressRouteLagIDInsensitively(input string) (*ExpressRouteLagId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExpressRouteLagId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExpressRouteLagId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ExpressRouteLagId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ExpressRouteLagName, ok = input.Parsed["expressRouteLagName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "expressRouteLagName", input)
	}

	return nil
}

// ValidateExpressRouteLagID checks that 'input' can be parsed as a Express Route Lag ID
func ValidateExpressRouteLagID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExpressRouteLagID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Express Route Lag ID
func (id ExpressRouteLagId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/expressRouteLags/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ExpressRouteLagName)
}

// Segments returns a slice of Resource ID Segments which comprise this Express Route Lag ID
func (id ExpressRouteLagId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticExpressRouteLags", "expressRouteLags", "expressRouteLags"),
		resourceids.UserSpecifiedSegment("expressRouteLagName", "expressRouteLagName"),
	}
}

// String returns a human-readable description of this Express Route Lag ID
func (id ExpressRouteLagId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Express Route Lag Name: %q", id.ExpressRouteLagName),
	}
	return fmt.Sprintf("Express Route Lag (%s)", strings.Join(components, "\n"))
}
