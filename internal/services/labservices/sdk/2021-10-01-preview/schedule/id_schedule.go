package schedule

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ScheduleId struct {
	SubscriptionId string
	ResourceGroup  string
	LabName        string
	Name           string
}

func NewScheduleID(subscriptionId, resourceGroup, labName, name string) ScheduleId {
	return ScheduleId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		LabName:        labName,
		Name:           name,
	}
}

func (id ScheduleId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Lab Name %q", id.LabName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Schedule", segmentsStr)
}

func (id ScheduleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.LabServices/labs/%s/schedules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LabName, id.Name)
}

// ParseScheduleID parses a Schedule ID into an ScheduleId struct
func ParseScheduleID(input string) (*ScheduleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ScheduleId{
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
	if resourceId.Name, err = id.PopSegment("schedules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseScheduleIDInsensitively parses an Schedule ID into an ScheduleId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseScheduleID method should be used instead for validation etc.
func ParseScheduleIDInsensitively(input string) (*ScheduleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ScheduleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'labs' segment
	labsKey := "labs"
	for key := range id.Path {
		if strings.EqualFold(key, labsKey) {
			labsKey = key
			break
		}
	}
	if resourceId.LabName, err = id.PopSegment(labsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'schedules' segment
	schedulesKey := "schedules"
	for key := range id.Path {
		if strings.EqualFold(key, schedulesKey) {
			schedulesKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(schedulesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
