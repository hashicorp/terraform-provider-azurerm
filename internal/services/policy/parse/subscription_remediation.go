package parse

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
)

type SubscriptionRemediationId struct {
	SubscriptionId  commonids.SubscriptionId
	RemediationName string
}

func NewSubscriptionRemediationID(subscriptionId commonids.SubscriptionId, remediationName string) SubscriptionRemediationId {
	return SubscriptionRemediationId{
		SubscriptionId:  subscriptionId,
		RemediationName: remediationName,
	}
}

func ParseSubscriptionRemediationID(input string) (*SubscriptionRemediationId, error) {
	parsed, err := remediations.ParseScopedRemediationID(input)
	if err != nil {
		return nil, err
	}

	subscriptionId, err := commonids.ParseSubscriptionID(parsed.ResourceId)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a Subscription ID: %+v", parsed.ResourceId, err)
	}

	return &SubscriptionRemediationId{
		SubscriptionId:  *subscriptionId,
		RemediationName: parsed.RemediationName,
	}, nil
}

func ParseSubscriptionRemediationIDInsensitively(input string) (*SubscriptionRemediationId, error) {
	parsed, err := remediations.ParseScopedRemediationIDInsensitively(input)
	if err != nil {
		return nil, err
	}

	subscriptionId, err := commonids.ParseSubscriptionIDInsensitively(parsed.ResourceId)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a Subscription ID: %+v", parsed.ResourceId, err)
	}

	return &SubscriptionRemediationId{
		SubscriptionId:  *subscriptionId,
		RemediationName: parsed.RemediationName,
	}, nil
}

func (id SubscriptionRemediationId) ToRemediationID() remediations.ScopedRemediationId {
	return remediations.ScopedRemediationId{
		ResourceId:      id.SubscriptionId.ID(),
		RemediationName: id.RemediationName,
	}
}
