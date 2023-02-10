package validate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/authorization/2020-10-01/rolemanagementpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
)

func ValidateRoleManagementPolicyId(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	ids := strings.Split(v, "|")
	if len(ids) != 2 {
		errors = append(errors, fmt.Errorf("expected %q to be in format roleManagementPolicyId|roleDefinitionId", v))
		return
	}

	_, err := rolemanagementpolicies.ParseScopedRoleManagementPolicyID(ids[0])
	if err != nil {
		errors = append(errors, fmt.Errorf("expected %s to be in format roleManagementPolicyId, %+v", ids[0], err))
		return
	}

	_, err = parse.RoleDefinitionId(ids[1])
	if err != nil {
		errors = append(errors, fmt.Errorf("expected %s to be in format roleManagementPolicyId, %+v", ids[1], err))
		return
	}

	return
}
