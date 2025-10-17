package virtualwans

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ConfigurationPolicyGroupId{})
}

var _ resourceids.ResourceId = &ConfigurationPolicyGroupId{}

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
	parser := resourceids.NewParserFromResourceIdType(&ConfigurationPolicyGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConfigurationPolicyGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConfigurationPolicyGroupIDInsensitively parses 'input' case-insensitively into a ConfigurationPolicyGroupId
// note: this method should only be used for API response data and not user input
func ParseConfigurationPolicyGroupIDInsensitively(input string) (*ConfigurationPolicyGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConfigurationPolicyGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConfigurationPolicyGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConfigurationPolicyGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.VpnServerConfigurationName, ok = input.Parsed["vpnServerConfigurationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "vpnServerConfigurationName", input)
	}

	if id.ConfigurationPolicyGroupName, ok = input.Parsed["configurationPolicyGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "configurationPolicyGroupName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("vpnServerConfigurationName", "vpnServerConfigurationName"),
		resourceids.StaticSegment("staticConfigurationPolicyGroups", "configurationPolicyGroups", "configurationPolicyGroups"),
		resourceids.UserSpecifiedSegment("configurationPolicyGroupName", "configurationPolicyGroupName"),
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
