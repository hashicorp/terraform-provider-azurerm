package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AutomanageConfigurationProfileHCIAssignmentId struct {
	SubscriptionId                     string
	ResourceGroup                      string
	ClusterName                        string
	ConfigurationProfileAssignmentName string
}

func NewAutomanageConfigurationProfileHCIAssignmentID(subscriptionId, resourceGroup, clusterName, configurationProfileAssignmentName string) AutomanageConfigurationProfileHCIAssignmentId {
	return AutomanageConfigurationProfileHCIAssignmentId{
		SubscriptionId:                     subscriptionId,
		ResourceGroup:                      resourceGroup,
		ClusterName:                        clusterName,
		ConfigurationProfileAssignmentName: configurationProfileAssignmentName,
	}
}

func (id AutomanageConfigurationProfileHCIAssignmentId) String() string {
	segments := []string{
		fmt.Sprintf("Configuration Profile Assignment Name %q", id.ConfigurationProfileAssignmentName),
		fmt.Sprintf("Cluster Name %q", id.ClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Automanage Configuration Profile H C I Assignment", segmentsStr)
}

func (id AutomanageConfigurationProfileHCIAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AzureStackHci/clusters/%s/providers/Microsoft.Automanage/configurationProfileAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName, id.ConfigurationProfileAssignmentName)
}

// AutomanageConfigurationProfileHCIAssignmentID parses a AutomanageConfigurationProfileHCIAssignment ID into an AutomanageConfigurationProfileHCIAssignmentId struct
func AutomanageConfigurationProfileHCIAssignmentID(input string) (*AutomanageConfigurationProfileHCIAssignmentId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AutomanageConfigurationProfileHCIAssignmentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ClusterName, err = id.PopSegment("clusters"); err != nil {
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
