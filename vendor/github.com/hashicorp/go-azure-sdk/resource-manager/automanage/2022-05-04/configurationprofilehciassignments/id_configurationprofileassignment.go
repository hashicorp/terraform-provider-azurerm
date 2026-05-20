package configurationprofilehciassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ConfigurationProfileAssignmentId{})
}

var _ resourceids.ResourceId = &ConfigurationProfileAssignmentId{}

// ConfigurationProfileAssignmentId is a struct representing the Resource ID for a Configuration Profile Assignment
type ConfigurationProfileAssignmentId struct {
	SubscriptionId                     string
	ResourceGroupName                  string
	ClusterName                        string
	ConfigurationProfileAssignmentName string
}

// NewConfigurationProfileAssignmentID returns a new ConfigurationProfileAssignmentId struct
func NewConfigurationProfileAssignmentID(subscriptionId string, resourceGroupName string, clusterName string, configurationProfileAssignmentName string) ConfigurationProfileAssignmentId {
	return ConfigurationProfileAssignmentId{
		SubscriptionId:                     subscriptionId,
		ResourceGroupName:                  resourceGroupName,
		ClusterName:                        clusterName,
		ConfigurationProfileAssignmentName: configurationProfileAssignmentName,
	}
}

// ParseConfigurationProfileAssignmentID parses 'input' into a ConfigurationProfileAssignmentId
func ParseConfigurationProfileAssignmentID(input string) (*ConfigurationProfileAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConfigurationProfileAssignmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConfigurationProfileAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConfigurationProfileAssignmentIDInsensitively parses 'input' case-insensitively into a ConfigurationProfileAssignmentId
// note: this method should only be used for API response data and not user input
func ParseConfigurationProfileAssignmentIDInsensitively(input string) (*ConfigurationProfileAssignmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConfigurationProfileAssignmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConfigurationProfileAssignmentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConfigurationProfileAssignmentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ClusterName, ok = input.Parsed["clusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "clusterName", input)
	}

	if id.ConfigurationProfileAssignmentName, ok = input.Parsed["configurationProfileAssignmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "configurationProfileAssignmentName", input)
	}

	return nil
}

// ValidateConfigurationProfileAssignmentID checks that 'input' can be parsed as a Configuration Profile Assignment ID
func ValidateConfigurationProfileAssignmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConfigurationProfileAssignmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Configuration Profile Assignment ID
func (id ConfigurationProfileAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AzureStackHCI/clusters/%s/providers/Microsoft.AutoManage/configurationProfileAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.ConfigurationProfileAssignmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Configuration Profile Assignment ID
func (id ConfigurationProfileAssignmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAzureStackHCI", "Microsoft.AzureStackHCI", "Microsoft.AzureStackHCI"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterName"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAutoManage", "Microsoft.AutoManage", "Microsoft.AutoManage"),
		resourceids.StaticSegment("staticConfigurationProfileAssignments", "configurationProfileAssignments", "configurationProfileAssignments"),
		resourceids.UserSpecifiedSegment("configurationProfileAssignmentName", "configurationProfileAssignmentName"),
	}
}

// String returns a human-readable description of this Configuration Profile Assignment ID
func (id ConfigurationProfileAssignmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Configuration Profile Assignment Name: %q", id.ConfigurationProfileAssignmentName),
	}
	return fmt.Sprintf("Configuration Profile Assignment (%s)", strings.Join(components, "\n"))
}
