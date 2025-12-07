package resourceanchors

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ResourceAnchorId{})
}

var _ resourceids.ResourceId = &ResourceAnchorId{}

// ResourceAnchorId is a struct representing the Resource ID for a Resource Anchor
type ResourceAnchorId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ResourceAnchorName string
}

// NewResourceAnchorID returns a new ResourceAnchorId struct
func NewResourceAnchorID(subscriptionId string, resourceGroupName string, resourceAnchorName string) ResourceAnchorId {
	return ResourceAnchorId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ResourceAnchorName: resourceAnchorName,
	}
}

// ParseResourceAnchorID parses 'input' into a ResourceAnchorId
func ParseResourceAnchorID(input string) (*ResourceAnchorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ResourceAnchorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ResourceAnchorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseResourceAnchorIDInsensitively parses 'input' case-insensitively into a ResourceAnchorId
// note: this method should only be used for API response data and not user input
func ParseResourceAnchorIDInsensitively(input string) (*ResourceAnchorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ResourceAnchorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ResourceAnchorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ResourceAnchorId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ResourceAnchorName, ok = input.Parsed["resourceAnchorName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceAnchorName", input)
	}

	return nil
}

// ValidateResourceAnchorID checks that 'input' can be parsed as a Resource Anchor ID
func ValidateResourceAnchorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseResourceAnchorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Resource Anchor ID
func (id ResourceAnchorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Oracle.Database/resourceAnchors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceAnchorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Resource Anchor ID
func (id ResourceAnchorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticResourceAnchors", "resourceAnchors", "resourceAnchors"),
		resourceids.UserSpecifiedSegment("resourceAnchorName", "resourceAnchorName"),
	}
}

// String returns a human-readable description of this Resource Anchor ID
func (id ResourceAnchorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Anchor Name: %q", id.ResourceAnchorName),
	}
	return fmt.Sprintf("Resource Anchor (%s)", strings.Join(components, "\n"))
}
