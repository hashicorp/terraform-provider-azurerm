package networkinterfaces

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RoleInstanceNetworkInterfaceId{})
}

var _ resourceids.ResourceId = &RoleInstanceNetworkInterfaceId{}

// RoleInstanceNetworkInterfaceId is a struct representing the Resource ID for a Role Instance Network Interface
type RoleInstanceNetworkInterfaceId struct {
	SubscriptionId       string
	ResourceGroupName    string
	CloudServiceName     string
	RoleInstanceName     string
	NetworkInterfaceName string
}

// NewRoleInstanceNetworkInterfaceID returns a new RoleInstanceNetworkInterfaceId struct
func NewRoleInstanceNetworkInterfaceID(subscriptionId string, resourceGroupName string, cloudServiceName string, roleInstanceName string, networkInterfaceName string) RoleInstanceNetworkInterfaceId {
	return RoleInstanceNetworkInterfaceId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		CloudServiceName:     cloudServiceName,
		RoleInstanceName:     roleInstanceName,
		NetworkInterfaceName: networkInterfaceName,
	}
}

// ParseRoleInstanceNetworkInterfaceID parses 'input' into a RoleInstanceNetworkInterfaceId
func ParseRoleInstanceNetworkInterfaceID(input string) (*RoleInstanceNetworkInterfaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RoleInstanceNetworkInterfaceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RoleInstanceNetworkInterfaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRoleInstanceNetworkInterfaceIDInsensitively parses 'input' case-insensitively into a RoleInstanceNetworkInterfaceId
// note: this method should only be used for API response data and not user input
func ParseRoleInstanceNetworkInterfaceIDInsensitively(input string) (*RoleInstanceNetworkInterfaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RoleInstanceNetworkInterfaceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RoleInstanceNetworkInterfaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RoleInstanceNetworkInterfaceId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.RoleInstanceName, ok = input.Parsed["roleInstanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "roleInstanceName", input)
	}

	if id.NetworkInterfaceName, ok = input.Parsed["networkInterfaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkInterfaceName", input)
	}

	return nil
}

// ValidateRoleInstanceNetworkInterfaceID checks that 'input' can be parsed as a Role Instance Network Interface ID
func ValidateRoleInstanceNetworkInterfaceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRoleInstanceNetworkInterfaceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Role Instance Network Interface ID
func (id RoleInstanceNetworkInterfaceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/cloudServices/%s/roleInstances/%s/networkInterfaces/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CloudServiceName, id.RoleInstanceName, id.NetworkInterfaceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Role Instance Network Interface ID
func (id RoleInstanceNetworkInterfaceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticCloudServices", "cloudServices", "cloudServices"),
		resourceids.UserSpecifiedSegment("cloudServiceName", "cloudServiceName"),
		resourceids.StaticSegment("staticRoleInstances", "roleInstances", "roleInstances"),
		resourceids.UserSpecifiedSegment("roleInstanceName", "roleInstanceName"),
		resourceids.StaticSegment("staticNetworkInterfaces", "networkInterfaces", "networkInterfaces"),
		resourceids.UserSpecifiedSegment("networkInterfaceName", "networkInterfaceName"),
	}
}

// String returns a human-readable description of this Role Instance Network Interface ID
func (id RoleInstanceNetworkInterfaceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cloud Service Name: %q", id.CloudServiceName),
		fmt.Sprintf("Role Instance Name: %q", id.RoleInstanceName),
		fmt.Sprintf("Network Interface Name: %q", id.NetworkInterfaceName),
	}
	return fmt.Sprintf("Role Instance Network Interface (%s)", strings.Join(components, "\n"))
}
