package configurationstores

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type PrivateEndpointConnectionId struct {
	SubscriptionId         string
	ResourceGroup          string
	ConfigurationStoreName string
	Name                   string
}

func NewPrivateEndpointConnectionID(subscriptionId, resourceGroup, configurationStoreName, name string) PrivateEndpointConnectionId {
	return PrivateEndpointConnectionId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		ConfigurationStoreName: configurationStoreName,
		Name:                   name,
	}
}

func (id PrivateEndpointConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Configuration Store Name %q", id.ConfigurationStoreName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Private Endpoint Connection", segmentsStr)
}

func (id PrivateEndpointConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppConfiguration/configurationStores/%s/privateEndpointConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ConfigurationStoreName, id.Name)
}

// PrivateEndpointConnectionID parses a PrivateEndpointConnection ID into an PrivateEndpointConnectionId struct
func PrivateEndpointConnectionID(input string) (*PrivateEndpointConnectionId, error) {
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

	if resourceId.ConfigurationStoreName, err = id.PopSegment("configurationStores"); err != nil {
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

// PrivateEndpointConnectionIDInsensitively parses an PrivateEndpointConnection ID into an PrivateEndpointConnectionId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the PrivateEndpointConnectionID method should be used instead for validation etc.
func PrivateEndpointConnectionIDInsensitively(input string) (*PrivateEndpointConnectionId, error) {
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

	// find the correct casing for the 'configurationStores' segment
	configurationStoresKey := "configurationStores"
	for key := range id.Path {
		if strings.EqualFold(key, configurationStoresKey) {
			configurationStoresKey = key
			break
		}
	}
	if resourceId.ConfigurationStoreName, err = id.PopSegment(configurationStoresKey); err != nil {
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
