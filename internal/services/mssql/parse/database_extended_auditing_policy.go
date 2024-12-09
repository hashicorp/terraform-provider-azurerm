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

type DatabaseExtendedAuditingPolicyId struct {
	SubscriptionId              string
	ResourceGroup               string
	ServerName                  string
	DatabaseName                string
	ExtendedAuditingSettingName string
}

func NewDatabaseExtendedAuditingPolicyID(subscriptionId, resourceGroup, serverName, databaseName, extendedAuditingSettingName string) DatabaseExtendedAuditingPolicyId {
	return DatabaseExtendedAuditingPolicyId{
		SubscriptionId:              subscriptionId,
		ResourceGroup:               resourceGroup,
		ServerName:                  serverName,
		DatabaseName:                databaseName,
		ExtendedAuditingSettingName: extendedAuditingSettingName,
	}
}

func (id DatabaseExtendedAuditingPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Extended Auditing Setting Name %q", id.ExtendedAuditingSettingName),
		fmt.Sprintf("Database Name %q", id.DatabaseName),
		fmt.Sprintf("Server Name %q", id.ServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Database Extended Auditing Policy", segmentsStr)
}

func (id DatabaseExtendedAuditingPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/servers/%s/databases/%s/extendedAuditingSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerName, id.DatabaseName, id.ExtendedAuditingSettingName)
}

// DatabaseExtendedAuditingPolicyID parses a DatabaseExtendedAuditingPolicy ID into an DatabaseExtendedAuditingPolicyId struct
func DatabaseExtendedAuditingPolicyID(input string) (*DatabaseExtendedAuditingPolicyId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an DatabaseExtendedAuditingPolicy ID: %+v", input, err)
	}

	resourceId := DatabaseExtendedAuditingPolicyId{
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
	if resourceId.DatabaseName, err = id.PopSegment("databases"); err != nil {
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
