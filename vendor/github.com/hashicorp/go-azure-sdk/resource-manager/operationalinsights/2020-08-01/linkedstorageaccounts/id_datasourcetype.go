package linkedstorageaccounts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DataSourceTypeId{}

// DataSourceTypeId is a struct representing the Resource ID for a Data Source Type
type DataSourceTypeId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	DataSourceType    DataSourceType
}

// NewDataSourceTypeID returns a new DataSourceTypeId struct
func NewDataSourceTypeID(subscriptionId string, resourceGroupName string, workspaceName string, dataSourceType DataSourceType) DataSourceTypeId {
	return DataSourceTypeId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		DataSourceType:    dataSourceType,
	}
}

// ParseDataSourceTypeID parses 'input' into a DataSourceTypeId
func ParseDataSourceTypeID(input string) (*DataSourceTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataSourceTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataSourceTypeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if v, ok := parsed.Parsed["dataSourceType"]; true {
		if !ok {
			return nil, resourceids.NewSegmentNotSpecifiedError(id, "dataSourceType", *parsed)
		}

		dataSourceType, err := parseDataSourceType(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.DataSourceType = *dataSourceType
	}

	return &id, nil
}

// ParseDataSourceTypeIDInsensitively parses 'input' case-insensitively into a DataSourceTypeId
// note: this method should only be used for API response data and not user input
func ParseDataSourceTypeIDInsensitively(input string) (*DataSourceTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataSourceTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataSourceTypeId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if v, ok := parsed.Parsed["dataSourceType"]; true {
		if !ok {
			return nil, resourceids.NewSegmentNotSpecifiedError(id, "dataSourceType", *parsed)
		}

		dataSourceType, err := parseDataSourceType(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.DataSourceType = *dataSourceType
	}

	return &id, nil
}

// ValidateDataSourceTypeID checks that 'input' can be parsed as a Data Source Type ID
func ValidateDataSourceTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataSourceTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Data Source Type ID
func (id DataSourceTypeId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/linkedStorageAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, string(id.DataSourceType))
}

// Segments returns a slice of Resource ID Segments which comprise this Data Source Type ID
func (id DataSourceTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticLinkedStorageAccounts", "linkedStorageAccounts", "linkedStorageAccounts"),
		resourceids.ConstantSegment("dataSourceType", PossibleValuesForDataSourceType(), "Alerts"),
	}
}

// String returns a human-readable description of this Data Source Type ID
func (id DataSourceTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Data Source Type: %q", string(id.DataSourceType)),
	}
	return fmt.Sprintf("Data Source Type (%s)", strings.Join(components, "\n"))
}
