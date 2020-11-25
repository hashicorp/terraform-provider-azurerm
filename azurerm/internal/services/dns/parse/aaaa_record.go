package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type AaaaRecordId struct {
	ResourceGroup string
	DnszoneName   string
	AAAAName      string
}

func AaaaRecordID(input string) (*AaaaRecordId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse DNS AAAA Record ID %q: %+v", input, err)
	}

	record := AaaaRecordId{
		ResourceGroup: id.ResourceGroup,
	}

	if record.DnszoneName, err = id.PopSegment("dnszones"); err != nil {
		return nil, err
	}

	if record.AAAAName, err = id.PopSegment("AAAA"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &record, nil
}
