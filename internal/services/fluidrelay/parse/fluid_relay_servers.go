package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FluidRelayServersId struct {
	SubscriptionId       string
	ResourceGroup        string
	FluidRelayServerName string
}

func NewFluidRelayServersID(subscriptionId, resourceGroup, fluidRelayServerName string) FluidRelayServersId {
	return FluidRelayServersId{
		SubscriptionId:       subscriptionId,
		ResourceGroup:        resourceGroup,
		FluidRelayServerName: fluidRelayServerName,
	}
}

func (id FluidRelayServersId) String() string {
	segments := []string{
		fmt.Sprintf("Fluid Relay Server Name %q", id.FluidRelayServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Fluid Relay Servers", segmentsStr)
}

func (id FluidRelayServersId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.FluidRelay/fluidRelayServers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FluidRelayServerName)
}

// FluidRelayServersID parses a FluidRelayServers ID into an FluidRelayServersId struct
func FluidRelayServersID(input string) (*FluidRelayServersId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FluidRelayServersId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.FluidRelayServerName, err = id.PopSegment("fluidRelayServers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
