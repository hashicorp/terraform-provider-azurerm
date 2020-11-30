package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DataConnectionId struct {
	SubscriptionId string
	ResourceGroup  string
	ClusterName    string
	DatabaseName   string
	Name           string
}

func NewDataConnectionID(subscriptionId, resourceGroup, clusterName, databaseName, name string) DataConnectionId {
	return DataConnectionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ClusterName:    clusterName,
		DatabaseName:   databaseName,
		Name:           name,
	}
}

func (id DataConnectionId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/Clusters/%s/Databases/%s/DataConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.Name)
}

// DataConnectionID parses a DataConnection ID into an DataConnectionId struct
func DataConnectionID(input string) (*DataConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DataConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ClusterName, err = id.PopSegment("Clusters"); err != nil {
		return nil, err
	}
	if resourceId.DatabaseName, err = id.PopSegment("Databases"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("DataConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
