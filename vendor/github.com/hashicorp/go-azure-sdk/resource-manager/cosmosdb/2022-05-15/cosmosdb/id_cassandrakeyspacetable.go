package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = CassandraKeyspaceTableId{}

// CassandraKeyspaceTableId is a struct representing the Resource ID for a Cassandra Keyspace Table
type CassandraKeyspaceTableId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	KeyspaceName      string
	TableName         string
}

// NewCassandraKeyspaceTableID returns a new CassandraKeyspaceTableId struct
func NewCassandraKeyspaceTableID(subscriptionId string, resourceGroupName string, accountName string, keyspaceName string, tableName string) CassandraKeyspaceTableId {
	return CassandraKeyspaceTableId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		KeyspaceName:      keyspaceName,
		TableName:         tableName,
	}
}

// ParseCassandraKeyspaceTableID parses 'input' into a CassandraKeyspaceTableId
func ParseCassandraKeyspaceTableID(input string) (*CassandraKeyspaceTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(CassandraKeyspaceTableId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CassandraKeyspaceTableId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.KeyspaceName, ok = parsed.Parsed["keyspaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'keyspaceName' was not found in the resource id %q", input)
	}

	if id.TableName, ok = parsed.Parsed["tableName"]; !ok {
		return nil, fmt.Errorf("the segment 'tableName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseCassandraKeyspaceTableIDInsensitively parses 'input' case-insensitively into a CassandraKeyspaceTableId
// note: this method should only be used for API response data and not user input
func ParseCassandraKeyspaceTableIDInsensitively(input string) (*CassandraKeyspaceTableId, error) {
	parser := resourceids.NewParserFromResourceIdType(CassandraKeyspaceTableId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CassandraKeyspaceTableId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.KeyspaceName, ok = parsed.Parsed["keyspaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'keyspaceName' was not found in the resource id %q", input)
	}

	if id.TableName, ok = parsed.Parsed["tableName"]; !ok {
		return nil, fmt.Errorf("the segment 'tableName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateCassandraKeyspaceTableID checks that 'input' can be parsed as a Cassandra Keyspace Table ID
func ValidateCassandraKeyspaceTableID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCassandraKeyspaceTableID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Cassandra Keyspace Table ID
func (id CassandraKeyspaceTableId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/cassandraKeyspaces/%s/tables/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.KeyspaceName, id.TableName)
}

// Segments returns a slice of Resource ID Segments which comprise this Cassandra Keyspace Table ID
func (id CassandraKeyspaceTableId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticCassandraKeyspaces", "cassandraKeyspaces", "cassandraKeyspaces"),
		resourceids.UserSpecifiedSegment("keyspaceName", "keyspaceValue"),
		resourceids.StaticSegment("staticTables", "tables", "tables"),
		resourceids.UserSpecifiedSegment("tableName", "tableValue"),
	}
}

// String returns a human-readable description of this Cassandra Keyspace Table ID
func (id CassandraKeyspaceTableId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Keyspace Name: %q", id.KeyspaceName),
		fmt.Sprintf("Table Name: %q", id.TableName),
	}
	return fmt.Sprintf("Cassandra Keyspace Table (%s)", strings.Join(components, "\n"))
}
