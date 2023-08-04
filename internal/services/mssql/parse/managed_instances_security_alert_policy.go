// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ManagedInstancesSecurityAlertPolicyId struct {
	SubscriptionId          string
	ResourceGroup           string
	ManagedInstanceName     string
	SecurityAlertPolicyName string
}

func NewManagedInstancesSecurityAlertPolicyID(subscriptionId, resourceGroup, managedInstanceName, securityAlertPolicyName string) ManagedInstancesSecurityAlertPolicyId {
	return ManagedInstancesSecurityAlertPolicyId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ManagedInstanceName:     managedInstanceName,
		SecurityAlertPolicyName: securityAlertPolicyName,
	}
}

func (id ManagedInstancesSecurityAlertPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Security Alert Policy Name %q", id.SecurityAlertPolicyName),
		fmt.Sprintf("Managed Instance Name %q", id.ManagedInstanceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Managed Instances Security Alert Policy", segmentsStr)
}

func (id ManagedInstancesSecurityAlertPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/managedInstances/%s/securityAlertPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName, id.SecurityAlertPolicyName)
}

// ManagedInstancesSecurityAlertPolicyID parses a ManagedInstancesSecurityAlertPolicy ID into an ManagedInstancesSecurityAlertPolicyId struct
func ManagedInstancesSecurityAlertPolicyID(input string) (*ManagedInstancesSecurityAlertPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ManagedInstancesSecurityAlertPolicy ID: %+v", input, err)
	}

	resourceId := ManagedInstancesSecurityAlertPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ManagedInstanceName, err = id.PopSegment("managedInstances"); err != nil {
		return nil, err
	}
	if resourceId.SecurityAlertPolicyName, err = id.PopSegment("securityAlertPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
