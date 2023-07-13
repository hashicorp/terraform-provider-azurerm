// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SubscriptionPolicyExemptionId struct {
	SubscriptionId      string
	PolicyExemptionName string
}

func NewSubscriptionPolicyExemptionID(subscriptionId, policyExemptionName string) SubscriptionPolicyExemptionId {
	return SubscriptionPolicyExemptionId{
		SubscriptionId:      subscriptionId,
		PolicyExemptionName: policyExemptionName,
	}
}

func (id SubscriptionPolicyExemptionId) String() string {
	segments := []string{
		fmt.Sprintf("Policy Exemption Name %q", id.PolicyExemptionName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Subscription Policy Exemption", segmentsStr)
}

func (id SubscriptionPolicyExemptionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Authorization/policyExemptions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.PolicyExemptionName)
}

// SubscriptionPolicyExemptionID parses a SubscriptionPolicyExemption ID into an SubscriptionPolicyExemptionId struct
func SubscriptionPolicyExemptionID(input string) (*SubscriptionPolicyExemptionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SubscriptionPolicyExemption ID: %+v", input, err)
	}

	resourceId := SubscriptionPolicyExemptionId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.PolicyExemptionName, err = id.PopSegment("policyExemptions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
