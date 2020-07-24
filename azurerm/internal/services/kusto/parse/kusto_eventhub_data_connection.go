package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type KustoEventHubDataConnectionId struct {
	ResourceGroup string
	Cluster       string
	Database      string
	Name          string
}

func KustoEventHubDataConnectionID(input string) (*KustoEventHubDataConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Kusto Event Hub Data Connection ID %q: %+v", input, err)
	}

	dataConnection := KustoEventHubDataConnectionId{
		ResourceGroup: id.ResourceGroup,
	}

	if dataConnection.Cluster, err = id.PopSegment("Clusters"); err != nil {
		return nil, err
	}

	if dataConnection.Database, err = id.PopSegment("Databases"); err != nil {
		return nil, err
	}

	if dataConnection.Name, err = id.PopSegment("DataConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &dataConnection, nil
}
