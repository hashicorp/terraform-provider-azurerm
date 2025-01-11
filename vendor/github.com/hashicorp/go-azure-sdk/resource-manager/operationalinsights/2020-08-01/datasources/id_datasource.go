package datasources

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DataSourceId{})
}

var _ resourceids.ResourceId = &DataSourceId{}

// DataSourceId is a struct representing the Resource ID for a Data Source
type DataSourceId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	DataSourceName    string
}

// NewDataSourceID returns a new DataSourceId struct
func NewDataSourceID(subscriptionId string, resourceGroupName string, workspaceName string, dataSourceName string) DataSourceId {
	return DataSourceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		DataSourceName:    dataSourceName,
	}
}

// ParseDataSourceID parses 'input' into a DataSourceId
func ParseDataSourceID(input string) (*DataSourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataSourceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataSourceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDataSourceIDInsensitively parses 'input' case-insensitively into a DataSourceId
// note: this method should only be used for API response data and not user input
func ParseDataSourceIDInsensitively(input string) (*DataSourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataSourceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataSourceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DataSourceId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.DataSourceName, ok = input.Parsed["dataSourceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dataSourceName", input)
	}

	return nil
}

// ValidateDataSourceID checks that 'input' can be parsed as a Data Source ID
func ValidateDataSourceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataSourceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Data Source ID
func (id DataSourceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/dataSources/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.DataSourceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Data Source ID
func (id DataSourceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceName"),
		resourceids.StaticSegment("staticDataSources", "dataSources", "dataSources"),
		resourceids.UserSpecifiedSegment("dataSourceName", "dataSourceName"),
	}
}

// String returns a human-readable description of this Data Source ID
func (id DataSourceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Data Source Name: %q", id.DataSourceName),
	}
	return fmt.Sprintf("Data Source (%s)", strings.Join(components, "\n"))
}
