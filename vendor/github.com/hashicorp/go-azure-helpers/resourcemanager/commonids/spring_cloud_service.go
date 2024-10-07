// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &SpringCloudServiceId{}

// SpringCloudServiceId is a struct representing the Resource ID for a Spring Cloud Service
type SpringCloudServiceId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
}

// NewSpringCloudServiceID returns a new SpringCloudServiceId struct
func NewSpringCloudServiceID(subscriptionId string, resourceGroupName string, serviceName string) SpringCloudServiceId {
	return SpringCloudServiceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
	}
}

// ParseSpringCloudServiceID parses 'input' into a SpringCloudServiceId
func ParseSpringCloudServiceID(input string) (*SpringCloudServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SpringCloudServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SpringCloudServiceId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSpringCloudServiceIDInsensitively parses 'input' case-insensitively into a SpringCloudServiceId
// note: this method should only be used for API response data and not user input
func ParseSpringCloudServiceIDInsensitively(input string) (*SpringCloudServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SpringCloudServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SpringCloudServiceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SpringCloudServiceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ServiceName, ok = input.Parsed["serviceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceName", input)
	}

	return nil
}

// ValidateSpringCloudServiceID checks that 'input' can be parsed as a Spring Cloud Service ID
func ValidateSpringCloudServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSpringCloudServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Spring Cloud Service ID
func (id SpringCloudServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Spring Cloud Service ID
func (id SpringCloudServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
	}
}

// String returns a human-readable description of this Spring Cloud Service ID
func (id SpringCloudServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
	}
	return fmt.Sprintf("Spring Cloud Service (%s)", strings.Join(components, "\n"))
}
