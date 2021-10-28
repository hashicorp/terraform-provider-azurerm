package mhsmprivateendpointconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ManagedHSMPrivateEndpointConnectionId struct {
	SubscriptionId                string
	ResourceGroup                 string
	ManagedHSMName                string
	PrivateEndpointConnectionName string
}

func NewManagedHSMPrivateEndpointConnectionID(subscriptionId, resourceGroup, managedHSMName, privateEndpointConnectionName string) ManagedHSMPrivateEndpointConnectionId {
	return ManagedHSMPrivateEndpointConnectionId{
		SubscriptionId:                subscriptionId,
		ResourceGroup:                 resourceGroup,
		ManagedHSMName:                managedHSMName,
		PrivateEndpointConnectionName: privateEndpointConnectionName,
	}
}

func (id ManagedHSMPrivateEndpointConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Private Endpoint Connection Name %q", id.PrivateEndpointConnectionName),
		fmt.Sprintf("Managed H S M Name %q", id.ManagedHSMName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Managed H S M Private Endpoint Connection", segmentsStr)
}

func (id ManagedHSMPrivateEndpointConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/managedHSMs/%s/privateEndpointConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ManagedHSMName, id.PrivateEndpointConnectionName)
}

// ParseManagedHSMPrivateEndpointConnectionID parses a ManagedHSMPrivateEndpointConnection ID into an ManagedHSMPrivateEndpointConnectionId struct
func ParseManagedHSMPrivateEndpointConnectionID(input string) (*ManagedHSMPrivateEndpointConnectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ManagedHSMPrivateEndpointConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ManagedHSMName, err = id.PopSegment("managedHSMs"); err != nil {
		return nil, err
	}
	if resourceId.PrivateEndpointConnectionName, err = id.PopSegment("privateEndpointConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseManagedHSMPrivateEndpointConnectionIDInsensitively parses an ManagedHSMPrivateEndpointConnection ID into an ManagedHSMPrivateEndpointConnectionId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseManagedHSMPrivateEndpointConnectionID method should be used instead for validation etc.
func ParseManagedHSMPrivateEndpointConnectionIDInsensitively(input string) (*ManagedHSMPrivateEndpointConnectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ManagedHSMPrivateEndpointConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'managedHSMs' segment
	managedHSMsKey := "managedHSMs"
	for key := range id.Path {
		if strings.EqualFold(key, managedHSMsKey) {
			managedHSMsKey = key
			break
		}
	}
	if resourceId.ManagedHSMName, err = id.PopSegment(managedHSMsKey); err != nil {
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
	if resourceId.PrivateEndpointConnectionName, err = id.PopSegment(privateEndpointConnectionsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
