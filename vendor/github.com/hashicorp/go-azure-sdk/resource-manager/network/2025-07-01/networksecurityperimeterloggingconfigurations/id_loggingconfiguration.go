package networksecurityperimeterloggingconfigurations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&LoggingConfigurationId{})
}

var _ resourceids.ResourceId = &LoggingConfigurationId{}

// LoggingConfigurationId is a struct representing the Resource ID for a Logging Configuration
type LoggingConfigurationId struct {
	SubscriptionId               string
	ResourceGroupName            string
	NetworkSecurityPerimeterName string
	LoggingConfigurationName     string
}

// NewLoggingConfigurationID returns a new LoggingConfigurationId struct
func NewLoggingConfigurationID(subscriptionId string, resourceGroupName string, networkSecurityPerimeterName string, loggingConfigurationName string) LoggingConfigurationId {
	return LoggingConfigurationId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		NetworkSecurityPerimeterName: networkSecurityPerimeterName,
		LoggingConfigurationName:     loggingConfigurationName,
	}
}

// ParseLoggingConfigurationID parses 'input' into a LoggingConfigurationId
func ParseLoggingConfigurationID(input string) (*LoggingConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LoggingConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LoggingConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseLoggingConfigurationIDInsensitively parses 'input' case-insensitively into a LoggingConfigurationId
// note: this method should only be used for API response data and not user input
func ParseLoggingConfigurationIDInsensitively(input string) (*LoggingConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&LoggingConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := LoggingConfigurationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *LoggingConfigurationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkSecurityPerimeterName, ok = input.Parsed["networkSecurityPerimeterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkSecurityPerimeterName", input)
	}

	if id.LoggingConfigurationName, ok = input.Parsed["loggingConfigurationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "loggingConfigurationName", input)
	}

	return nil
}

// ValidateLoggingConfigurationID checks that 'input' can be parsed as a Logging Configuration ID
func ValidateLoggingConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLoggingConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Logging Configuration ID
func (id LoggingConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkSecurityPerimeters/%s/loggingConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkSecurityPerimeterName, id.LoggingConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Logging Configuration ID
func (id LoggingConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkSecurityPerimeters", "networkSecurityPerimeters", "networkSecurityPerimeters"),
		resourceids.UserSpecifiedSegment("networkSecurityPerimeterName", "networkSecurityPerimeterName"),
		resourceids.StaticSegment("staticLoggingConfigurations", "loggingConfigurations", "loggingConfigurations"),
		resourceids.UserSpecifiedSegment("loggingConfigurationName", "loggingConfigurationName"),
	}
}

// String returns a human-readable description of this Logging Configuration ID
func (id LoggingConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Security Perimeter Name: %q", id.NetworkSecurityPerimeterName),
		fmt.Sprintf("Logging Configuration Name: %q", id.LoggingConfigurationName),
	}
	return fmt.Sprintf("Logging Configuration (%s)", strings.Join(components, "\n"))
}
