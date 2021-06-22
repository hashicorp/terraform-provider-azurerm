package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MaintenanceConfigurationId struct {
	SubscriptionId     string
	ResourceGroup      string
	ManagedClusterName string
	Name               string
}

func NewMaintenanceConfigurationID(subscriptionId, resourceGroup, managedClusterName, name string) MaintenanceConfigurationId {
	return MaintenanceConfigurationId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		ManagedClusterName: managedClusterName,
		Name:               name,
	}
}

func (id MaintenanceConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Managed Cluster Name %q", id.ManagedClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Maintenance Configuration", segmentsStr)
}

func (id MaintenanceConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/managedClusters/%s/maintenanceConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ManagedClusterName, id.Name)
}

// MaintenanceConfigurationID parses a MaintenanceConfiguration ID into an MaintenanceConfigurationId struct
func MaintenanceConfigurationID(input string) (*MaintenanceConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := MaintenanceConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ManagedClusterName, err = id.PopSegment("managedClusters"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("maintenanceConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
