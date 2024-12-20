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

type ServerExtendedAuditingPolicyId struct {
	SubscriptionId              string
	ResourceGroup               string
	ServerName                  string
	ExtendedAuditingSettingName string
}

func NewServerExtendedAuditingPolicyID(subscriptionId, resourceGroup, serverName, extendedAuditingSettingName string) ServerExtendedAuditingPolicyId {
	return ServerExtendedAuditingPolicyId{
		SubscriptionId:              subscriptionId,
		ResourceGroup:               resourceGroup,
		ServerName:                  serverName,
		ExtendedAuditingSettingName: extendedAuditingSettingName,
	}
}

func (id ServerExtendedAuditingPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Extended Auditing Setting Name %q", id.ExtendedAuditingSettingName),
		fmt.Sprintf("Server Name %q", id.ServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Server Extended Auditing Policy", segmentsStr)
}

func (id ServerExtendedAuditingPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/extendedAuditingSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerName, id.ExtendedAuditingSettingName)
}

// ServerExtendedAuditingPolicyID parses a ServerExtendedAuditingPolicy ID into an ServerExtendedAuditingPolicyId struct
func ServerExtendedAuditingPolicyID(input string) (*ServerExtendedAuditingPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ServerExtendedAuditingPolicy ID: %+v", input, err)
	}

	resourceId := ServerExtendedAuditingPolicyId{
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
	if resourceId.ExtendedAuditingSettingName, err = id.PopSegment("extendedAuditingSettings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
