package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

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

func (id DatabasePrincipalId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Cluster Name %q", id.ClusterName),
		fmt.Sprintf("Database Name %q", id.DatabaseName),
		fmt.Sprintf("Role Name %q", id.RoleName),
		fmt.Sprintf("F Q N Name %q", id.FQNName),
	}
	return strings.Join(segments, " / ")
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
