package mongorbacs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&MongodbUserDefinitionId{})
}

var _ resourceids.ResourceId = &MongodbUserDefinitionId{}

// MongodbUserDefinitionId is a struct representing the Resource ID for a Mongodb User Definition
type MongodbUserDefinitionId struct {
	SubscriptionId        string
	ResourceGroupName     string
	DatabaseAccountName   string
	MongoUserDefinitionId string
}

// NewMongodbUserDefinitionID returns a new MongodbUserDefinitionId struct
func NewMongodbUserDefinitionID(subscriptionId string, resourceGroupName string, databaseAccountName string, mongoUserDefinitionId string) MongodbUserDefinitionId {
	return MongodbUserDefinitionId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		DatabaseAccountName:   databaseAccountName,
		MongoUserDefinitionId: mongoUserDefinitionId,
	}
}

// ParseMongodbUserDefinitionID parses 'input' into a MongodbUserDefinitionId
func ParseMongodbUserDefinitionID(input string) (*MongodbUserDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MongodbUserDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MongodbUserDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMongodbUserDefinitionIDInsensitively parses 'input' case-insensitively into a MongodbUserDefinitionId
// note: this method should only be used for API response data and not user input
func ParseMongodbUserDefinitionIDInsensitively(input string) (*MongodbUserDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MongodbUserDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MongodbUserDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MongodbUserDefinitionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.MongoUserDefinitionId, ok = input.Parsed["mongoUserDefinitionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "mongoUserDefinitionId", input)
	}

	return nil
}

// ValidateMongodbUserDefinitionID checks that 'input' can be parsed as a Mongodb User Definition ID
func ValidateMongodbUserDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMongodbUserDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Mongodb User Definition ID
func (id MongodbUserDefinitionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/mongodbUserDefinitions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.MongoUserDefinitionId)
}

// Segments returns a slice of Resource ID Segments which comprise this Mongodb User Definition ID
func (id MongodbUserDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountName"),
		resourceids.StaticSegment("staticMongodbUserDefinitions", "mongodbUserDefinitions", "mongodbUserDefinitions"),
		resourceids.UserSpecifiedSegment("mongoUserDefinitionId", "mongoUserDefinitionId"),
	}
}

// String returns a human-readable description of this Mongodb User Definition ID
func (id MongodbUserDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Mongo User Definition: %q", id.MongoUserDefinitionId),
	}
	return fmt.Sprintf("Mongodb User Definition (%s)", strings.Join(components, "\n"))
}
