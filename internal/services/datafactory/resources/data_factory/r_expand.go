// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
)

func expandDataFactoryGlobalParameters(input []any) (*map[string]factories.GlobalParameterSpecification, error) {
	result := make(map[string]factories.GlobalParameterSpecification)
	if len(input) == 0 {
		return &result, nil
	}
	for _, item := range input {
		if item == nil {
			continue
		}
		v := item.(map[string]any)

		name := v["name"].(string)
		if _, ok := v[name]; ok {
			return nil, fmt.Errorf("duplicate parameter name")
		}

		result[name] = factories.GlobalParameterSpecification{
			Type:  factories.GlobalParameterType(v["type"].(string)),
			Value: v["value"].(string),
		}
	}
	return &result, nil
}

func expandGitHubRepoConfiguration(input []any) *factories.FactoryGitHubConfiguration {
	if len(input) == 0 {
		return nil
	}

	item := input[0].(map[string]any)
	return &factories.FactoryGitHubConfiguration{
		AccountName:         item["account_name"].(string),
		CollaborationBranch: item["branch_name"].(string),
		DisablePublish:      pointer.To(!item["publishing_enabled"].(bool)),
		HostName:            pointer.To(item["git_url"].(string)),
		RepositoryName:      item["repository_name"].(string),
		RootFolder:          item["root_folder"].(string),
	}
}

func expandVSTSRepoConfiguration(input []any) *factories.FactoryVSTSConfiguration {
	if len(input) == 0 {
		return nil
	}

	item := input[0].(map[string]any)
	return &factories.FactoryVSTSConfiguration{
		AccountName:         item["account_name"].(string),
		CollaborationBranch: item["branch_name"].(string),
		DisablePublish:      pointer.To(!item["publishing_enabled"].(bool)),
		ProjectName:         item["project_name"].(string),
		RepositoryName:      item["repository_name"].(string),
		RootFolder:          item["root_folder"].(string),
		TenantId:            pointer.To(item["tenant_id"].(string)),
	}
}
