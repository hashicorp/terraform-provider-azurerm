package vipswap

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CloudServiceId{})
}

var _ resourceids.ResourceId = &CloudServiceId{}

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
	parser := resourceids.NewParserFromResourceIdType(&CloudServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CloudServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCloudServiceIDInsensitively parses 'input' case-insensitively into a CloudServiceId
// note: this method should only be used for API response data and not user input
func ParseCloudServiceIDInsensitively(input string) (*CloudServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CloudServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CloudServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CloudServiceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CloudServiceName, ok = input.Parsed["cloudServiceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "cloudServiceName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("resourceGroupName", "resourceGroupName"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticCloudServices", "cloudServices", "cloudServices"),
		resourceids.UserSpecifiedSegment("cloudServiceName", "cloudServiceName"),
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
