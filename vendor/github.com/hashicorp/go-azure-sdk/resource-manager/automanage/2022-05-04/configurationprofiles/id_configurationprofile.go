package configurationprofiles

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ConfigurationProfileId{})
}

var _ resourceids.ResourceId = &ConfigurationProfileId{}

// ConfigurationProfileId is a struct representing the Resource ID for a Configuration Profile
type ConfigurationProfileId struct {
	SubscriptionId           string
	ResourceGroupName        string
	ConfigurationProfileName string
}

// NewConfigurationProfileID returns a new ConfigurationProfileId struct
func NewConfigurationProfileID(subscriptionId string, resourceGroupName string, configurationProfileName string) ConfigurationProfileId {
	return ConfigurationProfileId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		ConfigurationProfileName: configurationProfileName,
	}
}

// ParseConfigurationProfileID parses 'input' into a ConfigurationProfileId
func ParseConfigurationProfileID(input string) (*ConfigurationProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConfigurationProfileId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConfigurationProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConfigurationProfileIDInsensitively parses 'input' case-insensitively into a ConfigurationProfileId
// note: this method should only be used for API response data and not user input
func ParseConfigurationProfileIDInsensitively(input string) (*ConfigurationProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConfigurationProfileId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConfigurationProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConfigurationProfileId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ConfigurationProfileName, ok = input.Parsed["configurationProfileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "configurationProfileName", input)
	}

	return nil
}

// ValidateConfigurationProfileID checks that 'input' can be parsed as a Configuration Profile ID
func ValidateConfigurationProfileID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConfigurationProfileID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Configuration Profile ID
func (id ConfigurationProfileId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AutoManage/configurationProfiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ConfigurationProfileName)
}

// Segments returns a slice of Resource ID Segments which comprise this Configuration Profile ID
func (id ConfigurationProfileId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutoManage", "Microsoft.AutoManage", "Microsoft.AutoManage"),
		resourceids.StaticSegment("staticConfigurationProfiles", "configurationProfiles", "configurationProfiles"),
		resourceids.UserSpecifiedSegment("configurationProfileName", "configurationProfileName"),
	}
}

// String returns a human-readable description of this Configuration Profile ID
func (id ConfigurationProfileId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Configuration Profile Name: %q", id.ConfigurationProfileName),
	}
	return fmt.Sprintf("Configuration Profile (%s)", strings.Join(components, "\n"))
}
