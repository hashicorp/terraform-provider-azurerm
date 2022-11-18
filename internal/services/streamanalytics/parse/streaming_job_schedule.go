package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type StreamingJobScheduleId struct {
	SubscriptionId   string
	ResourceGroup    string
	StreamingjobName string
	ScheduleName     string
}

func NewStreamingJobScheduleID(subscriptionId, resourceGroup, streamingjobName, scheduleName string) StreamingJobScheduleId {
	return StreamingJobScheduleId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		StreamingjobName: streamingjobName,
		ScheduleName:     scheduleName,
	}
}

func (id StreamingJobScheduleId) String() string {
	segments := []string{
		fmt.Sprintf("Schedule Name %q", id.ScheduleName),
		fmt.Sprintf("Streamingjob Name %q", id.JobName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroupName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Streaming Job Schedule", segmentsStr)
}

func (id StreamingJobScheduleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StreamAnalytics/streamingjobs/%s/schedule/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.JobName, id.ScheduleName)
}

// StreamingJobScheduleID parses a StreamingJobSchedule ID into an StreamingJobScheduleId struct
func StreamingJobScheduleID(input string) (*StreamingJobScheduleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := StreamingJobScheduleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroupName,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceid.ResourceGroupName == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceid.JobName, err = id.PopSegment("streamingjobs"); err != nil {
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
