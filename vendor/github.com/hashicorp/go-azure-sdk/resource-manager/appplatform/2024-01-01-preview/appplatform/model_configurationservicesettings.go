package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationServiceSettings struct {
	GitProperty              *ConfigurationServiceGitProperty `json:"gitProperty,omitempty"`
	RefreshIntervalInSeconds *int64                           `json:"refreshIntervalInSeconds,omitempty"`
}
