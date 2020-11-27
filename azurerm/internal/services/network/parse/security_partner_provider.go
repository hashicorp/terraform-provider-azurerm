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
		return nil, fmt.Errorf("parsing virtualHubSecurityPartnerProvider ID %q: %+v", input, err)
	}

	vhubSecurityPartnerProvider := SecurityPartnerProviderId{
		ResourceGroup: id.ResourceGroup,
	}

	if vhubSecurityPartnerProvider.Name, err = id.PopSegment("securityPartnerProviders"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &vhubSecurityPartnerProvider, nil
}
