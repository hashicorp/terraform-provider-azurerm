package containerinstance

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NgroupId{})
}

var _ resourceids.ResourceId = &NgroupId{}

// NgroupId is a struct representing the Resource ID for a Ngroup
type NgroupId struct {
	SubscriptionId    string
	ResourceGroupName string
	NgroupName        string
}

// NewNgroupID returns a new NgroupId struct
func NewNgroupID(subscriptionId string, resourceGroupName string, ngroupName string) NgroupId {
	return NgroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		NgroupName:        ngroupName,
	}
}

// ParseNgroupID parses 'input' into a NgroupId
func ParseNgroupID(input string) (*NgroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NgroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NgroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNgroupIDInsensitively parses 'input' case-insensitively into a NgroupId
// note: this method should only be used for API response data and not user input
func ParseNgroupIDInsensitively(input string) (*NgroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NgroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NgroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NgroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NgroupName, ok = input.Parsed["ngroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "ngroupName", input)
	}

	return nil
}

// ValidateNgroupID checks that 'input' can be parsed as a Ngroup ID
func ValidateNgroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNgroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Ngroup ID
func (id NgroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerInstance/ngroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NgroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Ngroup ID
func (id NgroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerInstance", "Microsoft.ContainerInstance", "Microsoft.ContainerInstance"),
		resourceids.StaticSegment("staticNgroups", "ngroups", "ngroups"),
		resourceids.UserSpecifiedSegment("ngroupName", "ngroupName"),
	}
}

// String returns a human-readable description of this Ngroup ID
func (id NgroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Ngroup Name: %q", id.NgroupName),
	}
	return fmt.Sprintf("Ngroup (%s)", strings.Join(components, "\n"))
}
