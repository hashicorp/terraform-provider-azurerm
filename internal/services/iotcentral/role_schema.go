package iotcentral

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	iotcentralDataplane "github.com/tombuildsstuff/kermit/sdk/iotcentral/2022-10-31-preview/iotcentral"
)

func schemaRole() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"role": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"organization": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func convertToRoleAssignments(roles []interface{}) []iotcentralDataplane.RoleAssignment {
	roleAssignments := []iotcentralDataplane.RoleAssignment{}

	for _, role := range roles {
		roleMap := role.(map[string]interface{})

		roleValue := roleMap["role"].(string)
		organizationValue := roleMap["organization"].(string)

		roleAssignments = append(roleAssignments, iotcentralDataplane.RoleAssignment{
			Role:         &roleValue,
			Organization: &organizationValue,
		})
	}

	return roleAssignments
}

func convertFromRoleAssignments(roleAssignments []iotcentralDataplane.RoleAssignment) []interface{} {
	var roles []interface{}

	for _, roleAssignment := range roleAssignments {
		role := make(map[string]interface{})

		if roleAssignment.Role != nil {
			role["role"] = *roleAssignment.Role
		}

		if roleAssignment.Organization != nil {
			role["organization"] = *roleAssignment.Organization
		}

		roles = append(roles, role)
	}

	return roles
}
