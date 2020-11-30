package web

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AppServiceCertificateResourceID struct {
	ResourceGroup string
	Name          string
}

func ParseAppServiceCertificateID(input string) (*AppServiceCertificateResourceID, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse App Service Certificate ID %q: %+v", input, err)
	}

	certificate := AppServiceCertificateResourceID{
		ResourceGroup: id.ResourceGroup,
	}

	if certificate.Name, err = id.PopSegment("certificates"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &certificate, nil
}
