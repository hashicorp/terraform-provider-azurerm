package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AutomanageConfigurationProfileAssignmentId struct {
	SubscriptionId                     string
	ResourceGroup                      string
	VirtualMachineName                 string
	ConfigurationProfileAssignmentName string
}

func NewAutomanageConfigurationProfileAssignmentID(subscriptionId, resourceGroup, virtualMachineName, configurationProfileAssignmentName string) AutomanageConfigurationProfileAssignmentId {
	return AutomanageConfigurationProfileAssignmentId{
		SubscriptionId:                     subscriptionId,
		ResourceGroup:                      resourceGroup,
		VirtualMachineName:                 virtualMachineName,
		ConfigurationProfileAssignmentName: configurationProfileAssignmentName,
	}
}

func (id AutomanageConfigurationProfileAssignmentId) String() string {
	segments := []string{
		fmt.Sprintf("Configuration Profile Assignment Name %q", id.ConfigurationProfileAssignmentName),
		fmt.Sprintf("Virtual Machine Name %q", id.VirtualMachineName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Automanage Configuration Profile Assignment", segmentsStr)
}

func (id AutomanageConfigurationProfileAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s/providers/Microsoft.Automanage/configurationProfileAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName, id.ConfigurationProfileAssignmentName)
}

// AutomanageConfigurationProfileAssignmentID parses a AutomanageConfigurationProfileAssignment ID into an AutomanageConfigurationProfileAssignmentId struct
func AutomanageConfigurationProfileAssignmentID(input string) (*AutomanageConfigurationProfileAssignmentId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AutomanageConfigurationProfileAssignmentId{
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
	if resourceId.ConfigurationProfileAssignmentName, err = id.PopSegment("configurationProfileAssignments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
