package staticsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticSiteBasicAuthPropertiesARMResourceProperties struct {
	ApplicableEnvironmentsMode string    `json:"applicableEnvironmentsMode"`
	Environments               *[]string `json:"environments,omitempty"`
	Password                   *string   `json:"password,omitempty"`
	SecretState                *string   `json:"secretState,omitempty"`
	SecretURL                  *string   `json:"secretUrl,omitempty"`
}
