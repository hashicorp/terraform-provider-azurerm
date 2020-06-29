package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DnsMxRecordId struct {
	ResourceGroup string
	ZoneName      string
	Name          string
}

func DnsMxRecordID(input string) (*DnsMxRecordId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse DNS MX Record ID %q: %+v", input, err)
	}

	record := DnsMxRecordId{
		ResourceGroup: id.ResourceGroup,
	}

	if record.ZoneName, err = id.PopSegment("dnszones"); err != nil {
		return nil, err
	}

	if record.Name, err = id.PopSegment("MX"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &record, nil
}
