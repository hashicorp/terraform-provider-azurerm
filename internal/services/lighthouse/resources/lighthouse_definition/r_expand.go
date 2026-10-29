// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lighthouse_definition

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2022-10-01/registrationdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandLighthouseDefinitionAuthorization(input []any) []registrationdefinitions.Authorization {
	results := make([]registrationdefinitions.Authorization, 0)
	for _, item := range input {
		v := item.(map[string]any)
		delegatedRoleDefinitionIds := pluginsdk.ExpandStringSlice(v["delegated_role_definition_ids"].(*pluginsdk.Set).List())
		result := registrationdefinitions.Authorization{
			RoleDefinitionId:           v["role_definition_id"].(string),
			PrincipalId:                v["principal_id"].(string),
			PrincipalIdDisplayName:     pointer.To(v["principal_display_name"].(string)),
			DelegatedRoleDefinitionIds: delegatedRoleDefinitionIds,
		}
		results = append(results, result)
	}
	return results
}

func expandLighthouseDefinitionPlan(input []any) *registrationdefinitions.Plan {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)
	return &registrationdefinitions.Plan{
		Name:      raw["name"].(string),
		Publisher: raw["publisher"].(string),
		Product:   raw["product"].(string),
		Version:   raw["version"].(string),
	}
}

func expandLighthouseDefinitionEligibleAuthorization(input []any) *[]registrationdefinitions.EligibleAuthorization {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]registrationdefinitions.EligibleAuthorization, 0, len(input))
	for _, item := range input {
		v := item.(map[string]any)

		result := registrationdefinitions.EligibleAuthorization{
			PrincipalId:            v["principal_id"].(string),
			RoleDefinitionId:       v["role_definition_id"].(string),
			JustInTimeAccessPolicy: expandLighthouseDefinitionJustInTimeAccessPolicy(v["just_in_time_access_policy"].([]any)),
		}

		if principalDisplayName := v["principal_display_name"].(string); principalDisplayName != "" {
			result.PrincipalIdDisplayName = pointer.To(principalDisplayName)
		}

		results = append(results, result)
	}

	return &results
}

func expandLighthouseDefinitionJustInTimeAccessPolicy(input []any) *registrationdefinitions.JustInTimeAccessPolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	justInTimeAccessPolicy := input[0].(map[string]any)

	result := registrationdefinitions.JustInTimeAccessPolicy{
		MaximumActivationDuration: pointer.To(justInTimeAccessPolicy["maximum_activation_duration"].(string)),
		ManagedByTenantApprovers:  expandLighthouseDefinitionApprover(justInTimeAccessPolicy["approver"].(*pluginsdk.Set).List()),
	}

	multiFactorAuthProvider := registrationdefinitions.MultiFactorAuthProviderNone
	if v := justInTimeAccessPolicy["multi_factor_auth_provider"].(string); v != "" {
		multiFactorAuthProvider = registrationdefinitions.MultiFactorAuthProvider(v)
	}
	result.MultiFactorAuthProvider = multiFactorAuthProvider

	return &result
}

func expandLighthouseDefinitionApprover(input []any) *[]registrationdefinitions.EligibleApprover {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]registrationdefinitions.EligibleApprover, 0)
	for _, v := range input {
		eligibleApprover := v.(map[string]any)

		result := registrationdefinitions.EligibleApprover{
			PrincipalId: eligibleApprover["principal_id"].(string),
		}

		if principalDisplayName := eligibleApprover["principal_display_name"].(string); principalDisplayName != "" {
			result.PrincipalIdDisplayName = pointer.To(principalDisplayName)
		}

		results = append(results, result)
	}

	return &results
}
