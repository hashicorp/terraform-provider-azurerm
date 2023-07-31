package vipswap

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = CloudServiceId{}

// CloudServiceId is a struct representing the Resource ID for a Cloud Service
type CloudServiceId struct {
	SubscriptionId    string
	ResourceGroupName string
	CloudServiceName  string
}

// NewCloudServiceID returns a new CloudServiceId struct
func NewCloudServiceID(subscriptionId string, resourceGroupName string, cloudServiceName string) CloudServiceId {
	return CloudServiceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		CloudServiceName:  cloudServiceName,
	}
}

// ParseCloudServiceID parses 'input' into a CloudServiceId
func ParseCloudServiceID(input string) (*CloudServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(CloudServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CloudServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CloudServiceName, ok = parsed.Parsed["cloudServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "cloudServiceName", *parsed)
	}

	return &id, nil
}

// ParseCloudServiceIDInsensitively parses 'input' case-insensitively into a CloudServiceId
// note: this method should only be used for API response data and not user input
func ParseCloudServiceIDInsensitively(input string) (*CloudServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(CloudServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CloudServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CloudServiceName, ok = parsed.Parsed["cloudServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "cloudServiceName", *parsed)
	}

	return &id, nil
}

// ValidateCloudServiceID checks that 'input' can be parsed as a Cloud Service ID
func ValidateCloudServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCloudServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Cloud Service ID
func (id CloudServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/cloudServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CloudServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Cloud Service ID
func (id CloudServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.UserSpecifiedSegment("resourceGroupName", "resourceGroupValue"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticCloudServices", "cloudServices", "cloudServices"),
		resourceids.UserSpecifiedSegment("cloudServiceName", "cloudServiceValue"),
	}
}

// String returns a human-readable description of this Cloud Service ID
func (id CloudServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cloud Service Name: %q", id.CloudServiceName),
	}
	return fmt.Sprintf("Cloud Service (%s)", strings.Join(components, "\n"))
}
