// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SmartDetectorAlertRuleId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewSmartDetectorAlertRuleID(subscriptionId, resourceGroup, name string) SmartDetectorAlertRuleId {
	return SmartDetectorAlertRuleId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id SmartDetectorAlertRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Smart Detector Alert Rule", segmentsStr)
}

func (id SmartDetectorAlertRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AlertsManagement/smartDetectorAlertRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// SmartDetectorAlertRuleID parses a SmartDetectorAlertRule ID into an SmartDetectorAlertRuleId struct
func SmartDetectorAlertRuleID(input string) (*SmartDetectorAlertRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SmartDetectorAlertRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.Name, err = id.PopSegment("smartDetectorAlertRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SmartDetectorAlertRuleIDInsensitively parses an SmartDetectorAlertRule ID into an SmartDetectorAlertRuleId struct, insensitively
// This should only be used to parse an ID for rewriting, the SmartDetectorAlertRuleID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SmartDetectorAlertRuleIDInsensitively(input string) (*SmartDetectorAlertRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SmartDetectorAlertRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'smartDetectorAlertRules' segment
	smartDetectorAlertRulesKey := "smartDetectorAlertRules"
	for key := range id.Path {
		if strings.EqualFold(key, smartDetectorAlertRulesKey) {
			smartDetectorAlertRulesKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(smartDetectorAlertRulesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
