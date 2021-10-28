package privateendpointconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VaultPrivateEndpointConnectionId struct {
	SubscriptionId                string
	ResourceGroup                 string
	VaultName                     string
	PrivateEndpointConnectionName string
}

func NewVaultPrivateEndpointConnectionID(subscriptionId, resourceGroup, vaultName, privateEndpointConnectionName string) VaultPrivateEndpointConnectionId {
	return VaultPrivateEndpointConnectionId{
		SubscriptionId:                subscriptionId,
		ResourceGroup:                 resourceGroup,
		VaultName:                     vaultName,
		PrivateEndpointConnectionName: privateEndpointConnectionName,
	}
}

func (id VaultPrivateEndpointConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Private Endpoint Connection Name %q", id.PrivateEndpointConnectionName),
		fmt.Sprintf("Vault Name %q", id.VaultName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Vault Private Endpoint Connection", segmentsStr)
}

func (id VaultPrivateEndpointConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s/privateEndpointConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VaultName, id.PrivateEndpointConnectionName)
}

// ParseVaultPrivateEndpointConnectionID parses a VaultPrivateEndpointConnection ID into an VaultPrivateEndpointConnectionId struct
func ParseVaultPrivateEndpointConnectionID(input string) (*VaultPrivateEndpointConnectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VaultPrivateEndpointConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VaultName, err = id.PopSegment("vaults"); err != nil {
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

// ParseVaultPrivateEndpointConnectionIDInsensitively parses an VaultPrivateEndpointConnection ID into an VaultPrivateEndpointConnectionId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseVaultPrivateEndpointConnectionID method should be used instead for validation etc.
func ParseVaultPrivateEndpointConnectionIDInsensitively(input string) (*VaultPrivateEndpointConnectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VaultPrivateEndpointConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'vaults' segment
	vaultsKey := "vaults"
	for key := range id.Path {
		if strings.EqualFold(key, vaultsKey) {
			vaultsKey = key
			break
		}
	}
	if resourceId.VaultName, err = id.PopSegment(vaultsKey); err != nil {
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
