// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_runtime_environment

type AutomationRuntimeEnvironmentResourceModel struct {
	Name                   string            `tfschema:"name"`
	AutomationAccountId    string            `tfschema:"automation_account_id"`
	RuntimeLanguage        string            `tfschema:"runtime_language"`
	RuntimeVersion         string            `tfschema:"runtime_version"`
	RuntimeDefaultPackages map[string]string `tfschema:"runtime_default_packages"`
	Location               string            `tfschema:"location"`
	Description            string            `tfschema:"description"`
	Tags                   map[string]string `tfschema:"tags"`
}

func (m Resource) ModelObject() any {
	return &AutomationRuntimeEnvironmentResourceModel{}
}
