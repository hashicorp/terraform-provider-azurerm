package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = MongodbDatabaseCollectionId{}

// MongodbDatabaseCollectionId is a struct representing the Resource ID for a Mongodb Database Collection
type MongodbDatabaseCollectionId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DatabaseAccountName string
	MongodbDatabaseName string
	CollectionName      string
}

// NewMongodbDatabaseCollectionID returns a new MongodbDatabaseCollectionId struct
func NewMongodbDatabaseCollectionID(subscriptionId string, resourceGroupName string, databaseAccountName string, mongodbDatabaseName string, collectionName string) MongodbDatabaseCollectionId {
	return MongodbDatabaseCollectionId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DatabaseAccountName: databaseAccountName,
		MongodbDatabaseName: mongodbDatabaseName,
		CollectionName:      collectionName,
	}
}

// ParseMongodbDatabaseCollectionID parses 'input' into a MongodbDatabaseCollectionId
func ParseMongodbDatabaseCollectionID(input string) (*MongodbDatabaseCollectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(MongodbDatabaseCollectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MongodbDatabaseCollectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.MongodbDatabaseName, ok = parsed.Parsed["mongodbDatabaseName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mongodbDatabaseName", *parsed)
	}

	if id.CollectionName, ok = parsed.Parsed["collectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "collectionName", *parsed)
	}

	return &id, nil
}

// ParseMongodbDatabaseCollectionIDInsensitively parses 'input' case-insensitively into a MongodbDatabaseCollectionId
// note: this method should only be used for API response data and not user input
func ParseMongodbDatabaseCollectionIDInsensitively(input string) (*MongodbDatabaseCollectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(MongodbDatabaseCollectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MongodbDatabaseCollectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.MongodbDatabaseName, ok = parsed.Parsed["mongodbDatabaseName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "mongodbDatabaseName", *parsed)
	}

	if id.CollectionName, ok = parsed.Parsed["collectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "collectionName", *parsed)
	}

	return &id, nil
}

// ValidateMongodbDatabaseCollectionID checks that 'input' can be parsed as a Mongodb Database Collection ID
func ValidateMongodbDatabaseCollectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMongodbDatabaseCollectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Mongodb Database Collection ID
func (id MongodbDatabaseCollectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/mongodbDatabases/%s/collections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.MongodbDatabaseName, id.CollectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Mongodb Database Collection ID
func (id MongodbDatabaseCollectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountValue"),
		resourceids.StaticSegment("staticMongodbDatabases", "mongodbDatabases", "mongodbDatabases"),
		resourceids.UserSpecifiedSegment("mongodbDatabaseName", "mongodbDatabaseValue"),
		resourceids.StaticSegment("staticCollections", "collections", "collections"),
		resourceids.UserSpecifiedSegment("collectionName", "collectionValue"),
	}
}

// String returns a human-readable description of this Mongodb Database Collection ID
func (id MongodbDatabaseCollectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Mongodb Database Name: %q", id.MongodbDatabaseName),
		fmt.Sprintf("Collection Name: %q", id.CollectionName),
	}
	return fmt.Sprintf("Mongodb Database Collection (%s)", strings.Join(components, "\n"))
}
