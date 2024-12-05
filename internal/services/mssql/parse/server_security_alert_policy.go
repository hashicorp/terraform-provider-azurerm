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

type ServerSecurityAlertPolicyId struct {
	SubscriptionId          string
	ResourceGroup           string
	ServerName              string
	SecurityAlertPolicyName string
}

func NewServerSecurityAlertPolicyID(subscriptionId, resourceGroup, serverName, securityAlertPolicyName string) ServerSecurityAlertPolicyId {
	return ServerSecurityAlertPolicyId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ServerName:              serverName,
		SecurityAlertPolicyName: securityAlertPolicyName,
	}
}

func (id ServerSecurityAlertPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Security Alert Policy Name %q", id.SecurityAlertPolicyName),
		fmt.Sprintf("Server Name %q", id.ServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Server Security Alert Policy", segmentsStr)
}

func (id ServerSecurityAlertPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/securityAlertPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerName, id.SecurityAlertPolicyName)
}

// ServerSecurityAlertPolicyID parses a ServerSecurityAlertPolicy ID into an ServerSecurityAlertPolicyId struct
func ServerSecurityAlertPolicyID(input string) (*ServerSecurityAlertPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ServerSecurityAlertPolicy ID: %+v", input, err)
	}

	resourceId := ServerSecurityAlertPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, errors.New("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, errors.New("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServerName, err = id.PopSegment("servers"); err != nil {
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
