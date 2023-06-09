package customipprefixes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = CustomIPPrefixId{}

// CustomIPPrefixId is a struct representing the Resource ID for a Custom I P Prefix
type CustomIPPrefixId struct {
	SubscriptionId     string
	ResourceGroupName  string
	CustomIPPrefixName string
}

// NewCustomIPPrefixID returns a new CustomIPPrefixId struct
func NewCustomIPPrefixID(subscriptionId string, resourceGroupName string, customIPPrefixName string) CustomIPPrefixId {
	return CustomIPPrefixId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		CustomIPPrefixName: customIPPrefixName,
	}
}

// ParseCustomIPPrefixID parses 'input' into a CustomIPPrefixId
func ParseCustomIPPrefixID(input string) (*CustomIPPrefixId, error) {
	parser := resourceids.NewParserFromResourceIdType(CustomIPPrefixId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CustomIPPrefixId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CustomIPPrefixName, ok = parsed.Parsed["customIPPrefixName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "customIPPrefixName", *parsed)
	}

	return &id, nil
}

// ParseCustomIPPrefixIDInsensitively parses 'input' case-insensitively into a CustomIPPrefixId
// note: this method should only be used for API response data and not user input
func ParseCustomIPPrefixIDInsensitively(input string) (*CustomIPPrefixId, error) {
	parser := resourceids.NewParserFromResourceIdType(CustomIPPrefixId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CustomIPPrefixId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CustomIPPrefixName, ok = parsed.Parsed["customIPPrefixName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "customIPPrefixName", *parsed)
	}

	return &id, nil
}

// ValidateCustomIPPrefixID checks that 'input' can be parsed as a Custom I P Prefix ID
func ValidateCustomIPPrefixID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCustomIPPrefixID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Custom I P Prefix ID
func (id CustomIPPrefixId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/customIPPrefixes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CustomIPPrefixName)
}

// Segments returns a slice of Resource ID Segments which comprise this Custom I P Prefix ID
func (id CustomIPPrefixId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticCustomIPPrefixes", "customIPPrefixes", "customIPPrefixes"),
		resourceids.UserSpecifiedSegment("customIPPrefixName", "customIPPrefixValue"),
	}
}

// String returns a human-readable description of this Custom I P Prefix ID
func (id CustomIPPrefixId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Custom I P Prefix Name: %q", id.CustomIPPrefixName),
	}
	return fmt.Sprintf("Custom I P Prefix (%s)", strings.Join(components, "\n"))
}
