package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AttestationProviderId struct {
	ResourceGroup string
	Name          string
}

func AttestationId(input string) (*AttestationProviderId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing attestation ID %q: %+v", input, err)
	}

	attestationProvider := AttestationProviderId{
		ResourceGroup: id.ResourceGroup,
	}
	if attestationProvider.Name, err = id.PopSegment("attestationProviders"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &attestationProvider, nil
}
