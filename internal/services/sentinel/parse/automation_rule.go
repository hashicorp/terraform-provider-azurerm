// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AutomationRuleId struct {
	SubscriptionId string
	ResourceGroup  string
	WorkspaceName  string
	Name           string
}

func NewAutomationRuleID(subscriptionId, resourceGroup, workspaceName, name string) AutomationRuleId {
	return AutomationRuleId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		WorkspaceName:  workspaceName,
		Name:           name,
	}
}

func (id AutomationRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Automation Rule", segmentsStr)
}

func (id AutomationRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/automationRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.Name)
}

// AutomationRuleID parses a AutomationRule ID into an AutomationRuleId struct
func AutomationRuleID(input string) (*AutomationRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an AutomationRule ID: %+v", input, err)
	}

	resourceId := AutomationRuleId{
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
	if resourceId.Name, err = id.PopSegment("automationRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// AutomationRuleIDInsensitively parses an AutomationRule ID into an AutomationRuleId struct, insensitively
// This should only be used to parse an ID for rewriting, the AutomationRuleID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func AutomationRuleIDInsensitively(input string) (*AutomationRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AutomationRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'workspaces' segment
	workspacesKey := "workspaces"
	for key := range id.Path {
		if strings.EqualFold(key, workspacesKey) {
			workspacesKey = key
			break
		}
	}
	if resourceId.WorkspaceName, err = id.PopSegment(workspacesKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'automationRules' segment
	automationRulesKey := "automationRules"
	for key := range id.Path {
		if strings.EqualFold(key, automationRulesKey) {
			automationRulesKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(automationRulesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
