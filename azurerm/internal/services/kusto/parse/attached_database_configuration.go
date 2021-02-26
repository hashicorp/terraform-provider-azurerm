package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

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

func (id AttachedDatabaseConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Cluster Name %q", id.ClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Attached Database Configuration", segmentsStr)
}

func (id AttachedDatabaseConfigurationId) ID() string {
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

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
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
