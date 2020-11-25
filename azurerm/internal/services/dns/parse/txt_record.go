package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type TxtRecordId struct {
	ResourceGroup string
	DnsZoneName   string
	TXTName       string
}

func TxtRecordID(input string) (*TxtRecordId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse DNS TXT Record ID %q: %+v", input, err)
	}

	record := TxtRecordId{
		ResourceGroup: id.ResourceGroup,
	}

	if record.DnsZoneName, err = id.PopSegment("dnszones"); err != nil {
		return nil, err
	}

	if record.TXTName, err = id.PopSegment("TXT"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &record, nil
}
