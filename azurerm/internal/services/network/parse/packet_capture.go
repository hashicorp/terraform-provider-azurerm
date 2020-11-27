package parse

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

type PacketCaptureId struct {
	ResourceGroup      string
	NetworkWatcherName string
	Name               string
}

func PacketCaptureID(input string) (*PacketCaptureId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	packetCapture := PacketCaptureId{
		ResourceGroup: id.ResourceGroup,
	}

	if packetCapture.NetworkWatcherName, err = id.PopSegment("networkWatchers"); err != nil {
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
