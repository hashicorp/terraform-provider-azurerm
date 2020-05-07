package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type BatchCertificateId struct {
	ResourceGroup string
	AccountName   string
	Name          string
}

func BatchCertificateID(input string) (*BatchCertificateId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Batch Certificate ID %q: %+v", input, err)
	}

	certificate := BatchCertificateId{
		ResourceGroup: id.ResourceGroup,
	}

	if certificate.AccountName, err = id.PopSegment("batchAccounts"); err != nil {
		return nil, err
	}

	if certificate.Name, err = id.PopSegment("certificates"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &certificate, nil
}
