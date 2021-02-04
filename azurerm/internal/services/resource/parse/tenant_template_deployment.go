package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type TenantTemplateDeploymentId struct {
	DeploymentName string
}

func NewTenantTemplateDeploymentID(deploymentName string) TenantTemplateDeploymentId {
	return TenantTemplateDeploymentId{
		DeploymentName: deploymentName,
	}
}

func (id TenantTemplateDeploymentId) String() string {
	segments := []string{
		fmt.Sprintf("Deployment Name %q", id.DeploymentName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Tenant Template Deployment", segmentsStr)
}

func (id TenantTemplateDeploymentId) ID() string {
	fmtString := "/providers/Microsoft.Resources/deployments/%s"
	return fmt.Sprintf(fmtString, id.DeploymentName)
}

// TenantTemplateDeploymentID parses a TenantTemplateDeployment ID into an TenantTemplateDeploymentId struct
func TenantTemplateDeploymentID(input string) (*TenantTemplateDeploymentId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := TenantTemplateDeploymentId{}

	if resourceId.DeploymentName, err = id.PopSegment("deployments"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
