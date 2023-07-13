// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SentinelAlertRuleTemplateId struct {
	SubscriptionId        string
	ResourceGroup         string
	WorkspaceName         string
	AlertRuleTemplateName string
}

func NewSentinelAlertRuleTemplateID(subscriptionId, resourceGroup, workspaceName, alertRuleTemplateName string) SentinelAlertRuleTemplateId {
	return SentinelAlertRuleTemplateId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		WorkspaceName:         workspaceName,
		AlertRuleTemplateName: alertRuleTemplateName,
	}
}

func (id SentinelAlertRuleTemplateId) String() string {
	segments := []string{
		fmt.Sprintf("Alert Rule Template Name %q", id.AlertRuleTemplateName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Sentinel Alert Rule Template", segmentsStr)
}

func (id SentinelAlertRuleTemplateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/alertRuleTemplates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.AlertRuleTemplateName)
}

// SentinelAlertRuleTemplateID parses a SentinelAlertRuleTemplate ID into an SentinelAlertRuleTemplateId struct
func SentinelAlertRuleTemplateID(input string) (*SentinelAlertRuleTemplateId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SentinelAlertRuleTemplate ID: %+v", input, err)
	}

	resourceId := SentinelAlertRuleTemplateId{
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
	if resourceId.AlertRuleTemplateName, err = id.PopSegment("alertRuleTemplates"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// SentinelAlertRuleTemplateIDInsensitively parses an SentinelAlertRuleTemplate ID into an SentinelAlertRuleTemplateId struct, insensitively
// This should only be used to parse an ID for rewriting, the SentinelAlertRuleTemplateID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func SentinelAlertRuleTemplateIDInsensitively(input string) (*SentinelAlertRuleTemplateId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SentinelAlertRuleTemplateId{
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

	// find the correct casing for the 'alertRuleTemplates' segment
	alertRuleTemplatesKey := "alertRuleTemplates"
	for key := range id.Path {
		if strings.EqualFold(key, alertRuleTemplatesKey) {
			alertRuleTemplatesKey = key
			break
		}
	}
	if resourceId.AlertRuleTemplateName, err = id.PopSegment(alertRuleTemplatesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
