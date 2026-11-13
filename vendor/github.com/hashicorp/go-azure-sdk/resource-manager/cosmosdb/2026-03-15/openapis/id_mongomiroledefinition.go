package openapis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&MongoMIRoleDefinitionId{})
}

var _ resourceids.ResourceId = &MongoMIRoleDefinitionId{}

// MongoMIRoleDefinitionId is a struct representing the Resource ID for a Mongo MI Role Definition
type MongoMIRoleDefinitionId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DatabaseAccountName string
	RoleDefinitionId    string
}

// NewMongoMIRoleDefinitionID returns a new MongoMIRoleDefinitionId struct
func NewMongoMIRoleDefinitionID(subscriptionId string, resourceGroupName string, databaseAccountName string, roleDefinitionId string) MongoMIRoleDefinitionId {
	return MongoMIRoleDefinitionId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DatabaseAccountName: databaseAccountName,
		RoleDefinitionId:    roleDefinitionId,
	}
}

// ParseMongoMIRoleDefinitionID parses 'input' into a MongoMIRoleDefinitionId
func ParseMongoMIRoleDefinitionID(input string) (*MongoMIRoleDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MongoMIRoleDefinitionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MongoMIRoleDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMongoMIRoleDefinitionIDInsensitively parses 'input' case-insensitively into a MongoMIRoleDefinitionId
// note: this method should only be used for API response data and not user input
func ParseMongoMIRoleDefinitionIDInsensitively(input string) (*MongoMIRoleDefinitionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MongoMIRoleDefinitionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MongoMIRoleDefinitionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MongoMIRoleDefinitionId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.RoleDefinitionId, ok = input.Parsed["roleDefinitionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "roleDefinitionId", input)
	}

	return nil
}

// ValidateMongoMIRoleDefinitionID checks that 'input' can be parsed as a Mongo MI Role Definition ID
func ValidateMongoMIRoleDefinitionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMongoMIRoleDefinitionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Mongo MI Role Definition ID
func (id MongoMIRoleDefinitionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/mongoMIRoleDefinitions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.RoleDefinitionId)
}

// Segments returns a slice of Resource ID Segments which comprise this Mongo MI Role Definition ID
func (id MongoMIRoleDefinitionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountName"),
		resourceids.StaticSegment("staticMongoMIRoleDefinitions", "mongoMIRoleDefinitions", "mongoMIRoleDefinitions"),
		resourceids.UserSpecifiedSegment("roleDefinitionId", "roleDefinitionId"),
	}
}

// String returns a human-readable description of this Mongo MI Role Definition ID
func (id MongoMIRoleDefinitionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Role Definition: %q", id.RoleDefinitionId),
	}
	return fmt.Sprintf("Mongo MI Role Definition (%s)", strings.Join(components, "\n"))
}
