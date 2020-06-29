package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DnsCNameRecordId struct {
	ResourceGroup string
	ZoneName      string
	Name          string
}

func DnsCNameRecordID(input string) (*DnsCNameRecordId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse DNS CNAME Record ID %q: %+v", input, err)
	}

	record := DnsCNameRecordId{
		ResourceGroup: id.ResourceGroup,
	}

	if record.ZoneName, err = id.PopSegment("dnszones"); err != nil {
		return nil, err
	}

	if record.Name, err = id.PopSegment("CNAME"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &record, nil
}
