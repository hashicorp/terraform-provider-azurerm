package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ContainerRegistryScopeMapId struct {
	SubscriptionId string
	ResourceGroup  string
	RegistryName   string
	ScopeMapName   string
}

func NewContainerRegistryScopeMapID(subscriptionId, resourceGroup, registryName, scopeMapName string) ContainerRegistryScopeMapId {
	return ContainerRegistryScopeMapId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		RegistryName:   registryName,
		ScopeMapName:   scopeMapName,
	}
}

func (id ContainerRegistryScopeMapId) String() string {
	segments := []string{
		fmt.Sprintf("Scope Map Name %q", id.ScopeMapName),
		fmt.Sprintf("Registry Name %q", id.RegistryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Container Registry Scope Map", segmentsStr)
}

func (id ContainerRegistryScopeMapId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/scopeMaps/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.ScopeMapName)
}

// ContainerRegistryScopeMapID parses a ContainerRegistryScopeMap ID into an ContainerRegistryScopeMapId struct
func ContainerRegistryScopeMapID(input string) (*ContainerRegistryScopeMapId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ContainerRegistryScopeMapId{
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
	if resourceId.ScopeMapName, err = id.PopSegment("scopeMaps"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
