package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
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
	return strings.Join(segments, " / ")
}

func (id SubscriptionTemplateDeploymentId) ID(_ string) string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Resources/deployments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.DeploymentName)
}

// SubscriptionTemplateDeploymentID parses a SubscriptionTemplateDeployment ID into an SubscriptionTemplateDeploymentId struct
func SubscriptionTemplateDeploymentID(input string) (*SubscriptionTemplateDeploymentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
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
