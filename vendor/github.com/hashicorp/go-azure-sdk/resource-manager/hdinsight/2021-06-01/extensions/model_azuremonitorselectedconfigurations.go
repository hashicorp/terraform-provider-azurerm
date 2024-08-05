package extensions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMonitorSelectedConfigurations struct {
	ConfigurationVersion *string                           `json:"configurationVersion,omitempty"`
	GlobalConfigurations *map[string]string                `json:"globalConfigurations,omitempty"`
	TableList            *[]AzureMonitorTableConfiguration `json:"tableList,omitempty"`
}
