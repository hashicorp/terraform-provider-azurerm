package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ContainerRegistryAgentPoolId struct {
	SubscriptionId string
	ResourceGroup  string
	RegistryName   string
	AgentPoolName  string
}

func NewContainerRegistryAgentPoolID(subscriptionId, resourceGroup, registryName, agentPoolName string) ContainerRegistryAgentPoolId {
	return ContainerRegistryAgentPoolId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		RegistryName:   registryName,
		AgentPoolName:  agentPoolName,
	}
}

func (id ContainerRegistryAgentPoolId) String() string {
	segments := []string{
		fmt.Sprintf("Agent Pool Name %q", id.AgentPoolName),
		fmt.Sprintf("Registry Name %q", id.RegistryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Container Registry Agent Pool", segmentsStr)
}

func (id ContainerRegistryAgentPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/agentPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.AgentPoolName)
}

// ContainerRegistryAgentPoolID parses a ContainerRegistryAgentPool ID into an ContainerRegistryAgentPoolId struct
func ContainerRegistryAgentPoolID(input string) (*ContainerRegistryAgentPoolId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ContainerRegistryAgentPoolId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.RegistryName, err = id.PopSegment("registries"); err != nil {
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
