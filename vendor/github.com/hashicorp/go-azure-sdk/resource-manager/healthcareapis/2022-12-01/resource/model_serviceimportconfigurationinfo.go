package resource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceImportConfigurationInfo struct {
	Enabled              *bool   `json:"enabled,omitempty"`
	InitialImportMode    *bool   `json:"initialImportMode,omitempty"`
	IntegrationDataStore *string `json:"integrationDataStore,omitempty"`
}
