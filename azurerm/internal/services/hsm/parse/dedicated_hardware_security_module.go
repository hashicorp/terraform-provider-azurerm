package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DedicatedHardwareSecurityModuleId struct {
	SubscriptionId   string
	ResourceGroup    string
	DedicatedHSMName string
}

func NewDedicatedHardwareSecurityModuleID(subscriptionId, resourceGroup, dedicatedHSMName string) DedicatedHardwareSecurityModuleId {
	return DedicatedHardwareSecurityModuleId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		DedicatedHSMName: dedicatedHSMName,
	}
}

func (id DedicatedHardwareSecurityModuleId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HardwareSecurityModules/dedicatedHSMs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DedicatedHSMName)
}

// DedicatedHardwareSecurityModuleID parses a DedicatedHardwareSecurityModule ID into an DedicatedHardwareSecurityModuleId struct
func DedicatedHardwareSecurityModuleID(input string) (*DedicatedHardwareSecurityModuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DedicatedHardwareSecurityModuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.DedicatedHSMName, err = id.PopSegment("dedicatedHSMs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
