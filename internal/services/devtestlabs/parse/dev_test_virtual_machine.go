package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type DevTestVirtualMachineId struct {
	SubscriptionId     string
	ResourceGroup      string
	LabName            string
	VirtualMachineName string
}

func NewDevTestVirtualMachineID(subscriptionId, resourceGroup, labName, virtualMachineName string) DevTestVirtualMachineId {
	return DevTestVirtualMachineId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		LabName:            labName,
		VirtualMachineName: virtualMachineName,
	}
}

func (id DevTestVirtualMachineId) String() string {
	segments := []string{
		fmt.Sprintf("Virtual Machine Name %q", id.VirtualMachineName),
		fmt.Sprintf("Lab Name %q", id.LabName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Dev Test Virtual Machine", segmentsStr)
}

func (id DevTestVirtualMachineId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevTestLab/labs/%s/virtualMachines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LabName, id.VirtualMachineName)
}

// DevTestVirtualMachineID parses a DevTestVirtualMachine ID into an DevTestVirtualMachineId struct
func DevTestVirtualMachineID(input string) (*DevTestVirtualMachineId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DevTestVirtualMachineId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.LabName, err = id.PopSegment("labs"); err != nil {
		return nil, err
	}
	if resourceId.VirtualMachineName, err = id.PopSegment("virtualMachines"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// DevTestVirtualMachineIDInsensitively parses an DevTestVirtualMachine ID into an DevTestVirtualMachineId struct, insensitively
// This should only be used to parse an ID for rewriting, the DevTestVirtualMachineID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func DevTestVirtualMachineIDInsensitively(input string) (*DevTestVirtualMachineId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DevTestVirtualMachineId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'labs' segment
	labsKey := "labs"
	for key := range id.Path {
		if strings.EqualFold(key, labsKey) {
			labsKey = key
			break
		}
	}
	if resourceId.LabName, err = id.PopSegment(labsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'virtualMachines' segment
	virtualMachinesKey := "virtualMachines"
	for key := range id.Path {
		if strings.EqualFold(key, virtualMachinesKey) {
			virtualMachinesKey = key
			break
		}
	}
	if resourceId.VirtualMachineName, err = id.PopSegment(virtualMachinesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
