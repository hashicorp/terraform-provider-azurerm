package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ContainerConnectedRegistryId struct {
	SubscriptionId        string
	ResourceGroup         string
	RegistryName          string
	ConnectedRegistryName string
}

func NewContainerConnectedRegistryID(subscriptionId, resourceGroup, registryName, connectedRegistryName string) ContainerConnectedRegistryId {
	return ContainerConnectedRegistryId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		RegistryName:          registryName,
		ConnectedRegistryName: connectedRegistryName,
	}
}

func (id ContainerConnectedRegistryId) String() string {
	segments := []string{
		fmt.Sprintf("Connected Registry Name %q", id.ConnectedRegistryName),
		fmt.Sprintf("Registry Name %q", id.RegistryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Container Connected Registry", segmentsStr)
}

func (id ContainerConnectedRegistryId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/connectedRegistries/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.ConnectedRegistryName)
}

// ContainerConnectedRegistryID parses a ContainerConnectedRegistry ID into an ContainerConnectedRegistryId struct
func ContainerConnectedRegistryID(input string) (*ContainerConnectedRegistryId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ContainerConnectedRegistryId{
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
	if resourceId.ConnectedRegistryName, err = id.PopSegment("connectedRegistries"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
