package dataconnectors

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DataConnectorId{}

// DataConnectorId is a struct representing the Resource ID for a Data Connector
type DataConnectorId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	DataConnectorId   string
}

// NewDataConnectorID returns a new DataConnectorId struct
func NewDataConnectorID(subscriptionId string, resourceGroupName string, workspaceName string, dataConnectorId string) DataConnectorId {
	return DataConnectorId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		DataConnectorId:   dataConnectorId,
	}
}

// ParseDataConnectorID parses 'input' into a DataConnectorId
func ParseDataConnectorID(input string) (*DataConnectorId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataConnectorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataConnectorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'workspaceName' was not found in the resource id %q", input)
	}

	if id.DataConnectorId, ok = parsed.Parsed["dataConnectorId"]; !ok {
		return nil, fmt.Errorf("the segment 'dataConnectorId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseDataConnectorIDInsensitively parses 'input' case-insensitively into a DataConnectorId
// note: this method should only be used for API response data and not user input
func ParseDataConnectorIDInsensitively(input string) (*DataConnectorId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataConnectorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataConnectorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'workspaceName' was not found in the resource id %q", input)
	}

	if id.DataConnectorId, ok = parsed.Parsed["dataConnectorId"]; !ok {
		return nil, fmt.Errorf("the segment 'dataConnectorId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateDataConnectorID checks that 'input' can be parsed as a Data Connector ID
func ValidateDataConnectorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataConnectorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Data Connector ID
func (id DataConnectorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/dataConnectors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.DataConnectorId)
}

// Segments returns a slice of Resource ID Segments which comprise this Data Connector ID
func (id DataConnectorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSecurityInsights", "Microsoft.SecurityInsights", "Microsoft.SecurityInsights"),
		resourceids.StaticSegment("staticDataConnectors", "dataConnectors", "dataConnectors"),
		resourceids.UserSpecifiedSegment("dataConnectorId", "dataConnectorIdValue"),
	}
}

// String returns a human-readable description of this Data Connector ID
func (id DataConnectorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Data Connector: %q", id.DataConnectorId),
	}
	return fmt.Sprintf("Data Connector (%s)", strings.Join(components, "\n"))
}
