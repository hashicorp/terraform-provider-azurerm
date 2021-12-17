package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type JobScheduleId struct {
	SubscriptionId        string
	ResourceGroup         string
	AutomationAccountName string
	Name                  string
}

func NewJobScheduleID(subscriptionId, resourceGroup, automationAccountName, name string) JobScheduleId {
	return JobScheduleId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		AutomationAccountName: automationAccountName,
		Name:                  name,
	}
}

func (id JobScheduleId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Automation Account Name %q", id.AutomationAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Job Schedule", segmentsStr)
}

func (id JobScheduleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/jobSchedules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AutomationAccountName, id.Name)
}

// JobScheduleID parses a JobSchedule ID into an JobScheduleId struct
func JobScheduleID(input string) (*JobScheduleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := JobScheduleId{
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
	if resourceId.Name, err = id.PopSegment("jobSchedules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
