package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SrvRecordId struct {
	ResourceGroup string
	DnszoneName   string
	SRVName       string
}

func SrvRecordID(input string) (*SrvRecordId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse DNS SRV Record ID %q: %+v", input, err)
	}

	record := SrvRecordId{
		ResourceGroup: id.ResourceGroup,
	}

	if record.DnszoneName, err = id.PopSegment("dnszones"); err != nil {
		return nil, err
	}

	if record.SRVName, err = id.PopSegment("SRV"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &record, nil
}
