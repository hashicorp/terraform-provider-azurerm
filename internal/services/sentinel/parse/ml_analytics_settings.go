// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type MLAnalyticsSettingsId struct {
	SubscriptionId                 string
	ResourceGroup                  string
	WorkspaceName                  string
	SecurityMLAnalyticsSettingName string
}

func NewMLAnalyticsSettingsID(subscriptionId, resourceGroup, workspaceName, securityMLAnalyticsSettingName string) MLAnalyticsSettingsId {
	return MLAnalyticsSettingsId{
		SubscriptionId:                 subscriptionId,
		ResourceGroup:                  resourceGroup,
		WorkspaceName:                  workspaceName,
		SecurityMLAnalyticsSettingName: securityMLAnalyticsSettingName,
	}
}

func (id MLAnalyticsSettingsId) String() string {
	segments := []string{
		fmt.Sprintf("Security M L Analytics Setting Name %q", id.SecurityMLAnalyticsSettingName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "M L Analytics Settings", segmentsStr)
}

func (id MLAnalyticsSettingsId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/securityMLAnalyticsSettings/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.SecurityMLAnalyticsSettingName)
}

// MLAnalyticsSettingsID parses a MLAnalyticsSettings ID into an MLAnalyticsSettingsId struct
func MLAnalyticsSettingsID(input string) (*MLAnalyticsSettingsId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an MLAnalyticsSettings ID: %+v", input, err)
	}

	resourceId := MLAnalyticsSettingsId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if resourceId.SecurityMLAnalyticsSettingName, err = id.PopSegment("securityMLAnalyticsSettings"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
