package networkinterfaces

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = RoleInstanceId{}

// RoleInstanceId is a struct representing the Resource ID for a Role Instance
type RoleInstanceId struct {
	SubscriptionId    string
	ResourceGroupName string
	CloudServiceName  string
	RoleInstanceName  string
}

// NewRoleInstanceID returns a new RoleInstanceId struct
func NewRoleInstanceID(subscriptionId string, resourceGroupName string, cloudServiceName string, roleInstanceName string) RoleInstanceId {
	return RoleInstanceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		CloudServiceName:  cloudServiceName,
		RoleInstanceName:  roleInstanceName,
	}
}

// ParseRoleInstanceID parses 'input' into a RoleInstanceId
func ParseRoleInstanceID(input string) (*RoleInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(RoleInstanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RoleInstanceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CloudServiceName, ok = parsed.Parsed["cloudServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "cloudServiceName", *parsed)
	}

	if id.RoleInstanceName, ok = parsed.Parsed["roleInstanceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "roleInstanceName", *parsed)
	}

	return &id, nil
}

// ParseRoleInstanceIDInsensitively parses 'input' case-insensitively into a RoleInstanceId
// note: this method should only be used for API response data and not user input
func ParseRoleInstanceIDInsensitively(input string) (*RoleInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(RoleInstanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RoleInstanceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CloudServiceName, ok = parsed.Parsed["cloudServiceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "cloudServiceName", *parsed)
	}

	if id.RoleInstanceName, ok = parsed.Parsed["roleInstanceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "roleInstanceName", *parsed)
	}

	return &id, nil
}

// ValidateRoleInstanceID checks that 'input' can be parsed as a Role Instance ID
func ValidateRoleInstanceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRoleInstanceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Role Instance ID
func (id RoleInstanceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/cloudServices/%s/roleInstances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CloudServiceName, id.RoleInstanceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Role Instance ID
func (id RoleInstanceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticCloudServices", "cloudServices", "cloudServices"),
		resourceids.UserSpecifiedSegment("cloudServiceName", "cloudServiceValue"),
		resourceids.StaticSegment("staticRoleInstances", "roleInstances", "roleInstances"),
		resourceids.UserSpecifiedSegment("roleInstanceName", "roleInstanceValue"),
	}
}

// String returns a human-readable description of this Role Instance ID
func (id RoleInstanceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cloud Service Name: %q", id.CloudServiceName),
		fmt.Sprintf("Role Instance Name: %q", id.RoleInstanceName),
	}
	return fmt.Sprintf("Role Instance (%s)", strings.Join(components, "\n"))
}
