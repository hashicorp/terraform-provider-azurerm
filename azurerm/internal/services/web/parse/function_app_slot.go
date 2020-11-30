package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FunctionAppSlotId struct {
	ResourceGroup string
	SiteName      string
	SlotName      string
}

func FunctionAppSlotID(input string) (*FunctionAppSlotId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Slot ID %q: %+v", input, err)
	}

	slot := FunctionAppSlotId{
		ResourceGroup: id.ResourceGroup,
		SiteName:      id.Path["sites"],
		SlotName:      id.Path["slots"],
	}

	if slot.SiteName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}

	if slot.SlotName, err = id.PopSegment("slots"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &slot, nil
}
