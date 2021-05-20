package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ContainerRegistryTokenId struct {
	SubscriptionId string
	ResourceGroup  string
	RegistryName   string
	TokenName      string
}

func NewContainerRegistryTokenID(subscriptionId, resourceGroup, registryName, tokenName string) ContainerRegistryTokenId {
	return ContainerRegistryTokenId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		RegistryName:   registryName,
		TokenName:      tokenName,
	}
}

func (id ContainerRegistryTokenId) String() string {
	segments := []string{
		fmt.Sprintf("Token Name %q", id.TokenName),
		fmt.Sprintf("Registry Name %q", id.RegistryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Container Registry Token", segmentsStr)
}

func (id ContainerRegistryTokenId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/tokens/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.TokenName)
}

// ContainerRegistryTokenID parses a ContainerRegistryToken ID into an ContainerRegistryTokenId struct
func ContainerRegistryTokenID(input string) (*ContainerRegistryTokenId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ContainerRegistryTokenId{
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
	if resourceId.TokenName, err = id.PopSegment("tokens"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
