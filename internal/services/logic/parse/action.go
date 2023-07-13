// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ActionId struct {
	SubscriptionId string
	ResourceGroup  string
	WorkflowName   string
	Name           string
}

func NewActionID(subscriptionId, resourceGroup, workflowName, name string) ActionId {
	return ActionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		WorkflowName:   workflowName,
		Name:           name,
	}
}

func (id ActionId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Workflow Name %q", id.WorkflowName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Action", segmentsStr)
}

func (id ActionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/workflows/%s/actions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkflowName, id.Name)
}

// ActionID parses a Action ID into an ActionId struct
func ActionID(input string) (*ActionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an Action ID: %+v", input, err)
	}

	resourceId := ActionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.WorkflowName, err = id.PopSegment("workflows"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("actions"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
