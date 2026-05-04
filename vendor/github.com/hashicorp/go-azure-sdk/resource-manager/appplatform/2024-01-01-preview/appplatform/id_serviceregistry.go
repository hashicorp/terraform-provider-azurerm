package appplatform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ServiceRegistryId{})
}

var _ resourceids.ResourceId = &ServiceRegistryId{}

// ServiceRegistryId is a struct representing the Resource ID for a Service Registry
type ServiceRegistryId struct {
	SubscriptionId      string
	ResourceGroupName   string
	SpringName          string
	ServiceRegistryName string
}

// NewServiceRegistryID returns a new ServiceRegistryId struct
func NewServiceRegistryID(subscriptionId string, resourceGroupName string, springName string, serviceRegistryName string) ServiceRegistryId {
	return ServiceRegistryId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		SpringName:          springName,
		ServiceRegistryName: serviceRegistryName,
	}
}

// ParseServiceRegistryID parses 'input' into a ServiceRegistryId
func ParseServiceRegistryID(input string) (*ServiceRegistryId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ServiceRegistryId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ServiceRegistryId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseServiceRegistryIDInsensitively parses 'input' case-insensitively into a ServiceRegistryId
// note: this method should only be used for API response data and not user input
func ParseServiceRegistryIDInsensitively(input string) (*ServiceRegistryId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ServiceRegistryId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ServiceRegistryId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ServiceRegistryId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.SpringName, ok = input.Parsed["springName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "springName", input)
	}

	if id.ServiceRegistryName, ok = input.Parsed["serviceRegistryName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "serviceRegistryName", input)
	}

	return nil
}

// ValidateServiceRegistryID checks that 'input' can be parsed as a Service Registry ID
func ValidateServiceRegistryID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseServiceRegistryID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Service Registry ID
func (id ServiceRegistryId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/spring/%s/serviceRegistries/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.ServiceRegistryName)
}

// Segments returns a slice of Resource ID Segments which comprise this Service Registry ID
func (id ServiceRegistryId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAppPlatform", "Microsoft.AppPlatform", "Microsoft.AppPlatform"),
		resourceids.StaticSegment("staticSpring", "spring", "spring"),
		resourceids.UserSpecifiedSegment("springName", "springName"),
		resourceids.StaticSegment("staticServiceRegistries", "serviceRegistries", "serviceRegistries"),
		resourceids.UserSpecifiedSegment("serviceRegistryName", "serviceRegistryName"),
	}
}

// String returns a human-readable description of this Service Registry ID
func (id ServiceRegistryId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Spring Name: %q", id.SpringName),
		fmt.Sprintf("Service Registry Name: %q", id.ServiceRegistryName),
	}
	return fmt.Sprintf("Service Registry (%s)", strings.Join(components, "\n"))
}
