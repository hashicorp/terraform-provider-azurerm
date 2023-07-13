package customresourceprovider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ResourceProviderId{}

// ResourceProviderId is a struct representing the Resource ID for a Resource Provider
type ResourceProviderId struct {
	SubscriptionId       string
	ResourceGroupName    string
	ResourceProviderName string
}

// NewResourceProviderID returns a new ResourceProviderId struct
func NewResourceProviderID(subscriptionId string, resourceGroupName string, resourceProviderName string) ResourceProviderId {
	return ResourceProviderId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		ResourceProviderName: resourceProviderName,
	}
}

// ParseResourceProviderID parses 'input' into a ResourceProviderId
func ParseResourceProviderID(input string) (*ResourceProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(ResourceProviderId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ResourceProviderId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ResourceProviderName, ok = parsed.Parsed["resourceProviderName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceProviderName", *parsed)
	}

	return &id, nil
}

// ParseResourceProviderIDInsensitively parses 'input' case-insensitively into a ResourceProviderId
// note: this method should only be used for API response data and not user input
func ParseResourceProviderIDInsensitively(input string) (*ResourceProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(ResourceProviderId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ResourceProviderId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ResourceProviderName, ok = parsed.Parsed["resourceProviderName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceProviderName", *parsed)
	}

	return &id, nil
}

// ValidateResourceProviderID checks that 'input' can be parsed as a Resource Provider ID
func ValidateResourceProviderID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseResourceProviderID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Resource Provider ID
func (id ResourceProviderId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.CustomProviders/resourceProviders/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceProviderName)
}

// Segments returns a slice of Resource ID Segments which comprise this Resource Provider ID
func (id ResourceProviderId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCustomProviders", "Microsoft.CustomProviders", "Microsoft.CustomProviders"),
		resourceids.StaticSegment("staticResourceProviders", "resourceProviders", "resourceProviders"),
		resourceids.UserSpecifiedSegment("resourceProviderName", "resourceProviderValue"),
	}
}

// String returns a human-readable description of this Resource Provider ID
func (id ResourceProviderId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Provider Name: %q", id.ResourceProviderName),
	}
	return fmt.Sprintf("Resource Provider (%s)", strings.Join(components, "\n"))
}
