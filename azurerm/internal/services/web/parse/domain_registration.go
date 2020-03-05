package parse

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DomainRegistration struct {
	ResourceGroup string
	Name          string
}

// /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DomainRegistration/domains/{domainName}

func DomainRegistrationID(input string) (*DomainRegistration, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Domain Registration ID %q: %+v", input, err)
	}

	domainRegistration := DomainRegistration {
		ResourceGroup: id.ResourceGroup,
	}

	if domainRegistration.Name, err = id.PopSegment("domains"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &domainRegistration, nil
}