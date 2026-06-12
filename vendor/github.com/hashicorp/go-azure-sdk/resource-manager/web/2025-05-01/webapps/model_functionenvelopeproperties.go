package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FunctionEnvelopeProperties struct {
	Config             *interface{}       `json:"config,omitempty"`
	ConfigHref         *string            `json:"config_href,omitempty"`
	Files              *map[string]string `json:"files,omitempty"`
	FunctionAppId      *string            `json:"function_app_id,omitempty"`
	Href               *string            `json:"href,omitempty"`
	InvokeUrlTemplate  *string            `json:"invoke_url_template,omitempty"`
	IsDisabled         *bool              `json:"isDisabled,omitempty"`
	Language           *string            `json:"language,omitempty"`
	ScriptHref         *string            `json:"script_href,omitempty"`
	ScriptRootPathHref *string            `json:"script_root_path_href,omitempty"`
	SecretsFileHref    *string            `json:"secrets_file_href,omitempty"`
	TestData           *string            `json:"test_data,omitempty"`
	TestDataHref       *string            `json:"test_data_href,omitempty"`
}
