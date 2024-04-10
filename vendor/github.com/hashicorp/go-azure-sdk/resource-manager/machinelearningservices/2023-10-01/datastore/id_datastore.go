package datastore

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &DataStoreId{}

// DataStoreId is a struct representing the Resource ID for a Data Store
type DataStoreId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	DataStoreName     string
}

// NewDataStoreID returns a new DataStoreId struct
func NewDataStoreID(subscriptionId string, resourceGroupName string, workspaceName string, dataStoreName string) DataStoreId {
	return DataStoreId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		DataStoreName:     dataStoreName,
	}
}

// ParseDataStoreID parses 'input' into a DataStoreId
func ParseDataStoreID(input string) (*DataStoreId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataStoreId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataStoreId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDataStoreIDInsensitively parses 'input' case-insensitively into a DataStoreId
// note: this method should only be used for API response data and not user input
func ParseDataStoreIDInsensitively(input string) (*DataStoreId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataStoreId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataStoreId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DataStoreId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DataStoreName, ok = input.Parsed["dataStoreName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dataStoreName", input)
	}

	return nil
}

// ValidateDataStoreID checks that 'input' can be parsed as a Data Store ID
func ValidateDataStoreID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataStoreID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Data Store ID
func (id DataStoreId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MachineLearningServices/workspaces/%s/dataStores/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.DataStoreName)
}

// Segments returns a slice of Resource ID Segments which comprise this Data Store ID
func (id DataStoreId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMachineLearningServices", "Microsoft.MachineLearningServices", "Microsoft.MachineLearningServices"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticDataStores", "dataStores", "dataStores"),
		resourceids.UserSpecifiedSegment("dataStoreName", "dataStoreValue"),
	}
}

// String returns a human-readable description of this Data Store ID
func (id DataStoreId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Data Store Name: %q", id.DataStoreName),
	}
	return fmt.Sprintf("Data Store (%s)", strings.Join(components, "\n"))
}
