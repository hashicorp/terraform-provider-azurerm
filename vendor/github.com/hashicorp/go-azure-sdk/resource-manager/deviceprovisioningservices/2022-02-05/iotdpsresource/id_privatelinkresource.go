package iotdpsresource

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PrivateLinkResourceId{})
}

var _ resourceids.ResourceId = &PrivateLinkResourceId{}

// PrivateLinkResourceId is a struct representing the Resource ID for a Private Link Resource
type PrivateLinkResourceId struct {
	SubscriptionId          string
	ResourceGroupName       string
	ProvisioningServiceName string
	GroupId                 string
}

// NewPrivateLinkResourceID returns a new PrivateLinkResourceId struct
func NewPrivateLinkResourceID(subscriptionId string, resourceGroupName string, provisioningServiceName string, groupId string) PrivateLinkResourceId {
	return PrivateLinkResourceId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		ProvisioningServiceName: provisioningServiceName,
		GroupId:                 groupId,
	}
}

// ParsePrivateLinkResourceID parses 'input' into a PrivateLinkResourceId
func ParsePrivateLinkResourceID(input string) (*PrivateLinkResourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrivateLinkResourceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateLinkResourceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePrivateLinkResourceIDInsensitively parses 'input' case-insensitively into a PrivateLinkResourceId
// note: this method should only be used for API response data and not user input
func ParsePrivateLinkResourceIDInsensitively(input string) (*PrivateLinkResourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PrivateLinkResourceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PrivateLinkResourceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PrivateLinkResourceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ProvisioningServiceName, ok = input.Parsed["provisioningServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "provisioningServiceName", input)
	}

	if id.GroupId, ok = input.Parsed["groupId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "groupId", input)
	}

	return nil
}

// ValidatePrivateLinkResourceID checks that 'input' can be parsed as a Private Link Resource ID
func ValidatePrivateLinkResourceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePrivateLinkResourceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Private Link Resource ID
func (id PrivateLinkResourceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Devices/provisioningServices/%s/privateLinkResources/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProvisioningServiceName, id.GroupId)
}

// Segments returns a slice of Resource ID Segments which comprise this Private Link Resource ID
func (id PrivateLinkResourceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevices", "Microsoft.Devices", "Microsoft.Devices"),
		resourceids.StaticSegment("staticProvisioningServices", "provisioningServices", "provisioningServices"),
		resourceids.UserSpecifiedSegment("provisioningServiceName", "provisioningServiceName"),
		resourceids.StaticSegment("staticPrivateLinkResources", "privateLinkResources", "privateLinkResources"),
		resourceids.UserSpecifiedSegment("groupId", "groupId"),
	}
}

// String returns a human-readable description of this Private Link Resource ID
func (id PrivateLinkResourceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Provisioning Service Name: %q", id.ProvisioningServiceName),
		fmt.Sprintf("Group: %q", id.GroupId),
	}
	return fmt.Sprintf("Private Link Resource (%s)", strings.Join(components, "\n"))
}
