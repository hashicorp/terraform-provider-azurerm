package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ContainerRegistryTaskId struct {
	SubscriptionId string
	ResourceGroup  string
	RegistryName   string
	TaskName       string
}

func NewContainerRegistryTaskID(subscriptionId, resourceGroup, registryName, taskName string) ContainerRegistryTaskId {
	return ContainerRegistryTaskId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		RegistryName:   registryName,
		TaskName:       taskName,
	}
}

func (id ContainerRegistryTaskId) String() string {
	segments := []string{
		fmt.Sprintf("Task Name %q", id.TaskName),
		fmt.Sprintf("Registry Name %q", id.RegistryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Container Registry Task", segmentsStr)
}

func (id ContainerRegistryTaskId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/tasks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.TaskName)
}

// ContainerRegistryTaskID parses a ContainerRegistryTask ID into an ContainerRegistryTaskId struct
func ContainerRegistryTaskID(input string) (*ContainerRegistryTaskId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ContainerRegistryTaskId{
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
	if resourceId.TaskName, err = id.PopSegment("tasks"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
