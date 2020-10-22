package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SecurityPartnerProviderId struct {
	ResourceGroup string
	Name          string
}

func SecurityPartnerProviderID(input string) (*SecurityPartnerProviderId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing securityPartnerProvider ID %q: %+v", input, err)
	}

	securityPartnerProvider := SecurityPartnerProviderId{
		ResourceGroup: id.ResourceGroup,
	}

	if securityPartnerProvider.Name, err = id.PopSegment("securityPartnerProviders"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &securityPartnerProvider, nil
}
