package linkedstorageaccounts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DataSourceTypeId{})
}

var _ resourceids.ResourceId = &DataSourceTypeId{}

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
	parser := resourceids.NewParserFromResourceIdType(&DataSourceTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataSourceTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDataSourceTypeIDInsensitively parses 'input' case-insensitively into a DataSourceTypeId
// note: this method should only be used for API response data and not user input
func ParseDataSourceTypeIDInsensitively(input string) (*DataSourceTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataSourceTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataSourceTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DataSourceTypeId) FromParseResult(input resourceids.ParseResult) error {
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

	if v, ok := input.Parsed["dataSourceType"]; true {
		if !ok {
			return resourceids.NewSegmentNotSpecifiedError(id, "dataSourceType", input)
		}

		dataSourceType, err := parseDataSourceType(v)
		if err != nil {
			return fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.DataSourceType = *dataSourceType
	}

	return nil
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
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceName"),
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
