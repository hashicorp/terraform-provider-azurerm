package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type HybridConnectionId struct {
	SubscriptionId                string
	ResourceGroup                 string
	SiteName                      string
	HybridConnectionNamespaceName string
	RelayName                     string
}

func NewHybridConnectionID(subscriptionId, resourceGroup, siteName, hybridConnectionNamespaceName, relayName string) HybridConnectionId {
	return HybridConnectionId{
		SubscriptionId:                subscriptionId,
		ResourceGroup:                 resourceGroup,
		SiteName:                      siteName,
		HybridConnectionNamespaceName: hybridConnectionNamespaceName,
		RelayName:                     relayName,
	}
}

func (id HybridConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Relay Name %q", id.RelayName),
		fmt.Sprintf("Hybrid Connection Namespace Name %q", id.HybridConnectionNamespaceName),
		fmt.Sprintf("Site Name %q", id.SiteName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Hybrid Connection", segmentsStr)
}

func (id HybridConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Web/sites/%s/hybridConnectionNamespaces/%s/relays/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName)
}

// HybridConnectionID parses a HybridConnection ID into an HybridConnectionId struct
func HybridConnectionID(input string) (*HybridConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
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

	if resourceId.SiteName, err = id.PopSegment("sites"); err != nil {
		return nil, err
	}
	if resourceId.HybridConnectionNamespaceName, err = id.PopSegment("hybridConnectionNamespaces"); err != nil {
		return nil, err
	}
	if resourceId.RelayName, err = id.PopSegment("relays"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
