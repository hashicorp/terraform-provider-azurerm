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
	recaser.RegisterResourceId(&MongodbDatabaseCollectionId{})
}

var _ resourceids.ResourceId = &MongodbDatabaseCollectionId{}

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
	parser := resourceids.NewParserFromResourceIdType(&MongodbDatabaseCollectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MongodbDatabaseCollectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMongodbDatabaseCollectionIDInsensitively parses 'input' case-insensitively into a MongodbDatabaseCollectionId
// note: this method should only be used for API response data and not user input
func ParseMongodbDatabaseCollectionIDInsensitively(input string) (*MongodbDatabaseCollectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MongodbDatabaseCollectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MongodbDatabaseCollectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MongodbDatabaseCollectionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.MongodbDatabaseName, ok = input.Parsed["mongodbDatabaseName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "mongodbDatabaseName", input)
	}

	if id.CollectionName, ok = input.Parsed["collectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "collectionName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountName"),
		resourceids.StaticSegment("staticMongodbDatabases", "mongodbDatabases", "mongodbDatabases"),
		resourceids.UserSpecifiedSegment("mongodbDatabaseName", "mongodbDatabaseName"),
		resourceids.StaticSegment("staticCollections", "collections", "collections"),
		resourceids.UserSpecifiedSegment("collectionName", "collectionName"),
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
