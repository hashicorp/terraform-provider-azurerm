package hybridconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type HybridConnectionId struct {
	SubscriptionId string
	ResourceGroup  string
	NamespaceName  string
	Name           string
}

func NewHybridConnectionID(subscriptionId, resourceGroup, namespaceName, name string) HybridConnectionId {
	return HybridConnectionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		NamespaceName:  namespaceName,
		Name:           name,
	}
}

func (id HybridConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Namespace Name %q", id.NamespaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Hybrid Connection", segmentsStr)
}

func (id HybridConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Relay/namespaces/%s/hybridConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NamespaceName, id.Name)
}

// ParseHybridConnectionID parses a HybridConnection ID into an HybridConnectionId struct
func ParseHybridConnectionID(input string) (*HybridConnectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := HybridConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.NamespaceName, err = id.PopSegment("namespaces"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("hybridConnections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseHybridConnectionIDInsensitively parses an HybridConnection ID into an HybridConnectionId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseHybridConnectionID method should be used instead for validation etc.
func ParseHybridConnectionIDInsensitively(input string) (*HybridConnectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := HybridConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'namespaces' segment
	namespacesKey := "namespaces"
	for key := range id.Path {
		if strings.EqualFold(key, namespacesKey) {
			namespacesKey = key
			break
		}
	}
	if resourceId.NamespaceName, err = id.PopSegment(namespacesKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'hybridConnections' segment
	hybridConnectionsKey := "hybridConnections"
	for key := range id.Path {
		if strings.EqualFold(key, hybridConnectionsKey) {
			hybridConnectionsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(hybridConnectionsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
