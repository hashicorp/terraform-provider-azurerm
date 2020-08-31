package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SpringCloudCertificateId struct {
	ResourceGroup string
	ServiceName   string
	Name          string
}

func SpringCloudCertificateID(input string) (*SpringCloudCertificateId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Spring Cloud Certificate ID %q: %+v", input, err)
	}

	cert := SpringCloudCertificateId{
		ResourceGroup: id.ResourceGroup,
	}

	if cert.ServiceName, err = id.PopSegment("Spring"); err != nil {
		return nil, err
	}

	if cert.Name, err = id.PopSegment("certificates"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &cert, nil
}
