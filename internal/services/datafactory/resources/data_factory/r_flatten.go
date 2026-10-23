// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
)

func flattenDataFactoryGlobalParameters(input *map[string]factories.GlobalParameterSpecification) (*[]any, error) {
	output := make([]any, 0)
	if input == nil || len(*input) == 0 {
		return &output, nil
	}

	for name, item := range *input {
		var valueResult string
		_, valueIsString := item.Value.(string)
		if (item.Type == factories.GlobalParameterTypeArray || item.Type == factories.GlobalParameterTypeObject) && !valueIsString {
			bytes, err := json.Marshal(item.Value)
			if err != nil {
				return nil, fmt.Errorf("marshalling value for global parameter %q (value %+v): %+v", name, item.Value, err)
			}
			valueResult = string(bytes)
		} else {
			valueResult = fmt.Sprintf("%v", item.Value)
		}

		output = append(output, map[string]any{
			"name":  name,
			"type":  string(item.Type),
			"value": valueResult,
		})
	}
	return &output, nil
}

func flattenGitHubRepoConfiguration(input factories.FactoryRepoConfiguration) []any {
	output := make([]any, 0)

	if v, ok := input.(factories.FactoryGitHubConfiguration); ok {
		publishingEnabled := true
		if v.DisablePublish != nil {
			publishingEnabled = !*v.DisablePublish
		}
		output = append(output, map[string]any{
			"account_name":       v.AccountName,
			"branch_name":        v.CollaborationBranch,
			"git_url":            pointer.From(v.HostName),
			"publishing_enabled": publishingEnabled,
			"repository_name":    v.RepositoryName,
			"root_folder":        v.RootFolder,
		})
	}

	return output
}

func flattenVSTSRepoConfiguration(input factories.FactoryRepoConfiguration) []any {
	output := make([]any, 0)

	if v, ok := input.(factories.FactoryVSTSConfiguration); ok {
		publishingEnabled := true
		if v.DisablePublish != nil {
			publishingEnabled = !*v.DisablePublish
		}
		output = append(output, map[string]any{
			"account_name":       v.AccountName,
			"branch_name":        v.CollaborationBranch,
			"project_name":       v.ProjectName,
			"publishing_enabled": publishingEnabled,
			"repository_name":    v.RepositoryName,
			"root_folder":        v.RootFolder,
			"tenant_id":          pointer.From(v.TenantId),
		})
	}

	return output
}
