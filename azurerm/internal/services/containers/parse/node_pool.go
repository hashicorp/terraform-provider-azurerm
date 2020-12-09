package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NodePoolId struct {
	SubscriptionId     string
	ResourceGroup      string
	ManagedClusterName string
	AgentPoolName      string
}

func NewNodePoolID(subscriptionId, resourceGroup, managedClusterName, agentPoolName string) NodePoolId {
	return NodePoolId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		ManagedClusterName: managedClusterName,
		AgentPoolName:      agentPoolName,
	}
}

func (id NodePoolId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Managed Cluster Name %q", id.ManagedClusterName),
		fmt.Sprintf("Agent Pool Name %q", id.AgentPoolName),
	}
	return strings.Join(segments, " / ")
}

func (id NodePoolId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/managedClusters/%s/agentPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ManagedClusterName, id.AgentPoolName)
}

// NodePoolID parses a NodePool ID into an NodePoolId struct
func NodePoolID(input string) (*NodePoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := NodePoolId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ManagedClusterName, err = id.PopSegment("managedClusters"); err != nil {
		return nil, err
	}
	if resourceId.AgentPoolName, err = id.PopSegment("agentPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
