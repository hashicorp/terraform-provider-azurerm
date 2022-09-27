package parse

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type AutomanageConfigurationProfileHCIAssignmentId struct {
	SubscriptionId string
	ResourceGroup  string
	ClusterName    string
	Name           string
}

func NewAutomanageConfigurationProfileHCIAssignmentID(subscriptionId string, resourcegroup string, clustername string, name string) AutomanageConfigurationProfileHCIAssignmentId {
	return AutomanageConfigurationProfileHCIAssignmentId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourcegroup,
		ClusterName:    clustername,
		Name:           name,
	}
}

func (id AutomanageConfigurationProfileHCIAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AzureStackHci/clusters/%s/providers/Microsoft.Automanage/configurationProfileAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName, id.Name)
}

func AutomanageConfigurationProfileHCIAssignmentID(input string) (*AutomanageConfigurationProfileHCIAssignmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing automanageConfigurationProfileHCIAssignment ID %q: %+v", input, err)
	}

	automanageConfigurationProfileHCIAssignment := AutomanageConfigurationProfileHCIAssignmentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}
	if automanageConfigurationProfileHCIAssignment.ClusterName, err = id.PopSegment("clusters"); err != nil {
		return nil, err
	}
	if automanageConfigurationProfileHCIAssignment.Name, err = id.PopSegment("configurationProfileAssignments"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &automanageConfigurationProfileHCIAssignment, nil
}
