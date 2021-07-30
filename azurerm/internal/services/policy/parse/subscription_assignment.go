package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SubscriptionAssignmentId struct {
	SubscriptionId       string
	PolicyAssignmentName string
}

func NewSubscriptionAssignmentID(subscriptionId, policyAssignmentName string) SubscriptionAssignmentId {
	return SubscriptionAssignmentId{
		SubscriptionId:       subscriptionId,
		PolicyAssignmentName: policyAssignmentName,
	}
}

func (id SubscriptionAssignmentId) String() string {
	segments := []string{
		fmt.Sprintf("Policy Assignment Name %q", id.PolicyAssignmentName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Subscription Assignment", segmentsStr)
}

func (id SubscriptionAssignmentId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Authorization/policyAssignments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.PolicyAssignmentName)
}

// SubscriptionAssignmentID parses a SubscriptionAssignment ID into an SubscriptionAssignmentId struct
func SubscriptionAssignmentID(input string) (*SubscriptionAssignmentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SubscriptionAssignmentId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.PolicyAssignmentName, err = id.PopSegment("policyAssignments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
