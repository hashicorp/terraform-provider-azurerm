package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AvsPrivateCloudId struct {
	ResourceGroup string
	Name          string
}

func AvsPrivateCloudID(input string) (*AvsPrivateCloudId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing avsPrivateCloud ID %q: %+v", input, err)
	}

	avsPrivateCloud := AvsPrivateCloudId{
		ResourceGroup: id.ResourceGroup,
	}
	if avsPrivateCloud.Name, err = id.PopSegment("privateClouds"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &avsPrivateCloud, nil
}
