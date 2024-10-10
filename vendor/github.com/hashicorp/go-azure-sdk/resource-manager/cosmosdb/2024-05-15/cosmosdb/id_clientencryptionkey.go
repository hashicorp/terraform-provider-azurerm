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
	recaser.RegisterResourceId(&ClientEncryptionKeyId{})
}

var _ resourceids.ResourceId = &ClientEncryptionKeyId{}

// ClientEncryptionKeyId is a struct representing the Resource ID for a Client Encryption Key
type ClientEncryptionKeyId struct {
	SubscriptionId          string
	ResourceGroupName       string
	DatabaseAccountName     string
	SqlDatabaseName         string
	ClientEncryptionKeyName string
}

// NewClientEncryptionKeyID returns a new ClientEncryptionKeyId struct
func NewClientEncryptionKeyID(subscriptionId string, resourceGroupName string, databaseAccountName string, sqlDatabaseName string, clientEncryptionKeyName string) ClientEncryptionKeyId {
	return ClientEncryptionKeyId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		DatabaseAccountName:     databaseAccountName,
		SqlDatabaseName:         sqlDatabaseName,
		ClientEncryptionKeyName: clientEncryptionKeyName,
	}
}

// ParseClientEncryptionKeyID parses 'input' into a ClientEncryptionKeyId
func ParseClientEncryptionKeyID(input string) (*ClientEncryptionKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ClientEncryptionKeyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ClientEncryptionKeyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseClientEncryptionKeyIDInsensitively parses 'input' case-insensitively into a ClientEncryptionKeyId
// note: this method should only be used for API response data and not user input
func ParseClientEncryptionKeyIDInsensitively(input string) (*ClientEncryptionKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ClientEncryptionKeyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ClientEncryptionKeyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ClientEncryptionKeyId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ClientEncryptionKeyName, ok = input.Parsed["clientEncryptionKeyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "clientEncryptionKeyName", input)
	}

	return nil
}

// ValidateClientEncryptionKeyID checks that 'input' can be parsed as a Client Encryption Key ID
func ValidateClientEncryptionKeyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseClientEncryptionKeyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Client Encryption Key ID
func (id ClientEncryptionKeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/sqlDatabases/%s/clientEncryptionKeys/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.SqlDatabaseName, id.ClientEncryptionKeyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Client Encryption Key ID
func (id ClientEncryptionKeyId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticClientEncryptionKeys", "clientEncryptionKeys", "clientEncryptionKeys"),
		resourceids.UserSpecifiedSegment("clientEncryptionKeyName", "clientEncryptionKeyName"),
	}
}

// String returns a human-readable description of this Client Encryption Key ID
func (id ClientEncryptionKeyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Sql Database Name: %q", id.SqlDatabaseName),
		fmt.Sprintf("Client Encryption Key Name: %q", id.ClientEncryptionKeyName),
	}
	return fmt.Sprintf("Client Encryption Key (%s)", strings.Join(components, "\n"))
}
