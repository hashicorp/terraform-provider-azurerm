package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DatabasePrincipalId struct {
	SubscriptionId string
	ResourceGroup  string
	ClusterName    string
	DatabaseName   string
	RoleName       string
	FQNName        string
}

func NewDatabasePrincipalID(subscriptionId, resourceGroup, clusterName, databaseName, roleName, fQNName string) DatabasePrincipalId {
	return DatabasePrincipalId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ClusterName:    clusterName,
		DatabaseName:   databaseName,
		RoleName:       roleName,
		FQNName:        fQNName,
	}
}

func (id DatabasePrincipalId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/Clusters/%s/Databases/%s/Role/%s/FQN/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName, id.DatabaseName, id.RoleName, id.FQNName)
}

// DatabasePrincipalID parses a DatabasePrincipal ID into an DatabasePrincipalId struct
func DatabasePrincipalID(input string) (*DatabasePrincipalId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DatabasePrincipalId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.ClusterName, err = id.PopSegment("Clusters"); err != nil {
		return nil, err
	}
	if resourceId.DatabaseName, err = id.PopSegment("Databases"); err != nil {
		return nil, err
	}
	if resourceId.RoleName, err = id.PopSegment("Role"); err != nil {
		return nil, err
	}
	if resourceId.FQNName, err = id.PopSegment("FQN"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
