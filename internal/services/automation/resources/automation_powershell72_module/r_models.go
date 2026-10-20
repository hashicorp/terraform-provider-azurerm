// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_powershell72_module

type ModuleLinkModel struct {
	Uri  string       `tfschema:"uri"`
	Hash []ModuleHash `tfschema:"hash"`
}

type AutomationPowerShell72ModuleModel struct {
	AutomationAccountID string            `tfschema:"automation_account_id"`
	Name                string            `tfschema:"name"`
	ModuleLink          []ModuleLinkModel `tfschema:"module_link"`
	Tags                map[string]any    `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &AutomationPowerShell72ModuleModel{}
}
