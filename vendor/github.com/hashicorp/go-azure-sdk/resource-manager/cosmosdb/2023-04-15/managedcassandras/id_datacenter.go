package managedcassandras

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DataCenterId{})
}

var _ resourceids.ResourceId = &DataCenterId{}

// DataCenterId is a struct representing the Resource ID for a Data Center
type DataCenterId struct {
	SubscriptionId       string
	ResourceGroupName    string
	CassandraClusterName string
	DataCenterName       string
}

// NewDataCenterID returns a new DataCenterId struct
func NewDataCenterID(subscriptionId string, resourceGroupName string, cassandraClusterName string, dataCenterName string) DataCenterId {
	return DataCenterId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		CassandraClusterName: cassandraClusterName,
		DataCenterName:       dataCenterName,
	}
}

// ParseDataCenterID parses 'input' into a DataCenterId
func ParseDataCenterID(input string) (*DataCenterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataCenterId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataCenterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDataCenterIDInsensitively parses 'input' case-insensitively into a DataCenterId
// note: this method should only be used for API response data and not user input
func ParseDataCenterIDInsensitively(input string) (*DataCenterId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DataCenterId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DataCenterId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DataCenterId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CassandraClusterName, ok = input.Parsed["cassandraClusterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "cassandraClusterName", input)
	}

	if id.DataCenterName, ok = input.Parsed["dataCenterName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dataCenterName", input)
	}

	return nil
}

// ValidateDataCenterID checks that 'input' can be parsed as a Data Center ID
func ValidateDataCenterID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataCenterID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Data Center ID
func (id DataCenterId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/cassandraClusters/%s/dataCenters/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CassandraClusterName, id.DataCenterName)
}

// Segments returns a slice of Resource ID Segments which comprise this Data Center ID
func (id DataCenterId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticCassandraClusters", "cassandraClusters", "cassandraClusters"),
		resourceids.UserSpecifiedSegment("cassandraClusterName", "cassandraClusterName"),
		resourceids.StaticSegment("staticDataCenters", "dataCenters", "dataCenters"),
		resourceids.UserSpecifiedSegment("dataCenterName", "dataCenterName"),
	}
}

// String returns a human-readable description of this Data Center ID
func (id DataCenterId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cassandra Cluster Name: %q", id.CassandraClusterName),
		fmt.Sprintf("Data Center Name: %q", id.DataCenterName),
	}
	return fmt.Sprintf("Data Center (%s)", strings.Join(components, "\n"))
}
