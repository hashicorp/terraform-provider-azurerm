// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SubscriptionTemplateDeploymentId struct {
	SubscriptionId string
	DeploymentName string
}

func NewSubscriptionTemplateDeploymentID(subscriptionId, deploymentName string) SubscriptionTemplateDeploymentId {
	return SubscriptionTemplateDeploymentId{
		SubscriptionId: subscriptionId,
		DeploymentName: deploymentName,
	}
}

func (id SubscriptionTemplateDeploymentId) String() string {
	segments := []string{
		fmt.Sprintf("Deployment Name %q", id.DeploymentName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Subscription Template Deployment", segmentsStr)
}

func (id SubscriptionTemplateDeploymentId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Resources/deployments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.DeploymentName)
}

// SubscriptionTemplateDeploymentID parses a SubscriptionTemplateDeployment ID into an SubscriptionTemplateDeploymentId struct
func SubscriptionTemplateDeploymentID(input string) (*SubscriptionTemplateDeploymentId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an SubscriptionTemplateDeployment ID: %+v", input, err)
	}

	resourceId := SubscriptionTemplateDeploymentId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.DeploymentName, err = id.PopSegment("deployments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
