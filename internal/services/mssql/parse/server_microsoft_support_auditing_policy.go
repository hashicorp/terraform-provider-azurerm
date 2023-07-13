// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ServerMicrosoftSupportAuditingPolicyId struct {
	SubscriptionId            string
	ResourceGroup             string
	ServerName                string
	DevOpsAuditingSettingName string
}

func NewServerMicrosoftSupportAuditingPolicyID(subscriptionId, resourceGroup, serverName, devOpsAuditingSettingName string) ServerMicrosoftSupportAuditingPolicyId {
	return ServerMicrosoftSupportAuditingPolicyId{
		SubscriptionId:            subscriptionId,
		ResourceGroup:             resourceGroup,
		ServerName:                serverName,
		DevOpsAuditingSettingName: devOpsAuditingSettingName,
	}
}

func (id ServerMicrosoftSupportAuditingPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Dev Ops Auditing Setting Name %q", id.DevOpsAuditingSettingName),
		fmt.Sprintf("Server Name %q", id.ServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Server Microsoft Support Auditing Policy", segmentsStr)
}

func (id ServerMicrosoftSupportAuditingPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/devOpsAuditingSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerName, id.DevOpsAuditingSettingName)
}

// ServerMicrosoftSupportAuditingPolicyID parses a ServerMicrosoftSupportAuditingPolicy ID into an ServerMicrosoftSupportAuditingPolicyId struct
func ServerMicrosoftSupportAuditingPolicyID(input string) (*ServerMicrosoftSupportAuditingPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ServerMicrosoftSupportAuditingPolicy ID: %+v", input, err)
	}

	resourceId := ServerMicrosoftSupportAuditingPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}
	if resourceId.DevOpsAuditingSettingName, err = id.PopSegment("devOpsAuditingSettings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
