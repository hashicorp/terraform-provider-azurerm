package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CassandraKeyspaceId struct {
	SubscriptionId      string
	ResourceGroup       string
	DatabaseAccountName string
	Name                string
}

func NewCassandraKeyspaceID(subscriptionId, resourceGroup, databaseAccountName, name string) CassandraKeyspaceId {
	return CassandraKeyspaceId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		DatabaseAccountName: databaseAccountName,
		Name:                name,
	}
}

func (id CassandraKeyspaceId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/cassandraKeyspaces/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DatabaseAccountName, id.Name)
}

// CassandraKeyspaceID parses a CassandraKeyspace ID into an CassandraKeyspaceId struct
func CassandraKeyspaceID(input string) (*CassandraKeyspaceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := CassandraKeyspaceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.DatabaseAccountName, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("cassandraKeyspaces"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
