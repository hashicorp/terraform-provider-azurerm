package guestconfigurationassignments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationSetting struct {
	ActionAfterReboot              *ActionAfterReboot `json:"actionAfterReboot,omitempty"`
	AllowModuleOverwrite           *bool              `json:"allowModuleOverwrite,omitempty"`
	ConfigurationMode              *ConfigurationMode `json:"configurationMode,omitempty"`
	ConfigurationModeFrequencyMins *float64           `json:"configurationModeFrequencyMins,omitempty"`
	RebootIfNeeded                 *bool              `json:"rebootIfNeeded,omitempty"`
	RefreshFrequencyMins           *float64           `json:"refreshFrequencyMins,omitempty"`
}
