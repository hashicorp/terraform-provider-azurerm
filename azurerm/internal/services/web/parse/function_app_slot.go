package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FunctionAppSlotId struct {
	SubscriptionId string
	ResourceGroup  string
	SiteName       string
	SlotName       string
}

func NewFunctionAppSlotID(subscriptionId, resourceGroup, siteName, slotName string) FunctionAppSlotId {
	return FunctionAppSlotId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SiteName:       siteName,
		SlotName:       slotName,
	}
}

func (id FunctionAppSlotId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/slots/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SiteName, id.SlotName)
}

// FunctionAppSlotID parses a FunctionAppSlot ID into an FunctionAppSlotId struct
func FunctionAppSlotID(input string) (*FunctionAppSlotId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FunctionAppSlotId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SiteName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}
	if resourceId.SlotName, err = id.PopSegment("slots"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
