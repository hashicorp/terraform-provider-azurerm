package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ClusterPrincipalAssignmentId struct {
	SubscriptionId          string
	ResourceGroup           string
	ClusterName             string
	PrincipalAssignmentName string
}

func NewClusterPrincipalAssignmentID(subscriptionId, resourceGroup, clusterName, principalAssignmentName string) ClusterPrincipalAssignmentId {
	return ClusterPrincipalAssignmentId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ClusterName:             clusterName,
		PrincipalAssignmentName: principalAssignmentName,
	}
}

func (id ClusterPrincipalAssignmentId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Kusto/Clusters/%s/PrincipalAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName, id.PrincipalAssignmentName)
}

// ClusterPrincipalAssignmentID parses a ClusterPrincipalAssignment ID into an ClusterPrincipalAssignmentId struct
func ClusterPrincipalAssignmentID(input string) (*ClusterPrincipalAssignmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ClusterPrincipalAssignmentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ClusterName, err = id.PopSegment("Clusters"); err != nil {
		return nil, err
	}
	if resourceId.PrincipalAssignmentName, err = id.PopSegment("PrincipalAssignments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
