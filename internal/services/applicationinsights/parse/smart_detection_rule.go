// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SmartDetectionRuleId struct {
	SubscriptionId         string
	ResourceGroup          string
	ComponentName          string
	SmartDetectionRuleName string
}

func NewSmartDetectionRuleID(subscriptionId, resourceGroup, componentName, smartDetectionRuleName string) SmartDetectionRuleId {
	return SmartDetectionRuleId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		ComponentName:          componentName,
		SmartDetectionRuleName: smartDetectionRuleName,
	}
}

func (id SmartDetectionRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Smart Detection Rule Name %q", id.SmartDetectionRuleName),
		fmt.Sprintf("Component Name %q", id.ComponentName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Smart Detection Rule", segmentsStr)
}

func (id SmartDetectionRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/components/%s/smartDetectionRule/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ComponentName, id.SmartDetectionRuleName)
}

// SmartDetectionRuleID parses a SmartDetectionRule ID into an SmartDetectionRuleId struct
func SmartDetectionRuleID(input string) (*SmartDetectionRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SmartDetectionRule ID: %+v", input, err)
	}

	resourceId := SmartDetectionRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ComponentName, err = id.PopSegment("components"); err != nil {
		return nil, err
	}
	if resourceId.SmartDetectionRuleName, err = id.PopSegment("smartDetectionRule"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SmartDetectionRuleIDInsensitively parses an SmartDetectionRule ID into an SmartDetectionRuleId struct, insensitively
// This should only be used to parse an ID for rewriting, the SmartDetectionRuleID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SmartDetectionRuleIDInsensitively(input string) (*SmartDetectionRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SmartDetectionRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'components' segment
	componentsKey := "components"
	for key := range id.Path {
		if strings.EqualFold(key, componentsKey) {
			componentsKey = key
			break
		}
	}
	if resourceId.ComponentName, err = id.PopSegment(componentsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'smartDetectionRule' segment
	smartDetectionRuleKey := "smartDetectionRule"
	for key := range id.Path {
		if strings.EqualFold(key, smartDetectionRuleKey) {
			smartDetectionRuleKey = key
			break
		}
	}
	if resourceId.SmartDetectionRuleName, err = id.PopSegment(smartDetectionRuleKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
