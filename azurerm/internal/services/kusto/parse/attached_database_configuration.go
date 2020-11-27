package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AttachedDatabaseConfigurationId struct {
	SubscriptionId string
	ResourceGroup  string
	ClusterName    string
	Name           string
}

func NewAttachedDatabaseConfigurationID(subscriptionId, resourceGroup, clusterName, name string) AttachedDatabaseConfigurationId {
	return AttachedDatabaseConfigurationId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ClusterName:    clusterName,
		Name:           name,
	}
}

func (id AttachedDatabaseConfigurationId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/Clusters/%s/AttachedDatabaseConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName, id.Name)
}

// AttachedDatabaseConfigurationID parses a AttachedDatabaseConfiguration ID into an AttachedDatabaseConfigurationId struct
func AttachedDatabaseConfigurationID(input string) (*AttachedDatabaseConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AttachedDatabaseConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.ClusterName, err = id.PopSegment("Clusters"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("AttachedDatabaseConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
