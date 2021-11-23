package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type DevTestLabScheduleId struct {
	SubscriptionId string
	ResourceGroup  string
	LabName        string
	ScheduleName   string
}

func NewDevTestLabScheduleID(subscriptionId, resourceGroup, labName, scheduleName string) DevTestLabScheduleId {
	return DevTestLabScheduleId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		LabName:        labName,
		ScheduleName:   scheduleName,
	}
}

func (id DevTestLabScheduleId) String() string {
	segments := []string{
		fmt.Sprintf("Schedule Name %q", id.ScheduleName),
		fmt.Sprintf("Lab Name %q", id.LabName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Dev Test Lab Schedule", segmentsStr)
}

func (id DevTestLabScheduleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevTestLab/labs/%s/schedules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LabName, id.ScheduleName)
}

// DevTestLabScheduleID parses a DevTestLabSchedule ID into an DevTestLabScheduleId struct
func DevTestLabScheduleID(input string) (*DevTestLabScheduleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := DevTestLabScheduleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.LabName, err = id.PopSegment("labs"); err != nil {
		return nil, err
	}
	if resourceId.ScheduleName, err = id.PopSegment("schedules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
