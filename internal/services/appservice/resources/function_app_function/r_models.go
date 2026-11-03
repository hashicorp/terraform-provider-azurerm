// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package function_app_function

type FunctionAppFunctionModel struct {
	Name           string          `tfschema:"name"`
	AppID          string          `tfschema:"function_app_id"`
	Enabled        bool            `tfschema:"enabled"`
	ConfigJSON     string          `tfschema:"config_json"`
	Language       string          `tfschema:"language"`
	SecretsFileURL string          `tfschema:"secrets_file_url"`
	TestData       string          `tfschema:"test_data"`
	Files          []FunctionFiles `tfschema:"file"`

	ConfigURL         string `tfschema:"config_url"`
	FunctionURL       string `tfschema:"url"`
	InvokeURL         string `tfschema:"invocation_url"`
	ScriptURL         string `tfschema:"script_url"`
	ScriptRootPathURL string `tfschema:"script_root_path_url"`
	TestDataURL       string `tfschema:"test_data_url"`
}

func (r Resource) ModelObject() any {
	return &FunctionAppFunctionModel{}
}
