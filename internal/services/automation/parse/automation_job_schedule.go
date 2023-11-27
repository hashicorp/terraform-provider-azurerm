// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AutomationJobScheduleId struct {
	SubscriptionId        string
	ResourceGroup         string
	AutomationAccountName string
	RunBookName           string
	ScheduleName          string
}

func NewAutomationJobScheduleID(subscriptionId, resourceGroup, automationAccountName, runBookName, scheduleName string) AutomationJobScheduleId {
	return AutomationJobScheduleId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		AutomationAccountName: automationAccountName,
		RunBookName:           runBookName,
		ScheduleName:          scheduleName,
	}
}

func (id AutomationJobScheduleId) String() string {
	segments := []string{
		fmt.Sprintf("Schedule Name %q", id.ScheduleName),
		fmt.Sprintf("Run Book Name %q", id.RunBookName),
		fmt.Sprintf("Automation Account Name %q", id.AutomationAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Automation Job Schedule", segmentsStr)
}

func (id AutomationJobScheduleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/runBook/%s/schedule/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AutomationAccountName, id.RunBookName, id.ScheduleName)
}

// AutomationJobScheduleID parses a AutomationJobSchedule ID into an AutomationJobScheduleId struct
func AutomationJobScheduleID(input string) (*AutomationJobScheduleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an AutomationJobSchedule ID: %+v", input, err)
	}

	resourceId := AutomationJobScheduleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.AutomationAccountName, err = id.PopSegment("automationAccounts"); err != nil {
		return nil, err
	}
	if resourceId.RunBookName, err = id.PopSegment("runBook"); err != nil {
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
