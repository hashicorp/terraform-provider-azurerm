package web

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceSlotResourceID struct {
	ResourceGroup  string
	AppServiceName string
	Name           string
}

func ParseAppServiceSlotID(input string) (*AppServiceSlotResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Slot ID %q: %+v", input, err)
	}

	slot := AppServiceSlotResourceID{
		ResourceGroup:  id.ResourceGroup,
		AppServiceName: id.Path["sites"],
		Name:           id.Path["slots"],
	}

	if slot.AppServiceName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}

	if slot.Name, err = id.PopSegment("slots"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &slot, nil
}
