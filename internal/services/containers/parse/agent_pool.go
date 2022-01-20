package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AgentPoolId struct {
	SubscriptionId string
	ResourceGroup  string
	RegistryName   string
	Name           string
}

func NewAgentPoolID(subscriptionId, resourceGroup, registryName, name string) AgentPoolId {
	return AgentPoolId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		RegistryName:   registryName,
		Name:           name,
	}
}

func (id AgentPoolId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Registry Name %q", id.RegistryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Agent Pool", segmentsStr)
}

func (id AgentPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/agentPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.Name)
}

// AgentPoolID parses a AgentPool ID into an AgentPoolId struct
func AgentPoolID(input string) (*AgentPoolId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AgentPoolId{
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
	if resourceId.Name, err = id.PopSegment("agentPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
