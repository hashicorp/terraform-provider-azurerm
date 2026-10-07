package extensions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ExtensionId{})
}

var _ resourceids.ResourceId = &ExtensionId{}

// ExtensionId is a struct representing the Resource ID for a Extension
type ExtensionId struct {
	SubscriptionId    string
	ResourceGroupName string
	ClusterName       string
	ArcSettingName    string
	ExtensionName     string
}

// NewExtensionID returns a new ExtensionId struct
func NewExtensionID(subscriptionId string, resourceGroupName string, clusterName string, arcSettingName string, extensionName string) ExtensionId {
	return ExtensionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ClusterName:       clusterName,
		ArcSettingName:    arcSettingName,
		ExtensionName:     extensionName,
	}
}

// ParseExtensionID parses 'input' into a ExtensionId
func ParseExtensionID(input string) (*ExtensionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExtensionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExtensionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseExtensionIDInsensitively parses 'input' case-insensitively into a ExtensionId
// note: this method should only be used for API response data and not user input
func ParseExtensionIDInsensitively(input string) (*ExtensionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ExtensionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ExtensionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ExtensionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ArcSettingName, ok = input.Parsed["arcSettingName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "arcSettingName", input)
	}

	if id.ExtensionName, ok = input.Parsed["extensionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "extensionName", input)
	}

	return nil
}

// ValidateExtensionID checks that 'input' can be parsed as a Extension ID
func ValidateExtensionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseExtensionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Extension ID
func (id ExtensionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AzureStackHCI/clusters/%s/arcSettings/%s/extensions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.ArcSettingName, id.ExtensionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Extension ID
func (id ExtensionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAzureStackHCI", "Microsoft.AzureStackHCI", "Microsoft.AzureStackHCI"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterName"),
		resourceids.StaticSegment("staticArcSettings", "arcSettings", "arcSettings"),
		resourceids.UserSpecifiedSegment("arcSettingName", "arcSettingName"),
		resourceids.StaticSegment("staticExtensions", "extensions", "extensions"),
		resourceids.UserSpecifiedSegment("extensionName", "extensionName"),
	}
}

// String returns a human-readable description of this Extension ID
func (id ExtensionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Arc Setting Name: %q", id.ArcSettingName),
		fmt.Sprintf("Extension Name: %q", id.ExtensionName),
	}
	return fmt.Sprintf("Extension (%s)", strings.Join(components, "\n"))
}
