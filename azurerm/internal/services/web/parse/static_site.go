package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type StaticSiteResourceID struct {
	ResourceGroup string
	Name          string
}

func StaticSiteID(input string) (*StaticSiteResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Static Site ID %q: %+v", input, err)
	}

	staticSite := StaticSiteResourceID{
		ResourceGroup: id.ResourceGroup,
	}

	if staticSite.Name, err = id.PopSegment("staticSites"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &staticSite, nil
}
