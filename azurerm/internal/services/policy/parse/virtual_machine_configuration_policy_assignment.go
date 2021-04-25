package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VirtualMachineConfigurationPolicyAssignmentId struct {
	SubscriptionId                   string
	ResourceGroup                    string
	VirtualMachineName               string
	GuestConfigurationAssignmentName string
}

func NewVirtualMachineConfigurationPolicyAssignmentID(subscriptionId, resourceGroup, virtualMachineName, guestConfigurationAssignmentName string) VirtualMachineConfigurationPolicyAssignmentId {
	return VirtualMachineConfigurationPolicyAssignmentId{
		SubscriptionId:                   subscriptionId,
		ResourceGroup:                    resourceGroup,
		VirtualMachineName:               virtualMachineName,
		GuestConfigurationAssignmentName: guestConfigurationAssignmentName,
	}
}

func (id VirtualMachineConfigurationPolicyAssignmentId) String() string {
	segments := []string{
		fmt.Sprintf("Guest Configuration Assignment Name %q", id.GuestConfigurationAssignmentName),
		fmt.Sprintf("Virtual Machine Name %q", id.VirtualMachineName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Virtual Machine Configuration Policy Assignment", segmentsStr)
}

func (id VirtualMachineConfigurationPolicyAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName, id.GuestConfigurationAssignmentName)
}

// VirtualMachineConfigurationPolicyAssignmentID parses a VirtualMachineConfigurationPolicyAssignment ID into an VirtualMachineConfigurationPolicyAssignmentId struct
func VirtualMachineConfigurationPolicyAssignmentID(input string) (*VirtualMachineConfigurationPolicyAssignmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VirtualMachineConfigurationPolicyAssignmentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VirtualMachineName, err = id.PopSegment("virtualMachines"); err != nil {
		return nil, err
	}
	if resourceId.GuestConfigurationAssignmentName, err = id.PopSegment("guestConfigurationAssignments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// VirtualMachineConfigurationPolicyAssignmentIDInsensitively parses an VirtualMachineConfigurationPolicyAssignment ID into an VirtualMachineConfigurationPolicyAssignmentId struct, insensitively
// This should only be used to parse an ID for rewriting, the VirtualMachineConfigurationPolicyAssignmentID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func VirtualMachineConfigurationPolicyAssignmentIDInsensitively(input string) (*VirtualMachineConfigurationPolicyAssignmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VirtualMachineConfigurationPolicyAssignmentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
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

	// find the correct casing for the 'guestConfigurationAssignments' segment
	guestConfigurationAssignmentsKey := "guestConfigurationAssignments"
	for key := range id.Path {
		if strings.EqualFold(key, guestConfigurationAssignmentsKey) {
			guestConfigurationAssignmentsKey = key
			break
		}
	}
	if resourceId.GuestConfigurationAssignmentName, err = id.PopSegment(guestConfigurationAssignmentsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
