// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ApiOperationPolicyId struct {
	SubscriptionId string
	ResourceGroup  string
	ServiceName    string
	ApiName        string
	OperationName  string
	PolicyName     string
}

func NewApiOperationPolicyID(subscriptionId, resourceGroup, serviceName, apiName, operationName, policyName string) ApiOperationPolicyId {
	return ApiOperationPolicyId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ServiceName:    serviceName,
		ApiName:        apiName,
		OperationName:  operationName,
		PolicyName:     policyName,
	}
}

func (id ApiOperationPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Policy Name %q", id.PolicyName),
		fmt.Sprintf("Operation Name %q", id.OperationName),
		fmt.Sprintf("Api Name %q", id.ApiName),
		fmt.Sprintf("Service Name %q", id.ServiceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Api Operation Policy", segmentsStr)
}

func (id ApiOperationPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/apis/%s/operations/%s/policies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServiceName, id.ApiName, id.OperationName, id.PolicyName)
}

// ApiOperationPolicyID parses a ApiOperationPolicy ID into an ApiOperationPolicyId struct
func ApiOperationPolicyID(input string) (*ApiOperationPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ApiOperationPolicy ID: %+v", input, err)
	}

	resourceId := ApiOperationPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServiceName, err = id.PopSegment("service"); err != nil {
		return nil, err
	}
	if resourceId.ApiName, err = id.PopSegment("apis"); err != nil {
		return nil, err
	}
	if resourceId.OperationName, err = id.PopSegment("operations"); err != nil {
		return nil, err
	}
	if resourceId.PolicyName, err = id.PopSegment("policies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
