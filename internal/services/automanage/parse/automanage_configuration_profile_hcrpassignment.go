package parse

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type AutomanageConfigurationProfileHCRPAssignmentId struct {
	SubscriptionId string
	ResourceGroup  string
	MachineName    string
	Name           string
}

func NewAutomanageConfigurationProfileHCRPAssignmentID(subscriptionId string, resourcegroup string, machinename string, name string) AutomanageConfigurationProfileHCRPAssignmentId {
	return AutomanageConfigurationProfileHCRPAssignmentId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourcegroup,
		MachineName:    machinename,
		Name:           name,
	}
}

func (id AutomanageConfigurationProfileHCRPAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HybridCompute/machines/%s/providers/Microsoft.Automanage/configurationProfileAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.MachineName, id.Name)
}

func AutomanageConfigurationProfileHCRPAssignmentID(input string) (*AutomanageConfigurationProfileHCRPAssignmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing automanageConfigurationProfileHCRPAssignment ID %q: %+v", input, err)
	}

	automanageConfigurationProfileHCRPAssignment := AutomanageConfigurationProfileHCRPAssignmentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}
	if automanageConfigurationProfileHCRPAssignment.MachineName, err = id.PopSegment("machines"); err != nil {
		return nil, err
	}
	if automanageConfigurationProfileHCRPAssignment.Name, err = id.PopSegment("configurationProfileAssignments"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &automanageConfigurationProfileHCRPAssignment, nil
}
