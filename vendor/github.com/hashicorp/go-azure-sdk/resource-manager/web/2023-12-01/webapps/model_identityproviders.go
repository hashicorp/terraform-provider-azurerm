package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IdentityProviders struct {
	Apple                        *Apple                                  `json:"apple,omitempty"`
	AzureActiveDirectory         *AzureActiveDirectory                   `json:"azureActiveDirectory,omitempty"`
	AzureStaticWebApps           *AzureStaticWebApps                     `json:"azureStaticWebApps,omitempty"`
	CustomOpenIdConnectProviders *map[string]CustomOpenIdConnectProvider `json:"customOpenIdConnectProviders,omitempty"`
	Facebook                     *Facebook                               `json:"facebook,omitempty"`
	GitHub                       *GitHub                                 `json:"gitHub,omitempty"`
	Google                       *Google                                 `json:"google,omitempty"`
	LegacyMicrosoftAccount       *LegacyMicrosoftAccount                 `json:"legacyMicrosoftAccount,omitempty"`
	Twitter                      *Twitter                                `json:"twitter,omitempty"`
}
