package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ManagementGroupTemplateDeploymentId struct {
	ManagementGroupName string
	DeploymentName      string
}

func NewManagementGroupTemplateDeploymentID(managementGroupName, deploymentName string) ManagementGroupTemplateDeploymentId {
	return ManagementGroupTemplateDeploymentId{
		ManagementGroupName: managementGroupName,
		DeploymentName:      deploymentName,
	}
}

func (id ManagementGroupTemplateDeploymentId) String() string {
	segments := []string{
		fmt.Sprintf("Deployment Name %q", id.DeploymentName),
		fmt.Sprintf("Management Group Name %q", id.ManagementGroupName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Management Group Template Deployment", segmentsStr)
}

func (id ManagementGroupTemplateDeploymentId) ID() string {
	fmtString := "/providers/Microsoft.Management/managementGroups/%s/providers/Microsoft.Resources/deployments/%s"
	return fmt.Sprintf(fmtString, id.ManagementGroupName, id.DeploymentName)
}

// ManagementGroupTemplateDeploymentID parses a ManagementGroupTemplateDeployment ID into an ManagementGroupTemplateDeploymentId struct
func ManagementGroupTemplateDeploymentID(input string) (*ManagementGroupTemplateDeploymentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ManagementGroupTemplateDeploymentId{}

	if resourceId.ManagementGroupName, err = id.PopSegment("managementGroups"); err != nil {
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
