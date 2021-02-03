package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CassandraTableId struct {
	SubscriptionId        string
	ResourceGroup         string
	DatabaseAccountName   string
	CassandraKeyspaceName string
	TableName             string
}

func NewCassandraTableID(subscriptionId, resourceGroup, databaseAccountName, cassandraKeyspaceName, tableName string) CassandraTableId {
	return CassandraTableId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		DatabaseAccountName:   databaseAccountName,
		CassandraKeyspaceName: cassandraKeyspaceName,
		TableName:             tableName,
	}
}

func (id CassandraTableId) String() string {
	segments := []string{
		fmt.Sprintf("Table Name %q", id.TableName),
		fmt.Sprintf("Cassandra Keyspace Name %q", id.CassandraKeyspaceName),
		fmt.Sprintf("Database Account Name %q", id.DatabaseAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Cassandra Table", segmentsStr)
}

func (id CassandraTableId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/cassandraKeyspaces/%s/tables/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DatabaseAccountName, id.CassandraKeyspaceName, id.TableName)
}

// CassandraTableID parses a CassandraTable ID into an CassandraTableId struct
func CassandraTableID(input string) (*CassandraTableId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := CassandraTableId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.DatabaseAccountName, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}
	if resourceId.CassandraKeyspaceName, err = id.PopSegment("cassandraKeyspaces"); err != nil {
		return nil, err
	}
	if resourceId.TableName, err = id.PopSegment("tables"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
