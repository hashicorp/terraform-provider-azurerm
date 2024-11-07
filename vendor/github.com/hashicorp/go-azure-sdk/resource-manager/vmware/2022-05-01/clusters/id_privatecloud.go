package clusters

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PrivateCloudId{})
}

var _ resourceids.ResourceId = &PrivateCloudId{}

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
	parser := resourceids.NewParserFromResourceIdType(&PrivateCloudId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateCloudId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePrivateCloudIDInsensitively parses 'input' case-insensitively into a PrivateCloudId
// note: this method should only be used for API response data and not user input
func ParsePrivateCloudIDInsensitively(input string) (*PrivateCloudId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrivateCloudId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateCloudId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PrivateCloudId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.PrivateCloudName, ok = input.Parsed["privateCloudName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "privateCloudName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("privateCloudName", "privateCloudName"),
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
