package tables

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
)

// GetResourceManagerResourceID returns the Resource ID for the given Table
// This can be useful when, for example, you're using this as a unique identifier
func (c Client) GetResourceManagerResourceID(subscriptionID, resourceGroup, accountName, tableName string) string {
	fmtStr := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/tableServices/default/tables/%s"
	return fmt.Sprintf(fmtStr, subscriptionID, resourceGroup, accountName, tableName)
}

// TODO: update this to implement `resourceids.ResourceId` once
// https://github.com/hashicorp/go-azure-helpers/issues/187 is fixed
var _ resourceids.Id = TableId{}

type TableId struct {
	// AccountId specifies the ID of the Storage Account where this Table exists.
	AccountId accounts.AccountId

	// TableName specifies the name of this Table.
	TableName string
}

func NewTableID(accountId accounts.AccountId, tableName string) TableId {
	return TableId{
		AccountId: accountId,
		TableName: tableName,
	}
}

func (b TableId) ID() string {
	return fmt.Sprintf("%s/Tables('%s')", b.AccountId.ID(), b.TableName)
}

func (b TableId) String() string {
	components := []string{
		fmt.Sprintf("Account %q", b.AccountId.String()),
	}
	return fmt.Sprintf("Table %q (%s)", b.TableName, strings.Join(components, " / "))
}

// ParseTableID parses `input` into a Table ID using a known `domainSuffix`
func ParseTableID(input, domainSuffix string) (*TableId, error) {
	// example: https://foo.table.core.windows.net/Table('bar')
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

	// Tables and Table Entities are similar however Tables use a reserved namespace, for example:
	//   Table('tableName')
	// whereas Entities begin with the actual table name, for example:
	//   tableName(PartitionKey='samplepartition',RowKey='samplerow')
	// However, there was a period of time when Table IDs did not use the reserved namespace, so we attempt to parse
	// both forms for maximum compatibility.
	var tableName string
	slug := strings.TrimPrefix(uri.Path, "/")
	if strings.HasPrefix(slug, "Tables('") && strings.HasSuffix(slug, "')") {
		// Ensure both prefix and suffix are present before trimming them out
		tableName = strings.TrimSuffix(strings.TrimPrefix(slug, "Tables('"), "')")
	} else if !strings.Contains(slug, "(") && !strings.HasSuffix(slug, ")") {
		// Also accept a bare table name
		tableName = slug
	} else {
		return nil, fmt.Errorf("expected the path to a table name and not an entity name but got %q", tableName)
	}
	if tableName == "" {
		return nil, fmt.Errorf("expected the path to a table name but the path was empty")
	}

	return &TableId{
		AccountId: *account,
		TableName: tableName,
	}, nil
}
