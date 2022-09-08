package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FluidRelayId struct {
	SubscriptionId       string
	ResourceGroup        string
	FluidRelayServerName string
}

func NewFluidRelayID(subscriptionId, resourceGroup, fluidRelayServerName string) FluidRelayId {
	return FluidRelayId{
		SubscriptionId:       subscriptionId,
		ResourceGroup:        resourceGroup,
		FluidRelayServerName: fluidRelayServerName,
	}
}

func (id FluidRelayId) String() string {
	segments := []string{
		fmt.Sprintf("Fluid Relay Server Name %q", id.FluidRelayServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Fluid Relay", segmentsStr)
}

func (id FluidRelayId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.FluidRelay/fluidRelayServers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FluidRelayServerName)
}

// FluidRelayID parses a FluidRelay ID into an FluidRelayId struct
func FluidRelayID(input string) (*FluidRelayId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FluidRelayId{
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
