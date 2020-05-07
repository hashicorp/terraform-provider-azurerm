package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DnsAaaaRecordId struct {
	ResourceGroup string
	ZoneName      string
	Name          string
}

func DnsAaaaRecordID(input string) (*DnsAaaaRecordId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse DNS AAAA Record ID %q: %+v", input, err)
	}

	record := DnsAaaaRecordId{
		ResourceGroup: id.ResourceGroup,
	}

	if record.ZoneName, err = id.PopSegment("dnszones"); err != nil {
		return nil, err
	}

	if record.Name, err = id.PopSegment("AAAA"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &record, nil
}
