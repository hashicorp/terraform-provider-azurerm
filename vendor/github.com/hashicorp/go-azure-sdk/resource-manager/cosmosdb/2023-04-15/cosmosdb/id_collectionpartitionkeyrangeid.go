package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = CollectionPartitionKeyRangeIdId{}

// CollectionPartitionKeyRangeIdId is a struct representing the Resource ID for a Collection Partition Key Range Id
type CollectionPartitionKeyRangeIdId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DatabaseAccountName string
	RegionName          string
	DatabaseName        string
	CollectionName      string
	PartitionKeyRangeId string
}

// NewCollectionPartitionKeyRangeIdID returns a new CollectionPartitionKeyRangeIdId struct
func NewCollectionPartitionKeyRangeIdID(subscriptionId string, resourceGroupName string, databaseAccountName string, regionName string, databaseName string, collectionName string, partitionKeyRangeId string) CollectionPartitionKeyRangeIdId {
	return CollectionPartitionKeyRangeIdId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DatabaseAccountName: databaseAccountName,
		RegionName:          regionName,
		DatabaseName:        databaseName,
		CollectionName:      collectionName,
		PartitionKeyRangeId: partitionKeyRangeId,
	}
}

// ParseCollectionPartitionKeyRangeIdID parses 'input' into a CollectionPartitionKeyRangeIdId
func ParseCollectionPartitionKeyRangeIdID(input string) (*CollectionPartitionKeyRangeIdId, error) {
	parser := resourceids.NewParserFromResourceIdType(CollectionPartitionKeyRangeIdId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CollectionPartitionKeyRangeIdId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.RegionName, ok = parsed.Parsed["regionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "regionName", *parsed)
	}

	if id.DatabaseName, ok = parsed.Parsed["databaseName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseName", *parsed)
	}

	if id.CollectionName, ok = parsed.Parsed["collectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "collectionName", *parsed)
	}

	if id.PartitionKeyRangeId, ok = parsed.Parsed["partitionKeyRangeId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "partitionKeyRangeId", *parsed)
	}

	return &id, nil
}

// ParseCollectionPartitionKeyRangeIdIDInsensitively parses 'input' case-insensitively into a CollectionPartitionKeyRangeIdId
// note: this method should only be used for API response data and not user input
func ParseCollectionPartitionKeyRangeIdIDInsensitively(input string) (*CollectionPartitionKeyRangeIdId, error) {
	parser := resourceids.NewParserFromResourceIdType(CollectionPartitionKeyRangeIdId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CollectionPartitionKeyRangeIdId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.RegionName, ok = parsed.Parsed["regionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "regionName", *parsed)
	}

	if id.DatabaseName, ok = parsed.Parsed["databaseName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseName", *parsed)
	}

	if id.CollectionName, ok = parsed.Parsed["collectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "collectionName", *parsed)
	}

	if id.PartitionKeyRangeId, ok = parsed.Parsed["partitionKeyRangeId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "partitionKeyRangeId", *parsed)
	}

	return &id, nil
}

// ValidateCollectionPartitionKeyRangeIdID checks that 'input' can be parsed as a Collection Partition Key Range Id ID
func ValidateCollectionPartitionKeyRangeIdID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCollectionPartitionKeyRangeIdID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Collection Partition Key Range Id ID
func (id CollectionPartitionKeyRangeIdId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/region/%s/databases/%s/collections/%s/partitionKeyRangeId/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.RegionName, id.DatabaseName, id.CollectionName, id.PartitionKeyRangeId)
}

// Segments returns a slice of Resource ID Segments which comprise this Collection Partition Key Range Id ID
func (id CollectionPartitionKeyRangeIdId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountValue"),
		resourceids.StaticSegment("staticRegion", "region", "region"),
		resourceids.UserSpecifiedSegment("regionName", "regionValue"),
		resourceids.StaticSegment("staticDatabases", "databases", "databases"),
		resourceids.UserSpecifiedSegment("databaseName", "databaseValue"),
		resourceids.StaticSegment("staticCollections", "collections", "collections"),
		resourceids.UserSpecifiedSegment("collectionName", "collectionValue"),
		resourceids.StaticSegment("staticPartitionKeyRangeId", "partitionKeyRangeId", "partitionKeyRangeId"),
		resourceids.UserSpecifiedSegment("partitionKeyRangeId", "partitionKeyRangeIdValue"),
	}
}

// String returns a human-readable description of this Collection Partition Key Range Id ID
func (id CollectionPartitionKeyRangeIdId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Region Name: %q", id.RegionName),
		fmt.Sprintf("Database Name: %q", id.DatabaseName),
		fmt.Sprintf("Collection Name: %q", id.CollectionName),
		fmt.Sprintf("Partition Key Range: %q", id.PartitionKeyRangeId),
	}
	return fmt.Sprintf("Collection Partition Key Range Id (%s)", strings.Join(components, "\n"))
}
