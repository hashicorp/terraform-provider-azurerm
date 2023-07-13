// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ThreatIntelligenceIndicatorId struct {
	SubscriptionId         string
	ResourceGroup          string
	WorkspaceName          string
	ThreatIntelligenceName string
	IndicatorName          string
}

func NewThreatIntelligenceIndicatorID(subscriptionId, resourceGroup, workspaceName, threatIntelligenceName, indicatorName string) ThreatIntelligenceIndicatorId {
	return ThreatIntelligenceIndicatorId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		WorkspaceName:          workspaceName,
		ThreatIntelligenceName: threatIntelligenceName,
		IndicatorName:          indicatorName,
	}
}

func (id ThreatIntelligenceIndicatorId) String() string {
	segments := []string{
		fmt.Sprintf("Indicator Name %q", id.IndicatorName),
		fmt.Sprintf("Threat Intelligence Name %q", id.ThreatIntelligenceName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Threat Intelligence Indicator", segmentsStr)
}

func (id ThreatIntelligenceIndicatorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/threatIntelligence/%s/indicators/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.ThreatIntelligenceName, id.IndicatorName)
}

// ThreatIntelligenceIndicatorID parses a ThreatIntelligenceIndicator ID into an ThreatIntelligenceIndicatorId struct
func ThreatIntelligenceIndicatorID(input string) (*ThreatIntelligenceIndicatorId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an ThreatIntelligenceIndicator ID: %+v", input, err)
	}

	resourceId := ThreatIntelligenceIndicatorId{
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
	if resourceId.ThreatIntelligenceName, err = id.PopSegment("threatIntelligence"); err != nil {
		return nil, err
	}
	if resourceId.IndicatorName, err = id.PopSegment("indicators"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
