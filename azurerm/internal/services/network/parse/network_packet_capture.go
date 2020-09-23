package parse

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

type NetworkPacketCaptureId struct {
	ResourceGroup string
	WatcherName   string
	Name          string
}

func NetworkPacketCaptureID(input string) (*NetworkPacketCaptureId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	packetCapture := NetworkPacketCaptureId{
		ResourceGroup: id.ResourceGroup,
	}

	if packetCapture.WatcherName, err = id.PopSegment("networkWatchers"); err != nil {
		return nil, err
	}

	if packetCapture.Name, err = id.PopSegment("packetCaptures"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &packetCapture, nil
}
