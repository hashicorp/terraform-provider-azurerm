package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AnomalyAlertViewIdId struct {
	SubscriptionId string
	ViewName       string
}

func NewAnomalyAlertViewIdID(subscriptionId, viewName string) AnomalyAlertViewIdId {
	return AnomalyAlertViewIdId{
		SubscriptionId: subscriptionId,
		ViewName:       viewName,
	}
}

func (id AnomalyAlertViewIdId) String() string {
	segments := []string{
		fmt.Sprintf("View Name %q", id.ViewName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Anomaly Alert View Id", segmentsStr)
}

func (id AnomalyAlertViewIdId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.CostManagement/views/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ViewName)
}

// AnomalyAlertViewIdID parses a AnomalyAlertViewId ID into an AnomalyAlertViewIdId struct
func AnomalyAlertViewIdID(input string) (*AnomalyAlertViewIdId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AnomalyAlertViewIdId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ViewName, err = id.PopSegment("views"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
