package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AutomanageConfigurationProfileHCRPAssignmentId struct {
	SubscriptionId                     string
	ResourceGroup                      string
	MachineName                        string
	ConfigurationProfileAssignmentName string
}

func NewAutomanageConfigurationProfileHCRPAssignmentID(subscriptionId, resourceGroup, machineName, configurationProfileAssignmentName string) AutomanageConfigurationProfileHCRPAssignmentId {
	return AutomanageConfigurationProfileHCRPAssignmentId{
		SubscriptionId:                     subscriptionId,
		ResourceGroup:                      resourceGroup,
		MachineName:                        machineName,
		ConfigurationProfileAssignmentName: configurationProfileAssignmentName,
	}
}

func (id AutomanageConfigurationProfileHCRPAssignmentId) String() string {
	segments := []string{
		fmt.Sprintf("Configuration Profile Assignment Name %q", id.ConfigurationProfileAssignmentName),
		fmt.Sprintf("Machine Name %q", id.MachineName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Automanage Configuration Profile H C R P Assignment", segmentsStr)
}

func (id AutomanageConfigurationProfileHCRPAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.HybridCompute/machines/%s/providers/Microsoft.Automanage/configurationProfileAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.MachineName, id.ConfigurationProfileAssignmentName)
}

// AutomanageConfigurationProfileHCRPAssignmentID parses a AutomanageConfigurationProfileHCRPAssignment ID into an AutomanageConfigurationProfileHCRPAssignmentId struct
func AutomanageConfigurationProfileHCRPAssignmentID(input string) (*AutomanageConfigurationProfileHCRPAssignmentId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AutomanageConfigurationProfileHCRPAssignmentId{
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
	if resourceId.ConfigurationProfileAssignmentName, err = id.PopSegment("configurationProfileAssignments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
