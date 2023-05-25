package storageinsights

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = StorageInsightConfigId{}

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
	parser := resourceids.NewParserFromResourceIdType(StorageInsightConfigId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StorageInsightConfigId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.StorageInsightConfigName, ok = parsed.Parsed["storageInsightConfigName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "storageInsightConfigName", *parsed)
	}

	return &id, nil
}

// ParseStorageInsightConfigIDInsensitively parses 'input' case-insensitively into a StorageInsightConfigId
// note: this method should only be used for API response data and not user input
func ParseStorageInsightConfigIDInsensitively(input string) (*StorageInsightConfigId, error) {
	parser := resourceids.NewParserFromResourceIdType(StorageInsightConfigId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StorageInsightConfigId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.StorageInsightConfigName, ok = parsed.Parsed["storageInsightConfigName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "storageInsightConfigName", *parsed)
	}

	return &id, nil
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
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticStorageInsightConfigs", "storageInsightConfigs", "storageInsightConfigs"),
		resourceids.UserSpecifiedSegment("storageInsightConfigName", "storageInsightConfigValue"),
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
