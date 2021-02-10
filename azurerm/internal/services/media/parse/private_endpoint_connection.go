package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PrivateEndpointConnectionId struct {
	SubscriptionId   string
	ResourceGroup    string
	MediaserviceName string
	Name             string
}

func NewPrivateEndpointConnectionID(subscriptionId, resourceGroup, mediaserviceName, name string) PrivateEndpointConnectionId {
	return PrivateEndpointConnectionId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		MediaserviceName: mediaserviceName,
		Name:             name,
	}
}

func (id PrivateEndpointConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Mediaservice Name %q", id.MediaserviceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Private Endpoint Connection", segmentsStr)
}

func (id PrivateEndpointConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Media/mediaservices/%s/privateEndpointConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.MediaserviceName, id.Name)
}

// PrivateEndpointConnectionID parses a PrivateEndpointConnection ID into an PrivateEndpointConnectionId struct
func PrivateEndpointConnectionID(input string) (*PrivateEndpointConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := PrivateEndpointConnectionId{
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
	if resourceId.Name, err = id.PopSegment("privateEndpointConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
