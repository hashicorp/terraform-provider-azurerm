package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type MxRecordId struct {
	ResourceGroup string
	DnszoneName   string
	MXName        string
}

func MxRecordID(input string) (*MxRecordId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse DNS MX Record ID %q: %+v", input, err)
	}

	record := MxRecordId{
		ResourceGroup: id.ResourceGroup,
	}

	if record.DnszoneName, err = id.PopSegment("dnszones"); err != nil {
		return nil, err
	}

	if record.MXName, err = id.PopSegment("MX"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &record, nil
}
