package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type HybridMachineId struct {
	SubscriptionId string
	ResourceGroup  string
	MachineName    string
}

func NewHybridMachineID(subscriptionId, resourceGroup, machineName string) HybridMachineId {
	return HybridMachineId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		MachineName:    machineName,
	}
}

func (id HybridMachineId) String() string {
	segments := []string{
		fmt.Sprintf("Machine Name %q", id.MachineName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Hybrid Machine", segmentsStr)
}

func (id HybridMachineId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HybridCompute/machines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.MachineName)
}

// HybridMachineID parses a HybridMachine ID into an HybridMachineId struct
func HybridMachineID(input string) (*HybridMachineId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := HybridMachineId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.MachineName, err = id.PopSegment("machines"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
