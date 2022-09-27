package parse

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type AutomanageConfigurationProfileAssignmentId struct {
	SubscriptionId string
	ResourceGroup  string
	VMName         string
	Name           string
}

func NewAutomanageConfigurationProfileAssignmentID(subscriptionId string, resourcegroup string, vmname string, name string) AutomanageConfigurationProfileAssignmentId {
	return AutomanageConfigurationProfileAssignmentId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourcegroup,
		VMName:         vmname,
		Name:           name,
	}
}

func (id AutomanageConfigurationProfileAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s/providers/Microsoft.Automanage/configurationProfileAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VMName, id.Name)
}

func AutomanageConfigurationProfileAssignmentID(input string) (*AutomanageConfigurationProfileAssignmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing automanageConfigurationProfileAssignment ID %q: %+v", input, err)
	}

	automanageConfigurationProfileAssignment := AutomanageConfigurationProfileAssignmentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}
	if automanageConfigurationProfileAssignment.VMName, err = id.PopSegment("virtualMachines"); err != nil {
		return nil, err
	}
	if automanageConfigurationProfileAssignment.Name, err = id.PopSegment("configurationProfileAssignments"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &automanageConfigurationProfileAssignment, nil
}
