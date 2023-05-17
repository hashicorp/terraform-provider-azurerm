package configurationassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ResourceGroupProviderId{}

// ResourceGroupProviderId is a struct representing the Resource ID for a Resource Group Provider
type ResourceGroupProviderId struct {
	SubscriptionId     string
	ResourceGroupName  string
	ProviderName       string
	ResourceParentType string
	ResourceParentName string
	ResourceType       string
	ResourceName       string
}

// NewResourceGroupProviderID returns a new ResourceGroupProviderId struct
func NewResourceGroupProviderID(subscriptionId string, resourceGroupName string, providerName string, resourceParentType string, resourceParentName string, resourceType string, resourceName string) ResourceGroupProviderId {
	return ResourceGroupProviderId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		ProviderName:       providerName,
		ResourceParentType: resourceParentType,
		ResourceParentName: resourceParentName,
		ResourceType:       resourceType,
		ResourceName:       resourceName,
	}
}

// ParseResourceGroupProviderID parses 'input' into a ResourceGroupProviderId
func ParseResourceGroupProviderID(input string) (*ResourceGroupProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(ResourceGroupProviderId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ResourceGroupProviderId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "providerName", *parsed)
	}

	if id.ResourceParentType, ok = parsed.Parsed["resourceParentType"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceParentType", *parsed)
	}

	if id.ResourceParentName, ok = parsed.Parsed["resourceParentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceParentName", *parsed)
	}

	if id.ResourceType, ok = parsed.Parsed["resourceType"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceType", *parsed)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceName", *parsed)
	}

	return &id, nil
}

// ParseResourceGroupProviderIDInsensitively parses 'input' case-insensitively into a ResourceGroupProviderId
// note: this method should only be used for API response data and not user input
func ParseResourceGroupProviderIDInsensitively(input string) (*ResourceGroupProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(ResourceGroupProviderId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ResourceGroupProviderId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "providerName", *parsed)
	}

	if id.ResourceParentType, ok = parsed.Parsed["resourceParentType"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceParentType", *parsed)
	}

	if id.ResourceParentName, ok = parsed.Parsed["resourceParentName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceParentName", *parsed)
	}

	if id.ResourceType, ok = parsed.Parsed["resourceType"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceType", *parsed)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceName", *parsed)
	}

	return &id, nil
}

// ValidateResourceGroupProviderID checks that 'input' can be parsed as a Resource Group Provider ID
func ValidateResourceGroupProviderID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseResourceGroupProviderID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Resource Group Provider ID
func (id ResourceGroupProviderId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/%s/%s/%s/%s/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProviderName, id.ResourceParentType, id.ResourceParentName, id.ResourceType, id.ResourceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Resource Group Provider ID
func (id ResourceGroupProviderId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.UserSpecifiedSegment("providerName", "providerValue"),
		resourceids.UserSpecifiedSegment("resourceParentType", "resourceParentTypeValue"),
		resourceids.UserSpecifiedSegment("resourceParentName", "resourceParentValue"),
		resourceids.UserSpecifiedSegment("resourceType", "resourceTypeValue"),
		resourceids.UserSpecifiedSegment("resourceName", "resourceValue"),
	}
}

// String returns a human-readable description of this Resource Group Provider ID
func (id ResourceGroupProviderId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Provider Name: %q", id.ProviderName),
		fmt.Sprintf("Resource Parent Type: %q", id.ResourceParentType),
		fmt.Sprintf("Resource Parent Name: %q", id.ResourceParentName),
		fmt.Sprintf("Resource Type: %q", id.ResourceType),
		fmt.Sprintf("Resource Name: %q", id.ResourceName),
	}
	return fmt.Sprintf("Resource Group Provider (%s)", strings.Join(components, "\n"))
}
