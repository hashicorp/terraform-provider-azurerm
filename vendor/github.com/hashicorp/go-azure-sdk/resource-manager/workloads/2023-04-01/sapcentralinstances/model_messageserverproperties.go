package sapcentralinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MessageServerProperties struct {
	HTTPPort       *int64          `json:"httpPort,omitempty"`
	HTTPSPort      *int64          `json:"httpsPort,omitempty"`
	Health         *SAPHealthState `json:"health,omitempty"`
	Hostname       *string         `json:"hostname,omitempty"`
	IPAddress      *string         `json:"ipAddress,omitempty"`
	InternalMsPort *int64          `json:"internalMsPort,omitempty"`
	MsPort         *int64          `json:"msPort,omitempty"`
}
