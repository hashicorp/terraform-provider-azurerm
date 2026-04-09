package entities

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
)

// TODO: update this to implement `resourceids.ResourceId` once
// https://github.com/hashicorp/go-azure-helpers/issues/187 is fixed
var _ resourceids.Id = EntityId{}

type EntityId struct {
	// AccountId specifies the ID of the Storage Account where this Entity exists.
	AccountId accounts.AccountId

	// TableName specifies the name of the Table where this Entity exists.
	TableName string

	// PartitionKey specifies the Partition Key for this Entity.
	PartitionKey string

	// RowKey specifies the Row Key for this Entity.
	RowKey string
}

func NewEntityID(accountId accounts.AccountId, tableName, partitionKey, rowKey string) EntityId {
	return EntityId{
		AccountId:    accountId,
		TableName:    tableName,
		PartitionKey: partitionKey,
		RowKey:       rowKey,
	}
}

func (b EntityId) ID() string {
	return fmt.Sprintf("%s/%s(PartitionKey='%s',RowKey='%s')", b.AccountId.ID(), b.TableName, b.PartitionKey, b.RowKey)
}

func (b EntityId) String() string {
	components := []string{
		fmt.Sprintf("Partition Key %q", b.PartitionKey),
		fmt.Sprintf("Row Key %q", b.RowKey),
		fmt.Sprintf("Table Name %q", b.TableName),
		fmt.Sprintf("Account %q", b.AccountId.String()),
	}
	return fmt.Sprintf("Entity (%s)", strings.Join(components, " / "))
}

// ParseEntityID parses `input` into a Entity ID using a known `domainSuffix`
func ParseEntityID(input, domainSuffix string) (*EntityId, error) {
	// example: https://foo.table.core.windows.net/bar(PartitionKey='partition1',RowKey='row1')
	if input == "" {
		return nil, fmt.Errorf("`input` was empty")
	}

	account, err := accounts.ParseAccountID(input, domainSuffix)
	if err != nil {
		return nil, fmt.Errorf("parsing account %q: %+v", input, err)
	}

	if account.SubDomainType != accounts.TableSubDomainType {
		return nil, fmt.Errorf("expected the subdomain type to be %q but got %q", string(accounts.TableSubDomainType), string(account.SubDomainType))
	}

	uri, err := url.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a uri: %+v", input, err)
	}

	path := strings.TrimPrefix(uri.Path, "/")
	segments := strings.Split(path, "/")
	if len(segments) != 1 {
		return nil, fmt.Errorf("expected the path to contain 1 segment but got %d", len(segments))
	}

	// Tables and Table Entities are similar with table being `table1` and entities
	// being `table1(PartitionKey='samplepartition',RowKey='samplerow')` so we need to validate this is a table
	slug := strings.TrimPrefix(uri.Path, "/")
	if strings.HasPrefix(slug, "Tables('") && strings.HasSuffix(slug, "')") {
		// Ensure we do not parse a Table ID in the format: https://foo.table.core.windows.net/Table('foo')
		return nil, fmt.Errorf("expected the path to be an entity name but got a table name: %q", slug)
	} else if !strings.Contains(slug, "(") || !strings.HasSuffix(slug, ")") {
		// Ensure we do not try to parse a bare table name
		return nil, fmt.Errorf("expected the path to be an entity name but got an invalid format, possibly a table name: %q", slug)
	}

	indexOfFirstBracket := strings.Index(slug, "(")
	tableName := slug[0:indexOfFirstBracket]
	componentString := slug[indexOfFirstBracket:]
	componentString = strings.TrimPrefix(componentString, "(")
	componentString = strings.TrimSuffix(componentString, ")")
	components := strings.Split(componentString, ",")
	if len(components) != 2 {
		return nil, fmt.Errorf("expected the path to be an entity name but got %q", slug)
	}

	partitionKey := parseValueFromKey(components[0], "PartitionKey")
	rowKey := parseValueFromKey(components[1], "RowKey")

	return &EntityId{
		AccountId:    *account,
		TableName:    tableName,
		PartitionKey: *partitionKey,
		RowKey:       *rowKey,
	}, nil
}

func parseValueFromKey(input, expectedKey string) *string {
	components := strings.Split(input, "=")
	if len(components) != 2 {
		return nil
	}
	key := components[0]
	value := components[1]
	if key != expectedKey {
		return nil
	}

	// the value is surrounded in single quotes, remove those
	value = strings.TrimPrefix(value, "'")
	value = strings.TrimSuffix(value, "'")
	return &value
}
