package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PacketCaptureId struct {
	SubscriptionId     string
	ResourceGroup      string
	NetworkWatcherName string
	Name               string
}

func NewPacketCaptureID(subscriptionId, resourceGroup, networkWatcherName, name string) PacketCaptureId {
	return PacketCaptureId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		NetworkWatcherName: networkWatcherName,
		Name:               name,
	}
}

func (id PacketCaptureId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Network Watcher Name %q", id.NetworkWatcherName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Packet Capture", segmentsStr)
}

func (id PacketCaptureId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkWatchers/%s/packetCaptures/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NetworkWatcherName, id.Name)
}

// PacketCaptureID parses a PacketCapture ID into an PacketCaptureId struct
func PacketCaptureID(input string) (*PacketCaptureId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := PacketCaptureId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.NetworkWatcherName, err = id.PopSegment("networkWatchers"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("packetCaptures"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
