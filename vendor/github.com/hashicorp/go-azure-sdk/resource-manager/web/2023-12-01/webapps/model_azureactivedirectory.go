package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureActiveDirectory struct {
	Enabled           *bool                             `json:"enabled,omitempty"`
	IsAutoProvisioned *bool                             `json:"isAutoProvisioned,omitempty"`
	Login             *AzureActiveDirectoryLogin        `json:"login,omitempty"`
	Registration      *AzureActiveDirectoryRegistration `json:"registration,omitempty"`
	Validation        *AzureActiveDirectoryValidation   `json:"validation,omitempty"`
}
