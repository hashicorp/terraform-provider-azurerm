package appserviceenvironments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DefaultInstanceId{})
}

var _ resourceids.ResourceId = &DefaultInstanceId{}

// DefaultInstanceId is a struct representing the Resource ID for a Default Instance
type DefaultInstanceId struct {
	SubscriptionId         string
	ResourceGroupName      string
	HostingEnvironmentName string
	InstanceName           string
}

// NewDefaultInstanceID returns a new DefaultInstanceId struct
func NewDefaultInstanceID(subscriptionId string, resourceGroupName string, hostingEnvironmentName string, instanceName string) DefaultInstanceId {
	return DefaultInstanceId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		HostingEnvironmentName: hostingEnvironmentName,
		InstanceName:           instanceName,
	}
}

// ParseDefaultInstanceID parses 'input' into a DefaultInstanceId
func ParseDefaultInstanceID(input string) (*DefaultInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DefaultInstanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DefaultInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDefaultInstanceIDInsensitively parses 'input' case-insensitively into a DefaultInstanceId
// note: this method should only be used for API response data and not user input
func ParseDefaultInstanceIDInsensitively(input string) (*DefaultInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DefaultInstanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DefaultInstanceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DefaultInstanceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.HostingEnvironmentName, ok = input.Parsed["hostingEnvironmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "hostingEnvironmentName", input)
	}

	if id.InstanceName, ok = input.Parsed["instanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "instanceName", input)
	}

	return nil
}

// ValidateDefaultInstanceID checks that 'input' can be parsed as a Default Instance ID
func ValidateDefaultInstanceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDefaultInstanceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Default Instance ID
func (id DefaultInstanceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/hostingEnvironments/%s/multiRolePools/default/instances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HostingEnvironmentName, id.InstanceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Default Instance ID
func (id DefaultInstanceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftWeb", "Microsoft.Web", "Microsoft.Web"),
		resourceids.StaticSegment("staticHostingEnvironments", "hostingEnvironments", "hostingEnvironments"),
		resourceids.UserSpecifiedSegment("hostingEnvironmentName", "hostingEnvironmentName"),
		resourceids.StaticSegment("staticMultiRolePools", "multiRolePools", "multiRolePools"),
		resourceids.StaticSegment("staticDefault", "default", "default"),
		resourceids.StaticSegment("staticInstances", "instances", "instances"),
		resourceids.UserSpecifiedSegment("instanceName", "instanceName"),
	}
}

// String returns a human-readable description of this Default Instance ID
func (id DefaultInstanceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Hosting Environment Name: %q", id.HostingEnvironmentName),
		fmt.Sprintf("Instance Name: %q", id.InstanceName),
	}
	return fmt.Sprintf("Default Instance (%s)", strings.Join(components, "\n"))
}
