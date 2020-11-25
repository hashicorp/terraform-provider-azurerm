package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CaaRecordId struct {
	ResourceGroup string
	DnszoneName   string
	CAAName       string
}

func CaaRecordID(input string) (*CaaRecordId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse DNS CAA Record ID %q: %+v", input, err)
	}

	record := CaaRecordId{
		ResourceGroup: id.ResourceGroup,
	}

	if record.DnszoneName, err = id.PopSegment("dnszones"); err != nil {
		return nil, err
	}

	if record.CAAName, err = id.PopSegment("CAA"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &record, nil
}
