package signalr

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type PrivateEndpointConnectionId struct {
	SubscriptionId string
	ResourceGroup  string
	SignalRName    string
	Name           string
}

func NewPrivateEndpointConnectionID(subscriptionId, resourceGroup, signalRName, name string) PrivateEndpointConnectionId {
	return PrivateEndpointConnectionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SignalRName:    signalRName,
		Name:           name,
	}
}

func (id PrivateEndpointConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Signal R Name %q", id.SignalRName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Private Endpoint Connection", segmentsStr)
}

func (id PrivateEndpointConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SignalRService/SignalR/%s/privateEndpointConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SignalRName, id.Name)
}

// ParsePrivateEndpointConnectionID parses a PrivateEndpointConnection ID into an PrivateEndpointConnectionId struct
func ParsePrivateEndpointConnectionID(input string) (*PrivateEndpointConnectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
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

	if resourceId.SignalRName, err = id.PopSegment("SignalR"); err != nil {
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

// ParsePrivateEndpointConnectionIDInsensitively parses an PrivateEndpointConnection ID into an PrivateEndpointConnectionId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParsePrivateEndpointConnectionID method should be used instead for validation etc.
func ParsePrivateEndpointConnectionIDInsensitively(input string) (*PrivateEndpointConnectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
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

	// find the correct casing for the 'SignalR' segment
	SignalRKey := "SignalR"
	for key := range id.Path {
		if strings.EqualFold(key, SignalRKey) {
			SignalRKey = key
			break
		}
	}
	if resourceId.SignalRName, err = id.PopSegment(SignalRKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'privateEndpointConnections' segment
	privateEndpointConnectionsKey := "privateEndpointConnections"
	for key := range id.Path {
		if strings.EqualFold(key, privateEndpointConnectionsKey) {
			privateEndpointConnectionsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(privateEndpointConnectionsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
