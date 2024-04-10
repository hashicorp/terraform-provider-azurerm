package validate

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources" // nolint: staticcheck
)

func ResourceGroupTemplateDeploymentMode(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if ok && v == string(resources.DeploymentModeComplete) {
		warnings = append(warnings, fmt.Sprintf("If %q is set to `Complete` then resources within this Resource Group which are not defined in the ARM Template will be deleted, more info can be found at https://learn.microsoft.com/azure/azure-resource-manager/templates/deployment-modes#complete-mode", key))
		return
	}

	return
}
