package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type GuestConfigurationAssignmentId struct {
	SubscriptionId string
	ResourceGroup  string
	VMName         string
	Name           string
}

func NewGuestConfigurationAssignmentID(subscriptionId, resourceGroup, vmName, name string) GuestConfigurationAssignmentId {
	return GuestConfigurationAssignmentId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		VMName:         vmName,
		Name:           name,
	}
}

func (id GuestConfigurationAssignmentId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("VM Name %q", id.VMName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Guest Configuration Assignment", segmentsStr)
}

func (id GuestConfigurationAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VMName, id.Name)
}

func GuestConfigurationAssignmentID(input string) (*GuestConfigurationAssignmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing guestConfigurationAssignment ID %q: %+v", input, err)
	}

	guestConfigurationAssignment := GuestConfigurationAssignmentId{
		ResourceGroup: id.ResourceGroup,
	}
	if guestConfigurationAssignment.VMName, err = id.PopSegment("virtualMachines"); err != nil {
		return nil, err
	}
	if guestConfigurationAssignment.Name, err = id.PopSegment("guestConfigurationAssignments"); err != nil {
		return nil, err
	}
	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &guestConfigurationAssignment, nil
}
