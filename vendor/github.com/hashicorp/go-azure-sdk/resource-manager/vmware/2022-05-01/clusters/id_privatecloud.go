package clusters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PrivateCloudId{}

// PrivateCloudId is a struct representing the Resource ID for a Private Cloud
type PrivateCloudId struct {
	SubscriptionId    string
	ResourceGroupName string
	PrivateCloudName  string
}

// NewPrivateCloudID returns a new PrivateCloudId struct
func NewPrivateCloudID(subscriptionId string, resourceGroupName string, privateCloudName string) PrivateCloudId {
	return PrivateCloudId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		PrivateCloudName:  privateCloudName,
	}
}

// ParsePrivateCloudID parses 'input' into a PrivateCloudId
func ParsePrivateCloudID(input string) (*PrivateCloudId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateCloudId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateCloudId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PrivateCloudName, ok = parsed.Parsed["privateCloudName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateCloudName", *parsed)
	}

	return &id, nil
}

// ParsePrivateCloudIDInsensitively parses 'input' case-insensitively into a PrivateCloudId
// note: this method should only be used for API response data and not user input
func ParsePrivateCloudIDInsensitively(input string) (*PrivateCloudId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateCloudId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateCloudId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PrivateCloudName, ok = parsed.Parsed["privateCloudName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateCloudName", *parsed)
	}

	return &id, nil
}

// ValidatePrivateCloudID checks that 'input' can be parsed as a Private Cloud ID
func ValidatePrivateCloudID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePrivateCloudID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Private Cloud ID
func (id PrivateCloudId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AVS/privateClouds/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateCloudName)
}

// Segments returns a slice of Resource ID Segments which comprise this Private Cloud ID
func (id PrivateCloudId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAVS", "Microsoft.AVS", "Microsoft.AVS"),
		resourceids.StaticSegment("staticPrivateClouds", "privateClouds", "privateClouds"),
		resourceids.UserSpecifiedSegment("privateCloudName", "privateCloudValue"),
	}
}

// String returns a human-readable description of this Private Cloud ID
func (id PrivateCloudId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Cloud Name: %q", id.PrivateCloudName),
	}
	return fmt.Sprintf("Private Cloud (%s)", strings.Join(components, "\n"))
}
