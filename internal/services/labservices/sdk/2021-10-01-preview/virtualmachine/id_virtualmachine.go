package virtualmachine

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VirtualMachineId struct {
	SubscriptionId string
	ResourceGroup  string
	LabName        string
	Name           string
}

func NewVirtualMachineID(subscriptionId, resourceGroup, labName, name string) VirtualMachineId {
	return VirtualMachineId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		LabName:        labName,
		Name:           name,
	}
}

func (id VirtualMachineId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Lab Name %q", id.LabName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Virtual Machine", segmentsStr)
}

func (id VirtualMachineId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.LabServices/labs/%s/virtualMachines/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LabName, id.Name)
}

// ParseVirtualMachineID parses a VirtualMachine ID into an VirtualMachineId struct
func ParseVirtualMachineID(input string) (*VirtualMachineId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VirtualMachineId{
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
	if resourceId.Name, err = id.PopSegment("virtualMachines"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseVirtualMachineIDInsensitively parses an VirtualMachine ID into an VirtualMachineId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseVirtualMachineID method should be used instead for validation etc.
func ParseVirtualMachineIDInsensitively(input string) (*VirtualMachineId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VirtualMachineId{
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
	if resourceId.Name, err = id.PopSegment(virtualMachinesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
