package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MSDeployCore struct {
	AppOffline              *bool              `json:"appOffline,omitempty"`
	ConnectionString        *string            `json:"connectionString,omitempty"`
	DbType                  *string            `json:"dbType,omitempty"`
	PackageUri              *string            `json:"packageUri,omitempty"`
	SetParameters           *map[string]string `json:"setParameters,omitempty"`
	SetParametersXmlFileUri *string            `json:"setParametersXmlFileUri,omitempty"`
	SkipAppData             *bool              `json:"skipAppData,omitempty"`
}
