package cloudservicepublicipaddresses

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ProviderCloudServiceId{}

// ProviderCloudServiceId is a struct representing the Resource ID for a Provider Cloud Service
type ProviderCloudServiceId struct {
	SubscriptionId    string
	ResourceGroupName string
	CloudServiceName  string
}

// NewProviderCloudServiceID returns a new ProviderCloudServiceId struct
func NewProviderCloudServiceID(subscriptionId string, resourceGroupName string, cloudServiceName string) ProviderCloudServiceId {
	return ProviderCloudServiceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		CloudServiceName:  cloudServiceName,
	}
}

// ParseProviderCloudServiceID parses 'input' into a ProviderCloudServiceId
func ParseProviderCloudServiceID(input string) (*ProviderCloudServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderCloudServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderCloudServiceId{}

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

// ParseProviderCloudServiceIDInsensitively parses 'input' case-insensitively into a ProviderCloudServiceId
// note: this method should only be used for API response data and not user input
func ParseProviderCloudServiceIDInsensitively(input string) (*ProviderCloudServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderCloudServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderCloudServiceId{}

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

// ValidateProviderCloudServiceID checks that 'input' can be parsed as a Provider Cloud Service ID
func ValidateProviderCloudServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderCloudServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Cloud Service ID
func (id ProviderCloudServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/cloudServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CloudServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Cloud Service ID
func (id ProviderCloudServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticCloudServices", "cloudServices", "cloudServices"),
		resourceids.UserSpecifiedSegment("cloudServiceName", "cloudServiceValue"),
	}
}

// String returns a human-readable description of this Provider Cloud Service ID
func (id ProviderCloudServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cloud Service Name: %q", id.CloudServiceName),
	}
	return fmt.Sprintf("Provider Cloud Service (%s)", strings.Join(components, "\n"))
}
