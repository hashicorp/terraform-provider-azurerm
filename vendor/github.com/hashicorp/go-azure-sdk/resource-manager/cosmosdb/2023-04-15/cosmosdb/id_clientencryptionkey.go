package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ClientEncryptionKeyId{}

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
	parser := resourceids.NewParserFromResourceIdType(ClientEncryptionKeyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ClientEncryptionKeyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.SqlDatabaseName, ok = parsed.Parsed["sqlDatabaseName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sqlDatabaseName", *parsed)
	}

	if id.ClientEncryptionKeyName, ok = parsed.Parsed["clientEncryptionKeyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "clientEncryptionKeyName", *parsed)
	}

	return &id, nil
}

// ParseClientEncryptionKeyIDInsensitively parses 'input' case-insensitively into a ClientEncryptionKeyId
// note: this method should only be used for API response data and not user input
func ParseClientEncryptionKeyIDInsensitively(input string) (*ClientEncryptionKeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(ClientEncryptionKeyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ClientEncryptionKeyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "databaseAccountName", *parsed)
	}

	if id.SqlDatabaseName, ok = parsed.Parsed["sqlDatabaseName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "sqlDatabaseName", *parsed)
	}

	if id.ClientEncryptionKeyName, ok = parsed.Parsed["clientEncryptionKeyName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "clientEncryptionKeyName", *parsed)
	}

	return &id, nil
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
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountValue"),
		resourceids.StaticSegment("staticSqlDatabases", "sqlDatabases", "sqlDatabases"),
		resourceids.UserSpecifiedSegment("sqlDatabaseName", "sqlDatabaseValue"),
		resourceids.StaticSegment("staticClientEncryptionKeys", "clientEncryptionKeys", "clientEncryptionKeys"),
		resourceids.UserSpecifiedSegment("clientEncryptionKeyName", "clientEncryptionKeyValue"),
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
