package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SubscriptionCostManagementScheduledActionId struct {
	SubscriptionId      string
	ScheduledActionName string
}

func NewSubscriptionCostManagementScheduledActionID(subscriptionId, scheduledActionName string) SubscriptionCostManagementScheduledActionId {
	return SubscriptionCostManagementScheduledActionId{
		SubscriptionId:      subscriptionId,
		ScheduledActionName: scheduledActionName,
	}
}

func (id SubscriptionCostManagementScheduledActionId) String() string {
	segments := []string{
		fmt.Sprintf("Scheduled Action Name %q", id.ScheduledActionName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Subscription Cost Management Scheduled Action", segmentsStr)
}

func (id SubscriptionCostManagementScheduledActionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.CostManagement/scheduledActions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ScheduledActionName)
}

// SubscriptionCostManagementScheduledActionID parses a SubscriptionCostManagementScheduledAction ID into an SubscriptionCostManagementScheduledActionId struct
func SubscriptionCostManagementScheduledActionID(input string) (*SubscriptionCostManagementScheduledActionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SubscriptionCostManagementScheduledActionId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ScheduledActionName, err = id.PopSegment("scheduledActions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
