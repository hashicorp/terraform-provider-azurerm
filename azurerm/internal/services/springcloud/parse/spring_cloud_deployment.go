package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SpringCloudDeploymentId struct {
	SubscriptionId string
	ResourceGroup  string
	SpringName     string
	AppName        string
	DeploymentName string
}

func NewSpringCloudDeploymentID(subscriptionId, resourceGroup, springName, appName, deploymentName string) SpringCloudDeploymentId {
	return SpringCloudDeploymentId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SpringName:     springName,
		AppName:        appName,
		DeploymentName: deploymentName,
	}
}

func (id SpringCloudDeploymentId) String() string {
	segments := []string{
		fmt.Sprintf("Deployment Name %q", id.DeploymentName),
		fmt.Sprintf("App Name %q", id.AppName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Deployment", segmentsStr)
}

func (id SpringCloudDeploymentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/Spring/%s/apps/%s/deployments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.AppName, id.DeploymentName)
}

// SpringCloudDeploymentID parses a SpringCloudDeployment ID into an SpringCloudDeploymentId struct
func SpringCloudDeploymentID(input string) (*SpringCloudDeploymentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudDeploymentId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SpringName, err = id.PopSegment("Spring"); err != nil {
		return nil, err
	}
	if resourceId.AppName, err = id.PopSegment("apps"); err != nil {
		return nil, err
	}
	if resourceId.DeploymentName, err = id.PopSegment("deployments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
