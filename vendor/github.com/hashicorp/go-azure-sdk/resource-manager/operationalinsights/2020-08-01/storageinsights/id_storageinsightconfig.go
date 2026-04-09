package storageinsights

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&StorageInsightConfigId{})
}

var _ resourceids.ResourceId = &StorageInsightConfigId{}

// StorageInsightConfigId is a struct representing the Resource ID for a Storage Insight Config
type StorageInsightConfigId struct {
	SubscriptionId           string
	ResourceGroupName        string
	WorkspaceName            string
	StorageInsightConfigName string
}

// NewStorageInsightConfigID returns a new StorageInsightConfigId struct
func NewStorageInsightConfigID(subscriptionId string, resourceGroupName string, workspaceName string, storageInsightConfigName string) StorageInsightConfigId {
	return StorageInsightConfigId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		WorkspaceName:            workspaceName,
		StorageInsightConfigName: storageInsightConfigName,
	}
}

// ParseStorageInsightConfigID parses 'input' into a StorageInsightConfigId
func ParseStorageInsightConfigID(input string) (*StorageInsightConfigId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageInsightConfigId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageInsightConfigId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStorageInsightConfigIDInsensitively parses 'input' case-insensitively into a StorageInsightConfigId
// note: this method should only be used for API response data and not user input
func ParseStorageInsightConfigIDInsensitively(input string) (*StorageInsightConfigId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StorageInsightConfigId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StorageInsightConfigId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StorageInsightConfigId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.WorkspaceName, ok = input.Parsed["workspaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", input)
	}

	if id.StorageInsightConfigName, ok = input.Parsed["storageInsightConfigName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storageInsightConfigName", input)
	}

	return nil
}

// ValidateStorageInsightConfigID checks that 'input' can be parsed as a Storage Insight Config ID
func ValidateStorageInsightConfigID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStorageInsightConfigID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Storage Insight Config ID
func (id StorageInsightConfigId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/storageInsightConfigs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.StorageInsightConfigName)
}

// Segments returns a slice of Resource ID Segments which comprise this Storage Insight Config ID
func (id StorageInsightConfigId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceName"),
		resourceids.StaticSegment("staticStorageInsightConfigs", "storageInsightConfigs", "storageInsightConfigs"),
		resourceids.UserSpecifiedSegment("storageInsightConfigName", "storageInsightConfigName"),
	}
}

// String returns a human-readable description of this Storage Insight Config ID
func (id StorageInsightConfigId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Storage Insight Config Name: %q", id.StorageInsightConfigName),
	}
	return fmt.Sprintf("Storage Insight Config (%s)", strings.Join(components, "\n"))
}
