package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&PartitionKeyRangeIdId{})
}

var _ resourceids.ResourceId = &PartitionKeyRangeIdId{}

// PartitionKeyRangeIdId is a struct representing the Resource ID for a Partition Key Range Id
type PartitionKeyRangeIdId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DatabaseAccountName string
	DatabaseName        string
	CollectionName      string
	PartitionKeyRangeId string
}

// NewPartitionKeyRangeIdID returns a new PartitionKeyRangeIdId struct
func NewPartitionKeyRangeIdID(subscriptionId string, resourceGroupName string, databaseAccountName string, databaseName string, collectionName string, partitionKeyRangeId string) PartitionKeyRangeIdId {
	return PartitionKeyRangeIdId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DatabaseAccountName: databaseAccountName,
		DatabaseName:        databaseName,
		CollectionName:      collectionName,
		PartitionKeyRangeId: partitionKeyRangeId,
	}
}

// ParsePartitionKeyRangeIdID parses 'input' into a PartitionKeyRangeIdId
func ParsePartitionKeyRangeIdID(input string) (*PartitionKeyRangeIdId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PartitionKeyRangeIdId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PartitionKeyRangeIdId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParsePartitionKeyRangeIdIDInsensitively parses 'input' case-insensitively into a PartitionKeyRangeIdId
// note: this method should only be used for API response data and not user input
func ParsePartitionKeyRangeIdIDInsensitively(input string) (*PartitionKeyRangeIdId, error) {
	parser := resourceids.NewParserFromResourceIdType(&PartitionKeyRangeIdId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := PartitionKeyRangeIdId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *PartitionKeyRangeIdId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DatabaseAccountName, ok = input.Parsed["databaseAccountName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", input)
	}

	if id.DatabaseName, ok = input.Parsed["databaseName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "databaseName", input)
	}

	if id.CollectionName, ok = input.Parsed["collectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "collectionName", input)
	}

	if id.PartitionKeyRangeId, ok = input.Parsed["partitionKeyRangeId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "partitionKeyRangeId", input)
	}

	return nil
}

// ValidatePartitionKeyRangeIdID checks that 'input' can be parsed as a Partition Key Range Id ID
func ValidatePartitionKeyRangeIdID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePartitionKeyRangeIdID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Partition Key Range Id ID
func (id PartitionKeyRangeIdId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/databases/%s/collections/%s/partitionKeyRangeId/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.DatabaseName, id.CollectionName, id.PartitionKeyRangeId)
}

// Segments returns a slice of Resource ID Segments which comprise this Partition Key Range Id ID
func (id PartitionKeyRangeIdId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountName"),
		resourceids.StaticSegment("staticDatabases", "databases", "databases"),
		resourceids.UserSpecifiedSegment("databaseName", "databaseName"),
		resourceids.StaticSegment("staticCollections", "collections", "collections"),
		resourceids.UserSpecifiedSegment("collectionName", "collectionName"),
		resourceids.StaticSegment("staticPartitionKeyRangeId", "partitionKeyRangeId", "partitionKeyRangeId"),
		resourceids.UserSpecifiedSegment("partitionKeyRangeId", "partitionKeyRangeId"),
	}
}

// String returns a human-readable description of this Partition Key Range Id ID
func (id PartitionKeyRangeIdId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Database Name: %q", id.DatabaseName),
		fmt.Sprintf("Collection Name: %q", id.CollectionName),
		fmt.Sprintf("Partition Key Range: %q", id.PartitionKeyRangeId),
	}
	return fmt.Sprintf("Partition Key Range Id (%s)", strings.Join(components, "\n"))
}
