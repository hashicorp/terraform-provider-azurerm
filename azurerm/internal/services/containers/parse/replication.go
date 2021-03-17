package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ReplicationId struct {
	SubscriptionId string
	ResourceGroup  string
	RegistryName   string
	Name           string
}

func NewReplicationID(subscriptionId, resourceGroup, registryName, name string) ReplicationId {
	return ReplicationId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		RegistryName:   registryName,
		Name:           name,
	}
}

func (id ReplicationId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Registry Name %q", id.RegistryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Replication", segmentsStr)
}

func (id ReplicationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/replications/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.Name)
}

// ReplicationID parses a Replication ID into an ReplicationId struct
func ReplicationID(input string) (*ReplicationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ReplicationId{
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
	if resourceId.Name, err = id.PopSegment("replications"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
