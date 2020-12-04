package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LoadBalancerProbeId struct {
	SubscriptionId   string
	ResourceGroup    string
	LoadBalancerName string
	ProbeName        string
}

func NewLoadBalancerProbeID(subscriptionId, resourceGroup, loadBalancerName, probeName string) LoadBalancerProbeId {
	return LoadBalancerProbeId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		LoadBalancerName: loadBalancerName,
		ProbeName:        probeName,
	}
}

func (id LoadBalancerProbeId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/probes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.ProbeName)
}

// LoadBalancerProbeID parses a LoadBalancerProbe ID into an LoadBalancerProbeId struct
func LoadBalancerProbeID(input string) (*LoadBalancerProbeId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LoadBalancerProbeId{
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
	if resourceId.ProbeName, err = id.PopSegment("probes"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
