// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package static_web_app_function_app_registration

type StaticWebAppFunctionAppRegistrationModel struct {
	StaticWebAppID string `tfschema:"static_web_app_id"`
	FunctionAppID  string `tfschema:"function_app_id"`
}

func (r Resource) ModelObject() any {
	return &StaticWebAppFunctionAppRegistrationModel{}
}
