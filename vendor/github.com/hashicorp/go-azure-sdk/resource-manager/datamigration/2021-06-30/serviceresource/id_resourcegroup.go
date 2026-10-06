package serviceresource

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ResourceGroupId{})
}

var _ resourceids.ResourceId = &ResourceGroupId{}

// ResourceGroupId is a struct representing the Resource ID for a Resource Group
type ResourceGroupId struct {
	SubscriptionId    string
	ResourceGroupName string
}

// NewResourceGroupID returns a new ResourceGroupId struct
func NewResourceGroupID(subscriptionId string, resourceGroupName string) ResourceGroupId {
	return ResourceGroupId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
	}
}

// ParseResourceGroupID parses 'input' into a ResourceGroupId
func ParseResourceGroupID(input string) (*ResourceGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ResourceGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ResourceGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseResourceGroupIDInsensitively parses 'input' case-insensitively into a ResourceGroupId
// note: this method should only be used for API response data and not user input
func ParseResourceGroupIDInsensitively(input string) (*ResourceGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ResourceGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ResourceGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ResourceGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	return nil
}

// ValidateResourceGroupID checks that 'input' can be parsed as a Resource Group ID
func ValidateResourceGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseResourceGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Resource Group ID
func (id ResourceGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Resource Group ID
func (id ResourceGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.UserSpecifiedSegment("resourceGroupName", "resourceGroupName"),
	}
}

// String returns a human-readable description of this Resource Group ID
func (id ResourceGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
	}
	return fmt.Sprintf("Resource Group (%s)", strings.Join(components, "\n"))
}
