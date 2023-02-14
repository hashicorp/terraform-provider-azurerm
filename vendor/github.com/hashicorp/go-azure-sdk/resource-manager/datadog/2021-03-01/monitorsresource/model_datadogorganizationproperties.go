package monitorsresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatadogOrganizationProperties struct {
	ApiKey          *string `json:"apiKey,omitempty"`
	ApplicationKey  *string `json:"applicationKey,omitempty"`
	EnterpriseAppId *string `json:"enterpriseAppId,omitempty"`
	Id              *string `json:"id,omitempty"`
	LinkingAuthCode *string `json:"linkingAuthCode,omitempty"`
	LinkingClientId *string `json:"linkingClientId,omitempty"`
	Name            *string `json:"name,omitempty"`
	RedirectUri     *string `json:"redirectUri,omitempty"`
}
