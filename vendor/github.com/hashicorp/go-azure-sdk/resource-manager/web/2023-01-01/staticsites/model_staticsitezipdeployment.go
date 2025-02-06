package staticsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticSiteZipDeployment struct {
	ApiZipURL        *string `json:"apiZipUrl,omitempty"`
	AppZipURL        *string `json:"appZipUrl,omitempty"`
	DeploymentTitle  *string `json:"deploymentTitle,omitempty"`
	FunctionLanguage *string `json:"functionLanguage,omitempty"`
	Provider         *string `json:"provider,omitempty"`
}
