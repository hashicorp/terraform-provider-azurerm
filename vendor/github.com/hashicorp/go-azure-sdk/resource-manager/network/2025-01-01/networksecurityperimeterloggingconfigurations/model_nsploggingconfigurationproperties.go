package networksecurityperimeterloggingconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NspLoggingConfigurationProperties struct {
	EnabledLogCategories *[]string `json:"enabledLogCategories,omitempty"`
	Version              *string   `json:"version,omitempty"`
}
