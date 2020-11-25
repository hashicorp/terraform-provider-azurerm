package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ARecordId struct {
	ResourceGroup string
	DnszoneName   string
	AName         string
}

func ARecordID(input string) (*ARecordId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse DNS A Record ID %q: %+v", input, err)
	}

	record := ARecordId{
		ResourceGroup: id.ResourceGroup,
	}

	if record.DnszoneName, err = id.PopSegment("dnszones"); err != nil {
		return nil, err
	}

	if record.AName, err = id.PopSegment("A"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &record, nil
}
