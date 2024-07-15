// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ManagedInstanceStartStopScheduleId struct {
	SubscriptionId        string
	ResourceGroup         string
	ManagedInstanceName   string
	StartStopScheduleName string
}

func NewManagedInstanceStartStopScheduleID(subscriptionId, resourceGroup, managedInstanceName, startStopScheduleName string) ManagedInstanceStartStopScheduleId {
	return ManagedInstanceStartStopScheduleId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		ManagedInstanceName:   managedInstanceName,
		StartStopScheduleName: startStopScheduleName,
	}
}

func (id ManagedInstanceStartStopScheduleId) String() string {
	segments := []string{
		fmt.Sprintf("Start Stop Schedule Name %q", id.StartStopScheduleName),
		fmt.Sprintf("Managed Instance Name %q", id.ManagedInstanceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Managed Instance Start Stop Schedule", segmentsStr)
}

func (id ManagedInstanceStartStopScheduleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/managedInstances/%s/startStopSchedules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName, id.StartStopScheduleName)
}

// ManagedInstanceStartStopScheduleID parses a ManagedInstanceStartStopSchedule ID into an ManagedInstanceStartStopScheduleId struct
func ManagedInstanceStartStopScheduleID(input string) (*ManagedInstanceStartStopScheduleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ManagedInstanceStartStopSchedule ID: %+v", input, err)
	}

	resourceId := ManagedInstanceStartStopScheduleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ManagedInstanceName, err = id.PopSegment("managedInstances"); err != nil {
		return nil, err
	}
	if resourceId.StartStopScheduleName, err = id.PopSegment("startStopSchedules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
