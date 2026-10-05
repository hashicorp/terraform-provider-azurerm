package resourceguards

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ResourceGuardId{})
}

var _ resourceids.ResourceId = &ResourceGuardId{}

// ResourceGuardId is a struct representing the Resource ID for a Resource Guard
type ResourceGuardId struct {
	SubscriptionId    string
	ResourceGroupName string
	ResourceGuardName string
}

// NewResourceGuardID returns a new ResourceGuardId struct
func NewResourceGuardID(subscriptionId string, resourceGroupName string, resourceGuardName string) ResourceGuardId {
	return ResourceGuardId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ResourceGuardName: resourceGuardName,
	}
}

// ParseResourceGuardID parses 'input' into a ResourceGuardId
func ParseResourceGuardID(input string) (*ResourceGuardId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ResourceGuardId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ResourceGuardId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseResourceGuardIDInsensitively parses 'input' case-insensitively into a ResourceGuardId
// note: this method should only be used for API response data and not user input
func ParseResourceGuardIDInsensitively(input string) (*ResourceGuardId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ResourceGuardId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ResourceGuardId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ResourceGuardId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ResourceGuardName, ok = input.Parsed["resourceGuardName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGuardName", input)
	}

	return nil
}

// ValidateResourceGuardID checks that 'input' can be parsed as a Resource Guard ID
func ValidateResourceGuardID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseResourceGuardID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Resource Guard ID
func (id ResourceGuardId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/resourceGuards/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceGuardName)
}

// Segments returns a slice of Resource ID Segments which comprise this Resource Guard ID
func (id ResourceGuardId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticResourceGuards", "resourceGuards", "resourceGuards"),
		resourceids.UserSpecifiedSegment("resourceGuardName", "resourceGuardName"),
	}
}

// String returns a human-readable description of this Resource Guard ID
func (id ResourceGuardId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Guard Name: %q", id.ResourceGuardName),
	}
	return fmt.Sprintf("Resource Guard (%s)", strings.Join(components, "\n"))
}
