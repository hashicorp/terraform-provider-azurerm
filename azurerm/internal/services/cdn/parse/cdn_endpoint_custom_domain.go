package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CdnEndpointCustomDomainId struct {
	ResourceGroup string
	ProfileName   string
	EndpointName  string
	Name          string
}

func CdnEndpointCustomDomainID(input string) (*CdnEndpointCustomDomainId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse CDN Endpoint Custom Domain ID %q: %+v", input, err)
	}

	domainId := CdnEndpointCustomDomainId{
		ResourceGroup: id.ResourceGroup,
	}

	if domainId.ProfileName, err = id.PopSegment("profiles"); err != nil {
		return nil, err
	}

	if domainId.EndpointName, err = id.PopSegment("endpoints"); err != nil {
		return nil, err
	}

	if domainId.Name, err = id.PopSegment("customdomains"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &domainId, nil
}
