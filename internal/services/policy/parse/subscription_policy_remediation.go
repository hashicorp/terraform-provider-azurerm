// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SubscriptionPolicyRemediationId struct {
	SubscriptionId  string
	RemediationName string
}

func NewSubscriptionPolicyRemediationID(subscriptionId, remediationName string) SubscriptionPolicyRemediationId {
	return SubscriptionPolicyRemediationId{
		SubscriptionId:  subscriptionId,
		RemediationName: remediationName,
	}
}

func (id SubscriptionPolicyRemediationId) String() string {
	segments := []string{
		fmt.Sprintf("Remediation Name %q", id.RemediationName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Subscription Policy Remediation", segmentsStr)
}

func (id SubscriptionPolicyRemediationId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.PolicyInsights/remediations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.RemediationName)
}

// SubscriptionPolicyRemediationID parses a SubscriptionPolicyRemediation ID into an SubscriptionPolicyRemediationId struct
func SubscriptionPolicyRemediationID(input string) (*SubscriptionPolicyRemediationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SubscriptionPolicyRemediation ID: %+v", input, err)
	}

	resourceId := SubscriptionPolicyRemediationId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.RemediationName, err = id.PopSegment("remediations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
