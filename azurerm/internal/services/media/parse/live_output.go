package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LiveOutputId struct {
	SubscriptionId   string
	ResourceGroup    string
	MediaserviceName string
	LiveeventName    string
	Name             string
}

func NewLiveOutputID(subscriptionId, resourceGroup, mediaserviceName, liveeventName, name string) LiveOutputId {
	return LiveOutputId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		MediaserviceName: mediaserviceName,
		LiveeventName:    liveeventName,
		Name:             name,
	}
}

func (id LiveOutputId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Liveevent Name %q", id.LiveeventName),
		fmt.Sprintf("Mediaservice Name %q", id.MediaserviceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Live Output", segmentsStr)
}

func (id LiveOutputId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaservices/%s/liveevents/%s/liveoutputs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.MediaserviceName, id.LiveeventName, id.Name)
}

// LiveOutputID parses a LiveOutput ID into an LiveOutputId struct
func LiveOutputID(input string) (*LiveOutputId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LiveOutputId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.MediaserviceName, err = id.PopSegment("mediaservices"); err != nil {
		return nil, err
	}
	if resourceId.LiveeventName, err = id.PopSegment("liveevents"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("liveoutputs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
