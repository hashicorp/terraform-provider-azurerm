package virtualwans

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ConfigurationPolicyGroupId{}

// ConfigurationPolicyGroupId is a struct representing the Resource ID for a Configuration Policy Group
type ConfigurationPolicyGroupId struct {
	SubscriptionId               string
	ResourceGroupName            string
	VpnServerConfigurationName   string
	ConfigurationPolicyGroupName string
}

// NewConfigurationPolicyGroupID returns a new ConfigurationPolicyGroupId struct
func NewConfigurationPolicyGroupID(subscriptionId string, resourceGroupName string, vpnServerConfigurationName string, configurationPolicyGroupName string) ConfigurationPolicyGroupId {
	return ConfigurationPolicyGroupId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		VpnServerConfigurationName:   vpnServerConfigurationName,
		ConfigurationPolicyGroupName: configurationPolicyGroupName,
	}
}

// ParseConfigurationPolicyGroupID parses 'input' into a ConfigurationPolicyGroupId
func ParseConfigurationPolicyGroupID(input string) (*ConfigurationPolicyGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConfigurationPolicyGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConfigurationPolicyGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VpnServerConfigurationName, ok = parsed.Parsed["vpnServerConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnServerConfigurationName", *parsed)
	}

	if id.ConfigurationPolicyGroupName, ok = parsed.Parsed["configurationPolicyGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "configurationPolicyGroupName", *parsed)
	}

	return &id, nil
}

// ParseConfigurationPolicyGroupIDInsensitively parses 'input' case-insensitively into a ConfigurationPolicyGroupId
// note: this method should only be used for API response data and not user input
func ParseConfigurationPolicyGroupIDInsensitively(input string) (*ConfigurationPolicyGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConfigurationPolicyGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConfigurationPolicyGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.VpnServerConfigurationName, ok = parsed.Parsed["vpnServerConfigurationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "vpnServerConfigurationName", *parsed)
	}

	if id.ConfigurationPolicyGroupName, ok = parsed.Parsed["configurationPolicyGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "configurationPolicyGroupName", *parsed)
	}

	return &id, nil
}

// ValidateConfigurationPolicyGroupID checks that 'input' can be parsed as a Configuration Policy Group ID
func ValidateConfigurationPolicyGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConfigurationPolicyGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Configuration Policy Group ID
func (id ConfigurationPolicyGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnServerConfigurations/%s/configurationPolicyGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VpnServerConfigurationName, id.ConfigurationPolicyGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Configuration Policy Group ID
func (id ConfigurationPolicyGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVpnServerConfigurations", "vpnServerConfigurations", "vpnServerConfigurations"),
		resourceids.UserSpecifiedSegment("vpnServerConfigurationName", "vpnServerConfigurationValue"),
		resourceids.StaticSegment("staticConfigurationPolicyGroups", "configurationPolicyGroups", "configurationPolicyGroups"),
		resourceids.UserSpecifiedSegment("configurationPolicyGroupName", "configurationPolicyGroupValue"),
	}
}

// String returns a human-readable description of this Configuration Policy Group ID
func (id ConfigurationPolicyGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vpn Server Configuration Name: %q", id.VpnServerConfigurationName),
		fmt.Sprintf("Configuration Policy Group Name: %q", id.ConfigurationPolicyGroupName),
	}
	return fmt.Sprintf("Configuration Policy Group (%s)", strings.Join(components, "\n"))
}
