package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type BackendAddressPoolTunnelInterfaceId struct {
	SubscriptionId         string
	ResourceGroup          string
	LoadBalancerName       string
	BackendAddressPoolName string
	TunnelInterfaceName    string
}

func NewBackendAddressPoolTunnelInterfaceID(subscriptionId, resourceGroup, loadBalancerName, backendAddressPoolName, tunnelInterfaceName string) BackendAddressPoolTunnelInterfaceId {
	return BackendAddressPoolTunnelInterfaceId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		LoadBalancerName:       loadBalancerName,
		BackendAddressPoolName: backendAddressPoolName,
		TunnelInterfaceName:    tunnelInterfaceName,
	}
}

func (id BackendAddressPoolTunnelInterfaceId) String() string {
	segments := []string{
		fmt.Sprintf("Tunnel Interface Name %q", id.TunnelInterfaceName),
		fmt.Sprintf("Backend Address Pool Name %q", id.BackendAddressPoolName),
		fmt.Sprintf("Load Balancer Name %q", id.LoadBalancerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Backend Address Pool Tunnel Interface", segmentsStr)
}

func (id BackendAddressPoolTunnelInterfaceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/backendAddressPools/%s/tunnelInterfaces/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, id.TunnelInterfaceName)
}

// BackendAddressPoolTunnelInterfaceID parses a BackendAddressPoolTunnelInterface ID into an BackendAddressPoolTunnelInterfaceId struct
func BackendAddressPoolTunnelInterfaceID(input string) (*BackendAddressPoolTunnelInterfaceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := BackendAddressPoolTunnelInterfaceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.LoadBalancerName, err = id.PopSegment("loadBalancers"); err != nil {
		return nil, err
	}
	if resourceId.BackendAddressPoolName, err = id.PopSegment("backendAddressPools"); err != nil {
		return nil, err
	}
	if resourceId.TunnelInterfaceName, err = id.PopSegment("tunnelInterfaces"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
