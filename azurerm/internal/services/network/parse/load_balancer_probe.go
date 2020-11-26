package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LoadBalancerProbeId struct {
	ResourceGroup    string
	LoadBalancerName string
	Name             string
}

func LoadBalancerProbeID(input string) (*LoadBalancerProbeId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Load Balancer Probe ID %q: %+v", input, err)
	}

	probeId := LoadBalancerProbeId{
		ResourceGroup: id.ResourceGroup,
	}

	if probeId.LoadBalancerName, err = id.PopSegment("loadBalancers"); err != nil {
		return nil, err
	}

	if probeId.Name, err = id.PopSegment("probes"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &probeId, nil
}
