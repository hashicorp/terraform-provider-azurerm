package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PeerAsnId struct {
	Name string
}

func PeerAsnID(input string) (*PeerAsnId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Peer Asn ID %q: %+v", input, err)
	}

	asn := PeerAsnId{}

	if asn.Name, err = id.PopSegment("peerAsns"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &asn, nil
}
