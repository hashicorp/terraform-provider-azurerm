// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lighthouse_definition

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2022-10-01/registrationdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenLighthouseDefinitionAuthorization(input []registrationdefinitions.Authorization) []any {
	results := make([]any, 0)
	for _, item := range input {
		results = append(results, map[string]any{
			"role_definition_id":            item.RoleDefinitionId,
			"principal_id":                  item.PrincipalId,
			"principal_display_name":        pointer.From(item.PrincipalIdDisplayName),
			"delegated_role_definition_ids": pluginsdk.FlattenSlice(item.DelegatedRoleDefinitionIds),
		})
	}

	return results
}

func flattenLighthouseDefinitionPlan(input *registrationdefinitions.Plan) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"name":      input.Name,
			"publisher": input.Publisher,
			"product":   input.Product,
			"version":   input.Version,
		},
	}
}

func flattenLighthouseDefinitionEligibleAuthorization(input *[]registrationdefinitions.EligibleAuthorization) []any {
	if input == nil {
		return []any{}
	}

	results := make([]any, 0, len(*input))
	for _, item := range *input {
		result := map[string]any{
			"principal_id":       item.PrincipalId,
			"role_definition_id": item.RoleDefinitionId,
		}

		if item.JustInTimeAccessPolicy != nil {
			result["just_in_time_access_policy"] = flattenLighthouseDefinitionJustInTimeAccessPolicy(item.JustInTimeAccessPolicy)
		}

		if item.PrincipalIdDisplayName != nil {
			result["principal_display_name"] = *item.PrincipalIdDisplayName
		}

		results = append(results, result)
	}

	return results
}

func flattenLighthouseDefinitionJustInTimeAccessPolicy(input *registrationdefinitions.JustInTimeAccessPolicy) []any {
	if input == nil {
		return []any{}
	}

	var results []any

	result := map[string]any{}

	if v := input.MultiFactorAuthProvider; v != registrationdefinitions.MultiFactorAuthProviderNone {
		result["multi_factor_auth_provider"] = string(v)
	}

	if input.ManagedByTenantApprovers != nil {
		result["approver"] = flattenLighthouseDefinitionApprover(input.ManagedByTenantApprovers)
	}

	maximumActivationDuration := "PT8H"
	if input.MaximumActivationDuration != nil {
		maximumActivationDuration = *input.MaximumActivationDuration
	}
	result["maximum_activation_duration"] = maximumActivationDuration

	return append(results, result)
}

func flattenLighthouseDefinitionApprover(input *[]registrationdefinitions.EligibleApprover) []any {
	if input == nil {
		return []any{}
	}

	results := make([]any, 0, len(*input))
	for _, item := range *input {
		result := map[string]any{
			"principal_id": item.PrincipalId,
		}

		if item.PrincipalIdDisplayName != nil {
			result["principal_display_name"] = *item.PrincipalIdDisplayName
		}

		results = append(results, result)
	}

	return results
}
