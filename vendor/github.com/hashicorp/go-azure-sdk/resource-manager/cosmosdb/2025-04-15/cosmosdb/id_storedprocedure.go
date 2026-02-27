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
	recaser.RegisterResourceId(&StoredProcedureId{})
}

var _ resourceids.ResourceId = &StoredProcedureId{}

// StoredProcedureId is a struct representing the Resource ID for a Stored Procedure
type StoredProcedureId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DatabaseAccountName string
	SqlDatabaseName     string
	ContainerName       string
	StoredProcedureName string
}

// NewStoredProcedureID returns a new StoredProcedureId struct
func NewStoredProcedureID(subscriptionId string, resourceGroupName string, databaseAccountName string, sqlDatabaseName string, containerName string, storedProcedureName string) StoredProcedureId {
	return StoredProcedureId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DatabaseAccountName: databaseAccountName,
		SqlDatabaseName:     sqlDatabaseName,
		ContainerName:       containerName,
		StoredProcedureName: storedProcedureName,
	}
}

// ParseStoredProcedureID parses 'input' into a StoredProcedureId
func ParseStoredProcedureID(input string) (*StoredProcedureId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StoredProcedureId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StoredProcedureId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseStoredProcedureIDInsensitively parses 'input' case-insensitively into a StoredProcedureId
// note: this method should only be used for API response data and not user input
func ParseStoredProcedureIDInsensitively(input string) (*StoredProcedureId, error) {
	parser := resourceids.NewParserFromResourceIdType(&StoredProcedureId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := StoredProcedureId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *StoredProcedureId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.SqlDatabaseName, ok = input.Parsed["sqlDatabaseName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "sqlDatabaseName", input)
	}

	if id.ContainerName, ok = input.Parsed["containerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "containerName", input)
	}

	if id.StoredProcedureName, ok = input.Parsed["storedProcedureName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "storedProcedureName", input)
	}

	return nil
}

// ValidateStoredProcedureID checks that 'input' can be parsed as a Stored Procedure ID
func ValidateStoredProcedureID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStoredProcedureID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Stored Procedure ID
func (id StoredProcedureId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/sqlDatabases/%s/containers/%s/storedProcedures/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, id.StoredProcedureName)
}

// Segments returns a slice of Resource ID Segments which comprise this Stored Procedure ID
func (id StoredProcedureId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountName"),
		resourceids.StaticSegment("staticSqlDatabases", "sqlDatabases", "sqlDatabases"),
		resourceids.UserSpecifiedSegment("sqlDatabaseName", "sqlDatabaseName"),
		resourceids.StaticSegment("staticContainers", "containers", "containers"),
		resourceids.UserSpecifiedSegment("containerName", "containerName"),
		resourceids.StaticSegment("staticStoredProcedures", "storedProcedures", "storedProcedures"),
		resourceids.UserSpecifiedSegment("storedProcedureName", "storedProcedureName"),
	}
}

// String returns a human-readable description of this Stored Procedure ID
func (id StoredProcedureId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Sql Database Name: %q", id.SqlDatabaseName),
		fmt.Sprintf("Container Name: %q", id.ContainerName),
		fmt.Sprintf("Stored Procedure Name: %q", id.StoredProcedureName),
	}
	return fmt.Sprintf("Stored Procedure (%s)", strings.Join(components, "\n"))
}
