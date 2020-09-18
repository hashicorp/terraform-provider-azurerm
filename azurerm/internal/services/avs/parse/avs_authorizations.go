package parse

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AvsAuthorizationId struct {
	ResourceGroup    string
	PrivateCloudName string
	Name             string
}

func AvsAuthorizationID(input string) (*AvsAuthorizationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing avsAuthorization ID %q: %+v", input, err)
	}

	avsAuthorization := AvsAuthorizationId{
		ResourceGroup: id.ResourceGroup,
	}
	if avsAuthorization.PrivateCloudName, err = id.PopSegment("privateClouds"); err != nil {
		return nil, err
	}
	if avsAuthorization.Name, err = id.PopSegment("authorizations"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &avsAuthorization, nil
}
