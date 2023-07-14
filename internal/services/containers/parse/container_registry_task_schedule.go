// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ContainerRegistryTaskScheduleId struct {
	SubscriptionId string
	ResourceGroup  string
	RegistryName   string
	TaskName       string
	ScheduleName   string
}

func NewContainerRegistryTaskScheduleID(subscriptionId, resourceGroup, registryName, taskName, scheduleName string) ContainerRegistryTaskScheduleId {
	return ContainerRegistryTaskScheduleId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		RegistryName:   registryName,
		TaskName:       taskName,
		ScheduleName:   scheduleName,
	}
}

func (id ContainerRegistryTaskScheduleId) String() string {
	segments := []string{
		fmt.Sprintf("Schedule Name %q", id.ScheduleName),
		fmt.Sprintf("Task Name %q", id.TaskName),
		fmt.Sprintf("Registry Name %q", id.RegistryName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Container Registry Task Schedule", segmentsStr)
}

func (id ContainerRegistryTaskScheduleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s/tasks/%s/schedule/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.TaskName, id.ScheduleName)
}

// ContainerRegistryTaskScheduleID parses a ContainerRegistryTaskSchedule ID into an ContainerRegistryTaskScheduleId struct
func ContainerRegistryTaskScheduleID(input string) (*ContainerRegistryTaskScheduleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ContainerRegistryTaskSchedule ID: %+v", input, err)
	}

	resourceId := ContainerRegistryTaskScheduleId{
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
	if resourceId.ScheduleName, err = id.PopSegment("schedule"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
