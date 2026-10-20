// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_connection_type

type AutomationConnectionTypeModel struct {
	ResourceGroup         string  `json:"resource_group"          tfschema:"resource_group_name"`
	AutomationAccountName string  `json:"automation_account_name" tfschema:"automation_account_name"`
	Name                  string  `json:"name"                    tfschema:"name"`
	IsGlobal              bool    `json:"is_global"               tfschema:"is_global"`
	Field                 []Field `json:"field"                   tfschema:"field"`
}

func (m Resource) ModelObject() any {
	return &AutomationConnectionTypeModel{}
}
