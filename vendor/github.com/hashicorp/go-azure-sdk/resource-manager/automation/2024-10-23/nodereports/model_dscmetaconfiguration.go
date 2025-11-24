package nodereports

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DscMetaConfiguration struct {
	ActionAfterReboot              *string `json:"actionAfterReboot,omitempty"`
	AllowModuleOverwrite           *bool   `json:"allowModuleOverwrite,omitempty"`
	CertificateId                  *string `json:"certificateId,omitempty"`
	ConfigurationMode              *string `json:"configurationMode,omitempty"`
	ConfigurationModeFrequencyMins *int64  `json:"configurationModeFrequencyMins,omitempty"`
	RebootNodeIfNeeded             *bool   `json:"rebootNodeIfNeeded,omitempty"`
	RefreshFrequencyMins           *int64  `json:"refreshFrequencyMins,omitempty"`
}
